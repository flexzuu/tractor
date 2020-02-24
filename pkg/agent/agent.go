package agent

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sync"
	"time"
	"unsafe"

	"github.com/fsnotify/fsnotify"
	"github.com/manifold/tractor/pkg/agent/console"
	"github.com/manifold/tractor/pkg/config"
	"github.com/manifold/tractor/pkg/misc/daemon"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/misc/logging/null"
)

// Agent manages multiple workspaces in a directory (default: ~/.tractor).
type Agent struct {
	Path                 string `toml:"-"` // ~/.tractor
	ConfigPath           string `toml:"-"` // ~/.tractor/config.toml
	SocketPath           string // ~/.tractor/agent.sock
	WorkspacesPath       string // ~/.tractor/workspaces
	WorkspaceSocketsPath string // ~/.tractor/sockets
	WorkspaceBinPath     string // ~/.tractor/bin
	GoBin                string
	PreferredBrowser     string
	DevMode              bool

	Daemon  *daemon.Daemon   `toml:"-"`
	Console *console.Service `toml:"-"`
	Logger  logging.Logger   `toml:"-"`

	WorkspacesChanged chan struct{} `toml:"-"`
	workspaces        map[string]*Workspace
	mu                sync.RWMutex
}

// Open returns a new agent for the given path. If the given path is empty, a
// default of ~/.tractor will be used.
func Open(path string, console *console.Service, devMode bool) (*Agent, error) {
	bin, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}

	a := &Agent{
		DevMode:           devMode,
		Console:           console,
		Logger:            console,
		Path:              path,
		GoBin:             bin,
		workspaces:        make(map[string]*Workspace),
		WorkspacesChanged: make(chan struct{}),
	}

	if len(a.Path) == 0 {
		p, err := defaultPath()
		if err != nil {
			return nil, err
		}
		a.Path = p
	}

	a.SocketPath = filepath.Join(a.Path, "agent.sock")
	a.ConfigPath = filepath.Join(a.Path, "config.toml")
	a.WorkspacesPath = filepath.Join(a.Path, "workspaces")
	a.WorkspaceBinPath = filepath.Join(a.Path, "bin")
	a.WorkspaceSocketsPath = filepath.Join(a.Path, "sockets")
	if a.Logger == nil {
		a.Logger = &null.Logger{}
	}

	cfg, err := config.ParseFile(a.ConfigPath)
	if err != nil {
		return nil, err
	}

	a.PreferredBrowser = cfg.Agent.PreferredBrowser

	os.MkdirAll(a.WorkspacesPath, 0700)
	os.MkdirAll(a.WorkspaceSocketsPath, 0700)
	os.MkdirAll(a.WorkspaceBinPath, 0700)

	return a, nil
}

func (a *Agent) InitializeDaemon() (err error) {
	spaces, err := a.Workspaces()
	if err != nil {
		return err
	}
	for _, space := range spaces {
		a.Daemon.AddServices(space)
	}
	return nil
}

func (a *Agent) Serve(ctx context.Context) {
	a.Watch(ctx)
}

// Workspace returns a Workspace for the given path. The path must match
// either:
//   * the workspace symlink's basename in the agent's WorkspacesPath.
//   * the full path to the target of a workspace symlink in WorkspacesPath.
//   * full path to the workspace anywhere else. it will be symlinked to
//     the Workspaces path using the basename of the full path.
func (a *Agent) Workspace(path string) *Workspace {
	// check to see if the workspace is cached
	// cached=workspace is running through an agent QRPC call, or showing the
	// workspace in the systray.
	a.mu.RLock()
	ws := a.workspaces[path]
	a.mu.RUnlock()
	if ws != nil {
		return ws
	}

	// now look for a symlink in ~/.tractor/workspaces
	wss, err := a.Workspaces()
	if err != nil {
		panic(err)
	}
	for _, ws := range wss {
		if ws.Name == path || ws.TargetPath == path {
			return ws
		}
	}

	// if full path is a dir with workspace.go, symlink it
	basename, err := a.symlinkWorkspace(path)
	if err != nil {
		return nil
	}

	return a.Workspace(basename)
}

func (a *Agent) symlinkWorkspace(path string) (string, error) {
	fi, err := os.Lstat(filepath.Join(path, "workspace.go"))
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		return "", nil
	}

	basepath := filepath.Base(path)
	base := basepath
	i := 1
	for {
		err = os.Symlink(path, filepath.Join(a.WorkspacesPath, base))
		if err != nil && !os.IsExist(err) {
			return base, err
		}

		if err == nil {
			return base, nil
		}

		i++
		base = fmt.Sprintf("%s-%d", basepath, i)
	}
}

// Workspaces returns the workspaces under this agent's WorkspacesPath.
func (a *Agent) Workspaces() ([]*Workspace, error) {
	entries, err := ioutil.ReadDir(a.WorkspacesPath)
	if err != nil {
		return nil, err
	}

	workspaces := make([]*Workspace, 0, len(entries))
	a.mu.Lock()
	for _, entry := range entries {
		if !a.isWorkspaceDir(entry) {
			continue
		}

		n := entry.Name()
		ws := a.workspaces[n]
		if ws == nil {
			ws, err = OpenWorkspace(a, n)
			if err != nil {
				return nil, err
			}
			a.workspaces[n] = ws
		}
		workspaces = append(workspaces, ws)
	}
	a.mu.Unlock()
	return workspaces, nil
}

// Shutdown shuts all workspaces down and cleans up socket files.
func (a *Agent) Shutdown() {
	info(a.Logger, "[server] shutting down")
	os.RemoveAll(a.SocketPath)
	for _, ws := range a.workspaces {
		ws.Stop()
	}
}

func (a *Agent) Watch(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		info(a.Logger, "unable to create watcher:", err)
		return
	}
	watcher.Add(a.WorkspacesPath)
	debounce := Debounce(20 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			watcher.Close()
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}

			debounce(func() {
				a.WorkspacesChanged <- struct{}{}
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logErr(a.Logger, "watcher error:", err)
		}
	}
}

func (a *Agent) isWorkspaceDir(fi os.FileInfo) bool {
	if fi.IsDir() {
		return true
	}

	path := filepath.Join(a.WorkspacesPath, fi.Name())
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		logErr(a.Logger, err)
		return false
	}

	if resolved == path {
		return false
	}

	rfi, err := os.Lstat(resolved)
	if err != nil {
		logErr(a.Logger, err)
		return false
	}

	return rfi.IsDir()
}

func defaultPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, ".tractor"), nil
}

func info(l logging.Logger, args ...interface{}) {
	if !isNilValue(l) {
		l.Info(args...)
	}
}

func logErr(l logging.Logger, args ...interface{}) {
	if !isNilValue(l) {
		l.Error(args...)
	}
}

func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}

// New returns a debounced function that takes another functions as its argument.
// This function will be called when the debounced function stops being called
// for the given duration.
// The debounced function can be invoked with different functions, if needed,
// the last one will win.
func Debounce(after time.Duration) func(f func()) {
	d := &debouncer{after: after}

	return func(f func()) {
		d.add(f)
	}
}

type debouncer struct {
	mu    sync.Mutex
	after time.Duration
	timer *time.Timer
}

func (d *debouncer) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, f)
}
