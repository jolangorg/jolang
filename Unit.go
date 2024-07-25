package jolang2

import (
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
)

type Unit struct {
	Project    *Project
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
