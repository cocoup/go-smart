package api

import (
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/gogen"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var Cmd = &cobra.Command{
	Use:   "api",
	Short: "Generate api related files",
}

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Generate go files for provided api",
	//Example: "go-smart api new [options] service-name",
	RunE: gogen.DoCmd,
}

func init() {
	goCmd.Flags().StringVarP(&gogen.ApiFile, "file", "f", "", "api file")
	goCmd.Flags().StringVarP(&gogen.OutDir, "out", "o", "./", "The target dir")
	goCmd.Flags().StringVarP(&gogen.CodeStyle, "style", "s", config.DefaultFormat, "code style"+
		"  see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")

	Cmd.AddCommand(goCmd)
}
