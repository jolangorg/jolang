package jolang2

type NodeList []*Node
type NodeListMap map[string]NodeList

func (nl NodeList) Add(nodes ...*Node) NodeList {
	return append(nl, nodes...)
}
