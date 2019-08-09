package templates

import "html/template"

type ORMTplArg struct {
	ShortTable    string
	FullTable     string
	TableName     string
	SnakeProperty []ORMProperty
	CamelProperty []ORMProperty
}

type ORMProperty struct {
	Name string
	Type string
}

var ORMTemplate = template.Must(template.New("").Parse(ormTxt))