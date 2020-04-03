package agent

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/manifold/tractor/pkg/agent/console"
	"github.com/manifold/tractor/pkg/agent/rpc"
	"github.com/manifold/tractor/pkg/agent/selfdev"
	"github.com/manifold/tractor/pkg/agent/systray"
	"github.com/manifold/tractor/pkg/agent/workspace"
	"github.com/manifold/tractor/pkg/config"
	"github.com/manifold/tractor/pkg/misc/daemon"
	"github.com/manifold/tractor/pkg/misc/logging"
	"github.com/manifold/tractor/pkg/misc/logging/null"
	"github.com/manifold/tractor/pkg/workspace/supervisor"
)

// Agent manages multiple workspaces in a directory (default: ~/.tractor).
type Agent struct {
	Config     *config.Config
	DevMode    bool
	SocketPath string

	Daemon  *daemon.Daemon
	Systray *systray.Service
	Console *console.Service
	Logger  logging.Logger

	workspaces []*workspace.Entry

	mu sync.RWMutex
}

// New returns a new agent for the given path. If the given path is empty, a
// default of ~/.tractor will be used.
func New(path string, console *console.Service, devMode bool) (*Agent, error) {
	var cfg *config.Config
	var err error
	if path != "" {
		cfg, err = config.Open(path)
	} else {
		cfg, err = config.OpenDefault()
	}
	if err != nil {
		return nil, err
	}
	a := &Agent{
		Config:  cfg,
		DevMode: devMode,
		Console: console,
		Logger:  console,
	}
	if a.Logger == nil {
		a.Logger = &null.Logger{}
	}
	if devMode {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		os.Setenv("TRACTOR_SRC", wd)
	}
	return a, nil
}

func (a *Agent) InitializeDaemon() (err error) {
	spaces, err := a.Config.Workspaces()
	if err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, space := range spaces {
		sup := supervisor.New(space.Path, space.Name, a.Console.NewPipe("@"+space.Name))
		a.workspaces = append(a.workspaces, &workspace.Entry{
			Config:     space,
			Supervisor: sup,
			Service:    nil,
		})
		a.Daemon.AddServices(sup)
	}
	return nil
}

func (a *Agent) DaemonServices() []daemon.Service {
	a.Systray = &systray.Service{
		Config: a.Config,
		Agent:  a,
	}
	services := []daemon.Service{
		a,
		a.Console,
		a.Systray,
		&rpc.Service{
			Config: a.Config,
			Agent:  a,
		},
	}
	if a.DevMode {
		services = append(services, []daemon.Service{
			&selfdev.Service{},
		}...)
	}
	return services
}

func (a *Agent) Serve(ctx context.Context) {
	for {
		// TODO: listen for ctx finished

		time.Sleep(2 * time.Second) // TODO: no hardcode
		entries, err := a.LookupServices()
		if err != nil {
			panic(err)
		}

		reloadSystray := false
		a.mu.Lock()

		// delete non-supervised workspaces not in entries
		n := 0
		for _, s := range a.workspaces {
			if !s.Supervised() {
				for _, entry := range entries {
					if s.Path() == entry.Info {
						a.workspaces[n] = s
						n++
						break
					}
				}
			} else {
				a.workspaces[n] = s
				n++
			}
		}
		before := len(a.workspaces)
		a.workspaces = a.workspaces[:n]
		if len(a.workspaces) < before {
			reloadSystray = true
		}

		// update/add entries to workspaces
		for _, entry := range entries {
			parts := strings.Split(entry.Name, ".")
			name := parts[0]
			isnew := true
			for _, s := range a.workspaces {
				if s.Path() == entry.Info {
					isnew = false
					s.Service = entry
				}
			}
			if isnew {
				a.workspaces = append(a.workspaces, &workspace.Entry{
					Config: config.Workspace{
						Name: name,
						Path: entry.Info,
					},
					Supervisor: nil,
					Service:    entry,
				})
				reloadSystray = true
			}
		}
		a.mu.Unlock()

		if reloadSystray && a.Systray != nil {
			a.Systray.Restart()
		}
	}
}

func (a *Agent) Workspaces() []*workspace.Entry {
	a.mu.Lock()
	defer a.mu.Unlock()
	w := make([]*workspace.Entry, len(a.workspaces))
	copy(w, a.workspaces)
	return w
}

func (a *Agent) Supervisor(workspacePathOrName string) *supervisor.Supervisor {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, s := range a.workspaces {
		if s.Path() == workspacePathOrName || s.Name() == workspacePathOrName {
			return s.Supervisor
		}
	}
	return nil
}

// Shutdown shuts all workspaces down and cleans up socket files.
func (a *Agent) TerminateDaemon() {
	logging.Infof(a.Logger, "shutting down")
	os.RemoveAll(a.Config.Agent.SocketPath)
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, ws := range a.workspaces {
		if ws.Supervised() && ws.Supervisor.Daemon != nil {
			ws.Supervisor.Daemon.Stop()
		}
	}
}

func (a *Agent) LookupServices() ([]*mdns.ServiceEntry, error) {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	var entries []*mdns.ServiceEntry
	go func() {
		for entry := range entriesCh {
			if strings.HasSuffix(entry.Name, "_tractor._tcp.local.") {
				entries = append(entries, entry)
			}
		}
	}()
	if err := mdns.Lookup("_tractor._tcp", entriesCh); err != nil {
		return nil, err
	}
	close(entriesCh)
	return entries, nil
}

// func (a *Agent) Watch(ctx context.Context) {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		info(a.Logger, "unable to create watcher:", err)
// 		return
// 	}
// 	watcher.Add(a.Config.Agent.WorkspacesDir)
// 	debounce := debouncer.New(20 * time.Millisecond)
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			watcher.Close()
// 			return
// 		case event, ok := <-watcher.Events:
// 			if !ok {
// 				return
// 			}
// 			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
// 				continue
// 			}

// 			debounce(func() {
// 				a.WorkspacesChanged <- struct{}{}
// 			})

// 		case err, ok := <-watcher.Errors:
// 			if !ok {
// 				return
// 			}
// 			logErr(a.Logger, "watcher error:", err)
// 		}
// 	}
// }
