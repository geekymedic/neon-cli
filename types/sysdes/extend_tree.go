package sysdes

import (
	"context"
	"go/ast"

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
// topNode share memeory with sourcetree
func (topNode *ExtendTree) ReplaceExtraNode(sourceTree *xast.TopNode) (count int) {
	if sourceTree == nil {
		return
	}
	for _, targetNode := range topNode.FindNodesByFullNames([]string{sourceTree.TypeName}) {
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
