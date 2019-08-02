package mini_gateway

import (
	"github.com/geekymedic/neon-cli/services"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"os"
)

func GenerateMarkdown(sysDirNode, out types.DirNode, domain string) error {
	sys, err := sysdes.NewSystemDes(sysDirNode)
	if err != nil {
		return err
	}
	server := services.NewBaseCmdServer(func(server *services.BaseCmdServer) {
		server.GenerServer = services.NewGenerateServer(sys)
	})
	_, err = server.GenerServer.GenerateApiDoc(nil, &services.GenServerApiDocArg{Out: out, Domain: domain})
	return err
}

func GenerateAutomatedTest(sysDirNode, out types.DirNode) error {
	sys, err := sysdes.NewSystemDes(sysDirNode)
	if err != nil {
		return err
	}
	server := services.NewBaseCmdServer(func(server *services.BaseCmdServer) {
		server.GenerServer = services.NewGenerateServer(sys)
	})
	_, err = server.GenerServer.GenerateAutomatedTest(nil, &services.GenServerAutomatedTestArg{Out: out})
	return err
}

func GenerateSystem(sysDirNode types.DirNode, name string) error {
	server := services.NewBaseCmdServer()
	_, err := server.GenerServer.CreateSystem(nil, &services.GenServerCreeateSystemArg{Out: sysDirNode, Name: name})
	return err
}

func GenerateBff(sysDirNode types.DirNode, bffName string) error {
	sys, err := sysdes.NewSystemDes(sysDirNode)
	if err != nil {
		return err
	}
	server := services.NewBaseCmdServer(func(server *services.BaseCmdServer) {
		server.GenerServer = services.NewGenerateServer(sys)
	})
	bffBaseDir := sys.DirNode.Append("bff").(types.DirNode)
	implArg := services.GenServerCreateBffImplArg{
		Out:      bffBaseDir.Append(bffName, "impls", "demo").(types.DirNode),
		Name:     "demo",
		BffName:  bffName,
		SysName:  sys.Name,
		TagZh:    "健康检查",
		TagLogin: "N",
		TagPage:  "",
		TagURI:   "/api" + services.PacketRouter(bffName, "demo", sys.Name),
	}

	_, err = server.GenerServer.CreateBff(nil, &services.GenServerCreateBffArg{
		Out:     bffBaseDir,
		Name:    bffName,
		SysName: sys.ShortName,
		SysDir:  sys.DirNode,
		Impl:    &implArg,
	})

	return err
}

func GenerateService(sysDirNode types.DirNode, serviceName string) error {
	sys, err := sysdes.NewSystemDes(sysDirNode)
	if err != nil {
		return err
	}
	server := services.NewBaseCmdServer(func(server *services.BaseCmdServer) {
		server.GenerServer = services.NewGenerateServer(sys)
	})

	// create demo service
	{
		servicesBaseDir := sys.DirNode.Append("services").(types.DirNode)
		types.AssertNil(servicesBaseDir.Create(os.ModePerm))

		_, err := server.GenerServer.CreateService(nil, &services.GenServerServiceArg{
			Out:     servicesBaseDir,
			Name:    serviceName,
			SysName: sys.ShortName,
			SysDir:  sysDirNode,
			Impl: &services.GenServerCreateServiceImplArg{
				Out:         servicesBaseDir.Append(serviceName, "impls", "ping").(types.DirNode),
				Name:        "ping",
				ServiceName: "CheckHealth",
				SysName:     sys.ShortName,
				SubffixOpt:  "/ping",
			},
		})
		return err
	}
}
