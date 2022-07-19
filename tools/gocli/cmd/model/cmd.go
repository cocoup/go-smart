package model

import (
	"github.com/cocoup/go-smart/tools/gocli/cmd/config"
	"github.com/cocoup/go-smart/tools/gocli/cmd/model/mysql"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "model",
		Short: "Generate model code",
	}

	mysqlCmd = &cobra.Command{
		Use:   "mysql",
		Short: "Generate mysql model",
		RunE:  mysql.DoCmd,
	}
)

func init() {
	mysqlCmd.Flags().StringVarP(&mysql.DSN, "dsn", "d", "", `The data source of database,like "root:password@tcp(127.0.0.1:3306)/database"`)
	mysqlCmd.Flags().StringVarP(&mysql.Tables, "table", "t", "*", "The table or table globbing patterns in the database")
	mysqlCmd.Flags().BoolVarP(&mysql.UseCache, "cache", "c", false, "Generate code with cache [optional]")
	mysqlCmd.Flags().StringVarP(&mysql.OutDir, "out", "o", "./model", "The target dir")
	mysqlCmd.Flags().StringVarP(&mysql.CodeStyle, "style", "s", config.DefaultFormat, "The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	mysqlCmd.Flags().StringVarP(&mysql.JsonStyle, "json", "j", "goZero", "The model json name format.")

	Cmd.AddCommand(mysqlCmd)
}
