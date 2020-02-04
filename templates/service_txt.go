package templates

const (
	serviceConfigYmlTxt = `Name: "{{.Name}}"
`
	serviceConfigTxt = `package config

import (
	"sync"

	"github.com/spf13/viper"
)

var cfg = &Config{}

type Config struct {
	DB      DBConfig
	REDIS   RedisConfig
	Address string
	Metrics MetricsConfig
	Servers ServersConfig
	Log     LogConfig
}

type DBConfig struct {
	User DBUserConfig
}

type DBUserConfig struct {
	DSN     string
	MaxIdle int
	MaxOpen int
}

type RedisConfig struct {
	User RedisUserConfig
}

type RedisUserConfig struct {
	DSN string
}

type MetricsConfig struct {
	Address string
}

type ServersConfig struct {
	UserServer string
}

type LogConfig struct {
	Out   string
	Level string
	Dir   string
}

var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		err := viper.Unmarshal(cfg)
		if err != nil {
			panic(err)
		}
	})
	return cfg
}`

	serviceHookTxt = `package hook

import (
	{{.Alias}} "{{.ImportPackage}}"

	"github.com/geekymedic/neon/service"
)

// BeforeRun runs before run
func BeforeRun() {
	var opts = []func() error {func() error {
		config.GetConfig()
		return nil
	}}

	service.RegisterBeforeAppRunHook(opts...)
}

// BeforeExit is executed before process exits
func BeforeExit() {
}`

	serviceServerTxt = `package impls

import (
	{{.AliasName}} "{{.ImportPackage}}"

	"github.com/geekymedic/neon/service"
)

type {{.ServerName}} struct {
}

func init() {
	{{.SystemShortName}}.Register{{.ServerName}}(service.Server(), &{{.ServerName}}{})
}`

	serviceImplTxt  = `package impls

import (
	{{.AliasName}} "{{.ImportPackage}}"

	"golang.org/x/net/context"
)

func (m *{{.ServerName}}) {{.ImplName}}(ctx context.Context, in *{{.AliasName}}.{{.RequestName}}) (*{{.AliasName}}.{{.ResponseName}}, error) {
	return &{{.AliasName}}.{{.ResponseName}}{}, nil
}`

	serviceServerInitTxt = `package impls

import (
	_ "{{.ImportInit}}"
)
`

	serviceMainTxt = `package main

import (
	"fmt"

	{{range $i, $sube := $.List}}
	{{$sube}}
	{{end}}

	"github.com/geekymedic/neon/service"
	_ "github.com/geekymedic/neon/plugin/metrics"
)

func main() {
	hook.BeforeRun()
	defer hook.BeforeExit()

	if err := service.Main(); err != nil {
		fmt.Printf("fail to start {{.SystemName}}.{{.ServiceName}} service: %v\n", err)
	}
}`
)
