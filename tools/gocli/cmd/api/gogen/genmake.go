package gogen

import (
	_ "embed"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
)

//go:embed makefile.tpl
var makeTemplate string

func genMake(dir string, api *spec.ApiSpec) error {
	fileName := "Makefile"

	service := api.Service

	return utils.GenFile(utils.FileGenConfig{
		Dir:             dir,
		Subdir:          "./",
		Filename:        fileName,
		TemplateName:    "makeTemplate",
		Category:        category,
		TemplateFile:    makefileTemplateFile,
		BuiltinTemplate: makeTemplate,
		Data: map[string]string{
			"serviceName": service.Name,
		},
	})
}
