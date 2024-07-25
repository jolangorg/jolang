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

	printer := printers.NewPrinterJS(project)

	if false {
		for _, u := range project.UnitsByAbsName {
			filename := printer.Filename(u)
			filename = filepath.Join("output", filename)
			fmt.Println(filename)
		}
	}

	//unit, err := project.AddSource("examples/main/Mat33.java")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	unit, ok := project.UnitsByAbsName["org.jbox2d.dynamics.Profile"]
	if !ok {
		log.Println("not exists")
		return
	}

	//node := unit.FindNodeByType(unit.Root, "class_body")
	//fmt.Printf("Found: row: %d, column: %d", node.StartPoint().Row, node.StartPoint().Column)

	if false {
		unit.PrintAST()
	}

	if false {
		writeUnit(unit)
	}

	if true {
		writeAllUnits(project)
	}
}

func writeUnit(unit *jolang2.Unit) error {
	printer := printers.NewPrinterJS(unit.Project)

	content := printer.PrintUnit(unit)
	filename := printer.Filename(unit)
	filename = filepath.Join("output", filename)

	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}

	err = os.WriteFile(filename, []byte(content), os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func writeAllUnits(project *jolang2.Project) {
	var err error

	for _, unit := range project.Units {
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
	}

}
