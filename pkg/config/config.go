package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

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
	SocketPath           string // ~/.tractor/agent.sock
	WorkspacesPath       string // ~/.tractor/workspaces
	WorkspaceSocketsPath string // ~/.tractor/sockets
	WorkspaceBinPath     string // ~/.tractor/bin
	GoBin                string
	PreferredBrowser     string
}

// Config describes the user settings from ~/.tractor/config.toml
type Config struct {
	Agent Agent
}
