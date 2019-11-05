package templates

import (
	"os"
	"testing"

	"github.com/geekymedic/neon-cli/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseTemplate(t *testing.T) {
	var arg = MarkdownProperty{
		Login: "Y",
		Page:  []string{"a", "b"},
		Zh:    "登录",
		URI:   "https://geekymedic.com/api/amin/v1/boss/login",
		RequestTable: []MarkdownTable{
			{
				Title: "Session",
				Tables: []MarkdownReqRespTable{{
					FieldName:   "名称",
					FieldType:   "类型",
					FieldDesc:   "描述",
					FieldIgnore: "Y",
					DefValue:    "默认值",
					FieldRemark: "备注",
					//FieldValue:  interface{} // 值（根据不同的类型生成的值不一样）
				}},
			},
		},

		RequestJson: `{"token": "xxxxxxx"}`,
		RequestTypeScript: `
			export interface ChangeDrugPriceRequest {
				List: Array<ChangeDrugPriceRequestItem>;
			}
			interface ChangeDrugPriceRequestItem{
				Id: number;
				RetailPrice: number;
				MembershipPrice: number;
			}
		`,

		ResponseTable: []MarkdownTable{
			{
				Title: "Session",
				Tables: []MarkdownReqRespTable{{
					FieldName:   "名称",
					FieldType:   "类型",
					FieldDesc:   "描述",
					FieldIgnore: "Y",
					DefValue:    "默认值",
					FieldRemark: "备注",
					//FieldValue:  interface{} // 值（根据不同的类型生成的值不一样）
				}},
			},
		},
		ResposneJson: `{"token": "xxxxxxx"}`,
		ResposneTypeScript: `
			export interface ChangeDrugPriceRequest {
				List: Array<ChangeDrugPriceRequestItem>;
			}
			class ChangeDrugPriceRequestItem{
				Id: number;
				RetailPrice: number;
				MembershipPrice: number;
			}`,
		ErrCodeTable: []MarkdownErrCodeTable{

		},
	}

	Convey("", t, func() {
		s, err := ParseTemplate(InterfaceMarkdownTemplate, arg)
		So(err, ShouldBeNil)
		fp := types.NewBaseFile(os.TempDir() + "js.md")
		err = fp.Create(os.O_CREATE|os.O_WRONLY, types.DefPerm)
		So(err, ShouldBeNil)
		defer fp.Close()
		_, err = fp.WriteString(s)
		So(err, ShouldBeNil)
		t.Log(fp.Abs())
	})
}
