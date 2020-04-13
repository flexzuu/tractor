package workspace

import (
	"fmt"

	"github.com/hashicorp/mdns"
	"github.com/manifold/tractor/pkg/config"
	"github.com/manifold/tractor/pkg/workspace/supervisor"
)

type Entry struct {
	Config     config.Workspace
	Supervisor *supervisor.Supervisor
	Service    *mdns.ServiceEntry
}

func (w *Entry) Supervised() bool {
	return w.Supervisor != nil
}

func (w *Entry) Endpoint() string {
	if w.Service == nil {
		return ""
	}
	return fmt.Sprintf("ws://localhost:%d/rpc", w.Service.Port)
}

func (w *Entry) Name() string {
	return w.Config.Name
}

func (w *Entry) Path() string {
	return w.Config.Path
}

func (w *Entry) Status() supervisor.Status {
	if w.Supervisor != nil {
		return w.Supervisor.Status()
	}
	if w.Service != nil {
		return supervisor.StatusAvailable
	}
	return supervisor.StatusUnavailable
}

func (w *Entry) Observe(cb supervisor.Observer) {
	if w.Supervisor == nil {
		return
	}
	w.Supervisor.Observe(cb)
}
