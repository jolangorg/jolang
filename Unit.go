package jolang2

import (
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
	"log"
	"os"
)

type Unit struct {
	Project    *Project
	SourceCode []byte
	*sitter.Tree
	Root    *Node
	Package string
	Name    string

	SiblingUnits UnitsMap
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
		fmt.Println(node.Type(), node.GetName())
	})
}

func (u *Unit) WriteASTToFile(filename string) {
	s := ""
	u.TraverseAST(func(node *Node, level int) {
		for i := 0; i < level; i++ {
			s += "\t"
		}
		s += string(node.Type()) + " " + node.GetName() + "\n"
	})
	err := os.WriteFile(filename, []byte(s), os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func (u *Unit) NodeContent(node *sitter.Node) string {
	return string(u.SourceCode[node.StartByte():node.EndByte()])
}

func (u *Unit) WrapNode(node *sitter.Node) *Node {
	if node == nil {
		return nil
	}
	id := uint(node.ID())
	if result, ok := u.Project.NodesById[id]; ok {
		return result
	}

	result := &Node{node, u}
	u.Project.NodesById[id] = result

	return result
}

func (u *Unit) GetSiblingUnits() UnitsMap {
	if u.SiblingUnits == nil {
		u.SiblingUnits = make(UnitsMap)
		if unitMap, ok := u.Project.UnitsByPkg[u.Package]; ok {
			for name, unit := range unitMap {
				if u.Name != name {
					u.SiblingUnits[name] = unit
				}
			}
		}
	}

	return u.SiblingUnits
}
