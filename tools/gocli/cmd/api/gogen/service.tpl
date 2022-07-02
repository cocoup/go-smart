package {{.package}}

import (
	"github.com/gin-gonic/gin"

	{{.imports}}
)

type {{.service}} struct {
	ctx    *gin.Context
    svcCtx *service.Context
}

func New{{.service}}(ctx *gin.Context, svcCtx *service.Context) *{{.service}} {
	return &{{.service}}{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
