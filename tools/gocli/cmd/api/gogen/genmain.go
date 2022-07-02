package gogen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/util"
	"github.com/cocoup/go-smart/tools/gocli/util/format"
	"github.com/cocoup/go-smart/tools/gocli/util/pathx"
	"github.com/cocoup/go-smart/tools/gocli/vars"
)

//go:embed main.tpl
var mainTemplate string

func genMain(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.NamingFormat(cfg.NamingFormat, api.Service.Name)
	if err != nil {
		return err
	}

	configName := filename
	if strings.HasSuffix(filename, "-api") {
		filename = strings.ReplaceAll(filename, "-api", "")
	}

	return util.GenFile(util.FileGenConfig{
		Dir:             dir,
		Subdir:          "",
		Filename:        filename + ".go",
		TemplateName:    "mainTemplate",
		Category:        category,
		TemplateFile:    mainTemplateFile,
		BuiltinTemplate: mainTemplate,
		Data: map[string]string{
			"imports":     genMainImports(rootPkg),
			"serviceName": configName,
		},
	})
}

func genMainImports(parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, configDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, routeDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, serviceDir)))
	imports = append(imports, fmt.Sprintf("\"%s/core/conf\"", vars.ProjectOpenSourceURL))
	imports = append(imports, fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}
