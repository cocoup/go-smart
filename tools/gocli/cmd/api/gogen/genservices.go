package gogen

import (
	_ "embed"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/zeromicro/go-zero/core/logx"
	"path"
	"strings"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
	"github.com/cocoup/go-smart/tools/gocli/utils/format"
	"github.com/cocoup/go-smart/tools/gocli/utils/pathx"
)

const contextFilename = "context"

var (
	//go:embed service-context.tpl
	contextTemplate string
	//go:embed service.tpl
	serviceTemplate string
	//go:embed service-impl.tpl
	servicesImplTemplate string
)

func genServiceContext(dir, rootPkg string, cfg *config.Config) error {
	filename, err := format.NamingFormat(cfg.NamingFormat, contextFilename)
	if err != nil {
		return err
	}

	configImport := "\"" + pathx.JoinPackages(rootPkg, configDir) + "\""

	return utils.GenFile(utils.FileGenConfig{
		Dir:             dir,
		Subdir:          serviceDir,
		Filename:        filename + ".go",
		TemplateName:    "contextTemplate",
		Category:        category,
		TemplateFile:    serviceContextTemplateFile,
		BuiltinTemplate: contextTemplate,
		Data: map[string]string{
			"imports": configImport,
			"config":  "*config.Config",
		},
	})
}

func genServices(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	logx.Must(genServiceContext(dir, rootPkg, cfg))

	for _, group := range api.Service.Groups {
		var (
			builder  strings.Builder
			hasTypes bool
		)

		folder := group.GetAnnotation(groupProperty)

		folders := strings.Split(folder, "/")
		pkg := folders[0]
		entity := folders[len(folders)-1]
		service := entity + "Service"

		for _, route := range group.Routes {
			hasTypes = len(route.RequestTypeName()) > 0 || len(route.ResponseTypeName()) > 0

			handler := getHandlerName(route)
			handler = strings.TrimSuffix(handler, "Handler")

			filename, err := format.NamingFormat(cfg.NamingFormat, handler)
			if err != nil {
				fmt.Println(aurora.Red(fmt.Sprintf("format file name error, %s\n", err.Error())))
			}

			var (
				requestString  string
				responseString string
				returnString   string
			)
			if len(route.ResponseTypeName()) > 0 {
				resp := responseGoTypeName(route, typesPacket)
				responseString = "(resp " + resp + ", err error)"
				returnString = "return"
			} else {
				responseString = "error"
				returnString = "return nil"
			}
			if len(route.RequestTypeName()) > 0 {
				requestString = "req *" + requestGoTypeName(route, typesPacket)
			}

			err = utils.GenFile(utils.FileGenConfig{
				Dir:             path.Join(dir, serviceDir),
				Subdir:          pkg,
				Filename:        filename + ".go",
				TemplateName:    "servicesTemplate",
				Category:        category,
				TemplateFile:    servicesImplTemplateFile,
				BuiltinTemplate: servicesImplTemplate,
				Data: map[string]string{
					"package":  pkg,
					"imports":  genServiceImplImports(rootPkg, hasTypes),
					"service":  strings.Title(service),
					"handler":  strings.Title(handler),
					"pkg":      pkg,
					"entity":   strings.Title(entity),
					"request":  requestString,
					"response": responseString,
					"return":   returnString,
				},
			})
			if nil != err {
				fmt.Println(aurora.Red(fmt.Sprintf("generate service[%s] error, %s\n", filename, err.Error())))
			}
		}

		err := utils.GenFile(utils.FileGenConfig{
			Dir:             path.Join(dir, serviceDir),
			Subdir:          pkg,
			Filename:        "service.go",
			TemplateName:    "servicesTemplate",
			Category:        category,
			TemplateFile:    serviceTemplateFile,
			BuiltinTemplate: serviceTemplate,
			Data: map[string]string{
				"package":          pkg,
				"imports":          genServiceImports(rootPkg, pkg),
				"service":          strings.Title(service),
				"serviceAdditions": strings.TrimSpace(builder.String()),
			},
		})
		if nil != err {
			fmt.Println(aurora.Red(fmt.Sprintf("generate service[%s] error, %s\n", entity, err.Error())))
		}
	}
	return nil
}

func genServiceImports(parentPkg, subDir string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, serviceDir)))
	return strings.Join(imports, "\n\t")
}

func genServiceImplImports(parentPkg string, hasTypes bool) string {
	var imports []string
	if hasTypes {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}
	return strings.Join(imports, "\n\t")
}
