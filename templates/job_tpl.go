package templates

import "html/template"

type JobMainTplArg struct {
	ImportPackage string
}

type JobCmdSubTplArg struct {
	ImportPacket string
	CobraUse     string
	CmdName      string
	Schedule     string
}

var (
	JobMakefileTemplate = template.Must(template.New("").Parse(makeFileTxt))
	JobMainTpl          = template.Must(template.New("").Parse(jobMainTxt))
	JobScheduleTpl      = template.Must(template.New("").Parse(jobScheduleTxt))
	JobCmdRootTpl       = template.Must(template.New("").Parse(jobCmdTxt))
	JobCmdSubTpl        = template.Must(template.New("").Parse(jobCmdSubTxt))
)
