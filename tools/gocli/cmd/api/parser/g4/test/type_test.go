package test

import (
	ast2 "github.com/cocoup/go-smart/tools/gocli/cmd/api/parser/g4/ast"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/parser/g4/gen/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fieldAccept = func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
	return p.Field().Accept(visitor)
}

func TestField(t *testing.T) {
	t.Run("anonymous", func(t *testing.T) {
		v, err := parser.Accept(fieldAccept, `User`)
		assert.Nil(t, err)
		f := v.(*ast2.TypeField)
		assert.True(t, f.Equal(&ast2.TypeField{
			IsAnonymous: true,
			DataType:    &ast2.Literal{Literal: ast2.NewTextExpr("User")},
		}))

		v, err = parser.Accept(fieldAccept, `*User`)
		assert.Nil(t, err)
		f = v.(*ast2.TypeField)
		assert.True(t, f.Equal(&ast2.TypeField{
			IsAnonymous: true,
			DataType: &ast2.Pointer{
				PointerExpr: ast2.NewTextExpr("*User"),
				Star:        ast2.NewTextExpr("*"),
				Name:        ast2.NewTextExpr("User"),
			},
		}))

		v, err = parser.Accept(fieldAccept, `
		// anonymous user
		*User // pointer type`)
		assert.Nil(t, err)
		f = v.(*ast2.TypeField)
		assert.True(t, f.Equal(&ast2.TypeField{
			IsAnonymous: true,
			DataType: &ast2.Pointer{
				PointerExpr: ast2.NewTextExpr("*User"),
				Star:        ast2.NewTextExpr("*"),
				Name:        ast2.NewTextExpr("User"),
			},
			DocExpr: []ast2.Expr{
				ast2.NewTextExpr("// anonymous user"),
			},
			CommentExpr: ast2.NewTextExpr("// pointer type"),
		}))

		_, err = parser.Accept(fieldAccept, `interface`)
		assert.Error(t, err)

		_, err = parser.Accept(fieldAccept, `map`)
		assert.Error(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fieldAccept, `User int`)
		assert.Nil(t, err)
		f := v.(*ast2.TypeField)
		assert.True(t, f.Equal(&ast2.TypeField{
			Name:     ast2.NewTextExpr("User"),
			DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int")},
		}))
		v, err = parser.Accept(fieldAccept, `Foo Bar`)
		assert.Nil(t, err)
		f = v.(*ast2.TypeField)
		assert.True(t, f.Equal(&ast2.TypeField{
			Name:     ast2.NewTextExpr("Foo"),
			DataType: &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
		}))

		v, err = parser.Accept(fieldAccept, `Foo map[int]Bar`)
		assert.Nil(t, err)
		f = v.(*ast2.TypeField)
		assert.True(t, f.Equal(&ast2.TypeField{
			Name: ast2.NewTextExpr("Foo"),
			DataType: &ast2.Map{
				MapExpr: ast2.NewTextExpr("map[int]Bar"),
				Map:     ast2.NewTextExpr("map"),
				LBrack:  ast2.NewTextExpr("["),
				RBrack:  ast2.NewTextExpr("]"),
				Key:     ast2.NewTextExpr("int"),
				Value:   &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
			},
		}))
	})
}

func TestDataType_ID(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.DataType().Accept(visitor)
	}
	t.Run("Struct", func(t *testing.T) {
		v, err := parser.Accept(dt, `Foo`)
		assert.Nil(t, err)
		id := v.(ast2.DataType)
		assert.True(t, id.Equal(&ast2.Literal{Literal: ast2.NewTextExpr("Foo")}))
	})

	t.Run("basic", func(t *testing.T) {
		v, err := parser.Accept(dt, `int`)
		assert.Nil(t, err)
		id := v.(ast2.DataType)
		assert.True(t, id.Equal(&ast2.Literal{Literal: ast2.NewTextExpr("int")}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `map`)
		assert.Error(t, err)
	})
}

func TestDataType_Map(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.MapType().Accept(visitor)
	}
	t.Run("basicKey", func(t *testing.T) {
		v, err := parser.Accept(dt, `map[int]Bar`)
		assert.Nil(t, err)
		m := v.(ast2.DataType)
		assert.True(t, m.Equal(&ast2.Map{
			MapExpr: ast2.NewTextExpr("map[int]Bar"),
			Map:     ast2.NewTextExpr("map"),
			LBrack:  ast2.NewTextExpr("["),
			RBrack:  ast2.NewTextExpr("]"),
			Key:     ast2.NewTextExpr("int"),
			Value:   &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `map[var]Bar`)
		assert.Error(t, err)

		_, err = parser.Accept(dt, `map[*User]Bar`)
		assert.Error(t, err)

		_, err = parser.Accept(dt, `map[User]Bar`)
		assert.Error(t, err)
	})
}

func TestDataType_Array(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.ArrayType().Accept(visitor)
	}
	t.Run("basic", func(t *testing.T) {
		v, err := parser.Accept(dt, `[]int`)
		assert.Nil(t, err)
		array := v.(ast2.DataType)
		assert.True(t, array.Equal(&ast2.Array{
			ArrayExpr: ast2.NewTextExpr("[]int"),
			LBrack:    ast2.NewTextExpr("["),
			RBrack:    ast2.NewTextExpr("]"),
			Literal:   &ast2.Literal{Literal: ast2.NewTextExpr("int")},
		}))
	})

	t.Run("pointer", func(t *testing.T) {
		v, err := parser.Accept(dt, `[]*User`)
		assert.Nil(t, err)
		array := v.(ast2.DataType)
		assert.True(t, array.Equal(&ast2.Array{
			ArrayExpr: ast2.NewTextExpr("[]*User"),
			LBrack:    ast2.NewTextExpr("["),
			RBrack:    ast2.NewTextExpr("]"),
			Literal: &ast2.Pointer{
				PointerExpr: ast2.NewTextExpr("*User"),
				Star:        ast2.NewTextExpr("*"),
				Name:        ast2.NewTextExpr("User"),
			},
		}))
	})

	t.Run("interface{}", func(t *testing.T) {
		v, err := parser.Accept(dt, `[]interface{}`)
		assert.Nil(t, err)
		array := v.(ast2.DataType)
		assert.True(t, array.Equal(&ast2.Array{
			ArrayExpr: ast2.NewTextExpr("[]interface{}"),
			LBrack:    ast2.NewTextExpr("["),
			RBrack:    ast2.NewTextExpr("]"),
			Literal:   &ast2.Interface{Literal: ast2.NewTextExpr("interface{}")},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `[]var`)
		assert.Error(t, err)

		_, err = parser.Accept(dt, `[]interface`)
		assert.Error(t, err)
	})
}

func TestDataType_Interface(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.DataType().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(dt, `interface{}`)
		assert.Nil(t, err)
		inter := v.(ast2.DataType)
		assert.True(t, inter.Equal(&ast2.Interface{Literal: ast2.NewTextExpr("interface{}")}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `interface`)
		assert.Error(t, err)
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `interface{`)
		assert.Error(t, err)
	})
}

func TestDataType_Time(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.DataType().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(dt, `time.Time`)
		assert.Error(t, err)
	})
}

func TestDataType_Pointer(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.PointerType().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(dt, `*int`)
		assert.Nil(t, err)
		assert.True(t, v.(ast2.DataType).Equal(&ast2.Pointer{
			PointerExpr: ast2.NewTextExpr("*int"),
			Star:        ast2.NewTextExpr("*"),
			Name:        ast2.NewTextExpr("int"),
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `int`)
		assert.Error(t, err)
	})
}

func TestAlias(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.TypeAlias().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(fn, `Foo int`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `Foo=int`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `
		Foo int // comment`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `
		Foo int /**comment*/`)
		assert.Error(t, err)
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `Foo var`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `Foo 2`)
		assert.Error(t, err)
	})
}

func TestTypeStruct(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.TypeStruct().Accept(visitor)
	}

	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, "Foo {\n\t\t\tFoo string\n\t\t\tBar int `json:\"bar\"``\n\t\t}")
		assert.Nil(t, err)
		s := v.(*ast2.TypeStruct)
		assert.True(t, s.Equal(&ast2.TypeStruct{
			Name:   ast2.NewTextExpr("Foo"),
			LBrace: ast2.NewTextExpr("{"),
			RBrace: ast2.NewTextExpr("}"),
			Fields: []*ast2.TypeField{
				{
					Name:     ast2.NewTextExpr("Foo"),
					DataType: &ast2.Literal{Literal: ast2.NewTextExpr("string")},
				},
				{
					Name:     ast2.NewTextExpr("Bar"),
					DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int")},
					Tag:      ast2.NewTextExpr("`json:\"bar\"`"),
				},
			},
		}))

		v, err = parser.Accept(fn, "Foo struct{\n\t\t\tFoo string\n\t\t\tBar int `json:\"bar\"``\n\t\t}")
		assert.Nil(t, err)
		s = v.(*ast2.TypeStruct)
		assert.True(t, s.Equal(&ast2.TypeStruct{
			Name:   ast2.NewTextExpr("Foo"),
			LBrace: ast2.NewTextExpr("{"),
			RBrace: ast2.NewTextExpr("}"),
			Struct: ast2.NewTextExpr("struct"),
			Fields: []*ast2.TypeField{
				{
					Name:     ast2.NewTextExpr("Foo"),
					DataType: &ast2.Literal{Literal: ast2.NewTextExpr("string")},
				},
				{
					Name:     ast2.NewTextExpr("Bar"),
					DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int")},
					Tag:      ast2.NewTextExpr("`json:\"bar\"`"),
				},
			},
		}))
	})
}

func TestTypeBlock(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.TypeBlock().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(fn, `type(
			// doc
			Foo int
		)`)
		assert.Error(t, err)

		v, err := parser.Accept(fn, `type (
			// doc
			Foo {
				Bar int
			}
		)`)
		assert.Nil(t, err)
		st := v.([]ast2.TypeExpr)
		assert.True(t, st[0].Equal(&ast2.TypeStruct{
			Name:   ast2.NewTextExpr("Foo"),
			LBrace: ast2.NewTextExpr("{"),
			RBrace: ast2.NewTextExpr("}"),
			DocExpr: []ast2.Expr{
				ast2.NewTextExpr("// doc"),
			},
			Fields: []*ast2.TypeField{
				{
					Name:     ast2.NewTextExpr("Bar"),
					DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int")},
				},
			},
		}))
	})
}

func TestTypeLit(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.TypeLit().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(fn, `type Foo int`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `type Foo = int`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `
		// doc
		type Foo = int // comment`)
		assert.Error(t, err)

		v, err := parser.Accept(fn, `
		// doc
		type Foo {// comment
			Bar int
		}`)
		assert.Nil(t, err)
		st := v.(*ast2.TypeStruct)
		assert.True(t, st.Equal(&ast2.TypeStruct{
			Name: ast2.NewTextExpr("Foo"),
			Fields: []*ast2.TypeField{
				{
					Name:     ast2.NewTextExpr("Bar"),
					DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int")},
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("// comment"),
					},
				},
			},
			DocExpr: []ast2.Expr{
				ast2.NewTextExpr("// doc"),
			},
		}))

		v, err = parser.Accept(fn, `
		// doc
		type Foo {// comment
			Bar
		}`)
		assert.Nil(t, err)
		st = v.(*ast2.TypeStruct)
		assert.True(t, st.Equal(&ast2.TypeStruct{
			Name: ast2.NewTextExpr("Foo"),
			Fields: []*ast2.TypeField{
				{
					IsAnonymous: true,
					DataType:    &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("// comment"),
					},
				},
			},
			DocExpr: []ast2.Expr{
				ast2.NewTextExpr("// doc"),
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `type Foo`)
		assert.Error(t, err)
	})
}

func TestTypeUnExported(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.TypeSpec().Accept(visitor)
	}

	t.Run("type", func(t *testing.T) {
		_, err := parser.Accept(fn, `type foo {}`)
		assert.Nil(t, err)
	})

	t.Run("field", func(t *testing.T) {
		_, err := parser.Accept(fn, `type Foo {
			name int
		}`)
		assert.Nil(t, err)

		_, err = parser.Accept(fn, `type Foo {
			Name int
		}`)
		assert.Nil(t, err)
	})

	t.Run("filedDataType", func(t *testing.T) {
		_, err := parser.Accept(fn, `type Foo {
			Foo *foo
			Bar []bar
			FooBar map[int]fooBar
		}`)
		assert.Nil(t, err)
	})
}
