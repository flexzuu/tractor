package rpc

import (
	"bufio"
	"context"
	"fmt"
	"go/build"
	"go/scanner"
	"io"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/davecgh/go-spew/spew"
	qrpc "github.com/manifold/qtalk/golang/rpc"
	"github.com/manifold/tractor/pkg/manifold/library"
	"github.com/manifold/tractor/pkg/manifold/object"
	"github.com/manifold/tractor/pkg/manifold/prefab"
	"github.com/manifold/tractor/pkg/workspace/repl"
)

type AppendNodeParams struct {
	ID   string
	Name string
}

type SetValueParams struct {
	Path     string
	Value    interface{}
	IntValue *int
	RefValue *string
	Type     string
}

type RemoveComponentParams struct {
	ID        string
	Component string
}

type NodeParams struct {
	ID     string
	Name   *string
	Active *bool
}

type DelegateParams struct {
	ID       string
	Contents string
}

type MoveNodeParams struct {
	ID    string
	Index int
}

func (s *Service) Echo() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var s []string
		err := c.Decode(&s)
		if err != nil {
			r.Return(err)
			return
		}
		r.Return(s)
	}
}

func (s *Service) Console() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		console := s.Output.Pipe()

		ch, err := r.Hijack(nil)
		if err != nil {
			console.Close()
			r.Return(err)
			return
		}

		_, err = io.Copy(ch, console)
		ch.Close()
		console.Close()

		if err == io.ErrClosedPipe {
			r.Return(err)
			return
		}
	}
}

func (s *Service) Reload() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		// TODO
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) Repl() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params DelegateParams
		_ = c.Decode(&params)
		// ^^ TODO: make sure this isn't necessary before hijacking
		ch, err := r.Hijack(nil)
		if err != nil {
			log.Println(err)
		}

		i := interp.New(interp.Options{GoPath: build.Default.GOPATH})
		i.Use(stdlib.Symbols)
		i.Use(interp.Symbols)
		i.Use(repl.Symbols)
		i.Use(map[string]map[string]reflect.Value{
			"console": map[string]reflect.Value{
				"View":  reflect.ValueOf(s.viewState),
				"State": reflect.ValueOf(s.State),
			},
		})
		i.Eval("import \"console\"")
		i.Eval("import \"github.com/manifold/tractor/pkg/manifold\"")
		i.Eval("state := console.State")
		i.Eval("view := console.View")
		i.Eval("selected := func() manifold.Object { return state.Root.FindID(view.SelectedNode) }")
		func(i *interp.Interpreter, in io.Reader, out io.Writer) {
			scs := spew.ConfigState{
				MaxDepth: 2,
				Indent:   "  ",
			}
			s := bufio.NewScanner(in)
			src := ""
			for s.Scan() {
				src += s.Text() + "\n"
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				v, err := i.EvalWithContext(ctx, src)
				if err != nil {
					switch err.(type) {
					case scanner.ErrorList:
						// Early failure in the scanner: the source is incomplete
						// and no AST could be produced, neither compiled / run.
						// Get one more line, and retry
						continue
					default:
						fmt.Fprintln(out, err)
					}
				} else if v.IsValid() {
					scs.Fdump(out, v.Interface())
					fmt.Fprintf(out, "\r")
				}
				src = ""
			}
		}(i, ch, ch)

	}
}

func (s *Service) RemoveComponent() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params RemoveComponentParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		n := s.State.Root.FindID(params.ID)
		if n == nil {
			r.Return(fmt.Errorf("unable to find node: %s", params.ID))
			return
		}
		com := n.Component(params.Component)
		n.RemoveComponent(com)
		if com.ID() == n.ID() {
			if err := s.State.Image.DestroyObjectPackage(n); err != nil {
				fmt.Println(err)
			}
		}
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) RemoveValue() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params SetValueParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		n := s.State.Root.FindChild(params.Path)
		if n == nil {
			r.Return(fmt.Errorf("unable to find node: %s", params.Path))
			return
		}
		localPath := params.Path[len(n.Path())+1:]
		f, _, err := n.GetField(localPath)
		if err != nil {
			r.Return(fmt.Errorf("unable to get field: %s", localPath))
			return
		}
		rv := reflect.ValueOf(f)
		idx := *params.IntValue
		a := rv.Slice(0, idx)
		b := rv.Slice(idx+1, rv.Len())
		rv = reflect.AppendSlice(a, b)
		n.SetField(localPath, rv.Interface())
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) AddValue() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params SetValueParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		n := s.State.Root.FindChild(params.Path)
		if n == nil {
			r.Return(fmt.Errorf("unable to find node: %s", params.Path))
			return
		}
		localPath := params.Path[len(n.Path())+1:]
		f, t, err := n.GetField(localPath)
		if err != nil {
			r.Return(fmt.Errorf("unable to get field: %s", localPath))
			return
		}
		rv := reflect.ValueOf(f)
		idx := rv.Len()
		nv := reflect.New(t.Elem())
		rv = reflect.Append(rv, reflect.Indirect(nv))
		n.SetField(localPath, rv.Interface())
		// NOTE: this must also done in setValue
		v := params.Value
		switch params.Type {
		case "time":
			v, err = time.Parse("15:04", v.(string))
			if err != nil {
				r.Return(err)
				return
			}
		case "date":
			v, err = time.Parse("2006-01-02", v.(string))
			if err != nil {
				r.Return(err)
				return
			}
		}
		n.SetField(filepath.Join(localPath, strconv.Itoa(idx)), v)
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) RefreshObject() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params NodeParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		n := s.State.Root.FindID(params.ID)
		if n == nil {
			r.Return(fmt.Errorf("unable to find node: %s", params.ID))
			return
		}
		if err := n.Refresh(); err != nil {
			r.Return(err)
			return
		}
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) ReloadComponent() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params RemoveComponentParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		n := s.State.Root.FindID(params.ID)
		if n == nil {
			r.Return(fmt.Errorf("unable to find node: %s", params.ID))
			return
		}
		com := n.Component(params.Component)
		if com != nil {
			if err := com.Reload(); err != nil {
				r.Return(err)
				return
			}
		}
		n.UpdateRegistry()
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) AddDelegate() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params NodeParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		obj := s.State.Root.FindID(params.ID)
		if obj == nil {
			r.Return(nil)
			return
		}
		r.Return(s.State.Image.CreateObjectPackage(obj))
	}
}

func (s *Service) LoadPrefab() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params AppendNodeParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		obj := s.State.Root.FindID(params.ID)
		if obj == nil {
			r.Return(nil)
			return
		}
		pf := prefab.LookupID(params.Name)
		if pf != nil {
			child := pf.New()
			obj.AppendChild(child)
		}
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) SelectNode() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var id string
		err := c.Decode(&id)
		if err != nil {
			r.Return(err)
			return
		}
		s.viewState.SelectedNode = id
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) UpdateNode() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params NodeParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		n := s.State.Root.FindID(params.ID)
		if n == nil {
			return
		}
		if params.Name != nil {
			n.SetName(*params.Name)
		}
		// if params.Active != nil {
		// 	n.Active = *params.Active
		// }
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) CallMethod() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var path string
		err := c.Decode(&path)
		if err != nil {
			r.Return(err)
			return
		}
		if path == "" {
			return
		}
		n := s.State.Root.FindChild(path)
		localPath := path[len(n.Path())+1:]
		// TODO: support args+ret
		n.CallMethod(localPath, nil, nil)
		s.updateView()
		r.Return(nil)
	}
}

// func (s *Service) SetExpression() func(qrpc.Responder, *qrpc.Call) {
// 	return func(r qrpc.Responder, c *qrpc.Call) {
// 		var params SetValueParams
// 		err := c.Decode(&params)
// 		if err != nil {
// 			r.Return(err)
// 			return
// 		}
// 		n := s.State.Root.FindChild(params.Path)
// 		localPath := params.Path[len(n.Path())+1:]
// 		// n.SetExpression(localPath, params.Value.(string))
// 		s.updateView()
// 		r.Return(nil)
// 	}
// }

func (s *Service) SetValue() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params SetValueParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		n := s.State.Root.FindChild(params.Path)
		// fmt.Println(params.Type)
		localPath := params.Path[len(n.Path())+1:]
		switch {
		case params.IntValue != nil:
			n.SetField(localPath, *params.IntValue)
		case params.RefValue != nil:
			refPath := filepath.Dir(*params.RefValue) // TODO: support subfields
			refNode := s.State.Root.FindChild(refPath)
			parts := strings.SplitN(localPath, "/", 2)
			refType := n.Component(parts[0]).FieldType(parts[1])
			if refNode != nil {
				typeSelector := (*params.RefValue)[len(refNode.Path())+1:]
				c := refNode.Component(typeSelector)
				if c != nil {
					n.SetField(localPath, c)
				} else {
					// interface reference
					ptr := reflect.New(refType)
					refNode.ValueTo(ptr)
					if ptr.IsValid() {
						n.SetField(localPath, reflect.Indirect(ptr).Interface())
					}
				}
			}
		default:
			// NOTE: this must also done in addValue
			v := params.Value
			var err error
			switch params.Type {
			case "time":
				v, err = time.Parse("15:04", v.(string))
				if err != nil {
					r.Return(err)
					return
				}
			case "date":
				v, err = time.Parse("2006-01-02", v.(string))
				if err != nil {
					r.Return(err)
					return
				}
			}
			n.SetField(localPath, v)
		}
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) AppendComponent() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params AppendNodeParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		if params.Name == "" {
			return
		}
		p := s.State.Root.FindID(params.ID)
		if p == nil {
			p = s.State.Root
		}
		v := library.Lookup(params.Name).New()
		p.AppendComponent(v)
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) DeleteNode() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var id string
		err := c.Decode(&id)
		if err != nil {
			r.Return(err)
			return
		}
		if id == "" {
			return
		}
		s.State.Root.RemoveID(id)
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) AppendNode() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params AppendNodeParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		if params.Name == "" {
			return
		}
		p := s.State.Root.FindID(params.ID)
		if p == nil {
			p = s.State.Root
		}
		n := object.New(params.Name)
		p.AppendChild(n)
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) MoveNode() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var params MoveNodeParams
		err := c.Decode(&params)
		if err != nil {
			r.Return(err)
			return
		}
		n := s.State.Root.FindID(params.ID)
		if n == nil {
			return
		}
		n.SetSiblingIndex(params.Index)
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) Subscribe() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		s.clients[c.Caller] = "state"
		s.updateView()
		r.Return(nil)
	}
}

func (s *Service) SelectProject() func(qrpc.Responder, *qrpc.Call) {
	return func(r qrpc.Responder, c *qrpc.Call) {
		var name string
		err := c.Decode(&name)
		if err != nil {
			r.Return(err)
			return
		}
		s.viewState.CurrentProject = name
		s.updateView()
		r.Return(nil)
	}
}
