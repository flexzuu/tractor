package daemon

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/hashicorp/mdns"
	"github.com/manifold/qtalk/golang/mux"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/workspace/editor"
	"github.com/manifold/tractor/pkg/workspace/rpc"
	"github.com/miekg/dns"
	"golang.org/x/net/websocket"
)

type Service struct {
	ListenAddr string

	RPC    *rpc.Service
	Editor *editor.Service
	Log    logging.Logger

	l net.Listener
	s *http.Server
}

func (s *Service) InitializeDaemon() (err error) {
	s.l, err = net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.s = &http.Server{
		Handler: s,
	}
	return nil
}

func (s *Service) Records(q dns.Question) []dns.RR {
	_, p, _ := net.SplitHostPort(s.l.Addr().String())
	port, _ := strconv.Atoi(p)
	wd, _ := os.Getwd()
	zone, _ := mdns.NewMDNSService(path.Base(wd), "_tractor._tcp", "", "", port, nil, []string{wd})
	return zone.Records(q)
}

func (s *Service) Serve(ctx context.Context) {
	s.Log.Debugf("tractor listening at http://%s", s.l.Addr().String())
	log.Fatal(s.s.Serve(s.l))
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/rpc" {
		websocket.Handler(func(conn *websocket.Conn) {
			conn.PayloadType = websocket.BinaryFrame
			sess := mux.NewSession(conn, r.Context())
			s.Log.Debug("new tractor rpc session")
			s.RPC.Inbox <- sess
			sess.Wait()
		}).ServeHTTP(w, r)
		return
	}

	// if path exists against editors fs
	if s.Editor != nil && s.Editor.MatchHTTP(r) {
		s.Editor.ServeHTTP(w, r)
		return
	}

	// proxy to theia at 11010
	u, _ := url.Parse("http://localhost:11010/")
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.ServeHTTP(w, r)
}
