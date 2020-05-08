package templates

const (
	interfaceMarkdownTxt = `### {{.Zh}}
#### 请求方法

> POST

#### 是否需要登录

> {{.Login}}

#### 请求路径

> {{.URI}}

#### 请求格式
{{range $i,$c := .RequestTable}}
**{{$c.Title}}**
{{$fieldsSize := len $c.Tables}}
{{ if (gt $fieldsSize 0)}}
| 参数名称 |类型| 参数含义 |必填|默认值|备注|
| ------ | ------ |------ |------ |------ |------ |
{{range $i,$e := $c.Tables}}| {{$e.FieldName}} | {{$e.FieldType}} | {{$e.FieldDesc}} | {{$e.FieldIgnore}} | {{$e.DefValue}} | {{$e.FieldRemark}} |
{{end}}{{end}}
{{end}}
***Example***:
{{.RequestJson}}

***TypeScript Object***
{{.RequestTypeScript}}

***Java Object***
{{.RequestJava}}

***Swift Object***
{{.RequestSwift}}

***Mock***
{{.RequestCurl}}

#### 返回格式
{{range $i, $c := .ResponseTable}}
**{{$c.Title}}**
{{$fileSize := len $c.Tables}}
{{ if (gt $fileSize 0)}}
| 参数名称 |类型| 参数含义 |备注|
| ------ | ------ |------ |------ |
{{range $i,$e := $c.Tables}}| {{$e.FieldName}} | {{$e.FieldType}} | {{$e.FieldDesc}} | {{$e.FieldRemark}} |
{{end}}{{end}}
{{end}}
***Example***:
{{.ResponseJson}}

***TypeScript Object***
{{.ResponseTypeScript}}

***Java Object***
{{.ResponseJava}}

***Swift Object***
{{.ResponseSwift}}

### 错误码
{{$fileSize := len .ErrCodeTable}}
{{ if (gt $fileSize 0)}}
| 参数名称 |值| 描述 |备注|
| ------ | ------ |------ |------ |
{{range $i, $e := .ErrCodeTable}} | {{$e.Name}} | {{$e.Value}} | {{$e.Desc}} | {{$e.Remarks}} |
{{end}}
{{end}}
`

	apiListMarkdownTxt = `# {{.Title}}
{{range $i, $e := .List}}
### {{$e.Title}}
{{range $i, $sube := $e.List}}
- [{{$sube.Remarks}}]({{$sube.Link}})
{{end}}
{{end}}`

	intellijAutomatedTestTxt = `// {{.Title}}
POST http://{{"{{host}}"}}{{.URL}}?_uid={{"{{uid}}"}}&_token={{"{{token}}"}}&_storeId={{"{{storeId}}"}}
Content-Type: Application/json

{{.Data}}

###
`
)
