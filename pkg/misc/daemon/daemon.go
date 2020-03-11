package daemon

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/misc/registry"
)

// Initializer is initialized before services are started. Returning
// an error will cancel the start of daemon services.
type Initializer interface {
	InitializeDaemon() error
}

// Terminator is terminated when the daemon gets a stop signal.
type Terminator interface {
	TerminateDaemon() error
}

// Service is run after the daemon is initialized.
type Service interface {
	Serve(ctx context.Context)
}

// Daemon is a top-level daemon lifecycle manager runs services given to it.
type Daemon struct {
	Initializers []Initializer
	Services     []Service
	Terminators  []Terminator
	Logger       logging.Logger
	Context      context.Context
	OnFinished   func()
	running      int32
	cancel       context.CancelFunc
	termErrs     chan []error
}

// New builds a daemon configured to run a set of services. The services
// are populated with each other if they have fields that match anything
// that was passed in.
func New(services ...Service) *Daemon {
	d := &Daemon{}
	d.AddServices(services...)
	return d
}

// Run creates a daemon from services and runs it with a background context
func Run(services ...Service) error {
	d := New(services...)
	return d.Run(context.Background())
}

// AddServices appends Service and Terminators to daemon
func (d *Daemon) AddServices(services ...Service) {
	r, _ := registry.New(d)
	for _, s := range d.Services {
		r.Register(s)
	}
	for _, s := range services {
		r.Register(s)
		d.Services = append(d.Services, s)
		if t, ok := s.(Terminator); ok {
			d.Terminators = append(d.Terminators, t)
		}
	}
	r.SelfPopulate()
}

// Run executes the daemon lifecycle
func (d *Daemon) Run(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&d.running, 0, 1) {
		return errors.New("already running")
	}

	// call initializers
	for _, i := range d.Initializers {
		if err := i.InitializeDaemon(); err != nil {
			return err
		}
	}

	// finish if no services
	if len(d.Services) == 0 {
		return errors.New("no services to run")
	}

	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancelFunc := context.WithCancel(ctx)
	d.Context = ctx
	d.cancel = cancelFunc
	d.termErrs = make(chan []error)

	// setup terminators on stop signals
	go TerminateOnSignal(d)
	go TerminateOnContextDone(d)

	var wg sync.WaitGroup
	for _, service := range d.Services {
		wg.Add(1)
		go func(s Service) {
			s.Serve(d.Context)
			wg.Done()
		}(service)
	}

	finished := make(chan bool)
	go func() {
		wg.Wait()
		close(finished)
	}()

	var errs []error
	select {
	case <-finished:
		select {
		case errs = <-d.termErrs:
		default:
		}
	case errs = <-d.termErrs:
		select {
		case <-finished:
		case <-time.After(1 * time.Second):
			// TODO: show/track what servies
			log.Println("warning: unfinished services")
		}
	}

	if d.OnFinished != nil {
		d.OnFinished()
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Terminate cancels the daemon context and calls Terminators in reverse order
func (d *Daemon) Terminate() {
	if d == nil {
		// find these cases and prevent them!
		panic("daemon reference used to Terminate but daemon pointer is nil")
	}

	if !atomic.CompareAndSwapInt32(&d.running, 1, 0) {
		return
	}

	if d.cancel != nil {
		d.cancel()
	}

	var errs []error
	for i := len(d.Terminators) - 1; i >= 0; i-- {
		if err := d.Terminators[i].TerminateDaemon(); err != nil {
			errs = append(errs, err)
		}
	}
	d.termErrs <- errs
}

// TerminateOnSignal waits for SIGINT, SIGHUP, SIGTERM, SIGKILL(?) to terminate the daemon.
func TerminateOnSignal(d *Daemon) {
	termSigs := make(chan os.Signal, 1)
	signal.Notify(termSigs, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGTERM)
	<-termSigs
	d.Terminate()
}

// TerminateOnContextDone waits for the deamon's context to be canceled.
func TerminateOnContextDone(d *Daemon) {
	<-d.Context.Done()
	d.Terminate()
}
