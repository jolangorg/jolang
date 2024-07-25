package jolang2

type NodeList []*Node
type NodeListMap map[string]NodeList

func (nl NodeList) Add(nodes ...*Node) NodeList {
	return append(nl, nodes...)
}

func (nlm NodeListMap) AddNode(k string, node *Node) {
	if nl, ok := nlm[k]; ok {
		nlm[k] = append(nl, node)
	} else {
		nlm[k] = NodeList{node}
	}
}
