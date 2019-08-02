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
)

func (s *GenerateServer) CreateSystem(ctx context.Context, arg *GenServerCreeateSystemArg) (*EmptyReply, error) {
	// create system base directory
	sysBaseDir := arg.Out.Append(arg.Name + sysdes.SystemNameSubffix).(types.DirNode)
	if err := sysBaseDir.Create(os.ModePerm); err != nil {
		return NoReply, errors.Wrap(err)
	}
	util.StdoutOk("Create directory successfully: %s\n", sysBaseDir.Abs())

	// create neon.yml
	{
		fileNode := types.NewBaseFile(fmt.Sprintf("%s%s.neon.yml", sysBaseDir.Abs(), types.Separator))
		err := fileNode.Create(os.O_CREATE|os.O_WRONLY|os.O_EXCL, types.DefPerm)
		types.AssertNil(err)
		types.AssertNil(fileNode.Close())
	}

	// create go.mod
	{
		fileNode := types.NewBaseFile(fmt.Sprintf("%s%sgo.mod", sysBaseDir.Abs(), types.Separator))
		err := fileNode.Create(os.O_CREATE|os.O_WRONLY|os.O_EXCL, types.DefPerm)
		types.AssertNil(err)
		_, err = fileNode.WriteString(fmt.Sprintf("module git.gmtshenzhen.com/yaodao/%s%s\n\n", arg.Name, sysdes.SystemNameSubffix))
		types.AssertNil(err)
		_, err = fileNode.WriteString("go 1.12\n")
		types.AssertNil(err)
		types.AssertNil(fileNode.Close())
	}

	// create readme
	{
		fileNode := types.NewBaseFile(fmt.Sprintf("%s%sREADME.md", sysBaseDir.Abs(), types.Separator))
		err := fileNode.Create(os.O_CREATE|os.O_WRONLY|os.O_EXCL, types.DefPerm)
		types.AssertNil(err)
		_, err = fileNode.WriteString("# " + arg.Name + sysdes.SystemNameSubffix)
		types.AssertNil(err)
	}

	// create demo protocol
	{
		protocolBaseDir := sysBaseDir.Append("protocol", "demo").(types.DirNode)
		types.AssertNil(protocolBaseDir.Create(os.ModePerm))
		dst := types.NewBaseFile(protocolBaseDir.Append("ping.pb.go").(types.DirNode).Abs())
		err := dst.Create(os.O_WRONLY|os.O_CREATE|os.O_EXCL, types.DefPerm)
		types.AssertNil(err)
		_, err = dst.WriteString(templates.PingTxt)
		types.AssertNil(err)
		types.AssertNil(dst.Close())
	}

	// create demo service
	{
		servicesBaseDir := sysBaseDir.Append("services").(types.DirNode)
		types.AssertNil(servicesBaseDir.Create(os.ModePerm))

		_, err := s.CreateService(nil, &GenServerServiceArg{
			Out:     servicesBaseDir,
			Name:    "checkhealth",
			SysName: arg.Name,
			SysDir:  sysBaseDir,
			Impl: &GenServerCreateServiceImplArg{
				Out:         servicesBaseDir.Append("checkhealth", "impls", "ping").(types.DirNode),
				Name:        "ping",
				ServiceName: "CheckHealth",
				SysName:     arg.Name,
				SubffixOpt:  "/ping",
			},
		})
		types.AssertNil(err)
	}

	// create bff
	{
		bffBaseDir := sysBaseDir.Append("bff").(types.DirNode)
		types.AssertNil(bffBaseDir.Create(os.ModePerm))

		var createBffImplArg = &GenServerCreateBffImplArg{
			Out:      bffBaseDir.Append("demo", "impls", "demo").(types.DirNode),
			Name:     "demo",
			BffName:  "demo",
			SysName:  arg.Name,
			TagZh:    "健康检查",
			TagLogin: "N",
			TagPage:  "",
			TagURI:   "api" + "/" + PacketRouter("demo", "ping", arg.Name),
		}
		_, err := s.CreateBff(ctx, &GenServerCreateBffArg{
			Name:    "demo",
			Out:     bffBaseDir,
			SysName: arg.Name,
			SysDir:  sysBaseDir,
			Impl:    createBffImplArg})
		types.AssertNil(err)
	}
	return NoReply, nil
}
