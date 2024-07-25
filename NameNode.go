package jolang2

import "fmt"

type NameNode struct {
	Name     string
	Children map[string]*NameNode
}

func NewRootNameNode() *NameNode {
	return NewNameNode(".")
}

func NewNameNode(name string) *NameNode {
	return &NameNode{
		Name:     name,
		Children: make(map[string]*NameNode),
	}
}

func (nn *NameNode) AddChild(name string) *NameNode {
	if child, ok := nn.Children[name]; ok {
		return child
	}
	child := NewNameNode(name)
	nn.Children[name] = child
	return child
}

func (nn *NameNode) PrintNameTree(level int) {
	for i := 0; i < level; i++ {
		fmt.Print("\t")
	}
	fmt.Println(nn.Name)
	for _, child := range nn.Children {
		child.PrintNameTree(level + 1)
	}
}
