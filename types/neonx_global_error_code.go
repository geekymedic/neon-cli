package types

import "fmt"

//go:generate stringer -type=CodeType
type CodeType int

const (
	CodeSuccess                 CodeType = 0
	CodeVersionError            CodeType = 1
	CodeUpdating                CodeType = 101
	CodeNotFound                CodeType = 1000
	CodeRequestUrlParamError    CodeType = 1001
	CodeRequestQueryParamError  CodeType = 1002
	CodeRequestCommonParamError CodeType = 1003
	CodeRequestBodyError        CodeType = 1004
	CodeNotAllow                CodeType = 1005
	CodeServerError             CodeType = 1006
)

type Codes map[CodeType]string

var (
	_codes = Codes{
		CodeSuccess:                 "请求成功",
		CodeVersionError:            "客户端版本错误，请升级客户端",
		CodeUpdating:                "服务正在升级",
		CodeNotFound:                "找不到对于系统&模块",
		CodeRequestUrlParamError:    "请求的URL参数错误",
		CodeRequestQueryParamError:  "请求的查询参数错误",
		CodeRequestCommonParamError: "请求的查询参数错误",
		CodeRequestBodyError:        "请求的请求结构错误",
		CodeNotAllow:                "权限校验失败",
		CodeServerError:             "服务器错误",
	}
)

func GetMessage(code CodeType) string {
	return _codes[code]
}

func CodeIter(fn func(key CodeType, value string) bool) {
	for key, value := range _codes {
		if !fn(key, value) {
			return
		}
	}
}

//请在初始化阶段使用该函数，运行时使用该函数可能导致同步问题
func MergeCodes(codes Codes) {

	for code, describe := range codes {

		_, exists := _codes[code]

		if exists {
			panic(
				fmt.Sprintf("code %d[%s] already exists", code, GetMessage(code)),
			)
		}

		_codes[code] = describe

	}
}
