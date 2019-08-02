package sysdes

type MakeFileParam struct {
	SystemName   string
	BffNames     []string
	ServiceNames []string
}

type DepsEnv struct {
	ProjectRoot string
	GoRoot      string
	GoPath      string
	Protoc      string
	ProtocGenGo string
	Git         string
}
