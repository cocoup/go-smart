package gogen

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
	"github.com/cocoup/go-smart/tools/gocli/utils/format"
)

const typesFile = "types"

//go:embed types.tpl
var typesTemplate string

// BuildTypes gen types to string
func BuildTypes(types []spec.Type) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n\n")
		}
		if err := writeType(&builder, tp); err != nil {
			return "", utils.WrapErr(err, "Type "+tp.Name()+" generate error")
		}
	}

	return builder.String(), nil
}

func genTypes(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	val, err := BuildTypes(api.Types)
	if err != nil {
		return err
	}

	typeFilename, err := format.NamingFormat(cfg.NamingFormat, typesFile)
	if err != nil {
		return err
	}

	typeFilename = typeFilename + ".go"
	filename := path.Join(dir, typesDir, typeFilename)
	os.Remove(filename)

	return utils.GenFile(utils.FileGenConfig{
		Dir:             dir,
		Subdir:          typesDir,
		Filename:        typeFilename,
		TemplateName:    "typesTemplate",
		Category:        category,
		TemplateFile:    typesTemplateFile,
		BuiltinTemplate: typesTemplate,
		Data: map[string]interface{}{
			"types":        val,
			"containsTime": false,
		},
	})
}

func writeType(writer io.Writer, tp spec.Type) error {
	structType, ok := tp.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unspport struct type: %s", tp.Name())
	}

	fmt.Fprintf(writer, "type %s struct {\n", strings.Title(tp.Name()))
	for _, member := range structType.Members {
		if member.IsInline {
			if _, err := fmt.Fprintf(writer, "%s\n", strings.Title(member.Type.Name())); err != nil {
				return err
			}

			continue
		}

		if err := writeProperty(writer, member.Name, member.Tag, member.GetComment(), member.Type, 1); err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "}")
	return nil
}
