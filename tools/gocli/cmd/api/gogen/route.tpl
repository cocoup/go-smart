// Code generated by goctl. DO NOT EDIT.
package route

import ({{if .hasTimeout}}
    "time"{{end}}

    {{.imports}}
)

func RegisterHandlers(server *rest.Server, svcCtx *service.Context) {
    rootGroup := server.Engine.Group(server.Conf.RouteRoot)

    {{.routesAdditions}}
}