package new

import (
	_ "embed"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/gogen"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/util"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/util/pathx"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

//go:embed api.tpl
var apiTemplate string

var CodeStyle string

// CreateServiceCommand fast create service
func DoCmd(args []string) error {
	dirName := args[0]
	if len(CodeStyle) == 0 {
		CodeStyle = config.DefaultFormat
	}

	abs, err := filepath.Abs(dirName)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	dirName = filepath.Base(filepath.Clean(abs))
	filename := dirName + ".api"
	apiFilePath := filepath.Join(abs, filename)
	fp, err := os.Create(apiFilePath)
	if err != nil {
		return err
	}

	defer fp.Close()

	text, err := pathx.LoadTemplate(category, apiTemplateFile, apiTemplate)
	if err != nil {
		return err
	}

	t := template.Must(template.New("template").Parse(text))
	if err := t.Execute(fp, map[string]string{
		"gitUser":  util.GetGitName(),
		"gitEmail": util.GetGitEmail(),
		"name":     dirName,
		"handler":  strings.Title(dirName),
	}); err != nil {
		return err
	}

	err = gogen.GenProcjec(apiFilePath, abs, CodeStyle)
	return err
}
