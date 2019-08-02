package xast

import "go/ast"

var AstMetaKey = "xxxxx"

type AstMeta struct {
	VarName     string   // 声明
	SysType     string   // 系统类型名
	FullName    string   // 全名
	RawExpr     ast.Expr // array: *ast.ArrayType, map: *ast.MapType, struct: *ast.Ident
	Comment     *ast.CommentGroup
	Doc         *ast.CommentGroup
	CrossModule bool // 是否引用跨模块或者跨文件的类型
}

func NewAstMeta(varName, sysType, fullName string, comment, doc *ast.CommentGroup, crossModule bool, rawExpr ast.Expr) AstMeta {
	return AstMeta{
		VarName:     varName,
		SysType:     sysType,
		FullName:    fullName,
		Comment:     comment,
		Doc:         doc,
		RawExpr:     rawExpr,
		CrossModule: crossModule,
	}
}

func (meta *AstMeta) Copy() *AstMeta {
	if meta == nil {
		return nil
	}
	newMeta := NewAstMeta(meta.VarName, meta.SysType, meta.FullName, meta.Comment, meta.Doc, meta.CrossModule, meta.RawExpr)
	return &newMeta
}
