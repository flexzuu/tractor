package state

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/image"
	"github.com/manifold/tractor/pkg/misc/debouncer"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/misc/notify"
)

type Service struct {
	Protocol   string
	ListenAddr string

	Log   logging.Logger
	Root  manifold.Object
	Image *image.Image
}

func (s *Service) InitializeDaemon() (err error) {
	wd, err := os.Getwd() // TODO: override with env var
	if err != nil {
		return err
	}
	s.Image = image.New(wd)

	s.Root, err = s.Image.Load()
	if err != nil {
		return err
	}

	manifold.Walk(s.Root, func(n manifold.Object) {
		for _, com := range n.Components() {
			if initializer, ok := com.Pointer().(initializer); ok {
				if err := initializer.Initialize(); err != nil {
					log.Print(err)
				}
			}
		}
	})

	debounce := debouncer.New(2 * time.Second)
	notify.Observe(s.Root, notify.Func(func(event interface{}) {
		debounce(func() {
			// TODO: Log errors?
			log.Print("change triggered SNAPSHOT")
			s.Snapshot()
		})
	}))

	return nil
}

func (s *Service) TerminateDaemon() error {
	return s.Snapshot()
}

func (s *Service) Serve(ctx context.Context) {
	<-ctx.Done()
}

func (s *Service) Snapshot() error {
	return s.Image.Write(s.Root)
}

type preInitializer interface {
	PreInitialize()
}

type initializer interface {
	Initialize() error
}
