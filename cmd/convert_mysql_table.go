package cmd

import (
	mini_gateway "github.com/geekymedic/neon-cli/mini-gateway"
	"github.com/geekymedic/neon-cli/util"
	"github.com/laohanlinux/converter"
	"github.com/spf13/cobra"
)

var convertMySQLCmd = &cobra.Command{
	Use: "mysql",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			address string
			table   string
			err     error
		)
		if address, err = cmd.Flags().GetString("address"); err != nil {
			util.StdoutExit(-1, "Fail to convert mysql to Go struct: %v", err)
		}
		if table, err = cmd.Flags().GetString("table"); err != nil {
			util.StdoutExit(-1, "Fail to convert mysql to Go struct: %v", err)
		}
		t := converter.NewTable2Struct()
		err = t.TagKey("orm").
			Table(table).
			Dsn(address).
			Run()
		if err != nil {
			util.StdoutExit(-1, "Fail to convert mysql to Go struct: %v", err)
		}
		err = mini_gateway.ORM(t)
		if err != nil {
			util.StdoutExit(-1, "Fail to convert mysql to Go struct: %v", err)
		}
	},
}

func init() {
	convertMySQLCmd.Flags().String("address", "", "root:root@tcp(localhost:3306)/test?charset=utf8")
	convertMySQLCmd.Flags().String("table", "", "")
}
