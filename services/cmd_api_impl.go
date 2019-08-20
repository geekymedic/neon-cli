package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/geekymedic/neon-cli/templates"
	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"github.com/geekymedic/neon-cli/types/xast"
	"github.com/geekymedic/neon-cli/types/xast/astutil"

	"github.com/geekymedic/neon/errors"
	"github.com/geekymedic/neon/logger"
	"github.com/tidwall/pretty"
)

func (s *GenerateServer) GenerateApiDoc(_ context.Context, arg *GenServerApiDocArg) (*EmptyReply, error) {
	var (
		tables = templates.ApiListProperty{Title: "Api List"}
	)
	//var bffs sysdes.Bffs
	{
		s.sys.Bffs.BffIter(func(item *sysdes.BffItem) bool {
			var apiTable templates.ApiListTable
			apiTable.Title = item.DirNode.Name()
			item.ImplIter(func(impl *sysdes.BffImpl) bool {
				if _, err := s.convertInterfaceApiMarkdown(item.DirNode.Name(), impl, arg.Domain, arg.Out); err != nil {
					logger.With("file", impl.FileNode.Abs(), "err", err).Error("Fail to convert markdown")
				} else {
					apiTable.List = append(apiTable.List, &templates.ApiTable{
						Link:    fmt.Sprintf("%s%s%s.md", item.DirNode.Name(), types.Separator, impl.FileNode.Name()),
						Remarks: impl.AstTree.Interface.(*sysdes.BffInterfaceTree).Annotation.Zh + ":" + impl.FileNode.Name()})
					logger.With("path", impl.FileNode.Abs()).Info("Convert markdown successfully")
				}
				return true
			})
			if len(apiTable.List) > 0 {
				tables.List = append(tables.List, apiTable)
			}
			return true
		})
	}

	// Api List.md
	{
		txt, err := templates.ParseTemplate(templates.ApiListMarkdownTemplate, tables)
		if err != nil {
			return NoReply, err
		}
		dir := types.NewBaseFile(fmt.Sprintf("%s%s%s", arg.Out.Abs(), types.Separator, "api_list.md"))
		err = dir.Create(os.O_CREATE|os.O_WRONLY|os.O_TRUNC, types.DefPerm)
		if err != nil {
			return NoReply, err
		}
		defer func() {
			types.AssertNil(dir.Close())
		}()

		_, err = dir.WriteString(txt)
		if err != nil {
			return NoReply, err
		}
	}
	return NoReply, nil
}

func (s *GenerateServer) convertInterfaceApiMarkdown(bffName string, impl *sysdes.BffImpl, domain string, outDir types.DirNode) (templates.MarkdownProperty, error) {
	var (
		errCodeNode      = types.NewBaseDir(impl.Sys.DirNode.Abs()).Append("bff", bffName, "codes").(types.DirNode)
		errCodePath      = types.NewBaseFile(errCodeNode.Abs() + types.Separator + "error_code.go")
		sysName          = impl.Sys.Name
		bffInterfaceTree = impl.AstTree.Interface.(*sysdes.BffInterfaceTree)
		bffReqTree       = impl.AstTree.Request.(*sysdes.BffRequestTree)
		bffReplyTree     = impl.AstTree.Response.(*sysdes.BffResponseTree)
		bffApi           = templates.MarkdownProperty{
			Login: bffInterfaceTree.Annotation.Login,
			Page:  bffInterfaceTree.Annotation.Page,
			Zh:    bffInterfaceTree.Annotation.Zh,
			URI:   bffInterfaceTree.Annotation.URI,
		}
	)

	// request table
	{
		if table := assemblyRequestOrResponseTable(bffReqTree.TopNode, true); len(table) > 0 {
			bffApi.RequestTable = table
			var data bytes.Buffer
			if err := json.NewEncoder(&data).Encode(InjectAstTree(bffReqTree.TopNode)); err != nil {
				return bffApi, err
			}
			bffApi.RequestJson = fmt.Sprintf("```json\n%s```\n", pretty.Pretty(data.Bytes()))

			ts, err := convertTS(bffReqTree.TopNode)
			if err != nil {
				return bffApi, errors.Wrap(err)
			}
			bffApi.RequestTypeScript = fmt.Sprintf("```typescript\n%s```\n", ts)
		}
	}

	// response table
	{
		if table := assemblyRequestOrResponseTable(bffReplyTree.TopNode, false); len(table) > 0 {
			bffApi.ResponseTable = table
		}

		var data bytes.Buffer
		var ret = map[string]interface{}{
			"Code":    0,
			"Message": "请求成功",
			"Data":    InjectAstTree(bffReplyTree.TopNode),
		}
		if err := json.NewEncoder(&data).Encode(ret); err != nil {
			return bffApi, err
		}
		bffApi.ResposneJson = fmt.Sprintf("```json\n%s```\n", pretty.Pretty(data.Bytes()))

		ts, err := convertTS(bffReplyTree.TopNode)
		if err != nil {
			return bffApi, errors.Wrap(err)
		}
		bffApi.ResposneTypeScript = fmt.Sprintf("```typescript\n%s```\n", ts)
	}

	// errcode
	{
		if tables, err := parseErrCode(impl, errCodePath); err == nil {
			bffApi.ErrCodeTable = tables
		} else {
			logger.Error("Fail parse errcode: %v", err)
		}
	}

	// fix uri
	bffApi.URI = FixURI(domain, bffName, sysName, impl.FileNode.Name(), bffApi.URI)

	txt, err := templates.ParseTemplate(templates.InterfaceMarkdownTemplate, bffApi)
	if err != nil {
		return bffApi, err
	}

	mkFp := types.NewBaseFile(outDir.Append(bffName, impl.FileNode.Name()+".md").(types.DirNode).Abs())
	mkFp.MustCreate(os.O_CREATE|os.O_WRONLY|os.O_TRUNC, types.DefPerm)
	if _, err = mkFp.WriteString(txt); err != nil {
		return bffApi, errors.Wrap(err)
	}
	if err = mkFp.Sync(); err != nil {
		return bffApi, errors.Wrap(err)
	}
	return bffApi, mkFp.Close()
}

func assemblyRequestOrResponseTable(topNode *xast.TopNode, isReq bool) []templates.MarkdownTable {
	var tables []templates.MarkdownTable
	topNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
		var requestTable []templates.MarkdownReqRespTable
		var title string
		switch typ := node.(type) {
		case *xast.TopNode:
			title = typ.TypeName
			requestTable = parseExtraNode(typ.ExtraNodes, isReq)
			requestTable = append(requestTable, parseLeafNode(typ.LeavesNodes, isReq)...)
		case *xast.ExtraNode:
			title = astutil.SimpleName(typ.Meta.(*xast.AstMeta).RawExpr)
			requestTable = parseExtraNode(typ.ExtraNodes, isReq)
			requestTable = append(requestTable, parseLeafNode(typ.LeavesNodes, isReq)...)
		case *xast.LeafNode:
		default:
			types.PanicSanityf("unsupport type %T", typ)
		}
		if len(requestTable) > 0 {
			tables = append(tables, templates.MarkdownTable{Title: title, Tables: requestTable})
		}
		return true
	})

	return tables
}

func parseErrCode(impl *sysdes.BffImpl, fileNode types.FileNode) ([]templates.MarkdownErrCodeTable, error) {
	curBffErrCode, err := LoadErrCode(fileNode.Abs())
	if err != nil {
		return nil, err
	}
	useErrCode, err := loadUseErrorCode(impl)
	if err != nil {
		return nil, err
	}

	var tables []templates.MarkdownErrCodeTable
	for _, errCode := range useErrCode {
		find := false
		for _, cur := range curBffErrCode {
			if cur.VarName == errCode {
				tables = append(tables, templates.MarkdownErrCodeTable{
					// errcode
					Name:    cur.VarName,
					Value:   cur.Value,
					Desc:    cur.Des,
					Remarks: cur.Remarks,
				})
				find = true
				break
			}
		}
		if !find {
			types.CodeIter(func(key types.CodeType, value string) bool {
				if key.String() == errCode {
					tables = append(tables, templates.MarkdownErrCodeTable{
						Name:  key.String(),
						Value: int(key),
						Desc:  value, // TODO add remark
					})
					return false
				}
				return true
			})
		}
	}

	return tables, nil
}
