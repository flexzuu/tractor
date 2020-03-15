package rpc

import (
	"fmt"

	qrpc "github.com/manifold/qtalk/golang/rpc"
	"github.com/manifold/tractor/pkg/misc/logging"
)

func (s *Service) Start() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var workspacePath string
		if err := c.Decode(&workspacePath); err != nil {
			r.Return(err)
			return
		}

		sup := s.Agent.Supervisor(workspacePath)
		if sup == nil {
			r.Return(fmt.Errorf("no supervisor for %q", workspacePath))
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
		var workspacePath string
		if err := c.Decode(&workspacePath); err != nil {
			r.Return(err)
			return
		}

		sup := s.Agent.Supervisor(workspacePath)
		if sup == nil {
			r.Return(fmt.Errorf("no supervisor for %q", workspacePath))
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
