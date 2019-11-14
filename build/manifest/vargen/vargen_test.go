package vargen

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/require"
)

var manifest = &model.Manifest{
	Description:      "This plugin serves as a starting point for writing a Mattermost plugin.",
	Id:               "my-id",
	Name:             "Plugin Starter Template",
	MinServerVersion: "5.12.0",
	Server: &model.ManifestServer{Executables: &model.ManifestExecutables{
		DarwinAmd64:  "server/dist/plugin-darwin-amd64",
		LinuxAmd64:   "server/dist/plugin-linux-amd64",
		WindowsAmd64: "server/dist/plugin-windows-amd64.exe",
	}},
	Version: "0.1.0",
	Webapp:  &model.ManifestWebapp{BundlePath: "webapp/dist/main.js"},
}

func TestManifestFile(t *testing.T) {
	t.Parallel()
	vg := Generate("manifest", manifest)
	require.Equal(t, trimCode(t, `
	var manifest = &model.Manifest{
		Description:      "This plugin serves as a starting point for writing a Mattermost plugin.",
		Id:               "my-id",
		MinServerVersion: "5.12.0",
		Name:             "Plugin Starter Template",
		Server: &model.ManifestServer{Executables: &model.ManifestExecutables{
			DarwinAmd64:  "server/dist/plugin-darwin-amd64",
			LinuxAmd64:   "server/dist/plugin-linux-amd64",
			WindowsAmd64: "server/dist/plugin-windows-amd64.exe",
		}},
		Version: "0.1.0",
		Webapp:  &model.ManifestWebapp{BundlePath: "webapp/dist/main.js"},
	}	
	`), trimCode(t, vg.String()))
}

type A struct {
	B *B
	c string
	P bool
	Q float32
	R uint
	S string
	T time.Duration
	U int
	V *B
	W []*B
	X *[]*B
	Y map[string]interface{}
	Z map[time.Duration]bool
}

var a = &A{
	B: nil,
	c: "unexported",
	P: true,
	Q: 1.2,
	R: 3,
	S: "hey",
	T: time.Second,
	U: 4,
	V: &B{},
	W: []*B{},
	X: &[]*B{
		{1, "b"},
	},
	Y: map[string]interface{}{
		"a": 1,
		"b": "two",
		"c": []B{},
		"d": []interface{}{
			map[int]bool{
				1: false,
				2: true,
			},
			struct {
				A string
				B uint32
			}{"a", 5},
		},
	},
	Z: map[time.Duration]bool{
		time.Nanosecond: false,
		time.Minute * 2: true,
	},
}

type B struct {
	A uint16
	B string
}

func TestComplicated(t *testing.T) {
	t.Parallel()
	vg := Generate("a", a)
	require.Equal(t, trimCode(t, `
	var a = &vargen.A{
		P: true,
		Q: 1.2,
		R: 3,
		S: "hey",
		T: 1000000000,
		U: 4,
		V: &vargen.B{},
		W: []*vargen.B{},
		X: &[]*vargen.B{&vargen.B{
			A: 1,
			B: "b",
		}},
		Y: map[string]interface{}{
			"a": 1,
			"b": "two",
			"c": []vargen.B{},
			"d": []interface{}{map[int]bool{
				1: false,
				2: true,
			}, struct {
				A string
				B uint32
			}{
				A: "a",
				B: 5,
			}},
		},
		Z: map[time.Duration]bool{
			1:            false,
			120000000000: true,
		},
	}
	`), trimCode(t, vg.String()))
}

// TODO(ilgooz): add more test cases.
var testAll = []struct {
	variableName    string
	variableContent interface{}
	expectedCode    string
	err             error
}{
	{"a", "x", `var a = "x"`, nil},
	{"b", map[string]bool{"x": true}, `var b = map[string]bool{"x": true}`, nil},
	{"b", map[string]interface{}{"x": func() {}}, ``, errors.New(`value kind "func" is not supported`)},
	{"b", map[string]interface{}{"x": struct{ A string }{"y"}}, `var b = map[string]interface{}{"x": struct{A string}{A: "y"}}`, nil},
}

func TestAll(t *testing.T) {
	t.Parallel()
	for _, tt := range testAll {
		t.Run(tt.variableName, func(t *testing.T) {
			var buf bytes.Buffer
			err := Generate(tt.variableName, tt.variableContent).Render(&buf)
			if tt.err == nil {
				require.NoError(t, err)
				require.Equal(t, trimCode(t, tt.expectedCode), trimCode(t, buf.String()))
			} else {
				require.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func trimCode(t *testing.T, code string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(code)), "")
}
