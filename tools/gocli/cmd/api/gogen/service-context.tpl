package service

import (
    {{.imports}}
)

type Context struct {
	Config {{.config}}
}

func NewContext(conf {{.config}}) *Context {
	return &Context{Config: conf}
}