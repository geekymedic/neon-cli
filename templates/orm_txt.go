package templates

const ormTxt = `
{{$table := .}}
func({{.ShortTable}} *{{.FullTable}}) Table() string {
	return {{.TableName}}
}

func({{.ShortTable}} *{{.FullTable}}) Columns() []string {
	var columns = []string {}
	{{range $i,$c := .SnakeProperty}}columns = append(columns, {{$table.ShortTable}}.{{$c.Name}})
	{{end}}
	return columns
}

func({{.ShortTable}} *{{.FullTable}}) Scan(rows *sql.Rows) error {
	var dest []interface{}
	{{range $i,$c := .SnakeProperty}}dest = append(dest, &{{$table.ShortTable}}.{{$c.Name}})
	{{end}}
	return rows.Scan(dest...)
}
`