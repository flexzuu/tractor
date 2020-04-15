package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/manifold/qtalk/golang/mux"
)

const HomeDir = ".tractor"
const FileName = "config.toml"
const WorkspaceFile = "tractor.go"

// Duration wraps time.Duration to parse from TOML configs.
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// ParseFile parses the given toml file into a Config.
func ParseFile(path string, cfg *Config) error {
	_, err := toml.DecodeFile(path, cfg)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Agent describes the user settings for the Tractor agent.
type Agent struct {
	SocketPath string // ~/.tractor/agent.sock

}

func (a *Agent) CheckSocket() bool {
	_, err := mux.DialUnix(a.SocketPath)
	if err != nil {
		os.RemoveAll(a.SocketPath)
		return false
	}
	return true
}

type Workspace struct {
	Name string
	Path string
}

// Config describes the user settings from ~/.tractor/config.toml
type Config struct {
	Agent Agent

	Dir           string // ~/.tractor
	WorkspacesDir string // ~/.tractor/workspaces
	BinDir        string // ~/.tractor/bin
	GoBin         string
	BrowserPref   string
	DevWatch      Duration `toml:"DevWatchInterval"` // time duration: "50ms", prefer DevWatchInterval() for usage
}

func Open(dir string) (*Config, error) {
	cfg := baseConfig(dir)
	if err := ParseFile(filepath.Join(dir, FileName), cfg); err != nil {
		return nil, err
	}
	os.MkdirAll(cfg.WorkspacesDir, 0700)
	os.MkdirAll(cfg.BinDir, 0700)
	return cfg, nil
}

func OpenDefault() (*Config, error) {
	return Open(defaultPath())
}

func defaultPath() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, HomeDir)
}

func baseConfig(dir string) *Config {
	bin, err := exec.LookPath("go")
	if err != nil {
		panic(err)
	}
	return &Config{
		Dir: dir,
		Agent: Agent{
			SocketPath: filepath.Join(dir, "agent.sock"),
		},
		WorkspacesDir: filepath.Join(dir, "workspaces"),
		BinDir:        filepath.Join(dir, "bin"),
		GoBin:         bin,
		DevWatch:      Duration{Duration: defaultDevInterval},
	}
}

var defaultDevInterval = time.Millisecond * 100

func (cfg *Config) resolveWorkspaceSymlink(fi os.FileInfo) (bool, string, error) {
	if fi.IsDir() {
		return false, "", nil
	}

	path := filepath.Join(cfg.WorkspacesDir, fi.Name())
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false, resolved, err
	}

	if resolved == path {
		return false, resolved, nil
	}

	rfi, err := os.Lstat(resolved)
	if err != nil {
		return false, resolved, err
	}

	return rfi.IsDir(), resolved, nil
}

func (cfg *Config) DevWatchInterval() time.Duration {
	if cfg == nil {
		return defaultDevInterval
	}
	return cfg.DevWatch.Duration
}

func (cfg *Config) AddWorkspace(path string) (string, error) {
	_, err := os.Lstat(filepath.Join(path, WorkspaceFile))
	if err != nil {
		return "", err // not a tractor workspace
	}

	basepath := filepath.Base(path)
	base := basepath
	i := 1
	for {
		err = os.Symlink(path, filepath.Join(cfg.WorkspacesDir, base))
		if err != nil && !os.IsExist(err) {
			return base, err
		}

		if err == nil {
			return base, nil
		}

		i++
		base = fmt.Sprintf("%s-%d", basepath, i)
	}
}

func (cfg *Config) Workspaces() ([]Workspace, error) {
	entries, err := ioutil.ReadDir(cfg.WorkspacesDir)
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	for _, entry := range entries {
		ok, resolved, err := cfg.resolveWorkspaceSymlink(entry)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		workspaces = append(workspaces, Workspace{
			Name: entry.Name(),
			Path: resolved,
		})
	}
	return workspaces, nil
}
