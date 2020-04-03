package rpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strconv"

	"github.com/hashicorp/mdns"
	"github.com/manifold/qtalk/golang/mux"
	qrpc "github.com/manifold/qtalk/golang/rpc"
	"github.com/manifold/tractor/pkg/misc/buffer"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/workspace/editor"
	"github.com/manifold/tractor/pkg/workspace/state"
	"github.com/manifold/tractor/pkg/workspace/view"
	"github.com/miekg/dns"
)

type Service struct {
	// Protocol   string
	ListenAddr string

	Output *buffer.Buffer
	Log    logging.Logger
	State  *state.Service
	Editor *editor.Service

	viewState *view.State
	clients   map[qrpc.Caller]string
	api       qrpc.API
	l         mux.Listener
}

func (s *Service) UpdateView() {
	s.updateView()
}

func (s *Service) updateView() {
	if s.viewState == nil || s.State == nil {
		return
	}
	// TODO: mutex, etc
	s.viewState.EditorsEndpoint = s.Editor.Endpoint()
	s.viewState.Update(s.State.Root)
	for client, callback := range s.clients {
		_, err := client.Call(callback, s.viewState, nil)
		if err != nil {
			delete(s.clients, client)
			log.Println(err)
		}
	}
}

func (s *Service) InitializeDaemon() (err error) {
	if s.l, err = mux.ListenWebsocket(s.ListenAddr); err != nil {
		return err
	}

	s.clients = make(map[qrpc.Caller]string)
	s.viewState = view.New(s.State.Root)

	s.api = qrpc.NewAPI()
	s.api.HandleFunc("echo", s.Echo())
	s.api.HandleFunc("console", s.Console())
	s.api.HandleFunc("reload", s.Reload())
	s.api.HandleFunc("refreshObject", s.RefreshObject())
	s.api.HandleFunc("repl", s.Repl())
	s.api.HandleFunc("selectNode", s.SelectNode())
	s.api.HandleFunc("removeComponent", s.RemoveComponent())
	s.api.HandleFunc("reloadComponent", s.ReloadComponent())
	s.api.HandleFunc("selectProject", s.SelectProject())
	s.api.HandleFunc("moveNode", s.MoveNode())
	s.api.HandleFunc("subscribe", s.Subscribe())
	s.api.HandleFunc("appendNode", s.AppendNode())
	s.api.HandleFunc("deleteNode", s.DeleteNode())
	s.api.HandleFunc("loadPrefab", s.LoadPrefab())
	s.api.HandleFunc("appendComponent", s.AppendComponent())
	s.api.HandleFunc("setValue", s.SetValue())
	// s.api.HandleFunc("setExpression", s.SetExpression())
	s.api.HandleFunc("callMethod", s.CallMethod())
	s.api.HandleFunc("updateNode", s.UpdateNode())
	s.api.HandleFunc("addDelegate", s.AddDelegate())

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
	server := &qrpc.Server{}
	s.Log.Infof("daemon listening at %s", s.l.Addr().String())
	if err := server.Serve(s.l, s.api); err != nil {
		fmt.Println(err)
	}
	// if s.Protocol == "unix" {
	// 	os.Remove(s.ListenAddr)
	// }
}

func (s *Service) TerminateDaemon() error {
	for client, _ := range s.clients {
		client.Call("shutdown", nil, nil)
	}
	// if s.Protocol == "unix" {
	// 	os.Remove(s.ListenAddr)
	// }
	return nil
}

func muxListenTo(proto, addr string) (mux.Listener, error) {
	switch proto {
	case "websocket":
		return mux.ListenWebsocket(addr)
	case "unix":
		return mux.ListenUnix(addr)
	}

	return nil, fmt.Errorf("cannot connect to %s, unknown protocol %q", addr, proto)
}
