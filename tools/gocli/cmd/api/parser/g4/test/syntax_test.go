package test

import (
	ast2 "github.com/cocoup/go-smart/tools/gocli/cmd/api/parser/g4/ast"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/parser/g4/gen/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

var syntaxAccept = func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
	return p.SyntaxLit().Accept(visitor)
}

func TestSyntax(t *testing.T) {
	t.Run("matched", func(t *testing.T) {
		v, err := parser.Accept(syntaxAccept, `syntax = "v1"`)
		assert.Nil(t, err)

		syntax := v.(*ast2.SyntaxExpr)
		assert.True(t, syntax.Equal(&ast2.SyntaxExpr{
			Syntax:  ast2.NewTextExpr("syntax"),
			Assign:  ast2.NewTextExpr("="),
			Version: ast2.NewTextExpr(`"v1"`),
		}))
	})

	t.Run("expecting syntax", func(t *testing.T) {
		_, err := parser.Accept(syntaxAccept, `= "v1"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syn = "v1"`)
		assert.Error(t, err)
	})

	t.Run("missing assign", func(t *testing.T) {
		_, err := parser.Accept(syntaxAccept, `syntax  "v1"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syntax + "v1"`)
		assert.Error(t, err)
	})

	t.Run("mismatched version", func(t *testing.T) {
		_, err := parser.Accept(syntaxAccept, `syntax="v0"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syntax = "v1a"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syntax = "vv1"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syntax = "1"`)
		assert.Error(t, err)
	})

	t.Run("with comment", func(t *testing.T) {
		v, err := parser.Accept(syntaxAccept, `
		// doc
		syntax="v1" // line comment`)
		assert.Nil(t, err)

		syntax := v.(*ast2.SyntaxExpr)
		assert.True(t, syntax.Equal(&ast2.SyntaxExpr{
			Syntax:  ast2.NewTextExpr("syntax"),
			Assign:  ast2.NewTextExpr("="),
			Version: ast2.NewTextExpr(`"v1"`),
			DocExpr: []ast2.Expr{
				ast2.NewTextExpr("// doc"),
			},
			CommentExpr: ast2.NewTextExpr("// line comment"),
		}))
	})
}
