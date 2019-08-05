package services

import (
	"fmt"
	"github.com/geekymedic/neon-cli/templates"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"math/rand"
	"reflect"
	"strings"

	"github.com/geekymedic/neon-cli/types/xast"

	uuid "github.com/satori/go.uuid"
)

// 解析扩展节点，参数注解在ast扩展节点上
func parseExtraNode(extraNodes map[string]*xast.ExtraNode, isReq bool) (apiField []templates.MarkdownReqRespTable) {
	if extraNodes == nil {
		return nil
	}
	for varName, extraNode := range extraNodes {
		var item = templates.MarkdownReqRespTable{}
		meta := extraNode.Meta.(*xast.AstMeta)
		item.FieldName = varName
		item.FieldType = meta.FullName
		if meta.Comment != nil {
			commentSplit := strings.SplitN(strings.TrimSpace(meta.Comment.Text()), "|", 4)
			commentSplit = append(commentSplit, "", "", "", "")
			item.FieldDesc = commentSplit[0]
			item.FieldIgnore = commentSplit[1]
			item.FieldValue = commentSplit[2]
			item.FieldRemark = commentSplit[3]
			apiField = append(apiField, item)
		} else {
			commentSplit := strings.SplitN(strings.TrimSpace(meta.Comment.Text()), "|", 2)
			commentSplit = append(commentSplit, "", "")
			item.FieldDesc = commentSplit[0]
			item.FieldRemark = commentSplit[1]
			apiField = append(apiField, item)
		}
	}
	return
}

// 解析叶子节点，参数注解在ast叶子节点上
func parseLeafNode(leafNodes map[string]*xast.LeafNode, isReq bool) (apiField []templates.MarkdownReqRespTable) {
	if leafNodes == nil {
		return nil
	}
	for _, leafNode := range leafNodes {
		var item = templates.MarkdownReqRespTable{}
		meta := leafNode.Meta.(*xast.AstMeta)
		item.FieldName = meta.VarName
		item.FieldType = meta.FullName
		if meta.Comment != nil {
			if isReq {
				commentSplit := strings.SplitN(strings.TrimSpace(meta.Comment.Text()), "|", 4)
				commentSplit = append(commentSplit, "", "", "", "")
				item.FieldDesc = commentSplit[0]
				item.FieldIgnore = commentSplit[1]
				item.FieldValue = commentSplit[2]
				item.FieldRemark = commentSplit[3]
				apiField = append(apiField, item)
			} else {
				commentSplit := strings.SplitN(strings.TrimSpace(meta.Comment.Text()), "|", 2)
				commentSplit = append(commentSplit, "", "")
				item.FieldDesc = commentSplit[0]
				item.FieldRemark = commentSplit[1]
				apiField = append(apiField, item)
			}
		}
	}
	return
}

func InjectAstTree(topNode *xast.TopNode) interface{} {
	if topNode == nil {
		return nil
	}

	var out = make(map[string]interface{})
	for varName, leaf := range topNode.LeavesNodes {
		out[varName] = defValue(leaf.TypeName)
	}
	for varName, extraNode := range topNode.ExtraNodes {
		if extraNode.TypeName == reflect.Array.String() || extraNode.TypeName == reflect.Slice.String() {
			out[varName] = []interface{}{InjectExtraNode(extraNode)}
		} else {
			out[varName] = InjectExtraNode(extraNode)
		}
	}
	return out
}

func InjectExtraNode(extraNode *xast.ExtraNode) interface{} {
	var out = make(map[string]interface{})
	for varName, leaf := range extraNode.LeavesNodes {
		out[varName] = defValue(leaf.TypeName)
	}

	for varName, extraNode := range extraNode.ExtraNodes {
		if extraNode.TypeName == reflect.Array.String() || extraNode.TypeName == reflect.Slice.String() {
			out[varName] = []interface{}{InjectExtraNode(extraNode)}
		} else {
			out[varName] = InjectExtraNode(extraNode)
		}
	}

	// map[int]string, []int
	// TODO add more level parse,eg [][][][]int, map[int]map[int][]int
	if len(extraNode.LeavesNodes)+len(extraNode.ExtraNodes) == 0 {
		switch extraNode.TypeName {
		case reflect.Slice.String(), reflect.Array.String():
			fullName := extraNode.Meta.(*xast.AstMeta).FullName
			idx := strings.Index(fullName, "]")
			typ := fullName[idx+1:]
			return defValue(typ)
		case reflect.Map.String():
			fullName := extraNode.Meta.(*xast.AstMeta).FullName
			kIdx, vIdx := strings.Index(fullName, "["), strings.Index(fullName, "]")
			out[fmt.Sprintf("%v", defValue(fullName[kIdx+1:vIdx]))] = defValue(fullName[vIdx+1:])
			return out
		}
	}
	return out
}

func FixURI(domain, bffName, sysName, implName, uri string) string {
	if uri == "" {
		uri = fmt.Sprintf("https://%s/api/%s", domain, PacketRouter(bffName, implName, sysName))
	} else if !strings.HasPrefix(uri, "http") && !strings.HasPrefix(uri, "https") {
		uri = fmt.Sprintf("https://%s%s", domain, uri)
	}
	return uri
}

func PacketRouter(bffName, implName, sysName string) string {
	idx := strings.LastIndex(sysName, sysdes.SystemNameSubffix)
	if idx > 0 {
		sysName = sysName[:idx]
	}
	return fmt.Sprintf("%s/v1/%s/%s", bffName, sysName, implName)
}

func defValue(typ string) interface{} {
	switch typ {
	case reflect.Int.String(), reflect.Int8.String(), reflect.Int16.String(), reflect.Int32.String(), reflect.Int64.String(),
		reflect.Uint.String(), reflect.Uint8.String(), reflect.Uint16.String(), reflect.Uint32.String(), reflect.Uint64.String():
		return rand.Intn(1<<8 - 1)
	case reflect.Float32.String(), reflect.Float64.String():
		return rand.Float32()
	case reflect.Bool.String():
		return (rand.Intn(1<<8-1) % 2) == 0
	case reflect.String.String():
		return uuid.Must(uuid.NewV4(), nil).String()
	default:
		return "Not Supper Type"
	}
}
