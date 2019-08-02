package templates

import (
	"html/template"
)

type MarkdownProperty struct {
	Login              string
	Page               []string
	Zh                 string
	URI                string
	RequestTable       []MarkdownTable
	RequestJson        interface{} // 请求参数示例
	RequestTypeScript  interface{} // 请求参数-typescript 对象映射
	ResponseTable      []MarkdownTable
	ResposneJson       interface{} // 应答参数示例
	ResposneTypeScript interface{} // 应答参数-typescript 对象映射
	ErrCodeTable       []MarkdownErrCodeTable
}

type MarkdownTable struct {
	Title  string
	Tables []MarkdownReqRespTable
}

type MarkdownReqRespTable struct {
	FieldName   string      // 参数名称
	FieldType   string      // 类型
	FieldDesc   string      // 参数含义
	FieldIgnore string      // 必填
	DefValue    interface{} // 默认值
	FieldRemark interface{} // 备注
	FieldValue  interface{} // 值（根据不同的类型生成的值不一样）
}

type MarkdownErrCodeTable struct {
	Name    string // 名称
	Value   int    // 值
	Desc    string // 描述
	Remarks string // 备注
}

var InterfaceMarkdownTemplate = template.Must(template.New("").Parse(interfaceMarkdownTxt))

type ApiListProperty struct {
	Title string
	List  []ApiListTable
}

type ApiListTable struct {
	Title string
	List  []*ApiTable
}

type ApiTable struct {
	Link    string
	Remarks string
	Emoji   string
}

var ApiListMarkdownTemplate = template.Must(template.New("").Parse(apiListMarkdownTxt))

type IntellijAutomatedProperty struct {
	Title string
	URL   string
	Data  string
}

var IntellijAutomatedTemplate = template.Must(template.New("").Parse(intellijAutomatedTestTxt))
