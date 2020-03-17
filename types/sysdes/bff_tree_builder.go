package sysdes

import (
	"container/list"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"os"
	"strings"

	"github.com/geekymedic/neon/logger"

	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/xast"
	"github.com/geekymedic/neon-cli/types/xast/astutil"
	"github.com/geekymedic/neon-cli/util"

	"github.com/geekymedic/neon/errors"
)

func ParseBffInterfaceRequestAstTree(fileNode types.FileNode) (*BffTree, error) {
	return NewBffTree(fileNode)
}

type BffAnnotationType = string

const (
	TypeBffAnnotationType      = "type"
	NameBffAnnotationType      = "name"
	LoginBffAnnotationType     = "login"
	PageBffAnnotationType      = "page"
	URIBffAnnotationType       = "uri"
	DesBffAnnotationType       = "describe"
	InterfaceBffAnnotationType = "interface"
)

type BffInterfaceAnnotation struct {
	Typ   string `validate:"required,nx_contains=b.i-bff.interface"`
	Zh    string `validate:"required"`
	Login string `validate:"required"`
	URI   string
	Des   string
	Page  []string
}

type BffRequestAnnotation struct {
	Typ       string `validate:"required,nx_contains=b.i.rt-bff.interface.request"`
	Interface string `validate:"required"`
}

type BffResponseAnnotation struct {
	Typ       string `validate:"required,nx_contains=b.i.re-bff.interface.response"`
	Interface string `validate:"required"`
}

type BffTree struct {
	Interface AstTree
	Request   AstTree
	Response  AstTree
	rawAst    ast.Expr
}

func NewBffTree(fileNode types.FileNode) (*BffTree, error) {
	var bffTree = &BffTree{}

	// BuildInterface Tree
	{
		interfaceTree := &BffInterfaceTree{
			FileNode: fileNode,
		}
		if err := interfaceTree.Parse(); err != nil {
			return nil, err
		}
		bffTree.Interface = interfaceTree
	}

	// BuildRequest Tree
	{
		requestTree := &BffRequestTree{
			FileNode: fileNode,
		}
		if err := requestTree.Parse(); err != nil {
			return nil, err
		}
		if err := requestTree.FillCrossStructs(); err != nil {
			return nil, err
		}
		bffTree.Request = requestTree
	}

	// BuildResponse Tree
	{
		responseTree := &BffResponseTree{
			FileNode: fileNode,
		}
		if err := responseTree.Parse(); err != nil {
			return nil, err
		}
		if err := responseTree.FillCrossStructs(); err != nil {
			return nil, err
		}
		bffTree.Response = responseTree
	}
	return bffTree, nil
}

type BffInterfaceTree struct {
	FuncName   string `json:"-" yaml:"-"`
	FileNode   types.FileNode
	Annotation *BffInterfaceAnnotation `json:"-" yaml:"-"`
}

func (b *BffInterfaceTree) Parse() error {
	_, astFile := astutil.MustOpenAst(b.FileNode.Abs(), nil, parser.ParseComments|parser.AllErrors|parser.DeclarationErrors)
	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// parse interface name
		b.FuncName = astutil.VarType(genDecl.Name)
		// parse annotation
		err := b.ParseAnnotation(genDecl.Doc)
		if err != nil {
			continue
		}
		break
	}
	if b.Annotation == nil {
		return errors.Format("not found abount bff interface infomation, path: %s", b.FileNode.Abs())
	}
	return nil
}

func (b *BffInterfaceTree) ParseAnnotation(docs *ast.CommentGroup) error {
	var (
		annotation = &BffInterfaceAnnotation{}
	)
	if docs == nil {
		return errors.NewStackError("Not set comment")
	}

	for _, doc := range docs.List {
		doc.Text = strings.TrimSpace(doc.Text)
		idx := strings.Index(doc.Text, "@")
		if idx < 0 {
			continue
		}
		txt := doc.Text[idx+1:]
		txtList := strings.SplitN(txt, ":", 2)
		txtList = append(txtList, "", "")
		switch txtList[0] {
		case TypeBffAnnotationType:
			annotation.Typ = strings.TrimSpace(txtList[1])
		case NameBffAnnotationType:
			annotation.Zh = strings.TrimSpace(txtList[1])
		case LoginBffAnnotationType:
			annotation.Login = strings.TrimSpace(txtList[1])
		case PageBffAnnotationType:
			annotation.Page = util.SplitTrimSpace(txtList[1], "|")
		case DesBffAnnotationType:
			annotation.Des = strings.TrimSpace(txtList[1])
		case URIBffAnnotationType:
			annotation.URI = strings.TrimSpace(txtList[1])
		default:
			util.StdDebug("Not supper annotation: %v", txt)
		}
	}
	if err := bffValidate.Struct(annotation); err != nil {
		return errors.Wrap(err)
	}
	b.Annotation = annotation
	return nil
}

func (b *BffInterfaceTree) FillCrossStructs() error {
	return fmt.Errorf("unimplement function")
}

type BffRequestTree struct {
	Name       string
	FileNode   types.FileNode
	Annotation *BffRequestAnnotation
	TopNode    *xast.TopNode
}

func (b *BffRequestTree) Parse() error {
	_, astFile := astutil.MustOpenAst(b.FileNode.Abs(), nil, parser.ParseComments|parser.AllErrors|parser.DeclarationErrors)
	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			b.Name = astutil.VarType(typeSpec.Name)
			break
		}
		topNode, err := astutil.BuildStructTree(b.Name, b.FileNode.Abs(), nil)
		if err != nil {
			continue
		}
		if topNode == nil {
			continue
		}
		meta := topNode.Meta.(*xast.AstMeta)
		// b.Name = meta.VarName
		// parse annotation
		if err := b.ParseAnnotation(meta.Doc); err != nil {
			// logger.With("path", b.FileNode.Abs(), "err", err).Debug("Fail to parse request")
			continue
		}

		b.TopNode = topNode
		break
	}
	if b.Annotation == nil {
		return errors.Format("not found abount bff request infomation: %s", b.FileNode.Abs())
	}
	return nil
}

// TODO
// Opz walk performance
func (b *BffRequestTree) FillCrossStructs() error {
	types.AssertNotNil(b.TopNode)
	return fillCrossStructs(b.TopNode, b.FileNode)
}

func (b *BffRequestTree) ParseAnnotation(docs *ast.CommentGroup) error {
	var (
		annotation = &BffRequestAnnotation{}
	)
	if docs == nil {
		return errors.NewStackError("Not set comment")
	}
	for _, doc := range docs.List {
		doc.Text = strings.TrimSpace(doc.Text)
		idx := strings.Index(doc.Text, "@")
		if idx < 0 {
			continue
		}
		txt := doc.Text[idx+1:]
		txtList := strings.SplitN(txt, ":", 2)
		txtList = append(txtList, "", "")
		switch txtList[0] {
		case TypeBffAnnotationType:
			annotation.Typ = strings.TrimSpace(txtList[1])
		case InterfaceBffAnnotationType:
			annotation.Interface = strings.TrimSpace(txtList[1])
		default:
			util.StdDebug("Not supper annotation: %v", txt)
		}
	}
	if err := bffValidate.Struct(annotation); err != nil {
		return errors.Wrap(err)
	}
	b.Annotation = annotation
	return nil
}

type BffResponseTree struct {
	Name       string
	FileNode   types.FileNode
	Annotation *BffResponseAnnotation
	TopNode    *xast.TopNode
}

func (b *BffResponseTree) ParseAnnotation(docs *ast.CommentGroup) error {
	var (
		annotation = &BffResponseAnnotation{}
	)
	if docs == nil {
		return errors.NewStackError("Not set comment")
	}
	for _, doc := range docs.List {
		doc.Text = strings.TrimSpace(doc.Text)
		idx := strings.Index(doc.Text, "@")
		if idx < 0 {
			continue
		}
		txt := doc.Text[idx+1:]
		txtList := strings.SplitN(txt, ":", 2)
		txtList = append(txtList, "", "")
		switch txtList[0] {
		case TypeBffAnnotationType:
			annotation.Typ = strings.TrimSpace(txtList[1])
		case InterfaceBffAnnotationType:
			annotation.Interface = strings.TrimSpace(txtList[1])
		default:
			util.StdDebug("Not supper annotation: %v", txt)
		}
	}
	if err := bffValidate.Struct(annotation); err != nil {
		return errors.Wrap(err)
	}
	b.Annotation = annotation
	return nil
}

func (b *BffResponseTree) Parse() error {
	_, astFile := astutil.MustOpenAst(b.FileNode.Abs(), nil, parser.ParseComments|parser.AllErrors|parser.DeclarationErrors)
	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			b.Name = astutil.VarType(typeSpec.Name)
			break
		}
		topNode, err := astutil.BuildStructTree(b.Name, b.FileNode.Abs(), nil)
		if err != nil {
			continue
		}
		if topNode == nil {
			continue
		}
		meta := topNode.Meta.(*xast.AstMeta)
		// b.Name = meta.VarName
		// parse annotation
		if err := b.ParseAnnotation(meta.Doc); err != nil {
			// logger.With("path", b.FileNode.Abs(), "err", err).Debug("Fail to parse response")
			continue
		}
		b.TopNode = topNode
		break
	}
	if b.Annotation == nil {
		return errors.NewStackError("not found abount bff response infomation")
	}
	return nil
}

func (b *BffResponseTree) FillCrossStructs() error {
	types.AssertNotNil(b.TopNode)
	return fillCrossStructs(b.TopNode, b.FileNode)
}

func fillCrossStructs(topNode *xast.TopNode, fileNode types.FileNode) error {
	var (
		stack list.List
		links []*xast.TopNode
	)
	stack.PushBack(topNode)

	for {
		value := stack.Front()
		if value == nil {
			break
		}
		stack.Remove(value)
		topNode := value.Value.(*xast.TopNode)
		links = append(links, topNode)
		var crossModule []string
		topNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
			if extra, ok := node.(*xast.ExtraNode); ok && extra.Meta.(*xast.AstMeta).CrossModule {
				simpleName := astutil.SimpleName(extra.Meta.(*xast.AstMeta).RawExpr)
				if !astutil.IsSysInnerType(simpleName) {
					if fixCrossStructs(simpleName, fileNode.Abs(), nil) {
						logger.With("path", fileNode.Abs(), "typename", topNode.TypeName, "module", simpleName).Info("Find a cross struct")
						crossModule = append(crossModule, simpleName)
					}
				}
			}
			return true
		})
		if len(crossModule) == 0 {
			continue
		}
		for _, module := range crossModule {
			var targetErr = errors.NewStackError("found target")
			err := fileNode.Walk(func(path string, info os.FileInfo, err error) error {
				if info.IsDir() || path == fileNode.Abs() || strings.HasSuffix(path, ".go") == false {
					return nil
				}
				logger.With("path", path, "crossModule", crossModule).Info("Handle cross module")
				topNode, err := astutil.BuildStructTree(module, path, nil)
				if err == nil && topNode != nil && topNode.TypeName == module {
					logger.With("type-name", topNode.TypeName, "path", path, "module", module, "fileNode", fileNode.Abs()).Info("Fix cross file struct")
					stack.PushBack(topNode)
					return targetErr
				}
				return nil
			})
			if err != targetErr {
				logger.With("err", err).Error("fail to build struct tree")
				types.PanicSanity(fmt.Sprintf("not found '%s' object define", module))
			}
		}
	}

	// TODO
	for i := 0; i < len(links); i++ {
		extendTree := NewExtendTree(links[i])
		for j := i + 1; j < len(links); j++ {
			extendTree.ReplaceExtraNode(links[j])
		}
	}

	// TODO Why? need rebuild
	topNode.ReBuildWalkPath()
	return nil
}

func fixCrossStructs(structName string, filename string, src interface{}) bool {
	// logger.With("structName", structName, "filename", filename).Info("fix cross structs")
	topNode, err := astutil.BuildStructTree(structName, filename, src)
	return err == nil && topNode == nil
}
