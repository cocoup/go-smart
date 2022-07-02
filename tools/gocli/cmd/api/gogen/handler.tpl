package {{.package}}

import (
	"github.com/gin-gonic/gin"

	{{.imports}}
)

func {{.handler}}(svcCtx *service.Context) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        {{if .hasRequest}}var req types.{{.requestType}}
        if err := ctx.ShouldBindJSON(&req); err != nil {
            result.ParamError(ctx, err)
            return
        }

        {{end}}service := {{.pkg}}.New{{.entity}}Service(ctx, svcCtx)
        {{if .hasResp}}resp, {{end}}err := service.{{.call}}({{if .hasRequest}}&req{{end}})
        result.HttpResult(ctx, {{if .hasResp}}resp{{else}}nil{{end}}, err)
    }
}