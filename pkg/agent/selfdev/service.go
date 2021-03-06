package selfdev

import (
	"bytes"
	"context"
	"crypto/md5"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/manifold/tractor/pkg/agent/console"
	"github.com/manifold/tractor/pkg/config"
	"github.com/manifold/tractor/pkg/misc/daemon"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/misc/subcmd"
	"github.com/radovskyb/watcher"
)

type Service struct {
	Daemon  *daemon.Daemon
	Logger  logging.Logger
	Console *console.Service
	Config  *config.Config

	watcher *watcher.Watcher
	output  io.WriteCloser
}

func (s *Service) InitializeDaemon() (err error) {
	s.output = s.Console.NewPipe("selfdev")
	s.watcher = watcher.New()
	s.watcher.SetMaxEvents(1)
	s.watcher.IgnoreHiddenFiles(true)
	s.watcher.AddFilterHook(func(info os.FileInfo, fullPath string) error {
		allowedExt := []string{".go", ".ts", ".tsx", ".js", ".jsx", ".html"}
		ignoreSubstr := []string{"node_modules"}
		for _, substr := range ignoreSubstr {
			if strings.Contains(fullPath, substr) {
				return watcher.ErrSkip
			}
		}
		for _, ext := range allowedExt {
			if filepath.Ext(info.Name()) == ext {
				return nil
			}
		}
		return watcher.ErrSkip
	})

	s.watcher.AddRecursive("./pkg")
	s.watcher.AddRecursive("./studio")

	theiaOut := s.Console.NewPipe("theia")
	shellBuild := &cmdService{
		subcmd.New("yarn", "run", "theia", "build", "--watch", "--mode", "development"),
	}
	shellBuild.Setup = func(cmd *exec.Cmd) error {
		cmd.Dir = "./studio/shell"
		// theia watch barfs a lot of useless warnings with every change
		//cmd.Stdout = s.output
		cmd.Stderr = theiaOut
		return nil
	}
	s.Daemon.AddServices(shellBuild)

	studioOut := s.Console.NewPipe("studio")
	shellRun := &cmdService{
		subcmd.New("yarn", "run", "theia", "start", "--port", "11010", "--log-level", "warn"),
	}
	shellRun.Setup = func(cmd *exec.Cmd) error {
		cmd.Dir = "./studio/shell"
		cmd.Stdout = studioOut
		cmd.Stderr = studioOut
		return nil
	}
	s.Daemon.AddServices(shellRun)

	return err
}

func (s *Service) TerminateDaemon() error {
	s.watcher.Close()
	return nil
}

func (s *Service) Serve(ctx context.Context) {
	go s.handleLoop(ctx)
	if err := s.watcher.Start(s.Config.DevWatchInterval()); err != nil {
		logging.Error(s.Logger, err)
	}
}
func (s *Service) handleLoop(ctx context.Context) {
	// debounce := Debounce(20 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-s.watcher.Event:
			if !ok {
				return
			}
			if event.Op&watcher.Chmod == watcher.Chmod {
				continue
			}

			if filepath.Ext(event.Path) == ".ts" || filepath.Ext(event.Path) == ".tsx" {

				logging.Info(s.Logger, "ts file changed, compiling...")
				for _, ext := range []string{"tractor", "editorview"} {
					if strings.Contains(event.Path, "/studio/extensions/"+ext) {
						go func(ext string) {
							// theia extension
							cmd := exec.Command("yarn", "build")
							cmd.Dir = "studio/extensions/" + ext
							cmd.Stdout = s.output
							cmd.Stderr = s.output
							cmd.Run()
							logging.Info(s.Logger, "finished "+ext)
						}(ext)
					}
				}

			}

			if filepath.Ext(event.Path) == ".go" {
				logging.Info(s.Logger, "go file changed, testing/compiling...")
				errs := make(chan error)
				go func() {
					cmd := exec.Command("go", "build", "-o", "./local/bin/tractor.tmp", "./cmd/tractor")
					cmd.Stdout = s.output
					cmd.Stderr = s.output
					err := cmd.Run()
					errs <- err
					if exitStatus(err) > 0 {
						logging.Info(s.Logger, "ERROR")
					}
				}()
				go func() {
					cmd := exec.Command("go", "build", "-o", "./local/bin/tractor-agent.tmp", "./cmd/tractor-agent")
					cmd.Stdout = s.output
					cmd.Stderr = s.output
					err := cmd.Run()
					errs <- err
					if exitStatus(err) > 0 {
						logging.Info(s.Logger, "ERROR")
					}
				}()
				go func() {
					cmd := exec.Command("go", "test", "-race", "./pkg/...")
					cmd.Stdout = s.output
					cmd.Stderr = s.output
					err := cmd.Run()
					errs <- err
					if exitStatus(err) > 0 {
						logging.Info(s.Logger, "ERROR")
					}
				}()
				go func() {
					for i := 0; i < 3; i++ {
						if err := <-errs; err != nil {
							os.Remove("./local/bin/tractor-agent.tmp")
							os.Remove("./local/bin/tractor.tmp")
							return
						}
					}
					os.Rename("./local/bin/tractor.tmp", "./local/bin/tractor")

					// NOTE: this is useless since go doesn't make deterministic builds.
					// 		 just a reminder maybe someday we can restart more intelligently.
					if !checksumMatch("./local/bin/tractor-agent.tmp", "./local/bin/tractor-agent") {
						os.Rename("./local/bin/tractor-agent.tmp", "./local/bin/tractor-agent")
						s.Daemon.OnFinished = func() {
							err := syscall.Exec(os.Args[0], os.Args, os.Environ())
							if err != nil {
								panic(err)
							}
						}
						s.Daemon.Terminate()
					} else {
						os.Remove("./local/bin/tractor-agent.tmp")
					}

				}()
			}

		case err, ok := <-s.watcher.Error:
			if !ok {
				return
			}
			logging.Error(s.Logger, err)
		case <-s.watcher.Closed:
			return
		}
	}
}

func exitStatus(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func checksumMatch(bin1, bin2 string) bool {
	checksum := func(path string, ch chan []byte) {
		b, err := os.Open(path)
		if err != nil {
			return
		}
		defer b.Close()
		h := md5.New()
		io.Copy(h, b)
		ch <- h.Sum(nil)
	}
	chk1 := make(chan []byte)
	chk2 := make(chan []byte)
	go checksum(bin1, chk1)
	go checksum(bin2, chk2)
	return bytes.Equal(<-chk1, <-chk2)
}

type cmdService struct {
	*subcmd.Subcmd
}

func (s *cmdService) Serve(ctx context.Context) {
	s.Start()
	s.Wait()
}

func (s *cmdService) TerminateDaemon() error {
	return s.Stop()
}
