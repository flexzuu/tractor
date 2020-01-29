package library

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/progrium/prototypes/go-reflected"
)

var (
	registered []*RegisteredComponent
)

type RegisteredComponent struct {
	Type     reflected.Type
	Name     string
	Filepath string
	ID       string
}

func (rc *RegisteredComponent) New() manifold.Component {
	return newComponent(rc.Name, rc.NewValue(), rc.ID)
}

func (rc *RegisteredComponent) NewValue() interface{} {
	return reflected.New(rc.Type).Interface()
}

func Register(v interface{}, id, filepath string) {
	if filepath == "" {
		_, filepath, _, _ = runtime.Caller(1)
	}
	t := reflected.ValueOf(v).Type()
	registered = append(registered, &RegisteredComponent{
		Type:     t,
		Name:     fmt.Sprintf("%s.%s", path.Base(t.PkgPath()), t.Name()),
		Filepath: filepath,
		ID:       id,
	})
}

// deprecated
func Names() []string {
	var names []string
	for _, rc := range registered {
		if rc.ID != "" {
			continue
		}
		names = append(names, rc.Name)
	}
	return names
}

func Registered() []*RegisteredComponent {
	r := make([]*RegisteredComponent, len(registered))
	copy(r, registered)
	return r
}

func Lookup(name string) *RegisteredComponent {
	for _, rc := range registered {
		if rc.Name == name {
			return rc
		}
	}
	return nil
}

func LookupID(id string) *RegisteredComponent {
	for _, rc := range registered {
		if rc.ID == id {
			return rc
		}
	}
	return nil
}

func Related(c *RegisteredComponent) (related []*RegisteredComponent) {
	if c == nil {
		return
	}
	for _, rc := range registered {
		if rc.Type == c.Type {
			continue
		}
		if strings.HasPrefix(rc.Filepath, filepath.Dir(c.Filepath)) {
			related = append(related, rc)
		}
	}
	return
}
