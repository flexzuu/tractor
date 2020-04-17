package field

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/misc/reflectutil"
	"github.com/manifold/tractor/pkg/workspace/ui"
)

type Enumer interface {
	Enum() []string
}

func booleanField(path string, value reflect.Value) ui.Field {
	return ui.Field{
		Type:  "boolean",
		Name:  filepath.Base(path),
		Path:  path,
		Value: value.Interface(),
	}
}

func stringField(path string, value reflect.Value) ui.Field {
	v := value.Interface()
	var enum []string
	e, ok := v.(Enumer)
	if ok {
		enum = e.Enum()
	}
	return ui.Field{
		Type:  "string",
		Name:  filepath.Base(path),
		Path:  path,
		Value: v,
		Enum:  enum,
	}
}

func numberField(path string, value reflect.Value) ui.Field {
	return ui.Field{
		Type:  "number",
		Name:  filepath.Base(path),
		Path:  path,
		Value: value.Interface(),
	}
}

func invalidField(path string) ui.Field {
	return ui.Field{
		Type:  "string",
		Name:  filepath.Base(path),
		Path:  path,
		Value: "INVALID",
	}
}

func unknownField(path string) ui.Field {
	return ui.Field{
		Type:  "string",
		Name:  filepath.Base(path),
		Path:  path,
		Value: "UNKNOWN",
	}
}

func structField(path string, value reflect.Value, obj manifold.Object) ui.Field {
	var fields []ui.Field
	for _, f := range reflectutil.Fields(value.Type()) {
		fields = append(fields, exportField(value, f, path, obj))
	}
	return ui.Field{
		Type:   "struct",
		Name:   filepath.Base(path),
		Path:   path,
		Fields: fields,
	}
}

func mapField(path string, value reflect.Value, obj manifold.Object) ui.Field {
	var fields []ui.Field
	for _, f := range reflectutil.Keys(value) {
		fields = append(fields, exportField(value, f, path, obj))
	}
	return ui.Field{
		Type:   "map",
		Name:   filepath.Base(path),
		Path:   path,
		Fields: fields,
	}
}

func arrayField(path string, value reflect.Value, obj manifold.Object) ui.Field {
	var fields []ui.Field
	for idx, e := range reflectutil.Values(value) {
		f, ok := exportElem(e, path, idx, obj)
		if !ok {
			log.Println("sliceField: unsupported elements")
			fields = []ui.Field{}
			break
		}
		fields = append(fields, f)
	}
	return ui.Field{
		Type:   "array",
		Path:   path,
		Name:   filepath.Base(path),
		Fields: fields,
	}
}

func referenceField(path string, value reflect.Value, obj manifold.Object) ui.Field {
	if !value.IsValid() {
		panic("referenceField: invalid value used")
	}
	t := value.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var refPath string
	refNode := obj.Root().FindPointer(value.Interface())
	if refNode != nil {
		refPath = refNode.Path()
	}
	return ui.Field{
		Type:  fmt.Sprintf("reference:%s", t.Name()),
		Path:  path,
		Name:  filepath.Base(path),
		Value: refPath,
	}
}
