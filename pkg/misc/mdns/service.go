package mdns

import (
	"context"

	"github.com/hashicorp/mdns"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/miekg/dns"
)

type Service struct {
	Providers []mdns.Zone
	Log       logging.Logger

	srv *mdns.Server
}

func (s *Service) InitializeDaemon() (err error) {
	s.Log.Debug("discoverable via mdns")
	s.srv, err = mdns.NewServer(&mdns.Config{Zone: s})
	return
}

func (s *Service) Serve(ctx context.Context) {
	<-ctx.Done()
}

func (s *Service) Records(q dns.Question) []dns.RR {
	var r []dns.RR
	for _, contributor := range s.Providers {
		if contributor == s {
			continue
		}
		r = append(r, contributor.Records(q)...)
	}
	return r
}

func (s *Service) TerminateDaemon() error {
	return s.srv.Shutdown()
}
