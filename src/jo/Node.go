package jo

import "C"
import (
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
	"jolang2/src/jo/nodetype"
)

type Node struct {
	*sitter.Node
	Unit *Unit
}

func (n *Node) Child(idx int) *Node {
	if idx < 0 {
		idx = int(n.ChildCount()) + idx
	}
	child := n.Node.Child(idx)

	if child == nil {
		return nil
	}
	return n.Unit.WrapNode(child)
}

func (n *Node) PrevSibling() *Node {
	s := n.Node.PrevSibling()
	if s == nil {
		return nil
	}
	return n.Unit.WrapNode(s)
}

func (n *Node) PrevSiblings(types ...nodetype.NodeType) NodeList {
	result := NodeList{}
	prev := n.PrevSibling()
	for prev != nil {
		if prev.IsType(types...) {
			result = append(result, prev)
		}
		prev = prev.PrevSibling()
	}
	return result
}

func (n *Node) NextSibling() *Node {
	s := n.Node.NextSibling()
	if s == nil {
		return nil
	}
	return n.Unit.WrapNode(s)
}

func (n *Node) Content() string {
	return n.Node.Content(n.Unit.SourceCode)
	//s := n.Node.Content(n.Unit.SourceCode)
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

func (n *Node) FindNodeByType(types ...nodetype.NodeType) *Node {
	for _, child := range n.Children() {
		if child.IsType(types...) {
			return child
		}
	}
	return nil
}

func (n *Node) Contains(types ...nodetype.NodeType) bool {
	return n.FindNodeByType(types...) != nil
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
		if child.IsType(types...) {
			result = append(result, child)
		}
	}
	return result
}

func (n *Node) FindNodesByTypeRecursive(types ...nodetype.NodeType) NodeList {
	result := NodeList{}
	for _, child := range n.Children() {
		if child.IsType(types...) {
			result = append(result, child)
		}

		found := child.FindNodesByTypeRecursive(types...)
		result = append(result, found...)
	}
	return result
}

func (n *Node) IsType(types ...nodetype.NodeType) bool {
	for _, t := range types {
		if n.Type() == t {
			return true
		}
	}
	return false
}

func (n *Node) GetAbsName() string {
	name := n.GetName()
	parents := n.Parents()
	var parent *Node
	for _, parent = range parents {
		if parent.IsType(nodetype.DECLARATIONS...) {
			name = parent.GetName() + "." + name
		}
	}
	pkg := parent.FindNodeByType(nodetype.PACKAGE_DECLARATION)
	if pkg != nil {
		name = pkg.GetName() + "." + name
	}

	return name
}

func (n *Node) GetName() string {
	if n.Type() == nodetype.PACKAGE_DECLARATION {
		return n.Child(1).Content()
	}

	if n.Type() == nodetype.METHOD_DECLARATION {
		id := n.FindNodeByType(nodetype.IDENTIFIER)
		if id == nil {
			return ""
		}
		return id.Content()
	}

	id := n.FindNodeByTypeRecursive(nodetype.IDENTIFIER)
	if id == nil {
		return ""
	}
	return id.Content()
}

func (n *Node) Parent() *Node {
	return n.Unit.WrapNode(n.Node.Parent())
}

func (n *Node) FindParents(t nodetype.NodeType) NodeList {
	parents := n.Parents()
	result := NodeList{}
	for _, p := range parents {
		if p.Type() == t {
			result = append(result, p)
		}
	}
	return result
}

func (n *Node) FindParent(t nodetype.NodeType) *Node {
	parents := n.FindParents(t)
	if len(parents) > 0 {
		return parents[0]
	}
	return nil
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
	s := n.Content()
	firstParent := parents[0]

	//check if node n is method that will call now
	isMethodInvocation := firstParent.Type() == nodetype.METHOD_INVOCATION && len(firstParent.FindNodesByType(nodetype.IDENTIFIER)) == 1

	for _, parent := range parents {
		if isMethodInvocation {
			methodDeclarations := parent.FindNodesByType(nodetype.METHOD_DECLARATION)
			for _, methodDeclaration := range methodDeclarations {
				if methodDeclaration.GetName() == s {
					return methodDeclaration
				}
			}
			continue
		}

		fieldDeclarations := parent.FindNodesByType(nodetype.FIELD_DECLARATION)
		for _, fieldDeclaration := range fieldDeclarations {
			decls := fieldDeclaration.FindNodesByType(nodetype.VARIABLE_DECLARATOR)
			for _, decl := range decls {
				if decl.GetName() == s {
					return decl
				}
			}
		}

		localDeclarations := parent.FindNodesByType(nodetype.LOCAL_VARIABLE_DECLARATION)
		for _, localDeclaration := range localDeclarations {
			decls := localDeclaration.FindNodesByType(nodetype.VARIABLE_DECLARATOR)
			for _, decl := range decls {
				if decl.GetName() == s {
					return decl
				}
			}
		}

		formalParameters := parent.FindNodeByType(nodetype.FORMAL_PARAMETERS)
		if formalParameters == nil {
			continue
		}
		params := formalParameters.FindNodesByType(nodetype.FORMAL_PARAMETER)
		for _, param := range params {
			if param.GetName() == s {
				return param
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

func (n *Node) IsFinal() bool {
	modifiers := n.FindNodeByType(nodetype.MODIFIERS)
	if modifiers == nil {
		return false
	}
	return modifiers.FindNodeByType(nodetype.FINAL) != nil
}
