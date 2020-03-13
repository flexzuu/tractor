package agent

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/manifold/tractor/pkg/agent/console"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/workspace/supervisor"
)

type Workspace struct {
	*supervisor.Supervisor

	Name        string // base name of dir (~/.tractor/workspaces/{name})
	SymlinkPath string // absolute path to symlink file (~/.tractor/workspaces/{name})
	TargetPath  string // absolute path to target of symlink (actual workspace)

	log         logging.Logger
	consolePipe io.WriteCloser
}

func OpenWorkspace(a *Agent, name string) (*Workspace, error) {
	symlinkPath := filepath.Join(a.WorkspacesPath, name)
	binPath := filepath.Join(a.WorkspaceBinPath, name)
	targetPath, err := os.Readlink(symlinkPath)
	if err != nil {
		return nil, err
	}
	var consolePipe io.WriteCloser
	if svc, ok := a.Logger.(*console.Service); ok && svc != nil {
		consolePipe = svc.NewPipe("@" + name)
	}
	svr := supervisor.New(targetPath, name, consolePipe)
	svr.Log = a.Logger
	svr.GoBin = a.GoBin
	svr.DaemonBin = binPath
	return &Workspace{
		Supervisor: svr,

		Name:        name,
		SymlinkPath: symlinkPath,
		TargetPath:  targetPath,
	}, nil
}

// Start starts the workspace daemon. creates the symlink to the path if it does
// not exist, using the path basename as the symlink name
func (w *Workspace) Start() error {
	logging.Infof(w.log, "@%s: Start()", w.Name)
	return w.Supervisor.Daemon.Restart()
}

// Stop stops the workspace daemon, deleting the unix socket file.
func (w *Workspace) Stop() error {
	logging.Infof(w.log, "@%s: Stop()", w.Name)
	if w.Supervisor.Daemon != nil {
		return w.Supervisor.Daemon.Stop()
	}
	return nil
}

func (w *Workspace) Serve(ctx context.Context) {
	w.Supervisor.Serve(ctx)
}

func (w *Workspace) Status() supervisor.Status {
	return w.Supervisor.Status()
}

func (w *Workspace) Observe(cb supervisor.Observer) {
	w.Supervisor.Observe(cb)
}
