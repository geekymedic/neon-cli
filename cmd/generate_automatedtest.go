package cmd

import (
	mini_gateway "github.com/geekymedic/neon-cli/mini-gateway"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/util"
	"github.com/spf13/cobra"
	"os"
)

var generateAutomatedTestOpt = struct {
	SysDir    string
	Out       string
}{}

var generateAutomatedTestCmd = &cobra.Command{
	Use: "auto-test",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := util.AbsDir(generateAutomatedTestOpt.SysDir)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate automated test plugin: %v", err)
		}
		sysDirNode := types.NewBaseDir(dir)

		dir, err = util.AbsDir(generateAutomatedTestOpt.Out)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate automated test plugin: %v", err)
		}
		outDirNode := types.NewBaseDir(dir)
		err = mini_gateway.GenerateAutomatedTest(sysDirNode, outDirNode)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate automated test plugin: %v", err)
		}
		util.StdoutOk("Generate automated test plugin successfully")
	},
}

func init() {
	curDir, _ := os.Getwd()
	generateAutomatedTestCmd.Flags().StringVar(&generateAutomatedTestOpt.SysDir, "sys-dir", curDir, "system directory")
	generateAutomatedTestCmd.Flags().StringVar(&generateAutomatedTestOpt.Out, "out-dir", curDir + types.Separator + "integration-test", "dst directory")
}