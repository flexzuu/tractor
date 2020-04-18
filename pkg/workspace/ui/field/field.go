package field

import (
	"fmt"
	"math"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/misc/reflectutil"
	"github.com/manifold/tractor/pkg/workspace/ui"
)

type Enumer interface {
	Enum() []string
}

type field struct {
	ui.Field

	v        interface{}
	rv       reflect.Value
	parent   reflect.Value
	basepath string
	obj      manifold.Object
}

var typeBuilders map[string]func(*field)

func init() {
	typeBuilders = map[string]func(*field){
		"invalid": func(f *field) {
			f.Type = "string"
			f.Value = "INVALID"
		},
		"unknown": func(f *field) {
			f.Type = "string"
			f.Value = "UNKNOWN"
		},
		"boolean": func(f *field) {
			f.Value = f.v
		},
		"number": func(f *field) {
			f.Value = f.v
			applyMinMax(&(f.Field), f.rv)
		},
		"string": func(f *field) {
			e, ok := f.v.(Enumer)
			if ok {
				f.Enum = e.Enum()
			}
			f.Value = f.v
		},
		"time": func(f *field) {
			t := f.v.(time.Time)
			if t.IsZero() {
				t = time.Now()
			}
			f.Value = f.v.(time.Time).Format("15:04")
		},
		"date": func(f *field) {
			t := f.v.(time.Time)
			if t.IsZero() {
				t = time.Now()
			}
			f.Value = t.Format("2006-01-02")
		},
		"struct": func(f *field) {
			for _, fieldname := range reflectutil.Fields(f.rv.Type()) {
				f.Fields = append(f.Fields, subField(f, fieldname))
			}
		},
		"map": func(f *field) {
			f.SubType = collectionSubtypeField(f.rv)
			for _, keyname := range reflectutil.Keys(f.rv) {
				f.Fields = append(f.Fields, subField(f, keyname))
			}
		},
		"array": func(f *field) {
			f.SubType = collectionSubtypeField(f.rv)
			for idx, e := range reflectutil.Values(f.rv) {
				f.Fields = append(f.Fields, subFieldElem(f, idx, e))
			}
		},
		"reference": func(f *field) {
			if !f.rv.IsValid() {
				panic("referenceField: invalid value used")
			}
			t := f.rv.Type()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			var refPath string
			refNode := f.obj.Root().FindPointer(f.v)
			if refNode != nil {
				refPath = refNode.Path()
			}
			f.Type = fmt.Sprintf("reference:%s", t.Name()) // TODO: stop doing this
			f.SubType = &ui.Field{
				Type: t.Name(),
			}
			f.Value = refPath
		},
	}

}

func FromComponent(com manifold.Component) (fields []ui.Field) {
	obj := com.Container()
	rv := reflect.Indirect(reflect.ValueOf(com.Pointer()))
	path := filepath.Join(obj.Path(), com.Name())
	hiddenFields := reflectutil.FieldsTagged(rv.Type(), "tractor", "hidden")
	for _, fieldname := range reflectutil.Fields(rv.Type()) {
		if strInSlice(hiddenFields, fieldname) {
			continue
		}
		fields = append(fields, newField(obj, rv, path, fieldname))
	}
	return
}

func collectionSubtypeField(rv reflect.Value) *ui.Field {
	f := &ui.Field{
		Type: typeFromKind(rv.Type().Elem().Kind()),
	}
	applyMinMax(f, rv)
	e, ok := rv.Interface().(Enumer)
	if ok {
		f.Enum = e.Enum()
	}
	return f
}

func applyMinMax(f *ui.Field, rv reflect.Value) {
	switch rv.Type().Kind() {
	case reflect.Int, reflect.Int64:
		f.Min = math.MinInt64
		f.Max = math.MaxInt64
	case reflect.Int8:
		f.Min = math.MinInt8
		f.Max = math.MaxInt8
	case reflect.Int16:
		f.Min = math.MinInt16
		f.Max = math.MaxInt16
	case reflect.Int32:
		f.Min = math.MinInt32
		f.Max = math.MaxInt32
	case reflect.Uint, reflect.Uint64:
		f.Max = math.MaxUint64
	case reflect.Uint8:
		f.Max = math.MaxUint8
	case reflect.Uint16:
		f.Max = math.MaxUint16
	case reflect.Uint32:
		f.Max = math.MaxUint32
	case reflect.Float32:
		f.Max = math.MaxUint32 // math.MaxFloat32 someday
	case reflect.Float64:
		f.Max = math.MaxUint64 // math.MaxFloat64 someday
	}
}

func newField(obj manifold.Object, parent reflect.Value, basepath, fieldname string) ui.Field {
	return subField(&field{
		Field: ui.Field{
			Path: basepath,
		},
		rv:  parent,
		obj: obj,
	}, fieldname)
}

func subField(f *field, fieldname string) ui.Field {
	sf := &field{
		Field: ui.Field{
			Type: typeFromKind(reflectutil.MemberKind(f.rv, fieldname)),
			Path: filepath.Join(f.Path, fieldname),
		},
		parent: f.rv,
		obj:    f.obj,
		rv:     reflectutil.Get(f.rv, fieldname),
	}
	sf.Name = filepath.Base(sf.Path)
	sf.v = sf.rv.Interface()

	if f.rv.Kind() == reflect.Struct {
		structfield, _ := f.rv.Type().FieldByName(fieldname)
		if structfield.Tag.Get("field") != "" {
			sf.Type = structfield.Tag.Get("field")
		}
	}
	builder, ok := typeBuilders[sf.Type]
	if !ok {
		builder = typeBuilders["unknown"]
	}
	builder(sf)
	return sf.Field
}

func subFieldElem(f *field, idx int, value reflect.Value) ui.Field {
	sf := &field{
		Field: ui.Field{
			Type: typeFromKind(value.Type().Kind()),
			Path: filepath.Join(f.Path, strconv.Itoa(idx)),
		},
		parent: f.rv,
		obj:    f.obj,
		rv:     value,
		v:      value.Interface(),
	}
	sf.Name = filepath.Base(sf.Path)
	builder, ok := typeBuilders[sf.Type]
	if !ok {
		builder = typeBuilders["unknown"]
	}
	builder(sf)
	return sf.Field
}

func typeFromKind(kind reflect.Kind) string {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.String:
		return "string"
	case reflect.Struct:
		return "struct"
	case reflect.Map:
		return "map"
	case reflect.Slice:
		return "array"
	case reflect.Ptr, reflect.Interface:
		return "reference"
	case reflect.Invalid:
		return "invalid"
	default:
		return "unknown"
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
