package services

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/sysdes"
	"github.com/geekymedic/neon-cli/types/xast/astutil"
)

func LoadErrCode(filename string) ([]struct {
	Value   int
	VarName string
	Des     string
	Remarks string
}, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	obj := f.Scope.Lookup("_codes")
	if obj == nil {
		return nil, fmt.Errorf("not found _codes object")
	}
	var codesMapping map[string]string
	var constCodes map[string]int
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if _codeMapping, ok := parseCodesSet(genDecl); ok && len(codesMapping) <= 0 {
			codesMapping = _codeMapping
		}
		if _constCodes, ok := parseConstCodes(genDecl); ok && len(constCodes) <= 0 {
			constCodes = _constCodes
		}
		if len(codesMapping) > 0 && len(constCodes) > 0 {
			break
		}
	}
	var ret []struct {
		Value   int
		VarName string
		Des     string
		Remarks string
	}

	for varName, desc := range codesMapping {
		ret = append(ret, struct {
			Value   int
			VarName string
			Des     string
			Remarks string
		}{Value: constCodes[varName], VarName: varName, Des: desc})
	}

	return ret, nil
}

//_codes = Codes{
//	CodeSuccess:                 "请求成功",
//}
func parseCodesSet(genDecl *ast.GenDecl) (map[string]string, bool) {
	var (
		codesMapping = map[string]string{}
		ok           bool
		valueSpec    *ast.ValueSpec
	)
	for _, spec := range genDecl.Specs {
		valueSpec, ok = spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if astutil.VarType(valueSpec.Names) != "_codes" {
			continue
		}
		ok = true
		for _, expr := range valueSpec.Values {
			compositeLit, ok := expr.(*ast.CompositeLit)
			if !ok {
				continue
			}
			for _, expr := range compositeLit.Elts {
				keyValue, ok := expr.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key := astutil.VarType(keyValue.Key) // CodeSuccess
				if !strings.HasPrefix(key, "Code") {
					types.PanicSanity("it should be not happen")
				}
				value := keyValue.Value.(*ast.BasicLit).Value //"请求成功"
				codesMapping[key] = value
			}
		}
		break
	}
	return codesMapping, ok
}

//const (
//	CodeSuccess                 = 0
//)
func parseConstCodes(genDecl *ast.GenDecl) (map[string]int, bool) {
	var (
		constCodes = map[string]int{}
		ok         bool
		valueSpec  *ast.ValueSpec
	)
	for _, spec := range genDecl.Specs {
		valueSpec, ok = spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		varName := astutil.VarType(valueSpec.Names[0])
		if !strings.HasPrefix(varName, "Code") {
			continue
		}
		tokenType, value := astutil.VarValue(valueSpec.Values[0])
		if tokenType != token.INT {
			types.PanicSanity("code must be int type")
		}
		if !strings.HasPrefix(varName, "Code") {
			continue
		}
		constCodes[varName], _ = strconv.Atoi(value)
		ok = true
	}
	return constCodes, ok
}

func loadUseErrorCode(impl *sysdes.BffImpl) ([]string, error) {
	var (
		fset         = token.NewFileSet()
		targetName   = impl.AstTree.Interface.(*sysdes.BffInterfaceTree).FuncName
		fileBuf, err = ioutil.ReadFile(impl.FileNode.Abs())
	)
	if err != nil {
		return nil, err
	}
	f, err := parser.ParseFile(fset, impl.FileNode.Abs(), nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	obj := f.Scope.Lookup(targetName)
	if obj == nil {
		return nil, fmt.Errorf("not found %s interface object at %s", targetName, impl.FileNode.Abs())
	}
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		useCodes, ok := parseUsedErrcode(targetName, funcDecl, fset, bytes.NewBuffer(fileBuf))
		if !ok {
			continue
		}
		return useCodes, nil
	}
	return nil, nil
}

// return []string{"CodeRequestBodyError", "CodeServerError"}
func parseUsedErrcode(targetName string, genDecl *ast.FuncDecl, fset *token.FileSet, fileBuffer *bytes.Buffer) ([]string, bool) {
	var errCodes []string
	fnName := astutil.VarType(genDecl.Name)
	if fnName != targetName {
		return nil, false
	}

	startLine, endLine := fset.Position(genDecl.Body.Lbrace).Line, fset.Position(genDecl.Body.Rbrace).Line
	var regexp, err = regexp.Compile("state.Error\\(bff.Code.*,")
	if err != nil {
		types.PanicSanityf("%v", err)
	}
	var n = 0
	var repeat = map[string]struct{}{}
	for {
		line, err := fileBuffer.ReadString('\n')
		if err == io.EOF {
			break
		}
		n++
		if n < int(startLine) {
			continue
		}
		if n > int(endLine) {
			break
		}
		line = regexp.FindString(line)
		if line == "" {
			continue
		}
		startIdx, endIdx := strings.Index(line, ".Code"), strings.LastIndex(line, ",")
		line = line[startIdx+1 : endIdx]
		_, ok := repeat[line]
		if ok {
			continue
		}
		repeat[line] = struct{}{}
		errCodes = append(errCodes, line)
	}

	return errCodes, true
}
