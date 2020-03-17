package services

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/geekymedic/neon-cli/types/xast"
	"github.com/geekymedic/neon-cli/types/xast/astutil"
)

func convertTS(astTree *xast.TopNode) (string, error) {
	var ret []struct {
		Title string
		Items map[string]string
	}

	astTree.BreadthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
		structItem := struct {
			Title string
			Items map[string]string
		}{Title: "", Items: map[string]string{}}
		switch n := node.(type) {
		case *xast.TopNode:
			structItem.Title = n.TypeName
			if strings.Contains(astTree.TypeName, "Response") {
				ret = append(ret, struct {
					Title string
					Items map[string]string
				}{Title: n.TypeName, Items: map[string]string{
					"Code":    parseTypeToTS("int", "int"),
					"Message": parseTypeToTS("string", "string"),
					"Data":    "Body",
				}})

				structItem.Title = "Body"
			}
			for _, node := range n.LeavesNodes {
				meta := node.Meta.(*xast.AstMeta)
				structItem.Items[meta.VarName] = parseTypeToTS(node.TypeName, node.TypeName)
			}
			for _, node := range n.ExtraNodes {
				meta := node.Meta.(*xast.AstMeta)
				structItem.Items[meta.VarName] = parseTypeToTS(meta.SysType, meta.FullName)
			}
			ret = append(ret, structItem)
		case *xast.ExtraNode:
			structItem.Title = astutil.SimpleName(n.Meta.(*xast.AstMeta).RawExpr)
			for _, node := range n.LeavesNodes {
				meta := node.Meta.(*xast.AstMeta)
				structItem.Items[meta.VarName] = parseTypeToTS(node.TypeName, node.TypeName)
			}
			for _, node := range n.ExtraNodes {
				meta := node.Meta.(*xast.AstMeta)
				structItem.Items[meta.VarName] = parseTypeToTS(meta.SysType, meta.FullName)
			}
			if len(structItem.Items) > 0 {
				ret = append(ret, structItem)
			}
		case *xast.LeafNode:
		}
		return true
	})

	var txt = ""
	for i, obj := range ret {
		if i == 0 {
			txt = "export interface " + obj.Title + " {\n"
		} else {
			txt += "interface " + obj.Title + "{\n"
		}
		var kv []string
		for k, v := range obj.Items {
			kv = append(kv, fmt.Sprintf("\t%s: %s;", k, v))
		}
		txt += strings.Join(kv, "\n")
		txt += "\n}\n"
	}
	return txt, nil
}

func parseTypeToTS(shortName, fullName string) string {
	fmt.Println(shortName, fullName)
	switch shortName {
	case reflect.Int.String(), reflect.Int8.String(), reflect.Int16.String(), reflect.Int32.String(), reflect.Int64.String(),
		reflect.Uint.String(), reflect.Uint8.String(), reflect.Uint16.String(), reflect.Uint32.String(), reflect.Uint64.String():
		return "number"
	case reflect.Float32.String(), reflect.Float64.String():
		return "number"
	case reflect.Bool.String():
		return "boolean"
	case reflect.String.String():
		return "string"
	case reflect.Array.String(), reflect.Slice.String():
		var array string
		var blocks = strings.Split(fullName, "[]")
		for _, value := range blocks {
			if value == "" {
				array += "Array<"
			} else {
				array += parseTypeToTS(value, value)
			}
		}
		array += strings.Repeat(">", len(blocks)-1)
		return array
	case reflect.Map.String():
		return "Any"
	case reflect.Struct.String():
		return fullName
	default:
		return shortName
	}
}
