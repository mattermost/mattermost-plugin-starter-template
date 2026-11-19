package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"
)

func TestGenerateGoManifestEscapesBackticks(t *testing.T) {
	manifestJSON := `{"description":"Use ` + "`" + `code` + "`" + `"}`

	generated := generateGoManifest(manifestJSON)
	file, err := parser.ParseFile(token.NewFileSet(), "manifest.go", generated, parser.AllErrors)
	if err != nil {
		t.Fatalf("generated Go manifest is invalid: %v", err)
	}

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Names) == 0 || valueSpec.Names[0].Name != "manifestStr" {
				continue
			}
			literal, ok := valueSpec.Values[0].(*ast.BasicLit)
			if !ok {
				t.Fatal("manifestStr is not a string literal")
			}
			got, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Fatalf("failed to unquote generated manifest: %v", err)
			}
			if got != manifestJSON {
				t.Fatalf("generated manifest = %q, want %q", got, manifestJSON)
			}
			return
		}
	}

	t.Fatal("generated Go manifest has no manifestStr constant")
}

func TestGenerateJSManifestEscapesTemplateLiteralCharacters(t *testing.T) {
	manifestJSON := `{"description":"Say ` + `\"` + `hold` + `\"` + ` using ` + "`" + `code` + "`" + `; cost ${price}; newline \n"}`
	want := `{"description":"Say ` + `\\"` + `hold` + `\\"` + ` using ` + `\` + "`" + `code` + `\` + "`" + `; cost ` + `\${price}` + `; newline ` + `\\n` + `"}`

	generated := string(generateJSManifest(manifestJSON))
	if !strings.Contains(generated, want) {
		t.Fatalf("generated JavaScript does not contain escaped manifest JSON:\n%s", generated)
	}
}
