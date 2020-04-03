package editor

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Schobers/bindatafs"
	"github.com/manifold/tractor/pkg/data/editors"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/progrium/hotweb/pkg/hotweb"
	"github.com/spf13/afero"
)

type Service struct {
	//EditorsFs  afero.Fs

	Log logging.Logger

	l  net.Listener
	s  *http.Server
	hw *hotweb.Handler
}

func (s *Service) InitializeDaemon() (err error) {
	s.l, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	var fs afero.Fs = bindatafs.NewFs(editors.MustAsset, editors.AssetInfo, editors.AssetNames)
	if os.Getenv("TRACTOR_SRC") != "" {
		s.Log.Debugf("using source at %s", os.Getenv("TRACTOR_SRC"))
		fs = afero.NewBasePathFs(afero.NewOsFs(), os.Getenv("TRACTOR_SRC"))
	}
	s.hw = hotweb.New(fs, "studio/editors")
	s.s = &http.Server{
		Handler: s.hw,
	}
	return nil
}

func (s *Service) Endpoint() string {
	return s.l.Addr().String()
}

func (s *Service) Serve(ctx context.Context) {
	go func() {
		s.Log.Error(s.hw.Watch())
	}()
	s.Log.Debugf("editors listening at %s", s.l.Addr().String())
	log.Fatal(s.s.Serve(s.l))
}
