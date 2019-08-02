package sysdes

import (
	"go/ast"
)

type ServiceAnnotationType = string

const (
	TypeServiceAnnotationType = "type"
	PathServiceAnnotationType = "path"
	DesServiceAnnotationType  = "des"
)

type ServiceInterfaceAnnotation struct {
	Typ  string //`validate:"required,eq=b.i"`
	Path string //`validate:"required"`
}

type ServiceRequestAnnotation struct {
	Typ       string
	Interface string
}

type ServiceResponseAnnotation struct {
	Typ       string
	Interface string
}

type ServiceTree struct {
	Interface AstTree
	Request   AstTree
	Response  AstTree
	rawAst    ast.Expr
}

type ServiceInterfaceTree struct {
	FuncName string `json:"-" yaml:"-"`
	*BaseAstTree
}

func NewServiceInterfaceTree(funcName string) *ServiceInterfaceTree {
	return &ServiceInterfaceTree{FuncName: funcName, BaseAstTree: NewBaseAstTree()}
}

type ServiceRequestTree struct {
	VarName string
	*BaseAstTree
}

func NewServiceRequestTree(varName string) *ServiceRequestTree {
	return &ServiceRequestTree{VarName: varName, BaseAstTree: NewBaseAstTree()}
}

type ServiceResponseTree struct {
	VarName string
	*BaseAstTree
}

func NewServiceResponseTree(varName string) *ServiceResponseTree {
	return &ServiceResponseTree{VarName: varName, BaseAstTree: NewBaseAstTree()}
}
