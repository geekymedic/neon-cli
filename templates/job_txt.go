package templates

const (
	jobMainTxt = `package main

import "{{.ImportPackage}}"

func main() {
	cmd.Execute()
}`
	jobScheduleTxt = `package schedule

import "github.com/geekymedic/neon/logger"

type Schedule interface {
	Run() error
	Stop() error
}

func NewBaseSchedule() *BaseSchedule {
	return &BaseSchedule{}
} 

type BaseSchedule struct{}

func (schedule *BaseSchedule) Run() error {
	logger.Info("BaseSchedule run")
	return nil
}

func (schedule *BaseSchedule) Stop() error {
	logger.Info("BaseSchedule stop")
	return nil
}`
	jobCmdTxt = `package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "A brief description of your application",
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
    rootCmd.AddCommand({{.subCmd}})
	
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.syncorderstatus.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}`
	jobCmdSubTxt = `package cmd
	import (
	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/config"
	"github.com/geekymedic/neon/logger"
	"github.com/spf13/viper"
	"{{.ImportPacket}}"	
	
	"github.com/spf13/cobra"
)

	var {{.CmdName}} = &cobra.Command{
		Use: "{{.CobraUse}}",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.With("job-name", "sync-refund")
			if err := config.Load(&cfgFile); err != nil {
				log.Error(err)
				return
			}
			if err := neon.LoadPlugins(viper.GetViper()); err != nil {
				log.With("err", err).Error("fail to load plugin")
				return
			}

			var job schedule.Schedule = schedule.New{{.Schedule}}()
			defer job.Stop()
			if err := job.Run(); err != nil {
				log.Error(err)
			}
		},
	}
`
	)