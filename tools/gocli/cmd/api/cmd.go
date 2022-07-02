package api

import (
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/apigen"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/gogen"
	"github.com/cocoup/go-smart/tools/gocli/cmd/api/new"
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var Cmd = &cobra.Command{
	Use:   "api",
	Short: "Generate api related files",
	RunE:  apigen.DoCmd,
}

var newCmd = &cobra.Command{
	Use:     "new",
	Short:   "Fast create api service",
	Example: "gocli api new [options] service-name",
	Args:    cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return new.DoCmd(args)
	},
}

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Generate go files for provided api",
	//Example: "go-smart api new [options] service-name",
	RunE: gogen.DoCmd,
}

func init() {
	Cmd.Flags().StringVarP(&apigen.ApiFile, "out", "o", "", "Output a sample api file")

	newCmd.Flags().StringVarP(&new.CodeStyle, "style", "s", config.DefaultFormat, "The file naming format,"+
		" see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")

	goCmd.Flags().StringVarP(&gogen.ApiFile, "file", "f", "", "api file")
	goCmd.Flags().StringVarP(&gogen.OutDir, "out", "o", "./api", "The target dir")
	goCmd.Flags().StringVarP(&gogen.CodeStyle, "style", "s", config.DefaultFormat, "code style"+
		"  see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")

	Cmd.AddCommand(newCmd)
	Cmd.AddCommand(goCmd)
}
