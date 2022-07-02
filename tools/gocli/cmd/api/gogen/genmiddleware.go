package gogen

import (
	_ "embed"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/util"
	"github.com/cocoup/go-smart/tools/gocli/util/format"
	"strings"
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
		err = util.GenFile(util.FileGenConfig{
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
