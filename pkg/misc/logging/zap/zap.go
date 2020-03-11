package zap

import (
	"fmt"
	"io"
	"net/url"
	"time"

	api "github.com/manifold/tractor/pkg/misc/logging"
	"go.uber.org/zap"
)

func NewLogger(w io.WriteCloser, options ...zap.Option) *Logger {
	sinkName := fmt.Sprintf("logger-%d", time.Now().Unix())
	zap.RegisterSink(sinkName, func(u *url.URL) (zap.Sink, error) {
		return sink{w}, nil
	})
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{fmt.Sprintf("%s://", sinkName)}
	logger, _ := config.Build(options...)
	return &Logger{logger.Sugar()}
}

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) With(args ...interface{}) api.Logger {
	return l.With(args...)
}

type sink struct {
	io.WriteCloser
}

func (w sink) Sync() error {
	return nil
}
