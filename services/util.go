package services

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/geekymedic/neon-cli/templates"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"github.com/geekymedic/neon-cli/types/xast"
	"github.com/geekymedic/neon/utils/tool"

	"github.com/Pallinder/go-randomdata"
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
		out[varName] = defSmartValue(varName, leaf.TypeName)
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
		out[varName] = defSmartValue(varName, leaf.TypeName)
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
			return defSmartValue(extraNode.Meta.(*xast.AstMeta).VarName, typ)
		case reflect.Map.String():
			fullName := extraNode.Meta.(*xast.AstMeta).FullName
			kIdx, vIdx := strings.Index(fullName, "["), strings.Index(fullName, "]")
			out[fmt.Sprintf("%v", defSmartValue(extraNode.Meta.(*xast.AstMeta).VarName, fullName[kIdx+1:vIdx]))] = defSmartValue(extraNode.Meta.(*xast.AstMeta).VarName, fullName[vIdx+1:])
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

func defSmartValue(valName, typ string) interface{} {
	for _, fn := range []func(string, string) (interface{}, bool){
		localtionSmartDefValue,
		humanSmartDefValue,
		timeSmartDefValue,
		computerSmartDefValue,
		numericSmartDefValue,
	} {
		if ret, ok := fn(valName, typ); ok {
			return ret
		}
	}

	if strings.Contains(strings.ToLower(valName), "phone") && typ == reflect.String.String() {
		return fmt.Sprintf("137" + fmt.Sprintf("%d", tool.RangeBitsInt(10000000, 99999999)))
	}
	if strings.Contains(strings.ToLower(valName), "mail") {
		return randomdata.Email()
	}

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

func humanSmartDefValue(valName, typ string) (interface{}, bool) {
	if strings.Contains(strings.ToLower(valName), "age") && strings.Contains(typ, "int") {
		return tool.RangeBitsInt(16, 90), true
	}
	if strings.Contains(strings.ToLower(valName), "name") && typ == "string" {
		return randomdata.SillyName(), true
	}
	if strings.Contains(strings.ToLower(valName), "gender") && typ == "string" {
		return randomdata.FullName(randomdata.RandomGender), true
	}
	return nil, false
}

func localtionSmartDefValue(valName, typ string) (interface{}, bool) {
	if strings.Contains(strings.ToLower(valName), "country") && typ == "string" {
		return randomdata.Country(randomdata.FullCountry), true
	}
	if strings.Contains(strings.ToLower(valName), "city") && typ == "string" {
		return randomdata.City(), true
	}
	if strings.Contains(strings.ToLower(valName), "province") && typ == "string" {
		return randomdata.Country(randomdata.FullCountry), true
	}
	if strings.Contains(strings.ToLower(valName), "address") && typ == "string" {
		return randomdata.Address(), true
	}

	return nil, false
}

func timeSmartDefValue(valName, typ string) (interface{}, bool) {
	if strings.Contains(strings.ToLower(valName), "day") && typ == "string" {
		return randomdata.Day(), true
	}
	if strings.Contains(strings.ToLower(valName), "month") && typ == "string" {
		return randomdata.Month(), true
	}
	if strings.Contains(strings.ToLower(valName), "year") && typ == "string" {
		return randomdata.Month(), true
	}
	if strings.Contains(strings.ToLower(valName), "yesterday") && typ == "string" {
		return randomdata.Day(), true
	}
	if strings.Contains(strings.ToLower(valName), "tomorrow") && typ == "string" {
		return randomdata.Day(), true
	}
	if strings.Contains(strings.ToLower(valName), "date") && typ == "string" {
		tm, _ := time.Parse(randomdata.DateOutputLayout, randomdata.FullDate())
		return tm.Format("2006-01-02"), true
	}
	if strings.Contains(strings.ToLower(valName), "time") && typ == "string" {
		tm, _ := time.Parse(randomdata.DateOutputLayout, randomdata.FullDate())
		return tm.Format("2006-01-02"), true
	}
	if strings.Contains(strings.ToLower(valName), "hour") && typ == "int" {
		return tool.RangeBitsInt(1, 24), true
	}
	if strings.Contains(strings.ToLower(valName), "minute") && typ == "int" {
		return tool.RangeBitsInt(1, 60), true
	}
	if strings.Contains(strings.ToLower(valName), "second") && typ == "int" {
		return tool.RangeBitsInt(1, 60), true
	}
	return nil, false
}

func computerSmartDefValue(valName, typ string) (interface{}, bool) {
	if strings.Contains(strings.ToLower(valName), "ip") && typ == "string" {
		return randomdata.IpV4Address(), true
	}
	if strings.Contains(strings.ToLower(valName), "mac") && typ == "string" {
		return randomdata.MacAddress(), true
	}
	return nil, false
}

func numericSmartDefValue(valName, typ string) (interface{}, bool) {
	if strings.Contains(strings.ToLower(valName), "number") && strings.Contains(typ, "int") {
		return tool.RangeBitsInt(10, 100), true
	}
	return nil, false
}
