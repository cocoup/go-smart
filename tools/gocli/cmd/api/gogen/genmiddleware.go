package gogen

import (
	_ "embed"
	"strings"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
	"github.com/cocoup/go-smart/tools/gocli/utils/format"
)

//go:embed middleware.tpl
var middlewareImplementCode string

func genMiddleware(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	middlewares := getMiddleware(api)
	for _, item := range middlewares {
		filename, err := format.NamingFormat(cfg.NamingFormat, item)

		if err != nil {
			return err
		}

		name := strings.TrimSuffix(item, "Middleware")
		err = utils.GenFile(utils.FileGenConfig{
			Dir:             dir,
			Subdir:          middlewareDir,
			Filename:        filename + ".go",
			TemplateName:    "contextTemplate",
			Category:        category,
			TemplateFile:    middlewareImplementCodeFile,
			BuiltinTemplate: middlewareImplementCode,
			Data: map[string]string{
				"name": strings.Title(name),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
