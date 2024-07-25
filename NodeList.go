package jolang2

type NodeList []*Node

func (nl NodeList) Add(nodes ...*Node) NodeList {
	return append(nl, nodes...)
}
