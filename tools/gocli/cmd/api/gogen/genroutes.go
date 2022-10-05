package gogen

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/zeromicro/go-zero/core/collection"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
	"github.com/cocoup/go-smart/tools/gocli/utils/format"
	"github.com/cocoup/go-smart/tools/gocli/utils/pathx"
	"github.com/cocoup/go-smart/tools/gocli/vars"
)

var (
	//go:embed route.tpl
	routeTemplate string
)

const (
	jwtTransKey      = "jwtTransition"
	routesFilename   = "routes"
	timeoutThreshold = time.Millisecond
)

var mapping = map[string]string{
	"delete":  http.MethodDelete,
	"get":     http.MethodGet,
	"head":    http.MethodHead,
	"post":    http.MethodPost,
	"put":     http.MethodPut,
	"patch":   http.MethodPatch,
	"connect": http.MethodConnect,
	"options": http.MethodOptions,
	"trace":   http.MethodTrace,
}

type (
	group struct {
		routes           []route
		jwtEnabled       bool
		signatureEnabled bool
		authName         string
		timeout          string
		middlewares      []string
		prefix           string
		jwtTrans         string
		group            string
	}
	route struct {
		method  string
		path    string
		handler string
	}
)

func genRoutes(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	groups, err := getRoutes(api)
	if err != nil {
		return err
	}

	var (
		builder     strings.Builder
		hasTimeout  bool
		bMiddleware bool
		jwtEnabled  bool
	)
	for _, g := range groups {
		if len(g.group) <= 0 {
			fmt.Println(aurora.Red("service missing 'group' key, ignored generation"))
			continue
		}

		group := g.group
		if len(g.prefix) > 0 {
			group = fmt.Sprintf("%s/%s", g.prefix, group)
		}
		builder.WriteString("{\n")

		builder.WriteString(fmt.Sprintf("\tgroup := rootGroup.Group(\"%s\")\n", group))
		if g.jwtEnabled {
			jwtEnabled = g.jwtEnabled
			builder.WriteString(fmt.Sprintf("\tgroup.Use(restMid.JWTAuth(server.Conf.JWT.Secret))\n"))
		}
		if len(g.middlewares) > 0 {
			for _, m := range g.middlewares {
				builder.WriteString(fmt.Sprintf("\tgroup.Use(middleware.%s())\n", strings.Title(strings.TrimSpace(m))))
			}
			bMiddleware = true
		}
		builder.WriteString("\n")
		for _, route := range g.routes {
			handler := strings.TrimSuffix(route.handler, "Handler")
			path := strings.TrimPrefix(route.path, "/"+g.group)
			path = strings.TrimLeft(strings.TrimRight(path, "/"), "/")
			builder.WriteString(fmt.Sprintf("\tgroup.%s(\"%s\", %s.%s(svcCtx))\n",
				route.method, path, g.group, handler))
		}

		builder.WriteString("}\n")
	}

	routeFilename, err := format.NamingFormat(cfg.NamingFormat, routesFilename)
	if err != nil {
		return err
	}

	routeFilename = routeFilename + ".go"
	filename := path.Join(dir, routeDir, routeFilename)
	err = os.Remove(filename)

	return utils.GenFile(utils.FileGenConfig{
		Dir:             dir,
		Subdir:          routeDir,
		Filename:        routeFilename,
		TemplateName:    "routesTemplate",
		Category:        category,
		TemplateFile:    routeTemplateFile,
		BuiltinTemplate: routeTemplate,
		Data: map[string]interface{}{
			"hasTimeout":      hasTimeout,
			"imports":         genRoutesImports(rootPkg, groups, jwtEnabled, bMiddleware),
			"routesAdditions": strings.TrimSpace(builder.String()),
		},
	})
}

func genRoutesImports(parentPkg string, groups []group, jwtEnabled, bMiddleware bool) string {
	importSet := collection.NewSet()

	for _, group := range groups {
		if len(group.group) <= 0 {
			continue
		}
		folders := strings.Split(group.group, "/")
		folder := folders[0]
		importSet.AddStr(fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, handlerDir, folder)))
	}
	importSet.AddStr(fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, serviceDir)))
	if jwtEnabled {
		importSet.AddStr(fmt.Sprintf("restMid \"%s/rest/middleware\"", vars.ProjectOpenSourceURL))
	}
	if bMiddleware {
		//importSet.AddStr(fmt.Sprintf("\"%s/rest/middleware\"", vars.ProjectOpenSourceURL))
		importSet.AddStr(fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, middlewareDir)))
	}
	imports := importSet.KeysStr()
	sort.Strings(imports)
	projectSection := strings.Join(imports, "\n\t")
	depSection := fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL)
	return fmt.Sprintf("%s\n\n\t%s", projectSection, depSection)
}

func getRoutes(api *spec.ApiSpec) ([]group, error) {
	var routes []group

	for _, g := range api.Service.Groups {
		var groupedRoutes group
		for _, r := range g.Routes {
			handler := getHandlerName(r)
			//handler = handler + "(serverCtx)"
			folder := r.GetAnnotation(groupProperty)
			//if len(folder) > 0 {
			//	handler = toPrefix(folder) + "." + strings.ToUpper(handler[:1]) + handler[1:]
			//} else {
			//	folder = g.GetAnnotation(groupProperty)
			//	if len(folder) > 0 {
			//		handler = toPrefix(folder) + "." + strings.ToUpper(handler[:1]) + handler[1:]
			//	}
			//}

			if len(folder) <= 0 {
				folder = g.GetAnnotation(groupProperty)
			}
			groupedRoutes.routes = append(groupedRoutes.routes, route{
				method:  mapping[r.Method],
				path:    r.Path,
				handler: strings.Title(handler),
			})

			groupedRoutes.group = folder
		}
		groupedRoutes.timeout = g.GetAnnotation("timeout")

		jwt := g.GetAnnotation("jwt")
		if len(jwt) > 0 {
			groupedRoutes.authName = jwt
			groupedRoutes.jwtEnabled = true
		}
		jwtTrans := g.GetAnnotation(jwtTransKey)
		if len(jwtTrans) > 0 {
			groupedRoutes.jwtTrans = jwtTrans
		}

		signature := g.GetAnnotation("signature")
		if signature == "true" {
			groupedRoutes.signatureEnabled = true
		}
		middleware := g.GetAnnotation("middleware")
		if len(middleware) > 0 {
			groupedRoutes.middlewares = append(groupedRoutes.middlewares,
				strings.Split(middleware, ",")...)
		}
		prefix := g.GetAnnotation(spec.RoutePrefixKey)
		prefix = strings.ReplaceAll(prefix, `"`, "")
		prefix = strings.TrimSpace(prefix)
		if len(prefix) > 0 {
			prefix = path.Join("/", prefix)
			groupedRoutes.prefix = prefix
		}
		routes = append(routes, groupedRoutes)
	}

	return routes, nil
}

func toPrefix(folder string) string {
	return strings.ReplaceAll(folder, "/", "")
}

func getHandlerName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Handler"
}

func getHandlerBaseName(route spec.Route) (string, error) {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")
	return handler, nil
}
