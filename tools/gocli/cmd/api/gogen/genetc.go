package gogen

import (
	_ "embed"
	"fmt"
	"strconv"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/util"
	"github.com/cocoup/go-smart/tools/gocli/util/format"
)

const (
	defaultPort = 8888
	etcDir      = "etc"
	routeRoot   = "/"
)

//go:embed etc.tpl
var etcTemplate string

func genEtc(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	fileName, err := format.NamingFormat(cfg.NamingFormat, api.Service.Name)
	if nil != err {
		return nil
	}

	service := api.Service
	host := "0.0.0.0"
	port := strconv.Itoa(defaultPort)

	return util.GenFile(util.FileGenConfig{
		Dir:             dir,
		Subdir:          etcDir,
		Filename:        fmt.Sprintf("%v.yaml", fileName),
		TemplateName:    "etcTemplate",
		Category:        category,
		TemplateFile:    etcTemplateFile,
		BuiltinTemplate: etcTemplate,
		Data: map[string]string{
			"serviceName": service.Name,
			"host":        host,
			"port":        port,
			"routeRoot":   routeRoot,
		},
	})
}
