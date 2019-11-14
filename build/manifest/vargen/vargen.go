// Package vargen generates source code to represent contents of runtime variables.
// variables are usually custom structs with nested values. but they can be other
// built-in types too.
//
// exceptions:
// * variables' nested values also can be custom structs or built-in types but they
//   cannot contain other variables. for ex:
//     `var myVar = struct{ A string }{ myStringVar }` -> is not valid.
//     `var myVar = struct{ A string }{ "my-string" }` -> is valid.
//     `var myVar = struct{ B *B }{ &B{} }` -> is valid.
//
// * funcs and private struct fields are not supported.
//     `var myVar = struct{ C func() }{ func(){} }` -> is not valid.
//     `var myVar = struct{ d string }{ "my-other-string" }` -> is not valid.
package vargen

import (
	"bytes"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

// Vargen is a variable code generator.
type Vargen struct {
	// name of the generated variable.
	name string

	// err is the last happened error.
	err error

	st *jen.Statement
}

// Generate generates a source code for v with variable name.
func Generate(name string, v interface{}) *Vargen {
	vg := &Vargen{name: name}
	vg.st, vg.err = vg.generate(v)
	return vg
}

// Render renders code generation to w.
func (v *Vargen) Render(w io.Writer) error {
	if v.err != nil {
		return v.err
	}
	return v.st.Render(w)
}

// String satisfies fmt.Stringer interface and returns the generated code as a string.
// String() will return a string error if something went wrong during the code generation
// or rendering.
func (v *Vargen) String() string {
	var buf bytes.Buffer
	if err := v.Render(&buf); err != nil {
		return errors.Wrap(err, "error while rendering").Error()
	}
	return buf.String()
}
