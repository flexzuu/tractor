package rpc

import (
	"context"
	"fmt"
	"os"

	"github.com/manifold/qtalk/golang/mux"
	qrpc "github.com/manifold/qtalk/golang/rpc"
	"github.com/manifold/tractor/pkg/config"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/workspace/supervisor"
)

// Service provides a QRPC server to start, restart, and stop running
// workspace supervisors.
type Service struct {
	Agent      agent
	Config     *config.Config
	Log        logging.Logger
	api        qrpc.API
	l          mux.Listener
	socketPath string
}

// only dep on agent, avoids dep cycle
type agent interface {
	Supervisor(path string) *supervisor.Supervisor
}

func (s *Service) InitializeDaemon() (err error) {
	if s.l, err = mux.ListenUnix(s.Config.Agent.SocketPath); err != nil {
		return err
	}

	s.api = qrpc.NewAPI()
	s.api.HandleFunc("start", s.Start())
	s.api.HandleFunc("stop", s.Stop())
	return nil
}

func (s *Service) Serve(ctx context.Context) {
	server := &qrpc.Server{}

	logging.Infof(s.Log, "agent listening at unix://%s", s.Config.Agent.SocketPath)
	if err := server.Serve(s.l, s.api); err != nil {
		fmt.Println(err)
	}
	os.Remove(s.Config.Agent.SocketPath)
}
