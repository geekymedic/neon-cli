package services

import (
	"context"
	"fmt"
	"os"

	"github.com/geekymedic/neon-cli/templates"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"github.com/geekymedic/neon-cli/util"

	"github.com/geekymedic/neon/errors"
	"github.com/iancoleman/strcase"
)

func (s *GenerateServer) CreateBff(ctx context.Context, arg *GenServerCreateBffArg) (*EmptyReply, error) {
	util.StdoutOk("Start create bff\n")
	bffBaseDir := arg.Out.Append(arg.Name).(types.DirNode)
	if err := bffBaseDir.IsExist(); err == os.ErrExist {
		return NoReply, errors.NewStackError("bff has exist")
	}

	if err := bffBaseDir.Create(os.ModePerm); err != nil {
		return NoReply, err
	}

	// create makefile
	{
		makefileFp := types.NewBaseFile(bffBaseDir.Append("Makefile").(types.DirNode).Abs())
		types.AssertNil(makefileFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.BffMakefileTemplate,
			map[string]interface{}{"Name": arg.Name, "Typ": "bff", "System": arg.SysName})
		types.AssertNil(err)
		_, err = makefileFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(makefileFp.Close())
		util.StdoutOk("Create bff makefile successfully\n")
	}

	// create docker file
	{
		dockerFp := types.NewBaseFile(bffBaseDir.Append("Dockfile").(types.DirNode).Abs())
		types.AssertNil(dockerFp.Create(types.DefFlag, types.DefPerm))
		types.AssertNil(dockerFp.Close())
		util.StdoutOk("Create bff dockfile successfully\n")
	}

	// create k8s
	{
		k8sFp := types.NewBaseFile(bffBaseDir.Append(".k8s.yml").(types.DirNode).Abs())
		types.AssertNil(k8sFp.Create(types.DefFlag, types.DefPerm))
		types.AssertNil(k8sFp.Close())
		util.StdoutOk("Create bff k8s config successfully\n")
	}

	// create config
	{
		configBaseDir := bffBaseDir.Append("config").(types.DirNode)
		types.AssertNil(configBaseDir.Create(os.ModePerm))
		ymlFp := types.NewBaseFile(configBaseDir.Append("config.yml").(types.DirNode).Abs())
		types.AssertNil(ymlFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.BffConfigYmlTemplate, map[string]interface{}{
			"Name": fmt.Sprintf("%s%s-bff-%s", arg.SysName, sysdes.SystemNameSubffix, arg.Name)})
		types.AssertNil(err)
		_, err = ymlFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(ymlFp.Close())

		configFp := types.NewBaseFile(configBaseDir.Append("config.go").(types.DirNode).Abs())
		types.AssertNil(configFp.Create(types.DefFlag, types.DefPerm))
		txt, err = templates.ParseTemplate(templates.BffConfigTemplate, nil)
		types.AssertNil(err)
		_, err = configFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(configFp.Close())
	}

	// create codes
	{
		codesBaseDir := bffBaseDir.Append("codes").(types.DirNode)
		types.AssertNil(codesBaseDir.Create(os.ModePerm))

		codeFp := types.NewBaseFile(codesBaseDir.Append("error_code.go").(types.DirNode).Abs())
		types.AssertNil(codeFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.BffErrCodeTemplate, nil)
		types.AssertNil(err)
		_, err = codeFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(codeFp.Close())
	}

	{

		// create hook
		{
			hookDir := bffBaseDir.Append("hook").(types.DirNode)
			types.AssertNil(hookDir.Create(os.ModePerm))
			hookFp := types.NewBaseFile(hookDir.Append("hook.go").(types.DirNode).Abs())
			types.AssertNil(hookFp.Create(types.DefFlag, types.DefPerm))
			txt, err := templates.ParseTemplate(templates.BffHookTemplate, templates.BffHookTplArg{
				Alias:         "config",
				ImportPackage: fmt.Sprintf("%s/%s/bff/%s/config", sysdes.GoHostPreffix, arg.SysDir.Name(), arg.Name),
			})
			types.AssertNil(err)
			_, err = hookFp.WriteString(txt)
			types.AssertNil(err)
		}

		// create impls
		{
			_, err := s.CreateBffImpl(nil, arg.Impl)
			types.AssertNil(err)
		}

		// create router
		{
			var createBffRouterArg GenServerCreateBffRouterArg
			createBffRouterArg.SysName = arg.SysName
			createBffRouterArg.Out = types.NewBaseDir(bffBaseDir.Append("router").(types.DirNode).Abs())
			createBffRouterArg.BffName = arg.Impl.BffName
			createBffRouterArg.ImplName = arg.Impl.Name
			createBffRouterArg.SubffixOpt = "/demo" // TODO
			_, err := s.CreateBffRouter(nil, &createBffRouterArg)
			types.AssertNil(err)
			util.StdoutOk("Create router successfully\n")
		}

		// create main.go
		{
			mainFp := types.NewBaseFile(bffBaseDir.Append("main.go").(types.DirNode).Abs())
			types.AssertNil(mainFp.Create(types.DefFlag, types.DefPerm))
			txt, err := templates.ParseTemplate(templates.BffMainTemplate, templates.BffMainTplArg{
				HookImport:   fmt.Sprintf("%s/%s%s/bff/%s/hook", sysdes.GoHostPreffix, arg.SysName, sysdes.SystemNameSubffix, arg.Name),
				RouterImport: fmt.Sprintf("%s/%s%s/bff/%s/router", sysdes.GoHostPreffix, arg.SysName, sysdes.SystemNameSubffix, arg.Name),
			})
			types.AssertNil(err)
			_, err = mainFp.WriteString(txt)
			types.AssertNil(err)
			types.AssertNil(mainFp.Close())
		}
	}

	return NoReply, nil
}

func (s *GenerateServer) CreateBffImpl(ctx context.Context, arg *GenServerCreateBffImplArg) (*EmptyReply, error) {
	util.StdoutOk("Start create bff impl: %s\n", arg.Name)
	implBaseDir := arg.Out
	if err := implBaseDir.Create(os.ModePerm); err != nil {
		return NoReply, err
	}
	// create file
	implFp := types.NewBaseFile(implBaseDir.Append(fmt.Sprintf("%s.go", arg.Name)).(types.DirNode).Abs())
	types.AssertNil(implFp.Create(types.DefFlag, types.DefPerm))
	tplArg := templates.BffImplTplArg{
		PackageName:   implBaseDir.Name(),
		InterfaceName: strcase.ToCamel(arg.Name),
		TagZh:         arg.TagZh,
		TagLogin:      arg.TagLogin,
		TagPage:       arg.TagPage,
		TagURI:        arg.TagURI,
	}
	txt, err := templates.ParseTemplate(templates.BffImplTemplate, tplArg)
	types.AssertNil(err)
	_, err = implFp.WriteString(txt)
	types.AssertNil(err)
	types.AssertNil(implFp.Close())

	return NoReply, nil
}

func (s *GenerateServer) CreateBffRouter(_ context.Context, arg *GenServerCreateBffRouterArg) (*EmptyReply, error) {
	if err := arg.Out.Create(os.ModePerm); err != nil {
		return NoReply, err
	}

	var tplArg templates.BffRouterImplTplArg
	tplArg.CamelInterfaceName = strcase.ToCamel(arg.ImplName)
	tplArg.InterfaceName = arg.ImplName
	tplArg.GroupRouter = fmt.Sprintf("/%s", PacketRouter(arg.BffName, arg.ImplName, arg.SysName))
	tplArg.ImplsImport = fmt.Sprintf("git.gmtshenzhen.com/yaodao/%s%s/bff/%s/impls", arg.SysName, sysdes.SystemNameSubffix, arg.BffName)
	if arg.SubffixOpt != "" {
		tplArg.ImplsImport += arg.SubffixOpt
	}
	txt, err := templates.ParseTemplate(templates.BffRouterTemplate, tplArg)
	types.AssertNil(err)
	routerFp := types.NewBaseFile(arg.Out.Append("router.go").(types.DirNode).Abs())
	err = routerFp.Create(types.DefFlag, types.DefPerm)
	types.AssertNil(err)
	defer func() {
		types.AssertNil(routerFp.Close())
	}()
	_, err = routerFp.WriteString(txt)
	types.AssertNil(err)
	return NoReply, nil
}
