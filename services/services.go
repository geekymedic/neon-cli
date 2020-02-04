package services

import (
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"github.com/laohanlinux/converter"
)

type Server interface {
	CmdServer
}

type CmdServer interface{}

type EmptyArg struct{}

type EmptyReply struct{}

var NoReply = &EmptyReply{}

type GenServerApiDocArg struct {
	Out    types.DirNode
	Domain string
}

type GenServerCreeateSystemArg struct {
	Out  types.DirNode
	Name string
}

type GenServerCreateBffArg struct {
	Out     types.DirNode
	Name    string
	SysName string
	SysDir  types.DirNode
	Impl    *GenServerCreateBffImplArg
}

type GenServerCreateBffImplArg struct {
	Out     types.DirNode // eg: impls/ping or impls
	Name    string        // eg: ping
	BffName string        // eg: demo
	SysName string        // eg: demo

	//RpcPath  string
	TagZh    string
	TagLogin string
	TagPage  string
	TagURI   string
}

type GenServerCreateBffRouterArg struct {
	Out        types.DirNode
	BffName    string
	ImplName   string
	SysName    string
	SubffixOpt string // impls/demo/ping.go --> /demo
}

type GenServerServiceArg struct {
	Out     types.DirNode
	Name    string
	SysName string
	SysDir  types.DirNode
	Impl    *GenServerCreateServiceImplArg
}

type GenServerCreateServiceImplArg struct {
	Out         types.DirNode // eg: impls/ping or impls
	Name        string        // eg: ping
	ServiceName string        // eg: ping
	SysName     string        // eg: demo
	SubffixOpt  string        // impls/ping/ping.go --> /ping

	//TagZh    string
	//TagLogin string
	//TagPage  string
	//TagURI   string
}

type GenServerAutomatedTestArg struct {
	Out    types.DirNode
	Domain string
}

type GenServerORMArg struct {
	Table *converter.Table2Struct
}

type GenServerJobArg struct {
	Name    string
	CmdName string
	SysDir  types.DirNode
}

type OptsFunc func(*BaseCmdServer)

func NewBaseCmdServer(opts ...OptsFunc) *BaseCmdServer {
	s := &BaseCmdServer{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type BaseCmdServer struct {
	GenerServer *GenerateServer
}

type GenerateServer struct {
	sys *sysdes.SystemDes
}

func NewGenerateServer(sys *sysdes.SystemDes) *GenerateServer {
	return &GenerateServer{sys: sys}
}
