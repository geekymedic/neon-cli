package astutil

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/xast"
)

// ast parser
type WalkNodeFunc func(ctx context.Context, linkName string, leafIdent ast.Expr) (walkContinue bool)

func MustOpenAst(filename string, src interface{}, mode parser.Mode) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, mode)
	if err != nil {
		types.PanicSanity(err)
	}
	return fset, f
}

func doWalkNodeFuncs(ctx context.Context, linkName string, leafIdent ast.Expr, walkNodeFuncs ...WalkNodeFunc) (walkContinue bool) {
	for _, walkNodeFunc := range walkNodeFuncs {
		if ok := walkNodeFunc(ctx, linkName, leafIdent); !ok {
			return
		}
	}
	walkContinue = true
	return
}

func WalkTypeSpec(ctx context.Context, parentName string, typeSpec ast.Spec, walkNodeFuncs ...WalkNodeFunc) {
	if typeSpec == nil {
		return
	}
	// get type
	spec, ok := typeSpec.(*ast.TypeSpec)
	if !ok {
		return
	}
	identName := spec.Name.Name // type IdentName srtruct, skip it
	switch typ := spec.Type.(type) {
	case *ast.StructType:
		if typ.Fields == nil {
			return
		}
		// Outlive level, TODO Optimize
		if parentName == "" {
			parentName = identName
		}

		for _, field := range typ.Fields.List {
			walkAstField(ctx, parentName, field, walkNodeFuncs...)
		}
	case *ast.InterfaceType:
		if parentName == "" {
			parentName = identName
		}
		for _, field := range typ.Methods.List {
			walkAstField(ctx, parentName, field, walkNodeFuncs...)
		}
	default:
	}
}

func walkAstField(ctx context.Context, parentName string, astField *ast.Field, walkNodeFuncs ...WalkNodeFunc) {
	if astField.Names == nil {
		// anonymous
		// eg: type A struct {Location} or type A struct {*Location} or type A struct {time.Time}
		// A.Location == astField.type
		parentName = fmt.Sprintf("%s.%s", parentName, AnonymousFieldName(astField.Type))
	} else {
		// variable
		// eg: type A struct {L Location}
		// A.L == astField.Names[0].Name
		parentName = fmt.Sprintf("%s.%s", parentName, astField.Names[0].Name)
	}
	sysType, crossModule := SystemType(astField.Type)
	meta := xast.NewAstMeta(VarType(astField), sysType, FullName(astField.Type), astField.Comment, astField.Doc, crossModule, astField.Type)
	ctx = context.WithValue(ctx, xast.AstMetaKey, &meta)
	walkExpr(ctx, parentName, astField.Type, walkNodeFuncs...)
}

func AnonymousFieldName(expr ast.Expr) string {
	switch realType := expr.(type) {
	case *ast.Ident:
		return realType.Name
	case *ast.StarExpr:
		return realType.X.(*ast.Ident).Name
	case *ast.SelectorExpr:
		if realType.Sel != nil {
			return realType.Sel.Name // type A struct {time.Time} or type A struct {Time}
		}
		return realType.X.(*ast.Ident).Name // I don't known
	default:
		panic("unimplement")
	}
}

// System type
func SystemType(expr ast.Expr) (string, bool) {
	switch realType := expr.(type) {
	case *ast.Ident:
		//
		if realType.Obj == nil {
			if IsSysInnerType(realType.Name) {
				return realType.Name, false
			}
			// different file refernce, eg:
			// file a:
			// type Book struct {
			// 		Address Address
			// }
			// file b:
			// type Address struct {
			//
			// }
			// so, it should struct type
			return reflect.Struct.String(), true
		}
		if realType.Obj.Kind != ast.Typ { // type xxx struct
			types.PanicSanity(fmt.Sprintf("kind must be %v", ast.Typ))
		}
		var _ = realType.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
		return reflect.Struct.String(), false
	case *ast.MapType:
		return reflect.Map.String(), false
	case *ast.ArrayType:
		// logger.With("module", realType, "elt", realType.Elt, "type", fmt.Sprintf("%T", realType.Elt)).Info("Cross modules")
		return reflect.Array.String(), !IsSysInnerType(fmt.Sprintf("%v", realType.Elt))
	case *ast.StarExpr: // pointer
		return SystemType(realType.X)
	default:
		types.PanicSanity(fmt.Sprintf("unsupport type %T", expr))
	}
	return "", false
}

func VarType(expr interface{}) string {
	switch realVar := expr.(type) {
	case *ast.Field:
		if len(realVar.Names) <= 0 {
			// FIXME: 嵌套类型，如 struct Person {
			// 		Local
			// }
			// struct Local {}
			return VarType(realVar.Type)
		}
		return realVar.Names[0].Name
	case []*ast.Ident:
		return VarType(realVar[0])
	case *ast.Ident:
		return realVar.Name
	default:
		types.PanicSanity(fmt.Sprintf("unsupport type %T", expr))
	}
	return ""
}

// map[string]Foo parse to map[string]Foo
func FullName(expr ast.Expr) string {
	switch realType := expr.(type) {
	case *ast.Ident:
		return realType.Name
	case *ast.MapType:
		key := FullName(realType.Key)
		value := FullName(realType.Value)
		return fmt.Sprintf("map[%s]%s", key, value)
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", FullName(realType.Elt))
	case *ast.StarExpr: // pionter
		return FullName(realType.X)
	}
	types.PanicSanity(fmt.Sprintf("unsupport type %T", expr))
	return ""
}

// map[string]Foo parse to Foo, string parse to string
func SimpleName(expr ast.Expr) string {
	switch realType := expr.(type) {
	case *ast.Ident:
		return realType.Name
	case *ast.MapType:
		return SimpleName(realType.Value)
	case *ast.ArrayType:
		return SimpleName(realType.Elt)
	case *ast.StarExpr: // pointer
		return SimpleName(realType.X)
	}

	types.PanicSanity(fmt.Sprintf("unsupport type %T", expr))
	return ""
}

func VarValue(expr ast.Expr) (token.Token, string) {
	switch typ := expr.(type) {
	case *ast.BasicLit:
		return typ.Kind, typ.Value
	default:
		types.PanicSanity(fmt.Sprintf("unsupport type %T", expr))
		return token.EOF, ""
	}
}

func walkExpr(ctx context.Context, parentName string, expr ast.Expr, walkNodeFuncs ...WalkNodeFunc) {
	switch realType := expr.(type) {
	case *ast.Ident:
		walkContinue := doWalkNodeFuncs(ctx, parentName, realType, walkNodeFuncs...)
		if !walkContinue {
			return
		}
		if obj := realType.Obj; obj != nil && obj.Decl != nil {
			spec, ok := obj.Decl.(*ast.TypeSpec)
			if ok {
				WalkTypeSpec(ctx, parentName, spec, walkNodeFuncs...)
			}
		}
	case *ast.SelectorExpr: // time.Time struct var
		doWalkNodeFuncs(ctx, parentName, realType, walkNodeFuncs...)
	case *ast.MapType: // terminal it, skip next level
		walkContinue := doWalkNodeFuncs(ctx, parentName, realType, walkNodeFuncs...)
		if !walkContinue {
			return
		}
		walkExpr(ctx, parentName, realType.Value, walkNodeFuncs...)
	case *ast.ArrayType: // terminal it, skip next level
		walkContinue := doWalkNodeFuncs(ctx, parentName, realType, walkNodeFuncs...)
		if !walkContinue {
			return
		}
		walkExpr(ctx, parentName, realType.Elt, walkNodeFuncs...)
	case *ast.StarExpr: // pointer type
		walkExpr(ctx, parentName, realType.X, walkNodeFuncs...)
	case *ast.StructType:
		// doLeafFns(ctx, parentName, realType, leafFns...)
		for _, field := range realType.Fields.List {
			walkAstField(ctx, parentName, field, walkNodeFuncs...)
		}
	case *ast.FuncType:
		// param
		for _, field := range realType.Params.List {
			walkAstField(ctx, parentName, field, walkNodeFuncs...)
		}
		// result
		for _, field := range realType.Results.List {
			walkAstField(ctx, parentName, field, walkNodeFuncs...)
		}
	default:
		panic(fmt.Sprintf("unsupport type: %T", expr))
	}
}

func BuildStructTree(structName string, filename string, src interface{}) (*xast.TopNode, error) {
	var (
		topNode     *xast.TopNode
		structItems []string
		err         error
	)

	var fn WalkNodeFunc = func(ctx context.Context, linkName string, _ ast.Expr) bool {
		structItems = append(structItems, linkName)
		return true
	}

	var buildTreeFn WalkNodeFunc = func(ctx context.Context, linkName string, leafIdent ast.Expr) bool {
		meta := ctx.Value(xast.AstMetaKey).(*xast.AstMeta)
		idx := strings.LastIndex(linkName, ".")
		switch meta.SysType {
		case reflect.Map.String(), reflect.Array.String(), reflect.Struct.String():
			// fmt.Println(meta.SysType, linkName, meta)
			extraNode := xast.NewExtraNode(meta.SysType, linkName, meta, nil, nil)
			err := topNode.AfterInsertExtraNode(linkName[:idx], linkName[idx+1:], *extraNode)
			types.AssertNil(err)
		default:
			leafNode := xast.NewLeafNode(meta.SysType, linkName, meta)
			err := topNode.AfterInsertLeafNode(linkName[:idx], linkName[idx+1:], *leafNode)
			types.AssertNil(err)
		}
		return true
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	for _, decl := range f.Decls {
		genSpec, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, typeSpec := range genSpec.Specs {
			if spec, ok := typeSpec.(*ast.TypeSpec); ok {
				if spec.Name.Name == structName {
					meta := xast.AstMeta{Doc: genSpec.Doc, FullName: structName, SysType: reflect.Struct.String()}
					topNode = xast.NewTopNode(structName, nil, nil, &meta)
					WalkTypeSpec(context.TODO(), "", typeSpec, FilterRepeatMiddle(RecycleCheckMiddle(fn, buildTreeFn)))
					goto Done
				}
			}
		}
	}

Done:

	return topNode, nil
}

func RecycleCheckMiddle(walkNodeFuncs ...WalkNodeFunc) WalkNodeFunc {
	var links = map[string]string{}
	var i = 0 // TODO optimize
	return func(ctx context.Context, linkName string, leafIdent ast.Expr) (walkContinue bool) {
		idx := strings.LastIndex(linkName, ".")
		if i == 0 {
			links[linkName] = linkName
		}
		curTypeName := SimpleName(leafIdent)
		if idx < 0 {
			links[linkName] = curTypeName
			return true
		}
		parentLink := linkName[:idx]
		// check parent node
		for _, item := range strings.Split(links[parentLink], ".") {
			if item == curTypeName {
				return false
			}
		}
		links[linkName] = links[parentLink] + "." + curTypeName
		walkContinue = true
		for _, leafFn := range walkNodeFuncs {
			if ok := leafFn(ctx, linkName, leafIdent); !ok {
				walkContinue = false
			}
		}
		return walkContinue
	}
}

func FilterRepeatMiddle(walkNodeFuncs ...WalkNodeFunc) WalkNodeFunc {
	var flag = map[string]struct{}{}
	return func(ctx context.Context, linkName string, leafIdent ast.Expr) (walkContinue bool) {
		_, walkContinue = flag[linkName]
		if walkContinue {
			return
		}
		flag[linkName] = struct{}{}
		walkContinue = true
		for _, leafFn := range walkNodeFuncs {
			if ok := leafFn(ctx, linkName, leafIdent); !ok {
				walkContinue = false
			}
		}
		return
	}
}

func IsSysInnerType(typ string) bool {
	switch typ {
	case reflect.Bool.String(),
		reflect.Ptr.String(),
		reflect.Chan.String(),
		reflect.String.String(),
		reflect.Map.String(),
		reflect.Array.String(),
		reflect.Struct.String(),
		reflect.Slice.String(),
		reflect.Interface.String(),
		reflect.Func.String():
		return true
	case reflect.Uint8.String(), reflect.Uint16.String(), reflect.Uint32.String(), reflect.Uint64.String(), reflect.Uint.String(),
		reflect.Int8.String(), reflect.Int16.String(), reflect.Int32.String(), reflect.Int64.String(), reflect.Int.String():
		return true
	case reflect.Float32.String(), reflect.Float64.String():
		return true
	default:
		return false
	}
}
