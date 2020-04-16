package daemon

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/manifold/tractor/pkg/config"
	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/object"
	"github.com/manifold/tractor/pkg/misc/buffer"
	"github.com/manifold/tractor/pkg/misc/daemon"
	"github.com/manifold/tractor/pkg/misc/mdns"
	"github.com/manifold/tractor/pkg/stdlib"
	"github.com/manifold/tractor/pkg/workspace/editor"
	"github.com/manifold/tractor/pkg/workspace/rpc"
	"github.com/manifold/tractor/pkg/workspace/state"

	zapper "github.com/manifold/tractor/pkg/misc/logging/zap"
)

var (
	addr = flag.String("addr", "localhost:0", "server listener address")
)

func init() {
	stdlib.Load()
}

func captureOutput(size int64) (*buffer.Buffer, func()) {
	buf, err := buffer.NewBuffer(size)
	if err != nil {
		panic(err)
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	cancel := func() {
		writer.Close()
		os.Stdout = stdout
		os.Stderr = stderr
	}
	os.Stdout = writer
	os.Stderr = writer

	go io.Copy(io.MultiWriter(buf, stdout), reader)

	return buf, cancel
}

func Run() {
	cfg, err := config.OpenDefault()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	buf, undoCapture := captureOutput(1024 * 1024)
	defer undoCapture()

	logger, undoRedirect := zapper.NewRedirectedLogger(os.Stdout)
	defer undoRedirect()

	rpcSvc := &rpc.Service{
		Log:    logger,
		Output: buf,
	}

	// not sure why rpc service needs to ALWAYS be
	// in object registry. pls find out
	object.RegistryPreloader = func(o manifold.Object) []interface{} {
		return []interface{}{o, rpcSvc}
	}

	dm := daemon.New([]daemon.Service{
		&Service{
			ListenAddr: *addr,
			Log:        logger,
		},
		&state.Service{
			Log: logger,
		},
		rpcSvc,
		&mdns.Service{
			Log: logger,
		},
		&editor.Service{
			Log:    logger,
			Config: cfg,
		},
	}...)

	fatal(dm.Run(context.Background()))
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
