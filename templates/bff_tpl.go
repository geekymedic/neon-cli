package templates

import "html/template"

type BffImplTplArg struct {
	PackageName   string
	RPCPath       string
	InterfaceName string
	TagZh         string
	TagLogin      string
	TagPage       string
	TagURI        string
}

type BffRouterImplTplArg struct {
	ImplsImport        string
	GroupRouter        string
	CamelInterfaceName string
	InterfaceName      string
}

type BffMainTplArg struct {
	HookImport   string
	RouterImport string
}

type BffHookTplArg struct {
	Alias         string
	ImportPackage string
}

var (
	BffMakefileTemplate  = template.Must(template.New("").Parse(makeFileTxt))
	BffErrCodeTemplate   = template.Must(template.New("").Parse(errCodeTxt))
	BffConfigYmlTemplate = template.Must(template.New("").Parse(bffConfigYmlTxt))
	BffConfigTemplate    = template.Must(template.New("").Parse(bffConfigTxt))
	BffImplTemplate      = template.Must(template.New("").Parse(bffImplTxt))
	BffRouterTemplate    = template.Must(template.New("").Parse(bffRouterTxt))
	BffHookTemplate      = template.Must(template.New("").Parse(bffHookTxt))
	BffMainTemplate      = template.Must(template.New("").Parse(bffMainTxt))
)
