package test

import (
	_ "embed"
	ast2 "github.com/cocoup/go-smart/tools/gocli/cmd/api/parser/g4/ast"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/parser/g4/gen/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed apis/test.api
var testApi string

var parser = ast2.NewParser(ast2.WithParserPrefix("test.api"), ast2.WithParserDebug())

func TestApi(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.Api().Accept(visitor)
	}

	v, err := parser.Accept(fn, testApi)
	assert.Nil(t, err)
	api := v.(*ast2.Api)
	body := &ast2.Body{
		Lp:   ast2.NewTextExpr("("),
		Rp:   ast2.NewTextExpr(")"),
		Name: &ast2.Literal{Literal: ast2.NewTextExpr("FooBar")},
	}

	returns := ast2.NewTextExpr("returns")
	assert.True(t, api.Equal(&ast2.Api{
		Syntax: &ast2.SyntaxExpr{
			Syntax:  ast2.NewTextExpr("syntax"),
			Assign:  ast2.NewTextExpr("="),
			Version: ast2.NewTextExpr(`"v1"`),
			DocExpr: []ast2.Expr{
				ast2.NewTextExpr("// syntax doc"),
			},
			CommentExpr: ast2.NewTextExpr("// syntax comment"),
		},
		Import: []*ast2.ImportExpr{
			{
				Import: ast2.NewTextExpr("import"),
				Value:  ast2.NewTextExpr(`"foo.api"`),
				DocExpr: []ast2.Expr{
					ast2.NewTextExpr("// import doc"),
				},
				CommentExpr: ast2.NewTextExpr("// import comment"),
			},
			{
				Import: ast2.NewTextExpr("import"),
				Value:  ast2.NewTextExpr(`"bar.api"`),
				DocExpr: []ast2.Expr{
					ast2.NewTextExpr("// import group doc"),
				},
				CommentExpr: ast2.NewTextExpr("// import group comment"),
			},
		},
		Info: &ast2.InfoExpr{
			Info: ast2.NewTextExpr("info"),
			Lp:   ast2.NewTextExpr("("),
			Rp:   ast2.NewTextExpr(")"),
			Kvs: []*ast2.KvExpr{
				{
					Key:   ast2.NewTextExpr("author"),
					Value: ast2.NewTextExpr(`"songmeizi"`),
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("// author doc"),
					},
					CommentExpr: ast2.NewTextExpr("// author comment"),
				},
				{
					Key:   ast2.NewTextExpr("date"),
					Value: ast2.NewTextExpr(`2020-01-04`),
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("// date doc"),
					},
					CommentExpr: ast2.NewTextExpr("// date comment"),
				},
				{
					Key: ast2.NewTextExpr("desc"),
					Value: ast2.NewTextExpr(`"break line
    desc"`),
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("// desc doc"),
					},
					CommentExpr: ast2.NewTextExpr("// desc comment"),
				},
			},
		},
		Type: []ast2.TypeExpr{
			&ast2.TypeStruct{
				Name:   ast2.NewTextExpr("FooBar"),
				Struct: ast2.NewTextExpr("struct"),
				LBrace: ast2.NewTextExpr("{"),
				RBrace: ast2.NewTextExpr("}"),
				Fields: []*ast2.TypeField{
					{
						Name:     ast2.NewTextExpr("Foo"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int")},
					},
				},
			},
			&ast2.TypeStruct{
				Name:   ast2.NewTextExpr("Bar"),
				LBrace: ast2.NewTextExpr("{"),
				RBrace: ast2.NewTextExpr("}"),
				DocExpr: []ast2.Expr{
					ast2.NewTextExpr("// remove struct"),
				},
				Fields: []*ast2.TypeField{
					{
						Name:     ast2.NewTextExpr("VString"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("string")},
						Tag:      ast2.NewTextExpr("`json:\"vString\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vString"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VBool"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("bool")},
						Tag:      ast2.NewTextExpr("`json:\"vBool\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vBool"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VInt8"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int8")},
						Tag:      ast2.NewTextExpr("`json:\"vInt8\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vInt8"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VInt16"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int16")},
						Tag:      ast2.NewTextExpr("`json:\"vInt16\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vInt16"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VInt32"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int32")},
						Tag:      ast2.NewTextExpr("`json:\"vInt32\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vInt32"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VInt64"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int64")},
						Tag:      ast2.NewTextExpr("`json:\"vInt64\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vInt64"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VInt"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("int")},
						Tag:      ast2.NewTextExpr("`json:\"vInt\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vInt"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VUInt8"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("uint8")},
						Tag:      ast2.NewTextExpr("`json:\"vUInt8\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vUInt8"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VUInt16"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("uint16")},
						Tag:      ast2.NewTextExpr("`json:\"vUInt16\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vUInt16"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VUInt32"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("uint32")},
						Tag:      ast2.NewTextExpr("`json:\"vUInt32\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vUInt32"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VUInt64"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("uint64")},
						Tag:      ast2.NewTextExpr("`json:\"vUInt64\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vUInt64"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VFloat32"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("float32")},
						Tag:      ast2.NewTextExpr("`json:\"vFloat32\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vFloat32"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VFloat64"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("float64")},
						Tag:      ast2.NewTextExpr("`json:\"vFloat64\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vFloat64"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VByte"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("byte")},
						Tag:      ast2.NewTextExpr("`json:\"vByte\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vByte"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VRune"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("rune")},
						Tag:      ast2.NewTextExpr("`json:\"vRune\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vRune"),
						},
					},
					{
						Name: ast2.NewTextExpr("VMap"),
						DataType: &ast2.Map{
							MapExpr: ast2.NewTextExpr("map[string]int"),
							Map:     ast2.NewTextExpr("map"),
							LBrack:  ast2.NewTextExpr("["),
							RBrack:  ast2.NewTextExpr("]"),
							Key:     ast2.NewTextExpr("string"),
							Value:   &ast2.Literal{Literal: ast2.NewTextExpr("int")},
						},
						Tag: ast2.NewTextExpr("`json:\"vMap\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vMap"),
						},
					},
					{
						Name: ast2.NewTextExpr("VArray"),
						DataType: &ast2.Array{
							ArrayExpr: ast2.NewTextExpr("[]int"),
							LBrack:    ast2.NewTextExpr("["),
							RBrack:    ast2.NewTextExpr("]"),
							Literal:   &ast2.Literal{Literal: ast2.NewTextExpr("int")},
						},
						Tag: ast2.NewTextExpr("`json:\"vArray\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vArray"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VStruct"),
						DataType: &ast2.Literal{Literal: ast2.NewTextExpr("FooBar")},
						Tag:      ast2.NewTextExpr("`json:\"vStruct\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vStruct"),
						},
					},
					{
						Name: ast2.NewTextExpr("VStructPointer"),
						DataType: &ast2.Pointer{
							PointerExpr: ast2.NewTextExpr("*FooBar"),
							Star:        ast2.NewTextExpr("*"),
							Name:        ast2.NewTextExpr("FooBar"),
						},
						Tag: ast2.NewTextExpr("`json:\"vStructPointer\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vStructPointer"),
						},
					},
					{
						Name:     ast2.NewTextExpr("VInterface"),
						DataType: &ast2.Interface{Literal: ast2.NewTextExpr("interface{}")},
						Tag:      ast2.NewTextExpr("`json:\"vInterface\"`"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// vInterface"),
						},
					},
					{
						IsAnonymous: true,
						DataType:    &ast2.Literal{Literal: ast2.NewTextExpr("FooBar")},
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// inline"),
						},
					},
				},
			},
		},
		Service: []*ast2.Service{
			{
				AtServer: &ast2.AtServer{
					AtServerToken: ast2.NewTextExpr("@server"),
					Lp:            ast2.NewTextExpr("("),
					Rp:            ast2.NewTextExpr(")"),
					Kv: []*ast2.KvExpr{
						{
							Key:   ast2.NewTextExpr("host"),
							Value: ast2.NewTextExpr("0.0.0.0"),
						},
						{
							Key:   ast2.NewTextExpr("port"),
							Value: ast2.NewTextExpr("8080"),
						},
						{
							Key: ast2.NewTextExpr("annotation"),
							Value: ast2.NewTextExpr(`"break line
    desc"`),
						},
					},
				},
				ServiceApi: &ast2.ServiceApi{
					ServiceToken: ast2.NewTextExpr("service"),
					Name:         ast2.NewTextExpr("foo-api"),
					Lbrace:       ast2.NewTextExpr("{"),
					Rbrace:       ast2.NewTextExpr("}"),
					ServiceRoute: []*ast2.ServiceRoute{
						{
							AtDoc: &ast2.AtDoc{
								AtDocToken: ast2.NewTextExpr("@doc"),
								Lp:         ast2.NewTextExpr("("),
								Rp:         ast2.NewTextExpr(")"),
								LineDoc:    ast2.NewTextExpr(`"foo"`),
							},
							AtHandler: &ast2.AtHandler{
								AtHandlerToken: ast2.NewTextExpr("@handler"),
								Name:           ast2.NewTextExpr("postFoo"),
							},
							Route: &ast2.Route{
								Method:      ast2.NewTextExpr("post"),
								Path:        ast2.NewTextExpr("/foo"),
								Req:         body,
								ReturnToken: returns,
								Reply:       body,
								DocExpr: []ast2.Expr{
									ast2.NewTextExpr("// foo"),
								},
							},
						},
						{
							AtDoc: &ast2.AtDoc{
								AtDocToken: ast2.NewTextExpr("@doc"),
								Lp:         ast2.NewTextExpr("("),
								Rp:         ast2.NewTextExpr(")"),
								Kv: []*ast2.KvExpr{
									{
										Key:   ast2.NewTextExpr("summary"),
										Value: ast2.NewTextExpr("bar"),
									},
								},
							},
							AtServer: &ast2.AtServer{
								AtServerToken: ast2.NewTextExpr("@server"),
								Lp:            ast2.NewTextExpr("("),
								Rp:            ast2.NewTextExpr(")"),
								Kv: []*ast2.KvExpr{
									{
										Key:   ast2.NewTextExpr("handler"),
										Value: ast2.NewTextExpr("postBar"),
									},
								},
							},
							Route: &ast2.Route{
								Method: ast2.NewTextExpr("post"),
								Path:   ast2.NewTextExpr("/bar"),
								Req:    body,
							},
						},
						{
							AtDoc: &ast2.AtDoc{
								AtDocToken: ast2.NewTextExpr("@doc"),
								Lp:         ast2.NewTextExpr("("),
								Rp:         ast2.NewTextExpr(")"),
								LineDoc:    ast2.NewTextExpr(`"foobar"`),
							},
							AtHandler: &ast2.AtHandler{
								AtHandlerToken: ast2.NewTextExpr("@handler"),
								Name:           ast2.NewTextExpr("postFooBar"),
							},
							Route: &ast2.Route{
								Method:      ast2.NewTextExpr("post"),
								Path:        ast2.NewTextExpr("/foo/bar"),
								ReturnToken: returns,
								Reply:       body,
								DocExpr: []ast2.Expr{
									ast2.NewTextExpr(`/**
    * httpmethod: post
    * path: /foo/bar
    * reply: FooBar
    */`),
								},
							},
						},
						{
							AtDoc: &ast2.AtDoc{
								AtDocToken: ast2.NewTextExpr("@doc"),
								Lp:         ast2.NewTextExpr("("),
								Rp:         ast2.NewTextExpr(")"),
								LineDoc:    ast2.NewTextExpr(`"barfoo"`),
							},
							AtHandler: &ast2.AtHandler{
								AtHandlerToken: ast2.NewTextExpr("@handler"),
								Name:           ast2.NewTextExpr("postBarFoo"),
							},
							Route: &ast2.Route{
								Method:      ast2.NewTextExpr("post"),
								Path:        ast2.NewTextExpr("/bar/foo"),
								ReturnToken: returns,
								Reply:       body,
								CommentExpr: ast2.NewTextExpr("// post:/bar/foo"),
							},
						},
						{
							AtDoc: &ast2.AtDoc{
								AtDocToken: ast2.NewTextExpr("@doc"),
								Lp:         ast2.NewTextExpr("("),
								Rp:         ast2.NewTextExpr(")"),
								LineDoc:    ast2.NewTextExpr(`"barfoo"`),
							},
							AtHandler: &ast2.AtHandler{
								AtHandlerToken: ast2.NewTextExpr("@handler"),
								Name:           ast2.NewTextExpr("getBarFoo"),
							},
							Route: &ast2.Route{
								Method:      ast2.NewTextExpr("get"),
								Path:        ast2.NewTextExpr("/bar/foo"),
								ReturnToken: returns,
								Reply:       body,
							},
						},
					},
				},
			},
		},
	}))
}

func TestApiSyntax(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast2.ApiVisitor) interface{} {
		return p.Api().Accept(visitor)
	}
	parser.Accept(fn, `
	// doc 1
	// doc 2
	syntax = "v1" // comment 1
	// comment 2
	import "foo.api"
	`)
}
