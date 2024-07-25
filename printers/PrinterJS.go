package printers

import (
	"fmt"
	"jolang2"
	"jolang2/nodetype"
	"strings"
)

type PrinterJS struct {
	*BasePrinter
}

func NewPrinterJS(project *jolang2.Project) Printer {
	return &PrinterJS{
		BasePrinter: NewBasePrinter(),
	}
}

func (printer *PrinterJS) PrintUnit(unit *jolang2.Unit) string {
	root := unit.Root

	importDeclarations := root.FindNodesByType(nodetype.IMPORT_DECLARATION)
	for _, importDeclaration := range importDeclarations {
		printer.printImport(importDeclaration)
	}
	if len(importDeclarations) > 0 {
		printer.Println()
	}

	classDeclarations := root.FindNodesByType(nodetype.CLASS_DECLARATION)
	for _, classDeclaration := range classDeclarations {
		printer.printClass(classDeclaration)
	}

	return printer.Buffer
}

func (printer *PrinterJS) printImport(importDeclaration *jolang2.Node) {
	ids := importDeclaration.FindNodesByTypeRecursive(nodetype.IDENTIFIER)
	if len(ids) < 1 {
		return
	}
	name := ids[len(ids)-1].Content()
	path := importDeclaration.Child(1).Content()
	path = strings.ReplaceAll(path, ".", "/") + ".js"
	_, _ = fmt.Fprintf(printer, `import {%s} from "%s"`, name, path)
	printer.Println()
}

func (printer *PrinterJS) printExpr(exprNode *jolang2.Node) {
	printer.Print(exprNode.Content())
}

func (printer *PrinterJS) printFormalParams(params []*jolang2.Node) {
	printer.Print("(")
	for i, param := range params {
		if i != 0 {
			printer.Print(", ")
		}
		printer.Print(param.GetName())
	}
	printer.Print(")")
}

func (printer *PrinterJS) printMethods(classBody *jolang2.Node) {
	methodDeclarations := classBody.FindNodesByType(nodetype.METHOD_DECLARATION)
	for _, methodDeclaration := range methodDeclarations {
		printer.Println()

		name := methodDeclaration.GetName()
		if methodDeclaration.IsStatic() {
			printer.Print("static ")
		}
		printer.Print(name)
		params := methodDeclaration.FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)
		printer.printFormalParams(params)
		printer.Println("{}")
		//constructorDeclaration.PrintAST()
	}
}

func (printer *PrinterJS) printConstructors(classBody *jolang2.Node) {
	constructorDeclarations := classBody.FindNodesByType(nodetype.CONSTRUCTOR_DECLARATION)
	count := len(constructorDeclarations)
	if count == 0 {
		return
	}

	printer.Println()

	if count == 1 {
		printer.Print("constructor")
		params := constructorDeclarations[0].FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)
		printer.printFormalParams(params)
		printer.Println(" {")
		printer.Println("}")
	} else {
		printer.Println("constructor(){")
		printer.Println("}")
	}

	//for _, constructorDeclaration := range constructorDeclarations {
	//	printer.Println()
	//	printer.Println()
	//
	//	//name := constructorDeclaration.GetName()
	//
	//	//constructorDeclaration.PrintAST()
	//}
}

func (printer *PrinterJS) printFields(classBody *jolang2.Node) {
	fieldDeclarations := classBody.FindNodesByType(nodetype.FIELD_DECLARATION)
	for _, fieldDeclaration := range fieldDeclarations {
		variableDeclarators := fieldDeclaration.FindNodesByType(nodetype.VARIABLE_DECLARATOR)
		for _, variableDeclarator := range variableDeclarators {
			fieldName := variableDeclarator.FindNodeByType(nodetype.IDENTIFIER).Content()

			static := fieldDeclaration.FindNodeByType(nodetype.STATIC) != nil
			if static {
				printer.Print("static ")
			}

			eq := variableDeclarator.FindNodeByType(nodetype.EQUAL)
			if eq == nil {
				printer.Println(fieldName + ";")
			} else {
				expr := eq.NextSibling()
				printer.Print(fieldName, "= ")
				printer.printExpr(expr)
				printer.Println(";")
			}
		}

	}
}

func (printer *PrinterJS) printClass(classDeclaration *jolang2.Node) {
	className := classDeclaration.FindNodeByType(nodetype.IDENTIFIER).Content()
	classBody := classDeclaration.FindNodeByType(nodetype.CLASS_BODY)

	printer.Println("export class", className, "{")
	printer.Indent++
	printer.printFields(classBody)
	printer.printConstructors(classBody)
	printer.printMethods(classBody)
	printer.Indent--
	printer.Println("}")
}

func (printer *PrinterJS) Filename(unit *jolang2.Unit) string {
	return strings.ReplaceAll(unit.AbsName(), ".", "/") + ".js"
}