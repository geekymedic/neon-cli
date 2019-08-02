package templates

import "html/template"

type ServiceServerTplArg struct {
	ImportPackage   string
	AliasName       string
	SystemShortName string // eg: demo
	ServerName      string
}

type ServiceMainTplArg struct {
	List        []string
	SystemName  string
	ServiceName string
}

type ServiceImplArg struct {
	ImportPackage string
	AliasName     string
	ServerName    string
	ImplName      string
	RequestName   string
	ResponseName  string
}


type ServiceHookTplArg struct {
	Alias         string
	ImportPackage string
}

var (
	ServiceConfigYmlTpl  = template.Must(template.New("").Parse(serviceConfigYmlTxt))
	ServiceConfigTpl     = template.Must(template.New("").Parse(serviceConfigTxt))
	ServiceHookTpl       = template.Must(template.New("").Parse(serviceHookTxt))
	ServiceServerTpl     = template.Must(template.New("").Parse(serviceServerTxt))
	ServiceServerInitTpl = template.Must(template.New("").Parse(serviceServerInitTxt))
	ServiceMainTpl       = template.Must(template.New("").Parse(serviceMainTxt))
	ServiceImplTpl       = template.Must(template.New("").Parse(serviceImplTxt))
)
