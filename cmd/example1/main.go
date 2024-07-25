package main

import (
	"flag"
	"fmt"
	"jolang2"
	"jolang2/printers"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var srcPath string
	flag.StringVar(&srcPath, "src", "", "[required] dirs with java files (separated with ':')")

	var unitName string
	flag.StringVar(&unitName, "unit", "", "write specific unit e.g. `org.jbox2d.particle.ParticleSystem`")

	var writeAll bool
	flag.BoolVar(&writeAll, "write-all", false, "write all units")

	flag.Parse()

	if unitName == "" && !writeAll || srcPath == "" {
		flag.Usage()
		return
	}

	srcDirs := strings.Split(srcPath, ":")

	project := jolang2.NewProject()

	//Add src dirs
	for _, srcDir := range srcDirs {
		err := project.AddSourceDir(srcDir)
		if err != nil {
			log.Println(err)
			return
		}
	}

	{
		err := project.IndexDeclarations()
		if err != nil {
			log.Println(err)
			return
		}
	}

	//unit, err := project.AddSource("examples/main/Example1.java")

	printer := printers.NewPrinterJS(project)

	if false {
		printFilenames(project, printer)
	}

	if unitName != "" {
		unit, ok := project.UnitsByAbsName[unitName]
		if !ok {
			log.Println("unit `" + unitName + "` not exists")
			return
		}

		//unit.PrintAST()
		//unit.WriteASTToFile("txt/tree-World.txt")
		//unit.WriteASTToFile("txt/tree-" + unit.Name + ".txt")

		err := writeUnit(unit)
		if err != nil {
			log.Println(err)
			return
		}
	}

	if writeAll {
		writeAllUnits(project)
	}

	//node := unit.FindNodeByType(unit.Root, "class_body")
	//fmt.Printf("Found: row: %d, column: %d", node.StartPoint().Row, node.StartPoint().Column)
}

func writeUnit(unit *jolang2.Unit) error {
	printer := printers.NewPrinterJS(unit.Project)

	content := printer.PrintUnit(unit)
	filename := printer.Filename(unit)
	filenameAST := filepath.Join("output-ast", filename+".txt")
	filename = filepath.Join("output", filename)

	var err error

	//write ast
	{
		err = os.MkdirAll(filepath.Dir(filenameAST), os.ModePerm)
		if err != nil {
			return err
		}

		err = unit.WriteASTToFile(filenameAST)
		if err != nil {
			return err
		}
	}

	//write code
	{
		err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
		if err != nil {
			return err
		}

		err = os.WriteFile(filename, []byte(content), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func printFilenames(project *jolang2.Project, printer printers.Printer) {
	for _, u := range project.UnitsByAbsName {
		filename := printer.Filename(u)
		filename = filepath.Join("output", filename)
		fmt.Println(filename)
	}
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
