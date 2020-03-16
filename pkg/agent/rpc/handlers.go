package rpc

import (
	"fmt"

	qrpc "github.com/manifold/qtalk/golang/rpc"
	"github.com/manifold/tractor/pkg/misc/logging"
)

type workspaceInfo struct {
	Name     string
	Endpoint string
	Path     string
}

func (s *Service) List() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var list []workspaceInfo
		for _, e := range s.Agent.Workspaces() {
			list = append(list, workspaceInfo{
				Name:     e.Name(),
				Endpoint: e.Endpoint(),
				Path:     e.Path(),
			})
		}
		r.Return(list)
	}
}

func (s *Service) Start() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var workspacePathOrName string
		if err := c.Decode(&workspacePathOrName); err != nil {
			r.Return(err)
			return
		}

		sup := s.Agent.Supervisor(workspacePathOrName)
		if sup == nil {
			r.Return(fmt.Errorf("no supervisor for %q", workspacePathOrName))
			return
		}

		logging.Infof(s.Log, "@%s: Start()", sup.Name())
		if err := sup.Daemon.Restart(); err != nil {
			r.Return(err)
			return
		}
		r.Return(fmt.Sprintf("workspace %q started", sup.Name()))
	}
}

func (s *Service) Stop() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var workspacePathOrName string
		if err := c.Decode(&workspacePathOrName); err != nil {
			r.Return(err)
			return
		}

		sup := s.Agent.Supervisor(workspacePathOrName)
		if sup == nil {
			r.Return(fmt.Errorf("no supervisor for %q", workspacePathOrName))
			return
		}

		logging.Infof(s.Log, "@%s: Stop()", sup.Name())
		if sup.Daemon != nil {
			if err := sup.Daemon.Stop(); err != nil {
				r.Return(err)
				return
			}
		}

		r.Return(fmt.Sprintf("workspace %q stopped", sup.Name()))
	}
}
