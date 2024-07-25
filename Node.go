package jolang2

import "C"
import (
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
	"jolang2/nodetype"
)

type Node struct {
	*sitter.Node
	unit *Unit
}

func (n *Node) Child(idx int) *Node {
	if idx < 0 {
		idx = int(n.ChildCount()) + idx
	}
	child := n.Node.Child(idx)

	if child == nil {
		return nil
	}
	return n.unit.WrapNode(child)
}
func (n *Node) PrevSibling() *Node {
	s := n.Node.PrevSibling()
	if s == nil {
		return nil
	}
	return n.unit.WrapNode(s)
}

func (n *Node) NextSibling() *Node {
	s := n.Node.NextSibling()
	if s == nil {
		return nil
	}
	return n.unit.WrapNode(s)
}

func (n *Node) Content() string {
	return n.Node.Content(n.unit.SourceCode)
	//s := n.Node.Content(n.unit.SourceCode)
	//return strings.Trim(s, " ")
}

func (n *Node) Type() nodetype.NodeType {
	return nodetype.NodeType(n.Node.Type())
}

func (n *Node) Children() NodeList {
	count := int(n.ChildCount())
	children := make(NodeList, count)
	for i := 0; i < count; i++ {
		children[i] = n.Child(i)
	}
	return children
}

func (n *Node) Traverse(level int, handler NodeHandlerFunc) {
	handler(n, level)
	for _, child := range n.Children() {
		child.Traverse(level+1, handler)
	}
}

func (n *Node) PrintAST() {
	n.Traverse(0, func(node *Node, level int) {
		for i := 0; i < level; i++ {
			fmt.Print("\t")
		}
		fmt.Println(node.Type(), node.GetName())
	})
}

func (n *Node) FindNodeByType(t nodetype.NodeType) *Node {
	for _, child := range n.Children() {
		if child.Type() == t {
			return child
		}
	}
	return nil
}

func (n *Node) FindNodeByTypeRecursive(t nodetype.NodeType) *Node {
	for _, child := range n.Children() {
		if child.Type() == t {
			return child
		}
		if found := child.FindNodeByTypeRecursive(t); found != nil {
			return found
		}
	}
	return nil
}

func (n *Node) FindNodesByType(types ...nodetype.NodeType) NodeList {
	result := NodeList{}
	for _, child := range n.Children() {
		for _, t := range types {
			if child.Type() == t {
				result = append(result, child)
				break
			}
		}
	}
	return result
}

func (n *Node) FindNodesByTypeRecursive(t nodetype.NodeType) NodeList {
	result := NodeList{}
	for _, child := range n.Children() {
		if child.Type() == t {
			result = append(result, child)
		}
		found := child.FindNodesByTypeRecursive(t)
		result = append(result, found...)
	}
	return result
}

func (n *Node) GetName() string {
	id := n.FindNodeByTypeRecursive(nodetype.IDENTIFIER)
	if id == nil {
		return ""
	}
	return id.Content()
}

func (n *Node) Parent() *Node {
	return n.unit.WrapNode(n.Node.Parent())
}

func (n *Node) Parents() NodeList {
	parents := NodeList{}
	node := n
	for {
		parent := node.Parent()
		if parent == nil {
			break
		}
		parents = append(parents, parent)
		node = parent
	}
	return parents

}

func (n *Node) FindDeclaration() *Node {
	parents := n.Parents()
	for _, parent := range parents {
		fieldDeclarations := parent.FindNodesByType(nodetype.FIELD_DECLARATION)
		for _, fieldDeclaration := range fieldDeclarations {
			decls := fieldDeclaration.FindNodesByType(nodetype.VARIABLE_DECLARATOR)
			for _, decl := range decls {
				if decl.GetName() == n.Content() {
					return decl
				}
			}
		}

		methodDeclarations := parent.FindNodesByType(nodetype.METHOD_DECLARATION)
		for _, methodDeclaration := range methodDeclarations {
			if methodDeclaration.GetName() == n.Content() {
				return methodDeclaration
			}
		}
	}
	return nil
}

func (n *Node) IsStatic() bool {
	modifiers := n.FindNodeByType(nodetype.MODIFIERS)
	if modifiers == nil {
		return false
	}
	return modifiers.FindNodeByType(nodetype.STATIC) != nil
}
