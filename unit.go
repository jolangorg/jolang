package jolang2

import (
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
	"jolang2/nodetype"
)

type Unit struct {
	SourceCode []byte
	*sitter.Tree
	Root    *Node
	Package string
	Name    string
}

type NodeHandlerFunc func(node *Node, level int)

func (u *Unit) AbsName() string {
	return u.Package + "." + u.Name
}

func (u *Unit) TraverseNodeAST(node *Node, level int, handler NodeHandlerFunc) {
	handler(node, level)
	for _, child := range node.Children() {
		u.TraverseNodeAST(child, level+1, handler)
	}
}

func (u *Unit) TraverseAST(handler NodeHandlerFunc) {
	u.TraverseNodeAST(u.Root, 0, handler)
}

func (u *Unit) PrintAST() {
	u.TraverseAST(func(node *Node, level int) {
		for i := 0; i < level; i++ {
			fmt.Print("\t")
		}
		fmt.Println(node.Type())
	})
}

func (u *Unit) NodeContent(node *sitter.Node) string {
	return string(u.SourceCode[node.StartByte():node.EndByte()])
}

func (u *Unit) WrapNode(node *sitter.Node) *Node {
	return &Node{node, u}
}

func (u *Unit) FindNodeByType(node *Node, t nodetype.NodeType) *Node {
	for _, child := range node.Children() {
		if child.Type() == t {
			return child
		}
		found := u.FindNodeByType(child, t)
		if found != nil {
			return found
		}
	}
	return nil
}

func (u *Unit) FindNodesByType(node *Node, t nodetype.NodeType) []*Node {
	result := []*Node{}
	for _, child := range node.Children() {
		if child.Type() == t {
			result = append(result, child)
		}
		found := u.FindNodesByType(child, t)
		result = append(result, found...)
	}
	return result
}
