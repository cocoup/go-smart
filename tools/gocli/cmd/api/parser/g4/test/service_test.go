package test

import (
	ast2 "github.com/cocoup/go-smart/tools/gocli/cmd/api/parser/g4/ast"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/parser/g4/gen/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBody(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast2.ApiVisitor) interface{} {
		return p.Body().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `(Foo)`)
		assert.Nil(t, err)
		body := v.(*ast2.Body)
		assert.True(t, body.Equal(&ast2.Body{
			Lp:   ast2.NewTextExpr("("),
			Rp:   ast2.NewTextExpr(")"),
			Name: &ast2.Literal{Literal: ast2.NewTextExpr("Foo")},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `(var)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `()`)
		assert.Nil(t, err)
	})
}

func TestRoute(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast2.ApiVisitor) interface{} {
		return p.Route().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `post /foo/foo-bar/:bar (Foo) returns (Bar)`)
		assert.Nil(t, err)
		route := v.(*ast2.Route)
		assert.True(t, route.Equal(&ast2.Route{
			Method: ast2.NewTextExpr("post"),
			Path:   ast2.NewTextExpr("/foo/foo-bar/:bar"),
			Req: &ast2.Body{
				Lp:   ast2.NewTextExpr("("),
				Rp:   ast2.NewTextExpr(")"),
				Name: &ast2.Literal{Literal: ast2.NewTextExpr("Foo")},
			},
			ReturnToken: ast2.NewTextExpr("returns"),
			Reply: &ast2.Body{
				Lp:   ast2.NewTextExpr("("),
				Rp:   ast2.NewTextExpr(")"),
				Name: &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar (Foo)`)
		assert.Nil(t, err)
		route = v.(*ast2.Route)
		assert.True(t, route.Equal(&ast2.Route{
			Method: ast2.NewTextExpr("post"),
			Path:   ast2.NewTextExpr("/foo/foo-bar/:bar"),
			Req: &ast2.Body{
				Lp:   ast2.NewTextExpr("("),
				Rp:   ast2.NewTextExpr(")"),
				Name: &ast2.Literal{Literal: ast2.NewTextExpr("Foo")},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar returns (Bar)`)
		assert.Nil(t, err)
		route = v.(*ast2.Route)
		assert.True(t, route.Equal(&ast2.Route{
			Method:      ast2.NewTextExpr("post"),
			Path:        ast2.NewTextExpr("/foo/foo-bar/:bar"),
			ReturnToken: ast2.NewTextExpr("returns"),
			Reply: &ast2.Body{
				Lp:   ast2.NewTextExpr("("),
				Rp:   ast2.NewTextExpr(")"),
				Name: &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar returns ([]Bar)`)
		assert.Nil(t, err)
		route = v.(*ast2.Route)
		assert.True(t, route.Equal(&ast2.Route{
			Method:      ast2.NewTextExpr("post"),
			Path:        ast2.NewTextExpr("/foo/foo-bar/:bar"),
			ReturnToken: ast2.NewTextExpr("returns"),
			Reply: &ast2.Body{
				Lp: ast2.NewTextExpr("("),
				Rp: ast2.NewTextExpr(")"),
				Name: &ast2.Array{
					ArrayExpr: ast2.NewTextExpr("[]Bar"),
					LBrack:    ast2.NewTextExpr("["),
					RBrack:    ast2.NewTextExpr("]"),
					Literal:   &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
				},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar returns ([]*Bar)`)
		assert.Nil(t, err)
		route = v.(*ast2.Route)
		assert.True(t, route.Equal(&ast2.Route{
			Method:      ast2.NewTextExpr("post"),
			Path:        ast2.NewTextExpr("/foo/foo-bar/:bar"),
			ReturnToken: ast2.NewTextExpr("returns"),
			Reply: &ast2.Body{
				Lp: ast2.NewTextExpr("("),
				Rp: ast2.NewTextExpr(")"),
				Name: &ast2.Array{
					ArrayExpr: ast2.NewTextExpr("[]*Bar"),
					LBrack:    ast2.NewTextExpr("["),
					RBrack:    ast2.NewTextExpr("]"),
					Literal: &ast2.Pointer{
						PointerExpr: ast2.NewTextExpr("*Bar"),
						Star:        ast2.NewTextExpr("*"),
						Name:        ast2.NewTextExpr("Bar"),
					},
				},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar`)
		assert.Nil(t, err)
		route = v.(*ast2.Route)
		assert.True(t, route.Equal(&ast2.Route{
			Method: ast2.NewTextExpr("post"),
			Path:   ast2.NewTextExpr("/foo/foo-bar/:bar"),
		}))

		v, err = parser.Accept(fn, `
		// foo
		post /foo/foo-bar/:bar // bar`)
		assert.Nil(t, err)
		route = v.(*ast2.Route)
		assert.True(t, route.Equal(&ast2.Route{
			Method: ast2.NewTextExpr("post"),
			Path:   ast2.NewTextExpr("/foo/foo-bar/:bar"),
			DocExpr: []ast2.Expr{
				ast2.NewTextExpr("// foo"),
			},
			CommentExpr: ast2.NewTextExpr("// bar"),
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `posts /foo`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `gets /foo`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `post /foo/:`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `post /foo/`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `post foo/bar`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `post /foo/bar returns (Bar)`)
		assert.Nil(t, err)

		_, err = parser.Accept(fn, ` /foo/bar returns (Bar)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post   returns (Bar)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post /foo/bar returns (int)`)
		assert.Nil(t, err)

		_, err = parser.Accept(fn, ` post /foo/bar returns (*int)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post /foo/bar returns ([]var)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post /foo/bar returns (const)`)
		assert.Error(t, err)
	})
}

func TestAtHandler(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast2.ApiVisitor) interface{} {
		return p.AtHandler().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `@handler foo`)
		assert.Nil(t, err)
		atHandler := v.(*ast2.AtHandler)
		assert.True(t, atHandler.Equal(&ast2.AtHandler{
			AtHandlerToken: ast2.NewTextExpr("@handler"),
			Name:           ast2.NewTextExpr("foo"),
		}))

		v, err = parser.Accept(fn, `
		// foo
		@handler foo // bar`)
		assert.Nil(t, err)
		atHandler = v.(*ast2.AtHandler)
		assert.True(t, atHandler.Equal(&ast2.AtHandler{
			AtHandlerToken: ast2.NewTextExpr("@handler"),
			Name:           ast2.NewTextExpr("foo"),
			DocExpr: []ast2.Expr{
				ast2.NewTextExpr("// foo"),
			},
			CommentExpr: ast2.NewTextExpr("// bar"),
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, ``)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@handler`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@handler "foo"`)
		assert.Error(t, err)
	})
}

func TestAtDoc(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast2.ApiVisitor) interface{} {
		return p.AtDoc().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `@doc "foo"`)
		assert.Nil(t, err)
		atDoc := v.(*ast2.AtDoc)
		assert.True(t, atDoc.Equal(&ast2.AtDoc{
			AtDocToken: ast2.NewTextExpr("@doc"),
			LineDoc:    ast2.NewTextExpr(`"foo"`),
		}))

		v, err = parser.Accept(fn, `@doc("foo")`)
		assert.Nil(t, err)
		atDoc = v.(*ast2.AtDoc)
		assert.True(t, atDoc.Equal(&ast2.AtDoc{
			AtDocToken: ast2.NewTextExpr("@doc"),
			Lp:         ast2.NewTextExpr("("),
			Rp:         ast2.NewTextExpr(")"),
			LineDoc:    ast2.NewTextExpr(`"foo"`),
		}))

		v, err = parser.Accept(fn, `@doc(
			foo: bar
		)`)
		assert.Nil(t, err)
		atDoc = v.(*ast2.AtDoc)
		assert.True(t, atDoc.Equal(&ast2.AtDoc{
			AtDocToken: ast2.NewTextExpr("@doc"),
			Lp:         ast2.NewTextExpr("("),
			Rp:         ast2.NewTextExpr(")"),
			Kv: []*ast2.KvExpr{
				{
					Key:   ast2.NewTextExpr("foo"),
					Value: ast2.NewTextExpr("bar"),
				},
			},
		}))

		v, err = parser.Accept(fn, `@doc(
			// foo
			foo: bar // bar
		)`)
		assert.Nil(t, err)
		atDoc = v.(*ast2.AtDoc)
		assert.True(t, atDoc.Equal(&ast2.AtDoc{
			AtDocToken: ast2.NewTextExpr("@doc"),
			Lp:         ast2.NewTextExpr("("),
			Rp:         ast2.NewTextExpr(")"),
			Kv: []*ast2.KvExpr{
				{
					Key:   ast2.NewTextExpr("foo"),
					Value: ast2.NewTextExpr("bar"),
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("// foo"),
					},
					CommentExpr: ast2.NewTextExpr("// bar"),
				},
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `@doc("foo"`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@doc "foo")`)
		assert.Error(t, err)
	})
}

func TestServiceRoute(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast2.ApiVisitor) interface{} {
		return p.ServiceRoute().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `
		@doc("foo")
		// foo/bar
		// foo
		@handler foo // bar
		// foo/bar
		// foo
		post /foo (Foo) returns (Bar) // bar
		`)
		assert.Nil(t, err)
		sr := v.(*ast2.ServiceRoute)
		assert.True(t, sr.Equal(&ast2.ServiceRoute{
			AtDoc: &ast2.AtDoc{
				AtDocToken: ast2.NewTextExpr("@doc"),
				Lp:         ast2.NewTextExpr("("),
				Rp:         ast2.NewTextExpr(")"),
				LineDoc:    ast2.NewTextExpr(`"foo"`),
			},
			AtHandler: &ast2.AtHandler{
				AtHandlerToken: ast2.NewTextExpr("@handler"),
				Name:           ast2.NewTextExpr("foo"),
				DocExpr: []ast2.Expr{
					ast2.NewTextExpr("// foo"),
				},
				CommentExpr: ast2.NewTextExpr("// bar"),
			},
			Route: &ast2.Route{
				Method: ast2.NewTextExpr("post"),
				Path:   ast2.NewTextExpr("/foo"),
				Req: &ast2.Body{
					Lp:   ast2.NewTextExpr("("),
					Rp:   ast2.NewTextExpr(")"),
					Name: &ast2.Literal{Literal: ast2.NewTextExpr("Foo")},
				},
				ReturnToken: ast2.NewTextExpr("returns"),
				Reply: &ast2.Body{
					Lp:   ast2.NewTextExpr("("),
					Rp:   ast2.NewTextExpr(")"),
					Name: &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
				},
				DocExpr: []ast2.Expr{
					ast2.NewTextExpr("// foo"),
				},
				CommentExpr: ast2.NewTextExpr("// bar"),
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `post /foo (Foo) returns (Bar) // bar`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@handler foo`)
		assert.Error(t, err)
	})
}

func TestServiceApi(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast2.ApiVisitor) interface{} {
		return p.ServiceApi().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `
		service foo-api{
			@doc("foo")
			// foo/bar
			// foo
			@handler foo // bar
			// foo/bar
			// foo
			post /foo (Foo) returns (Bar) // bar
		}
		`)
		assert.Nil(t, err)
		api := v.(*ast2.ServiceApi)
		assert.True(t, api.Equal(&ast2.ServiceApi{
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
						Name:           ast2.NewTextExpr("foo"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// foo"),
						},
						CommentExpr: ast2.NewTextExpr("// bar"),
					},
					Route: &ast2.Route{
						Method: ast2.NewTextExpr("post"),
						Path:   ast2.NewTextExpr("/foo"),
						Req: &ast2.Body{
							Lp:   ast2.NewTextExpr("("),
							Rp:   ast2.NewTextExpr(")"),
							Name: &ast2.Literal{Literal: ast2.NewTextExpr("Foo")},
						},
						ReturnToken: ast2.NewTextExpr("returns"),
						Reply: &ast2.Body{
							Lp:   ast2.NewTextExpr("("),
							Rp:   ast2.NewTextExpr(")"),
							Name: &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
						},
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// foo"),
						},
						CommentExpr: ast2.NewTextExpr("// bar"),
					},
				},
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `services foo-api{}`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `service foo-api{`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `service foo-api{
		post /foo
		}`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `service foo-api{
		@handler foo
		}`)
		assert.Error(t, err)
	})
}

func TestAtServer(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast2.ApiVisitor) interface{} {
		return p.AtServer().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `
		@server(
			// foo
			foo1: bar1 // bar
			// foo
			foo2: "bar2" // bar
			/**foo*/
			foo3: "foo
			bar" /**bar*/		
		)
		`)
		assert.Nil(t, err)
		as := v.(*ast2.AtServer)
		assert.True(t, as.Equal(&ast2.AtServer{
			AtServerToken: ast2.NewTextExpr("@server"),
			Lp:            ast2.NewTextExpr("("),
			Rp:            ast2.NewTextExpr(")"),
			Kv: []*ast2.KvExpr{
				{
					Key:   ast2.NewTextExpr("foo1"),
					Value: ast2.NewTextExpr("bar1"),
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("// foo"),
					},
					CommentExpr: ast2.NewTextExpr("// bar"),
				},
				{
					Key:   ast2.NewTextExpr("foo2"),
					Value: ast2.NewTextExpr(`"bar2"`),
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("// foo"),
					},
					CommentExpr: ast2.NewTextExpr("// bar"),
				},
				{
					Key: ast2.NewTextExpr("foo3"),
					Value: ast2.NewTextExpr(`"foo
			bar"`),
					DocExpr: []ast2.Expr{
						ast2.NewTextExpr("/**foo*/"),
					},
					CommentExpr: ast2.NewTextExpr("/**bar*/"),
				},
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `server (
			foo:bar
		)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@server ()`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@server (
			foo: bar
		`)
		assert.Error(t, err)
	})
}

func TestServiceSpec(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast2.ApiVisitor) interface{} {
		return p.ServiceSpec().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(fn, `
		service foo-api{
			@handler foo
			post /foo returns ([]int)
		}
		`)
		assert.Nil(t, err)

		v, err := parser.Accept(fn, `
		@server(
			// foo
			foo1: bar1 // bar
			// foo
			foo2: "bar2" // bar
			/**foo*/
			foo3: "foo
			bar" /**bar*/		
		)
		service foo-api{
			@doc("foo")
			// foo/bar
			// foo
			@handler foo // bar
			// foo/bar
			// foo
			post /foo (Foo) returns (Bar) // bar
		}
		`)
		assert.Nil(t, err)
		service := v.(*ast2.Service)
		assert.True(t, service.Equal(&ast2.Service{
			AtServer: &ast2.AtServer{
				AtServerToken: ast2.NewTextExpr("@server"),
				Lp:            ast2.NewTextExpr("("),
				Rp:            ast2.NewTextExpr(")"),
				Kv: []*ast2.KvExpr{
					{
						Key:   ast2.NewTextExpr("foo1"),
						Value: ast2.NewTextExpr("bar1"),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// foo"),
						},
						CommentExpr: ast2.NewTextExpr("// bar"),
					},
					{
						Key:   ast2.NewTextExpr("foo2"),
						Value: ast2.NewTextExpr(`"bar2"`),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("// foo"),
						},
						CommentExpr: ast2.NewTextExpr("// bar"),
					},
					{
						Key: ast2.NewTextExpr("foo3"),
						Value: ast2.NewTextExpr(`"foo
			bar"`),
						DocExpr: []ast2.Expr{
							ast2.NewTextExpr("/**foo*/"),
						},
						CommentExpr: ast2.NewTextExpr("/**bar*/"),
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
							Name:           ast2.NewTextExpr("foo"),
							DocExpr: []ast2.Expr{
								ast2.NewTextExpr("// foo"),
							},
							CommentExpr: ast2.NewTextExpr("// bar"),
						},
						Route: &ast2.Route{
							Method: ast2.NewTextExpr("post"),
							Path:   ast2.NewTextExpr("/foo"),
							Req: &ast2.Body{
								Lp:   ast2.NewTextExpr("("),
								Rp:   ast2.NewTextExpr(")"),
								Name: &ast2.Literal{Literal: ast2.NewTextExpr("Foo")},
							},
							ReturnToken: ast2.NewTextExpr("returns"),
							Reply: &ast2.Body{
								Lp:   ast2.NewTextExpr("("),
								Rp:   ast2.NewTextExpr(")"),
								Name: &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
							},
							DocExpr: []ast2.Expr{
								ast2.NewTextExpr("// foo"),
							},
							CommentExpr: ast2.NewTextExpr("// bar"),
						},
					},
				},
			},
		}))

		v, err = parser.Accept(fn, `
		service foo-api{
			@doc("foo")
			// foo/bar
			// foo
			@handler foo // bar
			// foo/bar
			// foo
			post /foo (Foo) returns (Bar) // bar
		}
		`)
		assert.Nil(t, err)
		service = v.(*ast2.Service)
		assert.True(t, service.Equal(&ast2.Service{
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
							Name:           ast2.NewTextExpr("foo"),
							DocExpr: []ast2.Expr{
								ast2.NewTextExpr("// foo"),
							},
							CommentExpr: ast2.NewTextExpr("// bar"),
						},
						Route: &ast2.Route{
							Method: ast2.NewTextExpr("post"),
							Path:   ast2.NewTextExpr("/foo"),
							Req: &ast2.Body{
								Lp:   ast2.NewTextExpr("("),
								Rp:   ast2.NewTextExpr(")"),
								Name: &ast2.Literal{Literal: ast2.NewTextExpr("Foo")},
							},
							ReturnToken: ast2.NewTextExpr("returns"),
							Reply: &ast2.Body{
								Lp:   ast2.NewTextExpr("("),
								Rp:   ast2.NewTextExpr(")"),
								Name: &ast2.Literal{Literal: ast2.NewTextExpr("Bar")},
							},
							DocExpr: []ast2.Expr{
								ast2.NewTextExpr("// foo"),
							},
							CommentExpr: ast2.NewTextExpr("// bar"),
						},
					},
				},
			},
		}))
	})
}
