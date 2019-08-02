package cmd

import (
	mini_gateway "github.com/geekymedic/neon-cli/mini-gateway"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/util"
	"github.com/spf13/cobra"
	"os"
)

var generateOpt = struct {
	SysDir    string
	Out       string
	ApiDomain string
}{}

var generateMdCmd = &cobra.Command{
	Use: "md",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := util.AbsDir(generateOpt.SysDir)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate markdown: %v", err)
		}
		sysDirNode := types.NewBaseDir(dir)

		dir, err = util.AbsDir(generateOpt.Out)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate markdown: %v", err)
		}
		outDirNode := types.NewBaseDir(dir)
		err = mini_gateway.GenerateMarkdown(sysDirNode, outDirNode, generateOpt.ApiDomain)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate markdown: %v", err)
		}
		util.StdoutOk("Generate markdown successfully")
	},
}

func init() {
	curDir, _ := os.Getwd()
	generateMdCmd.Flags().StringVar(&generateOpt.SysDir, "sys-dir", curDir, "system directory")
	generateMdCmd.Flags().StringVar(&generateOpt.Out, "out-dir", curDir + types.Separator + "doc", "dst directory")
	generateMdCmd.Flags().StringVar(&generateOpt.ApiDomain, "domain", "api.geekymedic.com.cn", "api http domain")
}
