package sysdes

import (
	"container/list"
	"context"
	"fmt"
	"go/ast"
	"reflect"

	"github.com/geekymedic/neon-cli/types"
	"github.com/geekymedic/neon-cli/types/xast"

	"github.com/geekymedic/neon/logger"
)

type AstTree interface {
	Parse() error
	ParseAnnotation(doc *ast.CommentGroup) error
	FillCrossStructs() error
}

type BaseAstTree struct{}

func NewBaseAstTree() *BaseAstTree {
	return &BaseAstTree{}
}

func (base *BaseAstTree) Parse() error { return nil }

func (base *BaseAstTree) ParseAnnotation(doc *ast.CommentGroup) error { return nil }

func (base *BaseAstTree) FillCrossStructs() error { return nil }

type ExtendTree struct {
	*xast.TopNode
}

func NewExtendTree(tree *xast.TopNode) *ExtendTree {
	return &ExtendTree{tree}
}

// Note:
// topNode share memory with sourcetree
func (topNode *ExtendTree) ReplaceExtraNode(sourceTree *xast.TopNode, fullName ...string) (count int) {
	if sourceTree == nil {
		return
	}
	if len(fullName) == 0 {
		fullName = []string{sourceTree.TypeName}
	}
	// fullName := sourceTree.Meta.(*xast.AstMeta).FullName
	for _, targetNode := range topNode.FindNodesByFullNames(fullName) {
		node, ok := targetNode.(*xast.ExtraNode)
		if !ok {
			logger.Warnf("the node should extraNode, but actual it is %T", targetNode)
			continue
		}
		node.LeavesNodes = sourceTree.LeavesNodes
		node.ExtraNodes = sourceTree.ExtraNodes
		count++
	}
	if count > 0 {
		topNode.TopNode.ReBuildWalkPath()
		topNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
			return true
		})
	}
	return
}

// FIXME:
// struct {
// 	A int //
//  Local
// }
// struct {
// 	A int //
// }

func (base *ExtendTree) FlatNestedNodes() *xast.TopNode {
	var topNode = base.TopNode
	stack := list.New()
	stack.PushFront(topNode)
	for {
		// pop stack
		node := stack.Front()
		if node == nil {
			break
		}
		stack.Remove(node)

		switch typ := node.Value.(type) {
		case *xast.LeafNode:
		case *xast.ExtraNode:
			for _, leaf := range typ.LeavesNodes {
				stack.PushBack(leaf)
			}
			var flatStruct []string
			for varName, extraNode := range typ.ExtraNodes {
				meta := extraNode.Meta.(*xast.AstMeta)
				// find a nested node
				if meta.SysType == reflect.Struct.String() && meta.VarName == meta.FullName {
					for key, leaf := range extraNode.LeavesNodes {
						typ.LeavesNodes[key] = leaf.Copy()
					}
					for key, extraNode := range extraNode.ExtraNodes {
						typ.ExtraNodes[key] = extraNode.Copy()
					}
					flatStruct = append(flatStruct, varName)
				}
				stack.PushBack(extraNode)
			}
			for _, key := range flatStruct {
				delete(typ.ExtraNodes, key)
			}
		case *xast.TopNode:
			for varName, leaf := range typ.LeavesNodes {
				leaf.WalkPath = typ.TypeName + "." + varName
				stack.PushBack(leaf)
			}
			var flatStruct []string
			for varName, extraNode := range typ.ExtraNodes {
				extraNode.WalkPath = typ.TypeName + "." + varName
				meta := extraNode.Meta.(*xast.AstMeta)
				// find a nested node
				if meta.SysType == reflect.Struct.String() && meta.VarName == meta.FullName {
					for key, leaf := range extraNode.LeavesNodes {
						if typ.LeavesNodes == nil {
							typ.LeavesNodes = map[string]*xast.LeafNode{}
						}
						typ.LeavesNodes[key] = leaf.Copy()
					}
					for key, extraNode := range extraNode.ExtraNodes {
						typ.ExtraNodes[key] = extraNode.Copy()
					}
					flatStruct = append(flatStruct, varName)
				}
				stack.PushBack(extraNode)
			}
			for _, key := range flatStruct {
				delete(typ.ExtraNodes, key)
			}
		default:
			types.PanicSanity(fmt.Sprintf("not support the type: %v", typ))
		}
	}
	return topNode
}
