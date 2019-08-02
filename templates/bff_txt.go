package templates

const (
	bffImplTxt = `package {{.PackageName}}

import (	
	"github.com/geekymedic/neon/bff"
)

// @type: b.i.rt
// @interface: {{.InterfaceName}}Handler
type {{.InterfaceName}}Request struct {
}

// @type: b.i.re
// @interface: {{.InterfaceName}}Handler
type {{.InterfaceName}}Response struct {
}

// @type: b.i
// @name: {{.TagZh}}
// @login: {{.TagLogin}}
// @page: {{.TagPage}}
// @uri: {{.TagURI}}
func {{.InterfaceName}}Handler(state *bff.State) {
	var (
		ask = &{{.InterfaceName}}Request{}
		ack   = &{{.InterfaceName}}Response{}
	)
	if err := state.ShouldBindJSON(ask); err != nil {
		state.Error(bff.CodeRequestBodyError, err)
		return
	}
	
	state.Success(ack)
}
`

	errCodeTxt = `package codes

import (
	"github.com/geekymedic/neon/bff"
)

var (
	_codes = bff.Codes{}
)

func GetMessage(code int) string {
	return _codes[code]
}

func init() {
	bff.MergeCodes(_codes)
}`

	bffConfigYmlTxt = `Name: {{.Name}}

Address: ":2243"

Servers:
    UserServer: "127.0.0.1:50054"

Metrics:
    Address: "0.0.0.0:9090"

Log:
    Out: "consol"
    Level: "debug"
    Dir: ""`

	bffConfigTxt = `package config

import (
	"sync"

	"github.com/spf13/viper"
)

var cfg = &Config{}

type Config struct {
	Address string
	Metrics MetricsConfig` + "`yaml:\"Metrics\"`" + `
	Servers ServersConfig` + "`yaml:\"Servers\"`" + `
	Log     LogConfig ` + "`yaml:\"Log\"`" + `
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

	bffRouterTxt = `package router

import (
	impls "{{.ImplsImport}}"

	"github.com/geekymedic/neon/bff"
)

func init() {
	var (
		engine = bff.Engine()
		group  = engine.Group("{{.GroupRouter}}")
	)
	group.POST("/{{.InterfaceName}}", bff.HttpHandler(impls.{{.CamelInterfaceName}}Handler))
}`

	bffHookTxt = `package hook

import (
	{{.Alias}} "{{.ImportPackage}}"

	_ "github.com/geekymedic/neon/utils/validator"
	"github.com/geekymedic/neon/bff"
)

// BeforeRun runs before run
func BeforeRun() {
	var opts = []func() error {func() error {
		config.GetConfig()
		return nil
	}}

	bff.RegisterBeforeAppRunHook(opts...)
}

// BeforeExit is executed before process exits
func BeforeExit() {
}`

	bffMainTxt = `package main
import (
	"fmt"

	"{{.HookImport}}"
	_ "{{.RouterImport}}"

	"github.com/geekymedic/neon/bff"
	_ "github.com/geekymedic/neon/plugin/metrics"
)

func main() {
	hook.BeforeRun()
	defer hook.BeforeExit()

	err := bff.Main()
	if err != nil {
		fmt.Println(err)
	}
}`
)
