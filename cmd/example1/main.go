package main

import (
	"fmt"
	"jolang2"
	"jolang2/printers"
	"log"
)

func main() {
	project := jolang2.NewProject()
	err := project.AddSourceDir("~/Projects/jbox2d/jbox2d-library/src/main/java")
	if err != nil {
		log.Println(err)
		return
	}

	//unit, err := project.AddSource("examples/main/Example1.java")
	unit, err := project.AddSource("examples/main/Mat33.java")
	if err != nil {
		log.Println(err)
		return
	}

	//node := unit.FindNodeByType(unit.Root, "class_body")
	//fmt.Printf("Found: row: %d, column: %d", node.StartPoint().Row, node.StartPoint().Column)

	//unit.PrintAST()
	//return

	if false {
		printer := printers.NewPrinterJava(unit)
		printer.PrintNode(unit.Root)
		fmt.Println(printer.Buffer)
	}

	if true {
		printer := printers.NewPrinterJS(project)
		content := printer.PrintUnit(unit)
		filename := printer.Filename(unit)

		fmt.Println(filename)
		fmt.Println(content)
	}

	//
	//n := unit.Tree.RootNode()
	//unit.printNode(0, n)
}
