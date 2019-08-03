package services

import (
	"context"
	"fmt"
	"github.com/geekymedic/neon-cli/templates"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"github.com/geekymedic/neon-cli/util"
	"github.com/geekymedic/neon/errors"
	"github.com/iancoleman/strcase"
	"os"
)

func (s *GenerateServer) CreateService(_ context.Context, arg *GenServerServiceArg) (*EmptyReply, error) {
	util.StdoutOk("Start create service\n")
	serviceBaseDir := arg.Out.Append(arg.Name).(types.DirNode)
	if err := serviceBaseDir.IsExist(); err == os.ErrExist {
		return NoReply, errors.NewStackError("bff has exist")
	}

	if err := serviceBaseDir.Create(os.ModePerm); err != nil {
		return NoReply, err
	}

	// create makefile
	{
		makefileFp := types.NewBaseFile(serviceBaseDir.Append("Makefile").(types.DirNode).Abs())
		types.AssertNil(makefileFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.BffMakefileTemplate,
			map[string]interface{}{"Name": arg.Name, "Typ": "services", "System": arg.SysName})
		types.AssertNil(err)
		_, err = makefileFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(makefileFp.Close())
		util.StdoutOk("Create service makefile successfully\n")
	}

	// create docker file
	{
		dockerFp := types.NewBaseFile(serviceBaseDir.Append("Dockfile").(types.DirNode).Abs())
		types.AssertNil(dockerFp.Create(types.DefFlag, types.DefPerm))
		types.AssertNil(dockerFp.Close())
		util.StdoutOk("Create service dockfile successfully\n")
	}

	// create k8s
	{
		k8sFp := types.NewBaseFile(serviceBaseDir.Append(".k8s.yml").(types.DirNode).Abs())
		types.AssertNil(k8sFp.Create(types.DefFlag, types.DefPerm))
		types.AssertNil(k8sFp.Close())
		util.StdoutOk("Create service k8s config successfully\n")
	}

	// create config
	{
		configBaseDir := serviceBaseDir.Append("config").(types.DirNode)
		types.AssertNil(configBaseDir.Create(os.ModePerm))
		ymlFp := types.NewBaseFile(configBaseDir.Append("config.yml").(types.DirNode).Abs())
		types.AssertNil(ymlFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.ServiceConfigYmlTpl, map[string]interface{}{
			"Name": fmt.Sprintf("%s%s-services-%s", arg.SysName, sysdes.SystemNameSubffix, arg.Name)})
		types.AssertNil(err)
		_, err = ymlFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(ymlFp.Close())

		configFp := types.NewBaseFile(configBaseDir.Append("config.go").(types.DirNode).Abs())
		types.AssertNil(configFp.Create(types.DefFlag, types.DefPerm))
		txt, err = templates.ParseTemplate(templates.ServiceConfigTpl, nil)
		types.AssertNil(err)
		_, err = configFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(configFp.Close())
	}

	// create hook
	{
		hookDir := serviceBaseDir.Append("hook").(types.DirNode)
		types.AssertNil(hookDir.Create(os.ModePerm))
		hookFp := types.NewBaseFile(hookDir.Append("hook.go").(types.DirNode).Abs())
		types.AssertNil(hookFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.ServiceHookTpl, templates.ServiceHookTplArg{
			Alias:         "config",
			ImportPackage: fmt.Sprintf("%s/%s/services/%s/config", sysdes.GoHostPreffix, arg.SysDir.Name(), arg.Name),
		})
		types.AssertNil(err)
		_, err = hookFp.WriteString(txt)
		types.AssertNil(err)
	}

	// create impls
	{
		_, err := s.CreateServiceImpl(nil, arg.Impl)
		types.AssertNil(err)

		// create server init
		registerServerBaseDir := serviceBaseDir.Append("impls", "register_server").(types.DirNode)
		types.AssertNil(registerServerBaseDir.Create(os.ModePerm))
		serverInitFp := types.NewBaseFile(registerServerBaseDir.Append("init.go").(types.DirNode).Abs())
		types.AssertNil(serverInitFp.Create(types.DefFlag, types.DefPerm))
		var serviceServerInitTplArg = fmt.Sprintf("%s/%s%s/services/%s/impls",
			sysdes.GoHostPreffix, arg.SysName, sysdes.SystemNameSubffix,
			arg.Name)
		if arg.Impl.SubffixOpt != "" {
			serviceServerInitTplArg += arg.Impl.SubffixOpt
		}
		txt, err := templates.ParseTemplate(templates.ServiceServerInitTpl, map[string]interface{}{"ImportInit": serviceServerInitTplArg})
		types.AssertNil(err)
		_, err = serverInitFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(serverInitFp.Close())
	}

	// create main.go
	{
		mainFp := types.NewBaseFile(serviceBaseDir.Append("main.go").(types.DirNode).Abs())
		types.AssertNil(mainFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.ServiceMainTpl, templates.ServiceMainTplArg{
			ServiceName: arg.SysName,
			SystemName:  arg.Name,
			List: []string{
				fmt.Sprintf("\"%s/%s%s/services/%s/hook\"", sysdes.GoHostPreffix, arg.SysName, sysdes.SystemNameSubffix, arg.Name),
				fmt.Sprintf("_ \"%s/%s%s/services/%s/impls/register_server\"", sysdes.GoHostPreffix, arg.SysName, sysdes.SystemNameSubffix, arg.Name),
			},
		})
		types.AssertNil(err)
		_, err = mainFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(mainFp.Close())
	}

	return NoReply, nil
}

func (s *GenerateServer) CreateServiceImpl(ctx context.Context, arg *GenServerCreateServiceImplArg) (*EmptyReply, error) {
	util.StdoutOk("Start create service impl: %s\n", arg.Name)
	implBaseDir := arg.Out
	if err := implBaseDir.Create(os.ModePerm); err != nil {
		return NoReply, err
	}
	// create server file
	{
		serverFp := types.NewBaseFile(implBaseDir.Append(fmt.Sprintf("%s_server.go", strcase.ToSnake(arg.ServiceName))).(types.DirNode).Abs())
		types.AssertNil(serverFp.Create(types.DefFlag, types.DefPerm))
		tplArg := templates.ServiceServerTplArg{
			ImportPackage: fmt.Sprintf("%s/%s%s/%s/%s",
				sysdes.GoHostPreffix,
				arg.SysName,
				sysdes.SystemNameSubffix,
				sysdes.ProtocolName,
				arg.SysName),
			AliasName:       arg.SysName,
			SystemShortName: arg.SysName,
			ServerName:      fmt.Sprintf("%sServer", strcase.ToCamel(arg.ServiceName)),
		}
		txt, err := templates.ParseTemplate(templates.ServiceServerTpl, tplArg)
		types.AssertNil(err)
		_, err = serverFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(serverFp.Close())
	}

	// create impl file
	{
		serviceImplFp := types.NewBaseFile(implBaseDir.Append(fmt.Sprintf("%s.go", arg.Name)).(types.DirNode).Abs())
		types.AssertNil(serviceImplFp.Create(types.DefFlag, types.DefPerm))
		tplArg := templates.ServiceImplArg{
			ImportPackage: fmt.Sprintf("%s/%s%s/%s/%s",
				sysdes.GoHostPreffix,
				arg.SysName,
				sysdes.SystemNameSubffix,
				sysdes.ProtocolName,
				arg.SysName),
			AliasName:    arg.SysName,
			ServerName:   fmt.Sprintf("%sServer", strcase.ToCamel(arg.ServiceName)),
			ImplName:     strcase.ToCamel(arg.Name),
			RequestName:  fmt.Sprintf("%sRequest", strcase.ToCamel(arg.Name)),
			ResponseName: fmt.Sprintf("%sResponse", strcase.ToCamel(arg.Name)),
		}
		txt, err := templates.ParseTemplate(templates.ServiceImplTpl, &tplArg)
		types.AssertNil(err)
		_, err = serviceImplFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(serviceImplFp.Close())
	}
	return NoReply, nil
}
