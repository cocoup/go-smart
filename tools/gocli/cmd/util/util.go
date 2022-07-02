package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"

	"go/format"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cocoup/go-smart/tools/gocli/util/ctx"
	"github.com/cocoup/go-smart/tools/gocli/util/pathx"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileGenConfig struct {
	Dir             string
	Subdir          string
	Filename        string
	TemplateName    string
	Category        string
	TemplateFile    string
	BuiltinTemplate string
	Data            interface{}
}

func GenFile(c FileGenConfig) error {
	fp, created, err := MaybeCreateFile(c.Dir, c.Subdir, c.Filename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var text string
	if len(c.Category) == 0 || len(c.TemplateFile) == 0 {
		text = c.BuiltinTemplate
	} else {
		text, err = pathx.LoadTemplate(c.Category, c.TemplateFile, c.BuiltinTemplate)
		if err != nil {
			return err
		}
	}

	t := template.Must(template.New(c.TemplateName).Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, c.Data)
	if err != nil {
		return err
	}

	code := FormatCode(buffer.String())
	_, err = fp.WriteString(code)
	return err
}

func GetParentPackage(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	projectCtx, err := ctx.Prepare(abs)
	if err != nil {
		return "", err
	}

	// fix https://github.com/zeromicro/go-zero/issues/1058
	wd := projectCtx.WorkDir
	d := projectCtx.Dir
	same, err := pathx.SameFile(wd, d)
	if err != nil {
		return "", err
	}

	trim := strings.TrimPrefix(projectCtx.WorkDir, projectCtx.Dir)
	if same {
		trim = strings.TrimPrefix(strings.ToLower(projectCtx.WorkDir), strings.ToLower(projectCtx.Dir))
	}

	return filepath.ToSlash(filepath.Join(projectCtx.Path, trim)), nil
}

func FormatCode(code string) string {
	ret, err := format.Source([]byte(code))
	if err != nil {
		return code
	}

	return string(ret)
}

// MaybeCreateFile creates file if not exists
func MaybeCreateFile(dir, subdir, file string) (fp *os.File, created bool, err error) {
	logx.Must(pathx.MkdirIfNotExist(path.Join(dir, subdir)))
	fpath := path.Join(dir, subdir, file)
	if pathx.FileExists(fpath) {
		fmt.Println(aurora.Yellow(fmt.Sprintf("%s exists, ignored generation", fpath)))
		return nil, false, nil
	}

	fp, err = pathx.CreateIfNotExist(fpath)
	created = err == nil
	return
}

// WrapErr wraps an error with message
func WrapErr(err error, message string) error {
	return errors.New(message + ", " + err.Error())
}

// Copy calls io.Copy if the source file and destination file exists
func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
