package apigen

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/util"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"html/template"
	"path/filepath"
	"strings"
)

//go:embed api.tpl
var apiTemplate string

var ApiFile string

// CreateApiTemplate create api template file
func DoCmd(_ *cobra.Command, _ []string) error {
	if len(ApiFile) == 0 {
		return errors.New("missing -o")
	}

	fp, err := pathx.CreateIfNotExist(ApiFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	text, err := pathx.LoadTemplate(category, apiTemplateFile, apiTemplate)
	if err != nil {
		return err
	}

	baseName := pathx.FileNameWithoutExt(filepath.Base(ApiFile))
	if strings.HasSuffix(strings.ToLower(baseName), "-api") {
		baseName = baseName[:len(baseName)-4]
	} else if strings.HasSuffix(strings.ToLower(baseName), "api") {
		baseName = baseName[:len(baseName)-3]
	}

	t := template.Must(template.New("apiTemplate").Parse(text))
	if err := t.Execute(fp, map[string]string{
		"gitUser":     util.GetGitName(),
		"gitEmail":    util.GetGitEmail(),
		"serviceName": baseName,
	}); err != nil {
		return err
	}

	fmt.Println(aurora.Green("Done."))
	return nil
}
