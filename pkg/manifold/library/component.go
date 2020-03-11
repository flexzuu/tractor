package library

import (
	"fmt"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/misc/jsonpointer"
	"github.com/manifold/tractor/pkg/misc/notify"
	"github.com/mitchellh/mapstructure"
	reflected "github.com/progrium/prototypes/go-reflected"
)

type component struct {
	object  manifold.Object
	name    string
	id      string
	enabled bool
	value   interface{}
	typed   bool
	loaded  bool // if Reload has been called once
}

type ComponentEnabler interface {
	ComponentEnable()
}
type ComponentDisabler interface {
	ComponentDisable()
}

type ChildProvider interface {
	ChildNodes() []manifold.Object
}

func newComponent(name string, value interface{}, id string) *component {
	if id == "" {
		// TODO: make up mind on how to use component IDs
		//id = xid.New().String()
	}
	var typed bool
	if _, ok := value.(map[string]interface{}); !ok {
		typed = true
	}
	return &component{
		name:    name,
		enabled: true,
		value:   value,
		id:      id,
		typed:   typed,
	}
}

func NewComponent(name string, value interface{}, id string) manifold.Component {
	return newComponent(name, value, id)
}

func (c *component) GetField(path string) (interface{}, reflect.Type, error) {
	// TODO: check if field exists
	return jsonpointer.Reflect(c.Pointer(), path), c.FieldType(path), nil
}

func (c *component) SetField(path string, value interface{}) error {
	old, _, _ := c.GetField(path)
	if old == value {
		return nil
	}
	jsonpointer.SetReflect(c.value, path, value)
	notify.Send(c.object, manifold.ObjectChange{
		Object: c.object,
		Path:   fmt.Sprintf("%s/%s", c.name, path),
		Old:    old,
		New:    value,
	})
	return nil
}

func (c *component) FieldType(path string) reflect.Type {
	parts := strings.Split(path, "/")
	rt := reflect.TypeOf(c.Pointer())
	for _, part := range parts {
		switch rt.Kind() {
		case reflect.Slice:
			rt = rt.Elem()
		case reflect.Struct:
			field, _ := rt.FieldByName(part)
			rt = field.Type
		case reflect.Ptr:
			// TODO: might not be a struct!
			field, _ := rt.Elem().FieldByName(part)
			rt = field.Type
		default:
			panic("unhandled type: " + rt.String())
		}
	}
	return rt
}

func (c *component) CallMethod(path string, args []interface{}, reply interface{}) error {
	// TODO: support methods on sub paths / data structures
	rval := reflect.ValueOf(c.Pointer())
	method := rval.MethodByName(path)
	var params []reflect.Value
	for _, arg := range args {
		params = append(params, reflect.ValueOf(arg))
	}
	retVals := method.Call(params)
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	// assuming up to 2 return values, one being an error
	rreply := reflect.ValueOf(reply)
	var errVal error
	for _, v := range retVals {
		if v.Type().Implements(errorInterface) {
			if !v.IsNil() {
				errVal = v.Interface().(error)
			}
		} else {
			if reply != nil {
				rreply.Elem().Set(v)
			}
		}
	}
	return errVal
}

func (c *component) Index() int {
	for idx, com := range c.object.Components() {
		if com == c {
			return idx
		}
	}
	return 0
}

func (c *component) SetIndex(idx int) {
	if idx == -1 {
		idx = len(c.object.Components()) - 1
	}
	old := c.Index()
	if old == idx {
		return
	}
	c.object.RemoveComponent(c)
	c.object.InsertComponentAt(idx, c)
	notify.Send(c.object, manifold.ObjectChange{
		Object: c.object,
		Path:   fmt.Sprintf("%s/::Index", c.name),
		Old:    old,
		New:    idx,
	})
}

func (c *component) Name() string {
	return c.name
}

func (c *component) ID() string {
	return c.id
}

func (c *component) Enabled() bool {
	return c.enabled
}

func (c *component) SetEnabled(enable bool) {
	old := c.enabled
	if old == enable {
		return
	}
	c.enabled = enable
	if enable {
		if e, ok := c.Pointer().(ComponentEnabler); ok {
			e.ComponentEnable()
		}
	} else {
		if e, ok := c.Pointer().(ComponentDisabler); ok {
			e.ComponentDisable()
		}
	}
	notify.Send(c.object, manifold.ObjectChange{
		Object: c.object,
		Path:   fmt.Sprintf("%s/::Enabled", c.name),
		Old:    old,
		New:    enable,
	})

}

func (c *component) Container() manifold.Object {
	return c.object
}

func (c *component) SetContainer(obj manifold.Object) {
	c.object = obj
}

// TODO: rename to Value()?
func (c *component) Pointer() interface{} {
	if !c.typed {
		c.value = typedComponentValue(c.value, c.name, c.id)
		c.typed = true
	}
	return c.value
}

func (c *component) Type() reflect.Type {
	return reflect.TypeOf(c.Pointer())
}

func (c *component) Reload() error {
	if c.loaded && c.enabled {
		if e, ok := c.Pointer().(ComponentDisabler); ok {
			e.ComponentDisable()
		}
	}
	if e, ok := c.Pointer().(ComponentEnabler); ok {
		e.ComponentEnable()
	}
	if len(c.object.Children()) == 0 {
		if cp, ok := c.Pointer().(ChildProvider); ok {
			for _, obj := range cp.ChildNodes() {
				c.object.AppendChild(obj)
			}
		}
	}
	c.SetEnabled(true)
	c.loaded = true
	return nil
}

// TODO
func (c *component) Fields() {}

// TODO
func (c *component) Methods() {}

// TODO
func (c *component) RelatedPrefabs() {}

func (c *component) Snapshot() manifold.ComponentSnapshot {
	if !c.typed {
		panic("snapshot before component value is typed")
	}
	com := manifold.ComponentSnapshot{
		Name:    c.name,
		ID:      c.id,
		Enabled: c.enabled,
		Value:   c.value,
	}
	if c.object != nil {
		com.ObjectID = c.object.ID()
		com.Value, com.Refs = extractRefs(c.object, com.Name, com.Value)
	}
	return com
}

func extractRefs(obj manifold.Object, basePath string, v interface{}) (out map[string]interface{}, refs []manifold.SnapshotRef) {
	if obj.Root() == nil {
		return
	}
	out = make(map[string]interface{})
	rv := reflected.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return
	}
	rt := rv.Type()
	for _, field := range rt.Fields() {
		ft := rt.FieldType(field)
		fieldPath := path.Join(basePath, field)
		var subrefs []manifold.SnapshotRef
		switch ft.Kind() {
		case reflect.Slice:
			out[field], subrefs = extractRefsSlice(obj, fieldPath, rv.Get(field).Interface())
			refs = append(refs, subrefs...)
		case reflect.Struct, reflect.Map:
			out[field], subrefs = extractRefs(obj, fieldPath, rv.Get(field).Interface())
			refs = append(refs, subrefs...)
		case reflect.Ptr, reflect.Interface:
			if rv.Get(field).IsNil() {
				continue
			}
			target := obj.Root().FindPointer(rv.Get(field).Interface())
			if target != nil {
				refs = append(refs, manifold.SnapshotRef{
					ObjectID: obj.ID(),
					Path:     fieldPath,
					TargetID: target.ID(),
				})
				out[field] = nil
			}
		default:
			out[field] = rv.Get(field).Interface()
		}
	}
	return
}

func extractRefsSlice(obj manifold.Object, basePath string, v interface{}) (out []interface{}, refs []manifold.SnapshotRef) {
	if obj.Root() == nil {
		return
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return
	}
	for i := 0; i < rv.Len(); i++ {
		field := rv.Index(i)
		ft := field.Type()
		fieldPath := path.Join(basePath, strconv.Itoa(i))
		var subrefs []manifold.SnapshotRef
		var vv interface{}
		switch ft.Kind() {
		case reflect.Slice:
			vv, subrefs = extractRefsSlice(obj, fieldPath, field.Interface())
			out = append(out, vv)
			refs = append(refs, subrefs...)
		case reflect.Struct, reflect.Map:
			vv, subrefs = extractRefs(obj, fieldPath, field.Interface())
			out = append(out, vv)
			refs = append(refs, subrefs...)
		case reflect.Ptr, reflect.Interface:
			if field.IsNil() {
				continue
			}
			target := obj.Root().FindPointer(field.Interface())
			if target != nil {
				refs = append(refs, manifold.SnapshotRef{
					ObjectID: obj.ID(),
					Path:     fieldPath,
					TargetID: target.ID(),
				})
				out = append(out, nil)
			}
		default:
			out = append(out, field.Interface())
		}
	}
	return
}

func typedComponentValue(value interface{}, name, id string) interface{} {
	var typedValue interface{}
	if id == "" {
		if rc := Lookup(name); rc != nil {
			typedValue = rc.NewValue()
		}
	} else {
		if rc := LookupID(id); rc != nil {
			typedValue = rc.NewValue()
		}
	}
	if typedValue == nil {
		panic("unable to find registered component: " + name)
	}
	if err := mapstructure.Decode(value, typedValue); err == nil {
		return typedValue
	} else {
		panic(err)
	}
}
