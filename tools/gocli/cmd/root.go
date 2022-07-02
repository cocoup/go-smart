package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api"
	"github.com/cocoup/go-smart/tools/gocli/cmd/model"
	"github.com/cocoup/go-smart/tools/gocli/version"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gocli",
	Short: "A cli tool to generate frame code",
	Long: "A cli tool to generate api, grpc, model code\n" +
		"GitHub: https://github.com/cocoup/go-smart \n" +
		"Site: ***",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s %s/%s", version.BuildVersion, runtime.GOOS, runtime.GOARCH)
	//fmt.Println(rootCmd.UsageTemplate())

	rootCmd.AddCommand(api.Cmd)
	rootCmd.AddCommand(model.Cmd)
}
