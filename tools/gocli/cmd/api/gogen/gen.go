package gogen

import (
	"errors"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/parser"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/utils"
	"github.com/cocoup/go-smart/tools/gocli/utils/pathx"
)

var (
	ApiFile   string
	OutDir    string
	CodeStyle string
)

func DoCmd(_ *cobra.Command, _ []string) error {
	if 0 == len(ApiFile) {
		return errors.New("missing -f")
	}

	return GenProcjec(ApiFile, OutDir, CodeStyle)
}

func GenProcjec(apiFile, outDir, style string) error {
	defer func() {
		if err := recover(); err != nil {
			//输出panic信息
			logx.Must(errors.New(fmt.Sprintf("解析文件出错: %v", err)))
		}
	}()

	api, err := parser.Parse(apiFile)
	if err != nil {
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	cfg, err := config.NewConfig(style)
	if err != nil {
		return err
	}

	logx.Must(pathx.MkdirIfNotExist(outDir))
	rootPkg, err := utils.GetParentPackage(outDir)
	if err != nil {
		return err
	}

	fmt.Println(rootPkg)

	logx.Must(genEtc(outDir, cfg, api))
	logx.Must(genConfig(outDir, cfg, api))
	logx.Must(genMain(outDir, rootPkg, cfg, api))
	logx.Must(genTypes(outDir, cfg, api))
	logx.Must(genRoutes(outDir, rootPkg, cfg, api))
	logx.Must(genHandlers(outDir, rootPkg, cfg, api))
	logx.Must(genServices(outDir, rootPkg, cfg, api))
	logx.Must(genMiddleware(outDir, cfg, api))
	logx.Must(genMake(outDir, api))

	fmt.Println(aurora.Green("Done."))

	return nil
}
