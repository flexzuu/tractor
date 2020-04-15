package editor

import (
	"context"
	"net/http"
	"os"

	"github.com/Schobers/bindatafs"
	"github.com/manifold/tractor/pkg/config"
	"github.com/manifold/tractor/pkg/data/editors"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/progrium/hotweb/pkg/hotweb"
	"github.com/spf13/afero"
)

type Service struct {
	Log    logging.Logger
	Config *config.Config

	hw *hotweb.Handler
}

func (s *Service) InitializeDaemon() (err error) {
	var fs afero.Fs = bindatafs.NewFs(
		editors.MustAsset,
		editors.AssetInfo,
		editors.AssetNames,
	)
	if os.Getenv("TRACTOR_SRC") != "" {
		s.Log.Debugf("using source at %s", os.Getenv("TRACTOR_SRC"))
		fs = afero.NewBasePathFs(afero.NewOsFs(), os.Getenv("TRACTOR_SRC"))
	}
	s.hw = hotweb.New(fs, "studio/editors", "/views")
	s.hw.WatchInterval = s.Config.DevWatchInterval()
	return nil
}

func (s *Service) MatchHTTP(r *http.Request) bool {
	return s.hw.MatchHTTP(r)
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.hw.ServeHTTP(w, r)
}

// DELETE ME
// see also frontend session state
func (s *Service) Endpoint() string {
	return "localhost:11000"
}

func (s *Service) Serve(ctx context.Context) {
	s.Log.Error(s.hw.Watch())
}
