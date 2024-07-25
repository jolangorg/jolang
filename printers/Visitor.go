package printers

import "jolang2"

type Visitor interface {
	Visit(node *jolang2.Node)
}
