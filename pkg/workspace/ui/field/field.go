package field

import (
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/misc/reflectutil"
	"github.com/manifold/tractor/pkg/workspace/ui"
)

func FromComponent(com manifold.Component) (fields []ui.Field) {
	obj := com.Container()
	rc := reflect.Indirect(reflect.ValueOf(com.Pointer()))
	path := filepath.Join(obj.Path(), com.Name())
	hiddenFields := reflectutil.FieldsTagged(rc.Type(), "tractor", "hidden")
	for _, field := range reflectutil.Fields(rc.Type()) {
		if strInSlice(hiddenFields, field) {
			continue
		}
		fields = append(fields, exportField(rc, field, path, obj))
	}
	return
}

func exportField(o reflect.Value, field, path string, obj manifold.Object) ui.Field {
	path = filepath.Join(path, field)
	switch reflectutil.MemberKind(o, field) {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return numberField(path, reflectutil.Get(o, field))
	case reflect.Bool:
		return booleanField(path, reflectutil.Get(o, field))
	case reflect.String:
		return stringField(path, reflectutil.Get(o, field))
	case reflect.Struct:
		return structField(path, reflectutil.Get(o, field), obj)
	case reflect.Map:
		return mapField(path, reflectutil.Get(o, field), obj)
	case reflect.Slice:
		return arrayField(path, reflectutil.Get(o, field), obj)
	case reflect.Ptr, reflect.Interface:
		return referenceField(path, reflectutil.Get(o, field), obj)
	case reflect.Invalid:
		return invalidField(path)
	default:
		return unknownField(path)
	}
}

func exportElem(v reflect.Value, path string, idx int, n manifold.Object) (ui.Field, bool) {
	elemPath := filepath.Join(path, strconv.Itoa(idx))
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return numberField(elemPath, v), true
	case reflect.Bool:
		return booleanField(elemPath, v), true
	case reflect.String:
		return stringField(elemPath, v), true
	default:
		return ui.Field{}, false
	}
}

func strInSlice(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}
