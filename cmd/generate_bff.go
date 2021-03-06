/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	mini_gateway "github.com/geekymedic/neon-cli/mini-gateway"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/util"
	"github.com/spf13/cobra"
	"os"
)

var generatebffCmdOpt = struct {
	Name     string
	implName string
}{}

// bffCmd represents the bff command
var generateBffCmd = &cobra.Command{
	Use:   "bff",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		// build system
		dir, err := util.AbsDir(dir)
		if err != nil {
			util.StdoutExit(-1, "Fail to generate bff: %v", err)
		}
		sysDirNode := types.NewBaseDir(dir)
		err = mini_gateway.GenerateBff(sysDirNode, generatebffCmdOpt.Name)
		if err != nil {
			util.StdoutExit(-1, "Fail to create bff: %v", err)
		}
		util.StdoutOk("Create bff successfully\n")
	},
}

func init() {
	generateCmd.AddCommand(generateBffCmd)

	generateBffCmd.Flags().StringVar(&generatebffCmdOpt.Name, "name", "demo", "bff name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
