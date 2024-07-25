package jo

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

func (nl NodeList) Filter(fn func(node *Node) bool) NodeList {
	result := NodeList{}
	for _, node := range nl {
		if fn(node) {
			result = result.Add(node)
		}
	}
	return result
}

func (nl NodeList) Contains(node *Node) bool {
	for _, n := range nl {
		if n == node {
			return true
		}
	}
	return false
}

func (nl NodeList) Concat(list NodeList) NodeList {
	result := NodeList{}
	result = append(result, nl...)
	result = append(result, list...)
	return result
}
