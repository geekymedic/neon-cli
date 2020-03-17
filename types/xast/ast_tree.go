package xast

import (
	"container/list"
	"context"
	"fmt"
	"strings"

	. "github.com/geekymedic/neon-cli/types"
)

const (
	TypeNameMapPrefix   = "map["
	TypeNameArrayPrefix = "array["
	TypeNameSlicePrefix = "slice["
)

var (
	EmptyTopTree = TopNode{}
)

type WalkFunc func(ctx context.Context, walkPath string, node interface{}) bool

func NewTopNode(typeName string,
	Leaves map[string]*LeafNode,
	extraNodes map[string]*ExtraNode,
	meta interface{}) *TopNode {
	return &TopNode{TypeName: typeName, LeavesNodes: Leaves, ExtraNodes: extraNodes, Meta: meta}
}

func NewExtraNode(typeName string,
	walkPath string,
	meta interface{},
	extraNodes map[string]*ExtraNode,
	leaves map[string]*LeafNode) *ExtraNode {
	return &ExtraNode{TypeName: typeName, WalkPath: walkPath, Meta: meta, ExtraNodes: extraNodes, LeavesNodes: leaves}
}

func NewLeafNode(typeName string, walkPath string, meta interface{}) *LeafNode {
	return &LeafNode{TypeName: typeName, WalkPath: walkPath, Meta: meta}
}

// LeafNode include type: int8,uint8,int16,uint16,int32,uint32,int,int64,uint64,float32,float64,byte,string
type LeafNode struct {
	TypeName string
	Meta     interface{} `json:",omitempty"`
	WalkPath string      `json:",omitempty"`
}

func (leafNode *LeafNode) Copy() *LeafNode {
	var (
		copyNode = new(LeafNode)
	)
	if leafNode == nil {
		return nil
	}
	copyNode.WalkPath = leafNode.WalkPath
	copyNode.TypeName = leafNode.TypeName
	copyNode.Meta = leafNode.Meta
	return copyNode
}

// ExtraNode include type: map, array, struct
type ExtraNode struct {
	TypeName    string
	Meta        interface{}           `json:",omitempty"`
	LeavesNodes map[string]*LeafNode  `json:",omitempty"` // key => varName
	ExtraNodes  map[string]*ExtraNode `json:",omitempty"` // key => varName
	WalkPath    string                `json:",omitempty"`
}

func (extraNode *ExtraNode) Copy() *ExtraNode {
	var (
		copyNode = new(ExtraNode)
	)
	if extraNode == nil {
		return nil
	}
	copyNode.TypeName = extraNode.TypeName
	copyNode.WalkPath = extraNode.WalkPath
	copyNode.Meta = extraNode.Meta
	copyNode.LeavesNodes = extraNode.LeavesNodes
	copyNode.ExtraNodes = extraNode.ExtraNodes
	return copyNode
}

// TopNode: current only support struct
type TopNode struct {
	TypeName    string
	LeavesNodes map[string]*LeafNode  `json:",omitempty"` // key => varName
	ExtraNodes  map[string]*ExtraNode `json:",omitempty"` // key => varName
	Meta        interface{}
}

func (topNode *TopNode) BreadthFirst(ctx context.Context, walkFunc WalkFunc) {
	stack := list.New()
	stack.PushFront(topNode)
	for {
		// pop stack
		node := stack.Front()
		if node == nil {
			return
		}
		stack.Remove(node)

		switch typ := node.Value.(type) {
		case *LeafNode:
			if !walkFunc(ctx, typ.WalkPath, typ) {
				return
			}
		case *ExtraNode:
			if !walkFunc(ctx, typ.WalkPath, typ) {
				return
			}

			for _, leaf := range typ.LeavesNodes {
				stack.PushBack(leaf)
			}

			for _, extraNode := range typ.ExtraNodes {
				stack.PushBack(extraNode)
			}
		case *TopNode:
			ok := walkFunc(ctx, typ.TypeName, typ)
			if !ok {
				return
			}
			for varName, leaf := range typ.LeavesNodes {
				leaf.WalkPath = typ.TypeName + "." + varName
				stack.PushBack(leaf)
			}
			for varName, extraNode := range typ.ExtraNodes {
				extraNode.WalkPath = typ.TypeName + "." + varName
				stack.PushBack(extraNode)
			}
		default:
			PanicSanity(fmt.Sprintf("not support the type: %v", typ))
		}
	}
}

func (topNode *TopNode) DepthFirst(ctx context.Context, walkFunc WalkFunc) {
	stack := list.New()
	stack.PushFront(topNode)

	for {
		// pop stack
		node := stack.Front()
		if node == nil {
			return
		}
		stack.Remove(node)

		switch typ := node.Value.(type) {
		case *LeafNode:
			if !walkFunc(ctx, typ.WalkPath, typ) {
				return
			}
		case *ExtraNode:
			if !walkFunc(ctx, typ.WalkPath, typ) {
				return
			}
			for _, extraNode := range typ.ExtraNodes {
				stack.PushFront(extraNode)
			}
			for _, leaf := range typ.LeavesNodes {
				stack.PushFront(leaf)
			}
		case *TopNode:
			if !walkFunc(ctx, typ.TypeName, typ) {
				return
			}
			for varName, extraNode := range typ.ExtraNodes {
				extraNode.WalkPath = typ.TypeName + "." + varName
				stack.PushBack(extraNode)
			}
			for varName, leaf := range typ.LeavesNodes {
				leaf.WalkPath = typ.TypeName + "." + varName
				stack.PushBack(leaf)
			}
		default:
			PanicSanity(fmt.Sprintf("not support the type: %v", typ))
		}
	}
}

// Parent Node Must be topNode or ExtraNode
func (topNode *TopNode) AfterInsertExtraNode(parentWalkPath, varName string, extraNode ExtraNode) error {
	curNode, ok := topNode.FindNode(parentWalkPath)
	if !ok {
		return fmt.Errorf("not found parent: %s", parentWalkPath)
	}

	switch typ := curNode.(type) {
	case *LeafNode:
		return fmt.Errorf("parent node is a leaf node, can't insert other after it")
	case *TopNode:
		if typ.ExtraNodes == nil {
			typ.ExtraNodes = make(map[string]*ExtraNode)
		}
		_, ok := typ.ExtraNodes[varName]
		if ok {
			// fmt.Println("warning", parentWalkPath, "has exists")
			break
		}
		typ.ExtraNodes[varName] = &extraNode
	case *ExtraNode:
		if typ.ExtraNodes == nil {
			typ.ExtraNodes = make(map[string]*ExtraNode)
		}
		_, ok := typ.ExtraNodes[varName]
		if ok {
			// fmt.Println("warning", parentWalkPath, "has exists")
			break
		}
		typ.ExtraNodes[varName] = &extraNode
	default:
		PanicSanity(fmt.Sprintf("not support the type: %T", curNode))
	}
	// fmt.Printf("insert a new extra node, parentWalkPath:%s, type:%s\n", parentWalkPath, varName)
	return nil
}

// Parent Node Must be topNode or ExtraNode
func (topNode *TopNode) AfterInsertLeafNode(parentWalkPath, varName string, leafNode LeafNode) error {
	curNode, ok := topNode.FindNode(parentWalkPath)
	if !ok {
		return fmt.Errorf("not found parent: %s", parentWalkPath)
	}
	switch typ := curNode.(type) {
	case *LeafNode:
		return fmt.Errorf("parent node is a leaf node, can't insert other after it")
	case *TopNode:
		if typ.LeavesNodes == nil {
			typ.LeavesNodes = make(map[string]*LeafNode)
		}
		_, ok := typ.LeavesNodes[varName]
		if ok {
			// fmt.Println("warning", parentWalkPath, "has exists")
			break
		}
		typ.LeavesNodes[varName] = &leafNode
	case *ExtraNode:
		if typ.LeavesNodes == nil {
			typ.LeavesNodes = make(map[string]*LeafNode)
		}
		_, ok := typ.LeavesNodes[varName]
		if ok {
			// fmt.Println("warning", parentWalkPath, "has exists")
			break
		}
		typ.LeavesNodes[varName] = &leafNode
	default:
		PanicSanity(fmt.Sprintf("not support the type: %v", typ))
	}
	return nil
}

func (topNode *TopNode) DepthCount() int {
	if topNode == nil {
		return 0
	}
	var depth int
	var callback WalkFunc = func(ctx context.Context, walkPath string, node interface{}) bool {
		if count := len(strings.Split(walkPath, ".")); count > depth {
			depth = count
		}
		return true
	}

	topNode.BreadthFirst(context.TODO(), callback)
	return depth
}

func (topNode *TopNode) NodeCount() int {
	count := 0
	topNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
		count++
		return true
	})
	return count
}

func (topNode *TopNode) ReBuildWalkPath() {
	if _, ok := topNode.FindNode(".@."); ok {
		PanicSanity("fail to build ast tree")
	}
}

func (topNode *TopNode) FindNode(walkPath string) (targetNode interface{}, ok bool) {
	stack := list.New()
	stack.PushFront(topNode)
	var callback = findNodeWalkFunc(walkPath, &targetNode, &ok)

	for {
		// pop stack
		node := stack.Front()
		if node == nil {
			return
		}
		stack.Remove(node)
		// TODO
		if topNode.recycle() {
			continue
		}

		switch typ := node.Value.(type) {
		case *LeafNode:
			if !callback(nil, typ.WalkPath, typ) {
				return
			}
		case *ExtraNode:
			if !callback(nil, typ.WalkPath, typ) {
				return
			}
			for varName, extraNode := range typ.ExtraNodes {
				extraNode.WalkPath = typ.WalkPath + "." + varName
				stack.PushFront(extraNode)
			}
			for varName, leaf := range typ.LeavesNodes {
				leaf.WalkPath = typ.WalkPath + "." + varName
				stack.PushFront(leaf)
			}
		case *TopNode:
			if !callback(nil, typ.TypeName, typ) {
				return
			}
			for varName, leaf := range typ.LeavesNodes {
				leaf.WalkPath = typ.TypeName + "." + varName
				stack.PushBack(leaf)
			}
			for varName, extraNode := range typ.ExtraNodes {
				extraNode.WalkPath = typ.TypeName + "." + varName
				stack.PushBack(extraNode)
			}
		default:
			PanicSanity(fmt.Sprintf("not support the type: %v", typ))
		}
	}
}

func (topNode *TopNode) FindNodesBySimpleNames(simpleNames []string, simpleFn func(node interface{}) string) (targetNodes []interface{}) {
	if topNode == nil {
		return nil
	}

	findFn := func(simpleName string) bool {
		for _, _simpleName := range simpleNames {
			if simpleName == _simpleName {
				return true
			}
		}
		return false
	}
	topNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
		if findFn(simpleFn(node)) {
			targetNodes = append(targetNodes, node)
		}
		return true
	})
	return
}

// 注意map[key]value 和 value 以及 []xxxx 和 xxxx的比较
func (topNode *TopNode) FindNodesByFullNames(fullNames []string) (targetNodes []interface{}) {
	if topNode == nil {
		return nil
	}

	findFn := func(fullName string) bool {
		for _, _fullName := range fullNames {
			idx1 := strings.LastIndex(fullName, "]")
			idx2 := strings.LastIndex(_fullName, "]")
			if idx1 > 0 {
				fullName = fullName[idx1+1:]
			}
			if idx2 > 0 {
				_fullName = _fullName[idx2+1:]
			}
			if fullName == _fullName {
				return true
			}
		}
		return false
	}
	topNode.DepthFirst(nil, func(ctx context.Context, walkPath string, node interface{}) bool {
		switch typ := node.(type) {
		case *TopNode:
			if findFn(typ.Meta.(*AstMeta).FullName) {
				targetNodes = append(targetNodes, typ)
			}
		case *ExtraNode:
			if findFn(typ.Meta.(*AstMeta).FullName) {
				targetNodes = append(targetNodes, typ)
			}
		case *LeafNode:
			if findFn(typ.Meta.(*AstMeta).FullName) {
				targetNodes = append(targetNodes, typ)
			}
		}
		return true
	})
	return
}

func (topNode *TopNode) recycle() bool {
	return false
}

func findNodeWalkFunc(walkPath string, node *interface{}, ok *bool) WalkFunc {
	return func(_ context.Context, _walkPath string, _node interface{}) bool {
		if walkPath == _walkPath {
			*node = _node
			*ok = true
			return false
		}
		return true
	}
}

func IsMapType(typeName string) bool {
	return strings.HasPrefix(typeName, TypeNameMapPrefix)
}

func IsSliceType(typeName string) bool {
	return strings.HasPrefix(typeName, TypeNameSlicePrefix)
}

func IsArrayType(typeName string) bool {
	return strings.HasPrefix(typeName, TypeNameArrayPrefix)
}

func IsSliceOrArrayType(typeName string) bool {
	return IsSliceType(typeName) || IsArrayType(typeName)
}

func IsSpecBuildIn(typeName string) bool {
	return IsMapType(typeName) || IsSliceOrArrayType(typeName)
}
