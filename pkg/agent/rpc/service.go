package rpc

import (
	"context"
	"fmt"
	"os"

	"github.com/manifold/qtalk/golang/mux"
	qrpc "github.com/manifold/qtalk/golang/rpc"
	"github.com/manifold/tractor/pkg/agent"
	"github.com/manifold/tractor/pkg/misc/logging"
)

// Service provides a QRPC server to connect, restart, and stop running
// workspaces.
type Service struct {
	Agent *agent.Agent
	Log   logging.Logger
	api   qrpc.API
	l     mux.Listener
}

func (s *Service) InitializeDaemon() (err error) {
	if s.l, err = mux.ListenUnix(s.Agent.SocketPath); err != nil {
		return err
	}

	s.api = qrpc.NewAPI()
	s.api.HandleFunc("start", s.Start())
	s.api.HandleFunc("stop", s.Stop())
	return nil
}

func (s *Service) Serve(ctx context.Context) {
	server := &qrpc.Server{}

	logging.Infof(s.Log, "agent listening at unix://%s", s.Agent.SocketPath)
	if err := server.Serve(s.l, s.api); err != nil {
		fmt.Println(err)
	}
	os.Remove(s.Agent.SocketPath)
}

func (s *Service) TerminateDaemon() error {
	s.Agent.Shutdown()
	os.Remove(s.Agent.SocketPath)
	return nil
}

func findWorkspace(a *agent.Agent, call *qrpc.Call) (*agent.Workspace, error) {
	var workspacePath string
	if err := call.Decode(&workspacePath); err != nil {
		return nil, err
	}

	if ws := a.Workspace(workspacePath); ws != nil {
		return ws, nil
	}

	return nil, fmt.Errorf("no workspace found for %q", workspacePath)
}
