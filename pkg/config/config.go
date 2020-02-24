package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// ParseFile parses the given toml file into a Config.
func ParseFile(path string) (Config, error) {
	c := Config{}
	_, err := toml.DecodeFile(path, &c)
	if err != nil && !os.IsNotExist(err) {
		return c, err
	}
	return c, nil
}

// Config describes the user settings from ~/.tractor/config.toml
type Config struct {
	Agent struct {
		PreferredBrowser string
	}
}
