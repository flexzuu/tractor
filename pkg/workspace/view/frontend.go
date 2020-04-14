package view

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/library"
	"github.com/manifold/tractor/pkg/manifold/prefab"

	//"github.com/manifold/tractor/pkg/repl"

	reflected "github.com/progrium/prototypes/go-reflected"
)

type Field struct {
	Type       string      `msgpack:"type"`
	Name       string      `msgpack:"name"`
	Path       string      `msgpack:"path"`
	Value      interface{} `msgpack:"value"`
	Expression *string     `msgpack:"expression"`
	Fields     []Field     `msgpack:"fields"`
}

type Button struct {
	Name    string `msgpack:"name"`
	Path    string `msgpack:"path"`
	OnClick string `msgpack:"onclick"`
}

type Component struct {
	Name     string   `msgpack:"name"`
	Filepath string   `msgpack:"filepath"`
	Fields   []Field  `msgpack:"fields"`
	Buttons  []Button `msgpack:"buttons"`
	Related  []string `msgpack:"related"`
	CustomUI *Element `msgpack:"customUI"`
}

type Node struct {
	Name       string      `msgpack:"name"`
	Path       string      `msgpack:"path"`
	Dir        string      `msgpack:"dir"`
	ID         string      `msgpack:"id"`
	Index      int         `msgpack:"index"`
	Active     bool        `msgpack:"active"`
	Components []Component `msgpack:"components"`
}

type Project struct {
	Name string `msgpack:"name"`
	Path string `msgpack:"path"`
}

type Prefab struct {
	Name string `msgpack:"name"`
	ID   string `msgpack:"id"`
}

type ComponentType struct {
	Filepath string `msgpack:"filepath"`
	Name     string `msgpack:"name"`
}

type State struct {
	Projects        []Project         `msgpack:"projects"`
	CurrentProject  string            `msgpack:"currentProject"`
	Components      []ComponentType   `msgpack:"components"`
	Prefabs         []Prefab          `msgpack:"prefabs"`
	Hierarchy       []string          `msgpack:"hierarchy"`
	Nodes           map[string]Node   `msgpack:"nodes"`
	NodePaths       map[string]string `msgpack:"nodePaths"`
	SelectedNode    string            `msgpack:"selectedNode"`
	EditorsEndpoint string            `msgpack:"editorsEndpoint"`

	mu sync.Mutex
}

type Element struct {
	Name     string
	Attrs    Attrs
	Children []Element
}

func El(name string, attrs Attrs, children []Element) Element {
	return Element{name, attrs, children}
}

type Attrs map[string]string

func exportElem(v reflected.Value, path string, idx int, n manifold.Object) (Field, bool) {
	elemPath := path + "/" + strconv.Itoa(idx)
	switch v.Type().Kind() {
	case reflect.Bool:
		return Field{
			Path:  elemPath,
			Type:  "boolean",
			Value: v.Interface(),
		}, true
	case reflect.String:
		return Field{
			Path:  elemPath,
			Type:  "string",
			Value: v.Interface(),
		}, true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return Field{
			Path:  elemPath,
			Type:  "number",
			Value: v.Interface(),
		}, true
	default:
		return Field{}, false
	}
}

func exportField(o reflected.Value, field, path string, n manifold.Object) Field {
	var kind reflect.Kind
	if o.Type().Kind() == reflect.Struct {
		kind = o.Type().FieldType(field).Kind()
	} else {
		if !o.Get(field).IsValid() {
			kind = reflect.Invalid
		} else {
			kind = o.Get(field).Type().Kind()
		}
	}
	fieldPath := path + "/" + field
	// var expr *string
	// exprPath := fieldPath[len(n.FullPath())+1:]
	// if e := n.Expression(exprPath); e != "" {
	// 	expr = &e
	// }
	switch kind {
	case reflect.Invalid:
		return Field{
			Name: field,
			Path: fieldPath,
			// Expression: expr,
			Type:  "string",
			Value: "INVALID",
		}
	case reflect.Bool:
		return Field{
			Name: field,
			Path: fieldPath,
			// Expression: expr,
			Type:  "boolean",
			Value: o.Get(field).Interface(),
		}
	case reflect.String:
		return Field{
			Name: field,
			Path: fieldPath,
			// Expression: expr,
			Type:  "string",
			Value: o.Get(field).Interface(),
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return Field{
			Path: fieldPath,
			Name: field,
			// Expression: expr,
			Type:  "number",
			Value: o.Get(field).Interface(),
		}
	case reflect.Struct:
		var fields []Field
		v := o.Get(field)
		for _, f := range v.Type().Fields() {
			fields = append(fields, exportField(v, f, fieldPath, n))
		}
		return Field{
			Path: fieldPath,
			Name: field,
			// Expression: expr,
			Type:   "struct",
			Fields: fields,
		}
	case reflect.Map:
		var fields []Field
		v := o.Get(field)
		for _, f := range v.Keys() {
			fields = append(fields, exportField(v, f, fieldPath, n))
		}
		return Field{
			Path: fieldPath,
			Name: field,
			// Expression: expr,
			Type:   "map",
			Fields: fields,
		}
	case reflect.Slice:
		var fields []Field
		v := o.Get(field)
		for idx, e := range v.Iter() {
			f, ok := exportElem(e, fieldPath, idx, n)
			if !ok {
				return Field{
					Name: field,
					Path: fieldPath,
					// Expression: expr,
					Type:  "string",
					Value: "UNSUPPORTED SLICE",
				}
			}
			fields = append(fields, f)
		}
		return Field{
			Path: fieldPath,
			Name: field,
			// Expression: expr,
			Type:   "array",
			Fields: fields,
		}
	case reflect.Ptr, reflect.Interface:
		var v interface{}
		if o.Get(field).IsValid() {
			v = o.Get(field).Interface()
		}
		t := o.Type().FieldType(field)
		if kind == reflect.Ptr {
			t = reflected.Type{Type: t.Elem()}
		}
		var path string
		if v != nil {
			refNode := n.Root().FindPointer(v)
			if refNode != nil {
				path = refNode.Path()
			}
		}
		return Field{
			Path: fieldPath,
			Name: field,
			// Expression: expr,
			Type:  fmt.Sprintf("reference:%s", t.Name()),
			Value: path,
		}
	default:
		//log.Println("TODO: unknown type:", o, field)
		return Field{
			Name: field,
			Path: fieldPath,
			// Expression: expr,
			Type:  "string",
			Value: "INVALID",
		}
	}
}

type ButtonProvider interface {
	InspectorButtons() []Button
}

type UIProvider interface {
	InspectorUI() Element
}

func (s *State) Update(root manifold.Object) {
	// reset/clear nodes
	s.Hierarchy = []string{}
	s.Nodes = make(map[string]Node)
	// walk every object in the tree
	manifold.Walk(root, func(n manifold.Object) {
		// all the nodes paths
		s.Hierarchy = append(s.Hierarchy, n.Path())
		// start a node struct based on node passed in
		node := Node{
			Name:   n.Name(),
			Active: true,
			// Dir:        n.Dir,
			Path:       n.Path(),
			Index:      n.SiblingIndex(),
			ID:         n.ID(),
			Components: []Component{},
		}
		for _, com := range n.Components() {
			// get all the fields of the component
			var fields []Field
			c := reflected.ValueOf(com.Pointer())
			path := n.Path() + "/" + com.Name()
			hiddenFields := c.Type().FieldsTagged("tractor", "hidden")
			for _, field := range c.Type().Fields() {
				if strInSlice(hiddenFields, field) {
					continue
				}
				fields = append(fields, exportField(c, field, path, n))
			}
			// see if component provides custom ui
			var ui *Element
			uip, ok := com.Pointer().(UIProvider)
			if ok {
				v := uip.InspectorUI()
				ui = &v
			}
			// see if the component provides buttons
			var buttons []Button
			p, ok := com.Pointer().(ButtonProvider)
			if ok {
				buttons = p.InspectorButtons()
				for idx, button := range buttons {
					if button.OnClick != "" {
						continue
					}
					typ := reflect.ValueOf(com.Pointer()).Type()
					for i := 0; i < typ.NumMethod(); i++ {
						method := typ.Method(i)
						if method.Name != button.Name {
							continue
						}
						if method.Type.NumIn() == 1 {
							buttons[idx].Path = path + "/" + method.Name
							break
						}
					}
				}
			}

			// look up the filepath for this component
			var filepath string
			if com.ID() != "" {
				rc := library.LookupID(com.ID())
				if rc == nil {
					panic("component ID not in library: " + com.ID())
				}
				filepath = rc.Filepath
			} else {
				rc := library.Lookup(com.Name())
				if rc == nil {
					panic("component name not in library: " + com.Name())
				}
				filepath = rc.Filepath
			}

			// look up related components to this component
			var related []string
			for _, rc := range library.Related(library.Lookup(com.Name())) {
				related = append(related, rc.Name)
			}

			// add component to frontend node's components
			node.Components = append(node.Components, Component{
				Name:     com.Name(),
				Filepath: filepath,
				Fields:   fields,
				Buttons:  buttons,
				Related:  related,
				CustomUI: ui,
			})
		}
		// add the node to state
		s.mu.Lock()
		s.Nodes[n.ID()] = node
		s.NodePaths[n.Path()] = n.ID()
		s.mu.Unlock()
	})
}

func New(root manifold.Object) *State {
	state := &State{
		Projects:       []Project{},
		CurrentProject: "dev",
		Nodes:          make(map[string]Node),
		NodePaths:      make(map[string]string),
	}
	for _, com := range library.Registered() {
		state.Components = append(state.Components, ComponentType{
			Name:     com.Name,
			Filepath: com.Filepath,
		})
	}
	for _, pf := range prefab.Registered() {
		state.Prefabs = append(state.Prefabs, Prefab{
			Name: pf.Name,
			ID:   pf.ID,
		})
	}
	state.Update(root)
	return state
}

func strInSlice(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}
