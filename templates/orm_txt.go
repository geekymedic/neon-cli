package templates

import (
	"github.com/iancoleman/strcase"
	"strings"
)

const ormTxt = `
{{$table := .}}
func({{.ShortTable}} *{{.FullTable}}) Table() string {
	return {{.TableName}}
}

func({{.ShortTable}} *{{.FullTable}}) Columns() []string {
	var columns = []string {}
	{{range $i,$c := .SnakeProperty}}columns = append(columns, "{{$c.Name}}")
	{{end}}
	return columns
}

func({{.ShortTable}} *{{.FullTable}}) Scan(rows *sql.Rows) error {
	var dest []interface{}
	{{range $i,$c := .CamelProperty}}dest = append(dest, &{{$table.ShortTable}}.{{$c.Name}})
	{{end}}
	return rows.Scan(dest...)
}

Insert:
INSERT INTO {{.ShortTable}} ({{SqlColumn .SnakeProperty}}) VALUES()
`

func SqlColumn(list []ORMProperty) string {
	var columns []string
	for _, value := range list {
		columns = append(columns, strcase.ToSnake(value.Name))
	}
	return strings.Join(columns, ",\n")
}
