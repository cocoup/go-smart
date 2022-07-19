package gogen

import (
	"fmt"
	"github.com/cocoup/go-smart/tools/gocli/utils/pathx"
)

const (
	category                         = "api"
	configTemplateFile               = "config.tpl"
	serviceContextTemplateFile       = "service-context.tpl"
	serviceTemplateFile              = "service.tpl"
	servicesImplTemplateFile         = "service-impl.tpl"
	etcTemplateFile                  = "etc.tpl"
	makefileTemplateFile             = "makefile.tpl"
	handlerTemplateFile              = "handler.tpl"
	handlersAdditionTemplateFile     = "handler-addition.tpl"
	mainTemplateFile                 = "main.tpl"
	middlewareImplementCodeFile      = "middleware.tpl"
	routeTemplateFile                = "route.tpl"
	routeAdditionTemplateFile        = "route-addition.tpl"
	routePackageTemplateFile         = "route-package.tpl"
	routePackageAdditionTemplateFile = "route-package-addition.tpl"
	typesTemplateFile                = "types.tpl"
)

var templates = map[string]string{
	etcTemplateFile: etcTemplate,

	//configTemplateFile:          configTemplate,
	//contextTemplateFile:         contextTemplate,
	//handlerTemplateFile:         handlerTemplate,
	//logicTemplateFile:           logicTemplate,
	//mainTemplateFile:            mainTemplate,
	//middlewareImplementCodeFile: middlewareImplementCode,
	//routesTemplateFile:          routesTemplate,
	//routesAdditionTemplateFile:  routesAdditionTemplate,
	//typesTemplateFile:           typesTemplate,
}

// Category returns the category of the api files.
func Category() string {
	return category
}

// Clean cleans the generated deployment files.
func Clean() error {
	return pathx.Clean(category)
}

// GenTemplates generates api template files.
func GenTemplates() error {
	return pathx.InitTemplates(category, templates)
}

// RevertTemplate reverts the given template file to the default value.
func RevertTemplate(name string) error {
	content, ok := templates[name]
	if !ok {
		return fmt.Errorf("%s: no such file name", name)
	}
	return pathx.CreateTemplate(category, name, content)
}

// Update updates the template files to the templates built in current goctl.
func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return pathx.InitTemplates(category, templates)
}
