package sysdes

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/xast"
	"github.com/geekymedic/neon-cli/types/xast/astutil"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestBffInterfaceTree_Parse(t *testing.T) {
	Convey("", t, func() {
		var bffTree = BffInterfaceTree{
			FileNode: types.NewBaseFile(bffFile(t)),
		}
		err := bffTree.Parse()
		So(err, ShouldBeNil)
		So(bffTree.FuncName, ShouldEqual, "CreateHardwareHandler")
		So(bffTree.Annotation.Zh, ShouldEqual, "新增硬件")
		So(bffTree.Annotation.URI, ShouldEqual, "/api/boss/v1/admin/login")
		So(bffTree.Annotation.Typ, ShouldEqual, "b.i")
		assert.Equal(t, []string{"hardware", "software"}, bffTree.Annotation.Page)
	})

}

func TestBffRequestTree_Parse(t *testing.T) {
	var bffTree = BffRequestTree{
		FileNode: types.NewBaseFile(bffFile(t)),
	}
	err := bffTree.Parse()
	assert.Nil(t, err)
	// assert.NotNil(t, bffTree.TopNode)
	// assert.Equal(t, "CreateHardwareRequest", bffTree.Name)
	// assert.Equal(t, "b.i.rt", bffTree.Annotation.Typ)
	// assert.Equal(t, "CreateHardwareHandler", bffTree.Annotation.Interface)
}

func TestBffResponseTree_Parse(t *testing.T) {
	Convey("", t, func() {
		var bffTree = BffResponseTree{
			FileNode: types.NewBaseFile(bffFile(t)),
		}
		err := bffTree.Parse()
		So(err, ShouldBeNil)
		So(bffTree.TopNode, ShouldNotBeNil)
		So(bffTree.Name, ShouldEqual, "CreateHardwareResponse")
		So(bffTree.Annotation.Typ, ShouldEqual, "b.i.re")
		So(bffTree.Annotation.Interface, ShouldEqual, "CreateHardwareHandler")
	})
}

func TestNewBffTree(t *testing.T) {
	t.Run("", func(t *testing.T) {
		Convey("", t, func() {
			bffTree, err := NewBffTree(types.NewBaseFile(bffFile(t)))
			So(err, ShouldBeNil)
			So(bffTree, ShouldNotBeNil)
			So(bffTree.Interface, ShouldNotBeNil)
			So(bffTree.Request, ShouldNotBeNil)
			So(bffTree.Response, ShouldNotBeNil)

			b := false
			bffTree.Request.(*BffRequestTree).TopNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
				if node, ok := node.(*xast.ExtraNode); ok {
					if walkPath == "CreateHardwareRequest.User" {
						b = node.Meta.(*xast.AstMeta).CrossModule
					}
				}
				return true
			})
			assert.False(t, b)
		})
	})

	t.Run("cross file", func(t *testing.T) {
		Convey("cross file", t, func() {
			tmp := fmt.Sprintf("%s", os.TempDir())
			So(os.MkdirAll(tmp, os.ModePerm), ShouldBeNil)

			var fps = []types.FileNode{
				types.NewBaseFile(bffCrossFile(t, tmp)),
				types.NewBaseFile(bffCrossFile2(t, tmp)),
				types.NewBaseFile(bffCrossFile3(t, tmp)),
			}
			bffTree, err := NewBffTree(fps[0])
			So(err, ShouldBeNil)
			So(bffTree, ShouldNotBeNil)
			So(bffTree.Interface, ShouldNotBeNil)
			So(bffTree.Request, ShouldNotBeNil)
			So(bffTree.Response, ShouldNotBeNil)
			{
				bffInterface := bffTree.Interface.(*BffInterfaceTree)
				So(bffInterface.FuncName, ShouldEqual, "CreateHardwareHandler")
				So(bffInterface.Annotation.Zh, ShouldEqual, "新增硬件")
				So(bffInterface.Annotation.URI, ShouldEqual, "/api/boss/v1/admin/login")
				So(bffInterface.Annotation.Typ, ShouldEqual, "b.i")
				So(strings.Join([]string{"hardware", "software"}, ","),
					ShouldEqual,
					strings.Join(bffInterface.Annotation.Page, ","))
			}
		})
	})
}

func TestNewRequest(t *testing.T) {
	Convey("cross file", t, func() {
		tmp := fmt.Sprintf("%s", os.TempDir())
		assert.Nil(t, os.MkdirAll(tmp, os.ModePerm))
		var fps = []types.FileNode{
			types.NewBaseFile(bffCrossFile(t, tmp)),
			types.NewBaseFile(bffCrossFile2(t, tmp)),
			types.NewBaseFile(bffCrossFile3(t, tmp)),
		}
		var requestTree = &BffRequestTree{FileNode: fps[0]}
		So(requestTree.Parse(), ShouldBeNil)
		So(requestTree.FillCrossStructs(), ShouldBeNil)

		expect := []string{
			"CreateHardwareRequest",
			"CreateHardwareRequest.User",
			"CreateHardwareRequest.User.Age",
			"CreateHardwareRequest.Bo",
			"CreateHardwareRequest.Bo.Info",
			"CreateHardwareRequest.Third",
			"CreateHardwareRequest.Third.I",
			"CreateHardwareRequest.Sn",
			"CreateHardwareRequest.TypeId",
		}
		requestTree.TopNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
			So(expect, ShouldContain, walkPath)
			return true
		})
	})
}

func TestCrossFile(t *testing.T) {
	t.Parallel()
	topNode, err := astutil.BuildStructTree("CreateHardwareRequest", bffCrossFile(t), nil)
	Convey("", t, func() {
		So(err, ShouldBeNil)
		So(topNode, ShouldNotBeNil)
	})
}

func TestExtendTree_ReplaceExtraNode(t *testing.T) {
	Convey("", t, func() {
		topNode1, err := astutil.BuildStructTree("CreateHardwareRequest", bffCrossFile(t), nil)
		So(err, ShouldBeNil)
		So(topNode1, ShouldNotBeNil)

		topNode2, err := astutil.BuildStructTree("Person", bffCrossFile2(t), nil)
		So(err, ShouldBeNil)
		So(topNode2, ShouldNotBeNil)

		topNode3, err := astutil.BuildStructTree("Book", bffCrossFile2(t), nil)
		So(err, ShouldBeNil)
		So(topNode3, ShouldNotBeNil)

		topNode4, err := astutil.BuildStructTree("Three", bffCrossFile3(t), nil)
		So(err, ShouldBeNil)
		So(topNode4, ShouldNotBeNil)

		extendNode := NewExtendTree(topNode1)
		count := extendNode.ReplaceExtraNode(topNode2)
		So(count, ShouldEqual, 1)
		count = extendNode.ReplaceExtraNode(topNode3)
		So(count, ShouldEqual, 1)
		count = extendNode.ReplaceExtraNode(topNode4)
		So(count, ShouldEqual, 1)

		expect := []string{
			"CreateHardwareRequest",
			"CreateHardwareRequest.User",
			"CreateHardwareRequest.User.Age",
			"CreateHardwareRequest.Bo",
			"CreateHardwareRequest.Bo.Info",
			"CreateHardwareRequest.Third",
			"CreateHardwareRequest.Third.I",
			"CreateHardwareRequest.Sn",
			"CreateHardwareRequest.TypeId",
		}

		So(extendNode.NodeCount(), ShouldEqual, len(expect))
		extendNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
			So(expect, ShouldContain, walkPath)
			return true
		})
	})

}

func TestFlatNodes(t *testing.T) {
	Convey("", t, func() {
		var bffTree = BffRequestTree{
			FileNode: types.NewBaseFile(flatNodeFile(t)),
		}
		err := bffTree.Parse()
		assert.Nil(t, err)
		extendNode := NewExtendTree(bffTree.TopNode)
		extendNode.FlatNestedNodes()
		extendNode.BreadthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
			fmt.Println(walkPath)
			return true
		})
	})
}

func bffFile(t *testing.T) string {
	txt := `
package hardware

import (
	store "protocol/store_system/store"
	"store_system/bff/admin/rpc"

	"github.com/geekymedic/neon/bff"
)

type Person struct {
	Age int
}

// @type: b.i.rt
// @interface: CreateHardwareHandler
type CreateHardwareRequest struct {
	Sn     string // 硬件序列号 | Y | "" | 可支持多个添加，以英文逗号(,)隔开
	TypeId string // 硬件所属类型id | Y | "" |
	User Person // 用户 | Y | "" |
}

// @type: b.i.re
// @interface: CreateHardwareHandler
type CreateHardwareResponse struct {
	User Person // 用户 | Y | "" |
}

// @type: b.i
// @name: 新增硬件
// @login: Y
// @uri: /api/boss/v1/admin/login
// @page: hardware|software
func CreateHardwareHandler(state *bff.State) {
	var (
		ask = &CreateHardwareRequest{}
		ack = &CreateHardwareResponse{}
	)
	if err := state.ShouldBindJSON(ask); err != nil {
		state.Error(bff.CodeRequestBodyError, err)
		return
	}
	_, err := rpc.NewStoreHardwareServer().CreateHardware(state.Context(), &store.CreateHardwareRequest{Sn: ask.Sn, TypeId: ask.TypeId})
	if err != nil {
		state.Error(bff.CodeServerError, err)
		return
	}
	state.Success(ack)
}
`
	//	txt := `
	// package admin
	//
	// import (
	//	"github.com/geekymedic/neon/bff"
	// )
	//
	// // @type: bff.interface.request
	// // @interface: AddPortalHandler
	// // @des:
	// type AddPortalRequest struct {
	//	SupplierId  string // 供应商ID | Y | "" |
	//	Url         string // 接入商URL | Y | "" |
	//	Certificate string // 接入商证书 | Y | "PEM格式" |
	// }
	//
	// // @type: bff.interface.response
	// // @interface: AddPortalHandler
	// // @describe:
	// type AddPortalResponse struct {
	//	PortalID string // 供应商接入ID | Y | "" |
	// }
	//
	// // @type: bff.interface
	// // @name: 添加认证供应商接入
	// // @login: Y
	// // @page:
	// // @uri: /api/admin/v1/portals/add_portal
	// // @describe:
	// func AddPortalHandler(state *bff.State) {
	//	var (
	//		ask = &AddPortalRequest{}
	//		ack = &AddPortalResponse{}
	//	)
	//	if err := state.ShouldBindJSON(ask); err != nil {
	//		state.Error(bff.CodeRequestBodyError, err)
	//		return
	//	}
	//
	//	state.Success(ack)
	// }
	//
	// `
	fileNode := types.NewBaseFile(fmt.Sprintf("%s%d", os.TempDir(), time.Now().UnixNano()))
	err := fileNode.Create(os.O_CREATE|os.O_RDWR, types.DefPerm)
	assert.Nil(t, err)
	defer fileNode.Close()
	_, err = fileNode.WriteString(txt)
	assert.Nil(t, err)

	return fileNode.Abs()
}

func bffCrossFile(t *testing.T, abs ...string) string {
	txt := `
package hardware

import (
	store "protocol/store_system/store"
	"store_system/bff/admin/rpc"

	"github.com/geekymedic/neon/bff"
)

// @type: b.i.rt
// @interface: CreateHardwareHandler
type CreateHardwareRequest struct {
	Sn     string // 硬件序列号 | Y | "" | 可支持多个添加，以英文逗号(,)隔开
	TypeId string // 硬件所属类型id | Y | "" |
	User Person // 用户 | Y | "" |
	Bo Book  // 书籍 | Y | "" |
	Third Three // 三 | Y | "" |
}

// @type: bff.interface.response
// @interface: CreateHardwareHandler
type CreateHardwareResponse struct {
	User Person // 用户 | Y | "" |
}

// @type: b.i 
// @name: 新增硬件
// @login: Y
// @uri: /api/boss/v1/admin/login
// @page: hardware|software
func CreateHardwareHandler(state *bff.State) {
	var (
		ask = &CreateHardwareRequest{}
		ack = &CreateHardwareResponse{}
	)
	if err := state.ShouldBindJSON(ask); err != nil {
		state.Error(bff.CodeRequestBodyError, err)
		return
	}
	_, err := rpc.NewStoreHardwareServer().CreateHardware(state.Context(), &store.CreateHardwareRequest{Sn: ask.Sn, TypeId: ask.TypeId})
	if err != nil {
		state.Error(bff.CodeServerError, err)
		return
	}
	state.Success(ack)
}
`
	fileNode := types.NewBaseFile(fmt.Sprintf("%s%d.go", os.TempDir(), time.Now().UnixNano()))
	if len(abs) != 0 {
		fileNode = types.NewBaseFile(fmt.Sprintf("%s/%d.go", abs[0], time.Now().UnixNano()))
	}
	err := fileNode.Create(os.O_CREATE|os.O_RDWR, types.DefPerm)
	assert.Nil(t, err)
	defer fileNode.Close()
	_, err = fileNode.WriteString(txt)
	assert.Nil(t, err)

	return fileNode.Abs()
}

func bffCrossFile2(t *testing.T, abs ...string) string {
	txt := `
package hardware

import (
	store "protocol/store_system/store"
	"store_system/bff/admin/rpc"

	"github.com/geekymedic/neon/bff"
)

type Person struct {
	Age int // 年龄 | Y | "" |
}

type Book struct {
	Info string // 信息 | Y | "" |
}
`
	fileNode := types.NewBaseFile(fmt.Sprintf("%s%d.go", os.TempDir(), time.Now().UnixNano()))
	if len(abs) > 0 {
		fileNode = types.NewBaseFile(fmt.Sprintf("%s/%d.go", abs[0], time.Now().UnixNano()))
	}
	err := fileNode.Create(os.O_CREATE|os.O_RDWR, types.DefPerm)
	assert.Nil(t, err)
	defer fileNode.Close()
	_, err = fileNode.WriteString(txt)
	assert.Nil(t, err)

	return fileNode.Abs()
}

func bffCrossFile3(t *testing.T, abs ...string) string {
	txt := `
package hardware

import (
	store "protocol/store_system/store"
	"store_system/bff/admin/rpc"

	"github.com/geekymedic/neon/bff"
)

type Three struct {
	I int // 标识 | Y | "" |
}
`
	fileNode := types.NewBaseFile(fmt.Sprintf("%s%d.go", os.TempDir(), time.Now().UnixNano()))
	if len(abs) > 0 {
		fileNode = types.NewBaseFile(fmt.Sprintf("%s/%d.go", abs[0], time.Now().UnixNano()))
	}
	err := fileNode.Create(os.O_CREATE|os.O_RDWR, types.DefPerm)
	assert.Nil(t, err)
	defer fileNode.Close()
	_, err = fileNode.WriteString(txt)
	assert.Nil(t, err)

	return fileNode.Abs()
}

func flatNodeFile(t *testing.T, abs ...string) string {
	txt := `
package hardware

import (
	store "protocol/store_system/store"
	"store_system/bff/admin/rpc"

	"github.com/geekymedic/neon/bff"
)

type Person struct {
	Age int
	Loc Local 
}

type Local struct {
	Left int //
	Right int //
}

// @type: b.i.rt
// @interface: CreateHardwareHandler
type CreateHardwareRequest struct {
	Sn     string // 硬件序列号 | Y | "" | 可支持多个添加，以英文逗号(,)隔开
	TypeId string // 硬件所属类型id | Y | "" |
	Person // 用户 | Y | "" |
	P1 Person // 用户 | Y | "" |
}

// @type: b.i.re
// @interface: CreateHardwareHandler
type CreateHardwareResponse struct {
	User Person // 用户 | Y | "" |
}

// @type: b.i
// @name: 新增硬件
// @login: Y
// @uri: /api/boss/v1/admin/login
// @page: hardware|software
func CreateHardwareHandler(state *bff.State) {
	var (
		ask = &CreateHardwareRequest{}
		ack = &CreateHardwareResponse{}
	)
	if err := state.ShouldBindJSON(ask); err != nil {
		state.Error(bff.CodeRequestBodyError, err)
		return
	}
	_, err := rpc.NewStoreHardwareServer().CreateHardware(state.Context(), &store.CreateHardwareRequest{Sn: ask.Sn, TypeId: ask.TypeId})
	if err != nil {
		state.Error(bff.CodeServerError, err)
		return
	}
	state.Success(ack)
}
`
	fileNode := types.NewBaseFile(fmt.Sprintf("%s%d", os.TempDir(), time.Now().UnixNano()))
	err := fileNode.Create(os.O_CREATE|os.O_RDWR, types.DefPerm)
	assert.Nil(t, err)
	defer fileNode.Close()
	_, err = fileNode.WriteString(txt)
	assert.Nil(t, err)

	return fileNode.Abs()
}
