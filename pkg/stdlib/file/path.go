package file

import (
	"net/http"
	"os"

	"github.com/spf13/afero"
)

type Path struct {
	afero.Fs

	Filepath string
}

func (c *Path) exists() bool {
	if c.Filepath == "" {
		return false
	}
	if _, err := os.Stat(c.Filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Path) ComponentEnable() {
	c.Fs = nil
	if c.exists() {
		c.Fs = afero.NewBasePathFs(afero.NewOsFs(), c.Filepath)
	}
}

func (c *Path) Open(name string) (http.File, error) {
	return http.Dir(c.Filepath).Open(name)
}
