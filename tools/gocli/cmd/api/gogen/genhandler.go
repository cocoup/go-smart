package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"github.com/logrusorgru/aurora"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
	"github.com/cocoup/go-smart/tools/gocli/utils/format"
	"github.com/cocoup/go-smart/tools/gocli/utils/pathx"
	"github.com/cocoup/go-smart/tools/gocli/vars"
)

var (
	//go:embed handler.tpl
	handlerTemplate string
)

func genHandlers(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, group := range api.Service.Groups {
		var hasRequest bool

		g := group.GetAnnotation(groupProperty)

		for _, route := range group.Routes {
			hasRequest = len(route.RequestTypeName()) > 0

			handler := strings.TrimSuffix(getHandlerName(route), "Handler")
			subdir := group.GetAnnotation(groupProperty)
			filename, err := format.NamingFormat(cfg.NamingFormat, handler)
			if nil != err {
				fmt.Println(aurora.Red(fmt.Sprintf("format file name error, %s\n", err.Error())))
			}

			err = utils.GenFile(utils.FileGenConfig{
				Dir:             path.Join(dir, handlerDir),
				Subdir:          subdir,
				Filename:        filename + ".go",
				TemplateName:    "handlersTemplate",
				Category:        category,
				TemplateFile:    handlerTemplateFile,
				BuiltinTemplate: handlerTemplate,
				Data: map[string]interface{}{
					"imports":     genHandlerImports(rootPkg, g, hasRequest),
					"package":     g,
					"handler":     strings.Title(handler),
					"hasRequest":  hasRequest,
					"hasResp":     len(route.ResponseTypeName()) > 0,
					"requestType": strings.Title(route.RequestTypeName()),
					"pkg":         g,
					"entity":      strings.Title(g),
					"call":        strings.Title(handler),
				},
			})
			if nil != err {
				fmt.Println(aurora.Red(fmt.Sprintf("generate handler error, %s\n", err.Error())))
			}
		}

	}

	return nil
}

func genHandlerImports(parentPkg, subDir string, hasRequest bool) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s/rest/common/result\"", vars.ProjectOpenSourceURL))
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, serviceDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, serviceDir, subDir)))
	if hasRequest {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}
	return strings.Join(imports, "\n\t")
}
