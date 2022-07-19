package mysql

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/model/common"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
	cliUtil "github.com/cocoup/go-smart/tools/gocli/utils"
	"github.com/cocoup/go-smart/tools/gocli/utils/format"
	"github.com/cocoup/go-smart/tools/gocli/utils/pathx"
)

var (
	//go:embed model.tpl
	modelTemplate string

	//go:embed vars.tpl
	varsTemplate string
)

var (
	DSN       string
	OutDir    string
	CodeStyle string
	Tables    string
	UseCache  bool
	JsonStyle string //json名称格式
)

func DoCmd(_ *cobra.Command, _ []string) error {
	if 0 == len(DSN) {
		return errors.New("missing -d")
	}

	cfg, err := config.NewConfig(CodeStyle)
	if err != nil {
		return err
	}

	pkg := filepath.Base(OutDir)
	err = genVars(OutDir, pkg)
	if nil != err {
		fmt.Println(aurora.Red(fmt.Sprintf("generate vars error, %s", err.Error())))
		return err
	}

	tablesPattern := common.ParseTableList(Tables)
	return genTableModels(OutDir, pkg, cfg, tablesPattern)
}

func genVars(outDir, pkg string) error {
	err := pathx.MkdirIfNotExist(outDir)
	if err != nil {
		return err
	}

	varFilename := "vars"

	filename := filepath.Join(outDir, varFilename+".go")
	text, err := pathx.LoadTemplate(category, varsTemplateFile, varsTemplate)
	if err != nil {
		return err
	}

	err = cliUtil.With("vars").Parse(text).SaveTo(map[string]interface{}{
		"package": pkg,
	}, filename, false)
	if err != nil {
		return err
	}

	return nil
}

func genTableModels(dir, pkg string, cfg *config.Config, pattern common.Pattern) error {
	sql, err := NewSqlModel(DSN)
	if nil != err {
		return err
	}

	tables, err := sql.GetAllTables()
	if nil != err {
		return err
	}

	for _, table := range tables {
		if !pattern.Match(table) {
			continue
		}

		if tableData, err := sql.GetColumns(table); nil != err {
			return err
		} else {
			if err := genModel(dir, pkg, cfg, tableData); nil != err {
				fmt.Println(aurora.Red(fmt.Sprintf("%s, ignored generation", err.Error())))
				continue
			}
		}
	}

	fmt.Println(aurora.Green("Done."))
	return nil
}

func jsonName(s string) string {
	name, _ := format.NamingFormat(JsonStyle, s)
	return name
}

func gormName(s string) string {
	name, _ := format.NamingFormat("go_zero", s)
	return name
}

func genModel(dir, pkg string, cfg *config.Config, tableData *common.Table) error {
	fields, hasTime, err := getTableFields(tableData)
	if nil != err {
		return err
	}

	filename, err := format.NamingFormat(cfg.NamingFormat, tableData.Table)
	if err != nil {
		return err
	}

	modelName, _ := format.NamingFormat("GoZero", tableData.Table)
	snakeModelName, _ := format.NamingFormat("go_zero", tableData.Table)

	fp, created, err := utils.MaybeCreateFile(dir, "", filename+".go")
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var text string
	if len(modelTemplate) == 0 {
		text = modelTemplate
	} else {
		text, err = pathx.LoadTemplate(category, modelTemplateFile, modelTemplate)
		if err != nil {
			return err
		}
	}

	t := template.Must(template.New(modelTemplateFile).Funcs(template.FuncMap{
		"jsonName": jsonName,
		"gormName": gormName,
	}).Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]interface{}{
		"package":    pkg,
		"hasTime":    hasTime,
		"imports":    "",
		"model":      modelName,
		"fields":     fields,
		"snakeModel": snakeModelName,
	})
	if err != nil {
		return err
	}

	code := utils.FormatCode(buffer.String())
	_, err = fp.WriteString(code)
	return err
}

func getTableFields(table *common.Table) ([]common.Field, bool, error) {
	var fields = make([]common.Field, len(table.Columns))
	var hasId bool
	var hasTime bool
	for idx, column := range table.Columns {
		dt, err := DB2Go(column.DataType, false)
		if err != nil {
			return nil, hasTime, err
		}
		if strings.Index(dt, "time") >= 0 {
			hasTime = true
		}
		if "id" == strings.ToLower(column.Name) {
			hasId = true
		}
		fieldName, err := format.NamingFormat("GoZero", column.Name)
		if nil != err {
			return nil, hasTime, err
		}

		field := common.Field{
			Name:     fieldName,
			DataType: dt,
			Comment:  column.Comment,
		}
		fields[idx] = field
	}
	if !hasId {
		return nil, hasTime, errors.New(fmt.Sprintf("%s missing id field", table.Table))
	}
	return fields, hasTime, nil
}
