package gogen

import (
	"fmt"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
	"github.com/cocoup/go-smart/tools/gocli/utils/format"
	"github.com/cocoup/go-smart/tools/gocli/vars"
)

const (
	configFile = "config"

	configTemplate = `
		package config
		import {{.authImport}}
		
		type Config struct {
			rest.RestConf {{.bq}}yaml:",inline"{{.bq}}
		}
	`

	//jwtTemplate = ` struct {
	//	AccessSecret string
	//	AccessExpire int64
	//}
	//`
	//jwtTransTemplate = ` struct {
	//	Secret     string
	//	PrevSecret string
	//}
	//`
)

func genConfig(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.NamingFormat(cfg.NamingFormat, configFile)
	if err != nil {
		return err
	}

	//authNames := getAuths(api)
	//var auths []string
	//for _, item := range authNames {
	//	auths = append(auths, fmt.Sprintf("%s %s", item, jwtTemplate))
	//}
	//
	//jwtTransNames := getJwtTrans(api)
	//var jwtTransList []string
	//for _, item := range jwtTransNames {
	//	jwtTransList = append(jwtTransList, fmt.Sprintf("%s %s", item, jwtTransTemplate))
	//}
	authImportStr := fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL)

	return utils.GenFile(utils.FileGenConfig{
		Dir:             dir,
		Subdir:          configDir,
		Filename:        filename + ".go",
		TemplateName:    "configTemplate",
		Category:        category,
		TemplateFile:    configTemplateFile,
		BuiltinTemplate: configTemplate,
		Data: map[string]string{
			"authImport": authImportStr,
			//"auth":       strings.Join(auths, "\n"),
			//"jwtTrans":   strings.Join(jwtTransList, "\n"),
			"bq": "`", //反引号替换(解决嵌套)
		},
	})
}
