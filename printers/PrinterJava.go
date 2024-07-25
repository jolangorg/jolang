package printers

import (
	"jolang2"
	"jolang2/nodetype"
)

type PrinterJava struct {
	*BasePrinter
	unit   *jolang2.Unit
	Buffer string
	Indent int
}

func NewPrinterJava(unit *jolang2.Unit) *PrinterJava {
	return &PrinterJava{
		BasePrinter: NewBasePrinter(),
		unit:        unit,
		Buffer:      "",
		Indent:      0,
	}
}

func (u *PrinterJava) PrintChildren(node *jolang2.Node) {
	for _, child := range node.Children() {
		u.PrintNode(child)
	}
}

func (p *PrinterJava) PrintChildrenFromTo(node *jolang2.Node, from, to int) {
	if to < 0 {
		to = int(node.ChildCount()) + to
	}
	for i := from; i < to; i++ {
		p.PrintNode(node.Child(i))
	}
}

func (p *PrinterJava) PrintNode(node *jolang2.Node) {
	switch node.Type() {

	case "package_declaration":
		p.Print("package ")
		p.PrintNode(node.Child(1))
		p.Println(";")
		p.Println()
		return

	case "class_declaration":
		p.Print("class ")
		p.PrintNode(node.Child(1))
		p.Print(" ")
		p.PrintNode(node.Child(2))
		return

	case "method_declaration":
		p.PrintIndent()

	case "class_body":
		p.Println("{")
		p.Indent++
		p.PrintChildrenFromTo(node, 1, -1)
		p.Indent--
		p.Println("}")
		return

	case "public", "static":
		p.Print(node.Type().String())
		p.Print(" ")
	case ";", "{", "}", "(", ")", ".":
		p.Print(node.Type().String())

	case "void_type":
		p.Print("void ")

	case "block":
		p.Println("{")
		p.Indent++
		p.PrintIndent()
		p.PrintChildrenFromTo(node, 1, -1)
		p.Indent--
		p.Println("}")
		return

	case "decimal_integer_literal", "string_literal":
		p.Print(node.Content())

	case "identifier":
		p.Print(node.Content())

	case "type_identifier":
		p.Print(node.Content())
		if node.NextSibling() != nil {
			p.Print(" ")
		}
	}
	p.PrintChildren(node)
}

func (p *PrinterJava) PrintUnit() {
	packageDeclaration := p.unit.Root.FindNodeByType(nodetype.PACKAGE_DECLARATION)
	p.Printf("package %s;", packageDeclaration)

}
