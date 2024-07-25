package main

import (
	"fmt"
	"jolang2"
	"jolang2/printers"
	"log"
	"os"
	"path/filepath"
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

	unit.PrintAST()
	return

	if false {
		printer := printers.NewPrinterJava(unit)
		printer.PrintNode(unit.Root)
		fmt.Println(printer.Buffer)
	}

	if false {
		printer := printers.NewPrinterJS(project)
		content := printer.PrintUnit(unit)
		filename := printer.Filename(unit)
		filename = filepath.Join("output", filename)

		err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}

		err = os.WriteFile(filename, []byte(content), os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}

		//fmt.Println(filename)
		//fmt.Println(content)
	}

	project.PrintNameTree(0)

	//
	//n := unit.Tree.RootNode()
	//unit.printNode(0, n)
}
