package supervisor

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/manifold/tractor/pkg/data/icons"
	"github.com/manifold/tractor/pkg/misc/debouncer"
	L "github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/misc/subcmd"
	"github.com/radovskyb/watcher"
)

type Status string

const (
	StatusAvailable   Status = "Available"
	StatusPartially   Status = "Partially"
	StatusUnavailable Status = "Unavailable"

	WatchInterval = 50 * time.Millisecond
)

func (s Status) Icon() []byte {
	switch s {
	case StatusAvailable:
		return icons.Available
	case StatusPartially:
		return icons.Partially
	default:
		return icons.Unavailable
	}
}

func (s Status) String() string {
	return string(s)
}

type Observer func(*Supervisor, Status)

func watcherFilter(info os.FileInfo, fullPath string) error {
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
}

// inputs: path, bin path, console writer, logger, watch interval
// roles: manage process, alert observers, watch for changes to recompile and restart process
type Supervisor struct {
	Output     io.WriteCloser
	Log        L.Logger
	DaemonArgs []string
	DaemonBin  string // absolute path to compiled binary (~/.tractor/bin/{name})
	GoBin      string

	name string
	path string

	status Status
	statMu sync.Mutex

	observers []Observer
	obsMu     sync.Mutex

	Daemon   *subcmd.Subcmd
	starting sync.Mutex

	watcher *watcher.Watcher
}

func New(path string, name string, output io.WriteCloser) *Supervisor {
	if name == "" {
		name = filepath.Base(path)
	}
	return &Supervisor{
		name:       name,
		path:       path,
		status:     StatusPartially,
		Output:     output,
		GoBin:      "go",
		DaemonBin:  filepath.Join(path, name),
		DaemonArgs: []string{},
	}
}

func (s *Supervisor) Path() string {
	return s.path
}

func (s *Supervisor) Name() string {
	return s.name
}

func (s *Supervisor) Status() Status {
	s.statMu.Lock()
	defer s.statMu.Unlock()
	return s.status
}

func (s *Supervisor) Recompile() error {
	var cmd *exec.Cmd
	if s.DaemonBin != "" {
		cmd = exec.Command(s.GoBin, "build", "-o", s.DaemonBin, ".")
	} else {
		cmd = exec.Command(s.GoBin, "build", ".")
	}
	cmd.Dir = s.path
	cmd.Stdout = s.Output
	cmd.Stderr = s.Output
	return cmd.Run()
}

func (s *Supervisor) StartDaemon() error {
	if s.Daemon != nil {
		return errors.New("daemon already started")
	}

	if err := s.Recompile(); err != nil {
		return err
	}

	s.Daemon = subcmd.New(s.DaemonBin, s.DaemonArgs...)
	//s.Daemon.Log = s.Log
	s.Daemon.Setup = func(cmd *exec.Cmd) error {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Dir = s.path
		cmd.Stdout = s.Output
		cmd.Stderr = s.Output
		cmd.StdinPipe()
		return nil
	}

	s.Daemon.Observe(func(cmd *subcmd.Subcmd, status subcmd.Status) {
		switch status {
		case subcmd.StatusStarted:
			s.setStatus(StatusAvailable)
		case subcmd.StatusExited:
			if cmd.Error() != nil {
				s.Log.Info("exit error:", cmd.Error())
				s.setStatus(StatusUnavailable)
			} else {
				s.setStatus(StatusPartially)
			}
		case subcmd.StatusStopped:
			s.setStatus(StatusUnavailable)
		}
	})

	if err := s.Daemon.Start(); err != nil {
		s.setStatus(StatusUnavailable)
		return err
	}

	return nil
}

func (s *Supervisor) Serve(ctx context.Context) {
	if err := s.StartDaemon(); err != nil {
		panic(err)
	}

	s.watcher = watcher.New()
	s.watcher.IgnoreHiddenFiles(true)
	s.watcher.AddFilterHook(watcherFilter)

	if err := s.watcher.AddRecursive(s.path); err != nil {
		s.Log.Info("unable to watch path:", s.path, err)
		return
	}

	go s.handleWatcher(ctx)

	if err := s.watcher.Start(WatchInterval); err != nil {
		s.Log.Error("watcher error:", err)
	}
}

func (s *Supervisor) handleWatcher(ctx context.Context) {
	debounce := debouncer.New(WatchInterval)
	for {
		select {
		case <-ctx.Done():
			if err := s.Daemon.Stop(); err != nil {
				if err != subcmd.ErrNotRunning {
					s.Log.Error(err)
				}
			}
			s.watcher.Close()
			return
		case event, ok := <-s.watcher.Event:
			if !ok {
				return
			}
			if event.Op&watcher.Chmod == watcher.Chmod {
				continue
			}

			debounce(func() {
				s.Log.Infof("@%s: %s changed", s.name, event.Path)
				s.Log.Infof("@%s: recompiling", s.name)
				if err := s.Recompile(); err != nil {
					s.Log.Info(err)
					return
				}
				s.Log.Infof("@%s: reloading", s.name)
				if err := s.Daemon.Restart(); err != nil {
					s.Log.Info(err)
				}
			})

		case err, ok := <-s.watcher.Error:
			if !ok {
				return
			}
			s.Log.Error("watcher error:", err)
		case <-s.watcher.Closed:
			return
		}
	}
}

func (s *Supervisor) Observe(cb Observer) {
	s.obsMu.Lock()
	s.observers = append(s.observers, cb)
	s.obsMu.Unlock()
}

func (s *Supervisor) setStatus(status Status) {
	s.statMu.Lock()
	if s.status == status {
		s.statMu.Unlock()
		return
	}
	L.Infof(s.Log, "@%s: %s => %s", s.name, s.status, status)

	s.status = status
	s.obsMu.Lock()
	for _, cb := range s.observers {
		cb(s, status)
	}
	s.obsMu.Unlock()
	s.statMu.Unlock()
}
