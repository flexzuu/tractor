package daemon

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/object"
	"github.com/manifold/tractor/pkg/misc/daemon"
	"github.com/manifold/tractor/pkg/misc/mdns"
	"github.com/manifold/tractor/pkg/stdlib"
	"github.com/manifold/tractor/pkg/workspace/rpc"
	"github.com/manifold/tractor/pkg/workspace/state"

	zapper "github.com/manifold/tractor/pkg/misc/logging/zap"
)

var (
	addr = flag.String("addr", "localhost:4243", "server listener address")
	// proto = flag.String("proto", "websocket", "server listener protocol")
)

func init() {
	stdlib.Load()
}

func Run() {
	flag.Parse()
	logger, undo := zapper.NewRedirectedLogger(os.Stdout)
	defer undo()
	rpcSvc := &rpc.Service{
		// Protocol:   *proto,
		ListenAddr: *addr,
		Log:        logger,
	}
	object.RegistryPreloader = func(o manifold.Object) []interface{} {
		return []interface{}{o, rpcSvc}
	}
	dm := daemon.New([]daemon.Service{
		&state.Service{
			Log: logger,
		},
		rpcSvc,
		&mdns.Service{
			Log: logger,
		},
	}...)
	fatal(dm.Run(context.Background()))
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
