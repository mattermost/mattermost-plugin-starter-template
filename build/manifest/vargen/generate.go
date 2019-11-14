package vargen

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"
)

// generate creates a *jen.Statement to generate code for vv value.
func (v *Vargen) generate(vv interface{}) (st *jen.Statement, err error) {
	// if vv is nil there is nothing to generate.
	if vv == nil {
		return nil, errors.New("variable cannot be nil")
	}
	// panicing inside the walkX funcs and recovering here is better since
	// this makes it easier to build recursive funcs and work with chainable APIs.
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// initiate code generation by starting to write var keyword and variable's
	// name and walk through the variable's content to create a source code.
	st = jen.
		Var().                              // write keyword `var`.
		Id(v.name).                         // write var's name.
		Add(jen.Op("=")).                   // write operator `=`.
		Add(walkValue(reflect.ValueOf(vv))) // write var's content.
	return st, nil
}

// walkValue v value to do code generation.
func walkValue(v reflect.Value) *jen.Statement {
	switch v.Kind() {
	case reflect.Ptr:
		return jen.Op("&").
			Add(walkValue(v.Elem()))
	case reflect.Interface:
		return walkValue(v.Elem())
	case reflect.String:
		return walkString(v)
	case reflect.Bool:
		return walkBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return walkInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return walkUint(v)
	case reflect.Float32, reflect.Float64:
		return walkFloat(v)
	case reflect.Map:
		return walkMap(v)
	case reflect.Slice, reflect.Array:
		return walkSlice(v)
	case reflect.Struct:
		return walkStruct(v)
	}
	// only the above kinds are supported, remaning not.
	// for ex. if a struct has a field inside with a func() value, this panic will hit.
	// basically, we only support custom struct types and most of the built-in types.
	panic(fmt.Errorf("value kind %q is not supported", v.Kind()))
}

// walkBool walks through the string v and generates code for it.
func walkString(v reflect.Value) *jen.Statement {
	// jen.Lit() adds "" around string, this is why use Lit() here.
	return jen.Lit(v.String())
}

// walkBool walks through the boolean v and generates code for it.
func walkBool(v reflect.Value) *jen.Statement {
	return jen.Id(fmt.Sprint(v.Bool()))
}

// walkInt walks through the intX v and generates code for it.
// please note that it's better to not use v.String() here in order to not
// produce a code like `int(1)`. `1` is preferred which is better for human reading.
func walkInt(v reflect.Value) *jen.Statement {
	return jen.Id(fmt.Sprint(v.Int()))
}

// walkUint walks through the uintX v and generates code for it.
// please note that it's better to not use v.String() here in order to not
// produce a code like `uint(1)`. `1` is preferred which is better for human reading.
func walkUint(v reflect.Value) *jen.Statement {
	return jen.Id(fmt.Sprint(v.Uint()))
}

// walkFloat walks through the floatX v and generates code for it.
// please note that it's better to not use v.String() here in order to not
// produce a code like `float32(1.2)`. `1.2` is preferred which is better for human reading.
func walkFloat(v reflect.Value) *jen.Statement {
	// fmt.Sprint handles creating fixed point notation for us.
	return jen.Id(fmt.Sprint(v.Interface()))
}

// walkMap walks through the map type v and generates code for it.
func walkMap(v reflect.Value) *jen.Statement {
	t := v.Type()
	keyType := jen.Id(t.Key().String())
	valueType := t.Elem().String()
	keys := v.MapKeys()
	values := jen.Dict{}
	for i := 0; i < v.Len(); i++ {
		key := keys[i]
		value := v.MapIndex(key)
		values[walkValue(key)] = walkValue(value)
	}
	return jen.
		Map(keyType).  // writing a map.
		Id(valueType). // write map's type. for ex: map[string]string
		Values(values) // write map's value.
}

// walkSlice walks through the slice or array type v and generates code for it.
func walkSlice(v reflect.Value) *jen.Statement {
	items := []jen.Code{}
	for i := 0; i < v.Len(); i++ {
		vv := v.Index(i)
		items = append(items, walkValue(vv))
	}
	valueType := v.Type().Elem().String()
	return jen.
		Index().         // writing a slice or array.
		Id(valueType).   // write slice's type. for ex: []string.
		Values(items...) // write slice's values.
}

// walkStruct walks through the struct type v and generates code for it.
func walkStruct(v reflect.Value) *jen.Statement {
	t := v.Type()
	structType := v.Type().String()
	fields := jen.Dict{}
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		// skip field if unexported.
		if ft.PkgPath != "" {
			continue
		}
		fv := v.Field(i)
		// skip field if it has a zero value.
		if fv.IsZero() {
			continue
		}
		fields[jen.Id(ft.Name)] = walkValue(fv)
	}
	return jen.
		Id(structType). // start writing a struct and write it's type.
		Values(fields)  // write struct's values.
}
