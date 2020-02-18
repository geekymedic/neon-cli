package services

import (
	"context"
	"fmt"
	"os"

	"github.com/geekymedic/neon-cli/templates"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"github.com/geekymedic/neon-cli/util"
)

func (s *GenerateServer) CreateJob(ctx context.Context, arg *GenServerJobArg) (*EmptyReply, error) {
	jobBaseDir := fmt.Sprintf("%s%sjob%scronjob%s%s", arg.SysDir.Abs(), types.Separator, types.Separator, types.Separator, arg.Name)
	jobBaseNode := types.NewBaseDir(jobBaseDir)
	types.AssertNil(jobBaseNode.Create(os.ModePerm))
	util.StdoutOk("Create %s successfully\n", jobBaseNode.Abs())

	// create makefile
	{
		makefileFp := types.NewBaseFile(jobBaseNode.Append("Makefile").(types.DirNode).Abs())
		types.AssertNil(makefileFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.JobMakefileTemplate,
			map[string]interface{}{"Name": arg.Name, "Typ": "cronjob", "System": s.sys.ShortName})
		types.AssertNil(err)
		_, err = makefileFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(makefileFp.Close())
		util.StdoutOk("Create job makefile successfully\n")
	}

	// create docker file
	{
		dockerFp := types.NewBaseFile(jobBaseNode.Append("Dockfile").(types.DirNode).Abs())
		types.AssertNil(dockerFp.Create(types.DefFlag, types.DefPerm))
		types.AssertNil(dockerFp.Close())
		util.StdoutOk("Create service dockfile successfully\n")
	}

	// create k8s
	{
		k8sFp := types.NewBaseFile(jobBaseNode.Append(".k8s.yml").(types.DirNode).Abs())
		types.AssertNil(k8sFp.Create(types.DefFlag, types.DefPerm))
		types.AssertNil(k8sFp.Close())
		util.StdoutOk("Create service k8s config successfully\n")
	}

	// create config
	{
		configBaseDir := jobBaseNode.Append("config").(types.DirNode)
		types.AssertNil(configBaseDir.Create(os.ModePerm))
		ymlFp := types.NewBaseFile(configBaseDir.Append("config.yml").(types.DirNode).Abs())
		types.AssertNil(ymlFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.ServiceConfigYmlTpl, map[string]interface{}{
			"Name": fmt.Sprintf("%s%s-cronjob-%s", s.sys.ShortName, sysdes.SystemNameSubffix, arg.Name)})
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
		util.StdoutOk("Create config.go successfully\n")
	}

	// create cmd/root.go
	{
		types.AssertNil(jobBaseNode.Append("cmd").(types.DirNode).Create(os.ModePerm))
		scheduleRootFp := types.NewBaseFile(jobBaseNode.Append("cmd", "root.go").(types.DirNode).Abs())
		scheduleRootFp.MustCreate(types.DefFlag, types.DefPerm)
		txt, err := templates.ParseTemplate(templates.JobCmdRootTpl, map[string]interface{}{"subCmd": arg.CmdName})
		types.AssertNil(err)
		_, err = scheduleRootFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(scheduleRootFp.Close())
		util.StdoutOk("Create schedule root.go successfully\n")
	}

	// create cmd/sub.go
	{
		scheduleSubCmdFp := types.NewBaseFile(jobBaseNode.Append("cmd", arg.CmdName+".go").(types.DirNode).Abs())
		scheduleSubCmdFp.MustCreate(types.DefFlag, types.DefPerm)
		txt, err := templates.ParseTemplate(templates.JobCmdSubTpl, templates.JobCmdSubTplArg{
			ImportPacket: util.ConvertBreakLinePath(fmt.Sprintf("%s/job/cronjob/%s/schedule", s.sys.GoModel, arg.Name)),
			CobraUse:     arg.CmdName,
			CmdName:      arg.CmdName,
			Schedule:     "BaseSchedule",
		})
		types.AssertNil(err)
		_, err = scheduleSubCmdFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(scheduleSubCmdFp.Close())
		util.StdoutOk("Create schedule %s.go successfully\n", arg.CmdName)
	}

	// create schedule.go
	{
		types.AssertNil(jobBaseNode.Append("schedule").(types.DirNode).Create(os.ModePerm))
		scheduleFp := types.NewBaseFile(jobBaseNode.Append("schedule", "schedule.go").(types.DirNode).Abs())
		scheduleFp.MustCreate(types.DefFlag, types.DefPerm)
		txt, err := templates.ParseTemplate(templates.JobScheduleTpl, map[string]interface{}{})
		types.AssertNil(err)
		_, err = scheduleFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(scheduleFp.Close())
		util.StdoutOk("Create schedule schedule.go successfully\n")
	}

	// create main.go
	{
		mainFp := types.NewBaseFile(jobBaseNode.Append("main.go").(types.DirNode).Abs())
		types.AssertNil(mainFp.Create(types.DefFlag, types.DefPerm))
		txt, err := templates.ParseTemplate(templates.JobMainTpl, templates.JobMainTplArg{
			ImportPackage: fmt.Sprintf("%s/job/cronjob/%s/cmd", s.sys.GoModel, arg.Name),
		})
		types.AssertNil(err)
		_, err = mainFp.WriteString(txt)
		types.AssertNil(err)
		types.AssertNil(mainFp.Close())
		util.StdoutOk("Create main.go successfully\n")
	}
	return &EmptyReply{}, nil
}
