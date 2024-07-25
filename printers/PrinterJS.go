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
	printer.Println(";")
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
		printer.Println()

		name := methodDeclaration.GetName()
		if methodDeclaration.IsStatic() {
			printer.Print("static ")
		}
		printer.Print(name)
		params := methodDeclaration.FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)
		block := methodDeclaration.FindNodeByType(nodetype.BLOCK)
		printer.printFormalParams(params)
		printer.Indent++
		printer.Visit(block)
		printer.Indent--

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

			if fieldDeclaration.IsStatic() {
				printer.Print("static ")
			}

			eq := variableDeclarator.FindNodeByType(nodetype.EQUAL)
			if eq == nil {
				printer.Println(fieldName + ";")
			} else {
				expr := eq.NextSibling()
				printer.Print(fieldName, "= ")
				printer.Visit(expr)
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

func (printer *PrinterJS) VisitChildrenOf(node *jolang2.Node) {
	for _, child := range node.Children() {
		printer.Visit(child)
	}
}

func (printer *PrinterJS) printIntegerLiteral(node *jolang2.Node) {
	content := node.Content()
	content = strings.ReplaceAll(content, "L", "")
	printer.Print(content)
}

func (printer *PrinterJS) Visit(node *jolang2.Node) {
	switch node.Type() {

	case nodetype.FIELD_ACCESS, nodetype.IDENTIFIER:
		printer.Print(node.Content())

	case nodetype.EXPRESSION_STATEMENT:
		printer.Visit(node.Child(0))
		printer.Println(";")

	case nodetype.ASSIGNMENT_EXPRESSION:
		printer.Visit(node.Child(0))
		printer.Print(" = ")
		printer.Visit(node.Child(2))

	case nodetype.LEFT_BRACE:
		printer.Println(node.Content())

	case nodetype.RIGHT_BRACE, nodetype.SEMICOLON, nodetype.EQUAL:
		printer.Print(node.Content())

	case nodetype.DECIMAL_INTEGER_LITERAL:
		printer.printIntegerLiteral(node)

	case nodetype.LOCAL_VARIABLE_DECLARATION:
		printer.Print("let ")
		printer.VisitChildrenOf(node.FindNodeByType(nodetype.VARIABLE_DECLARATOR))
		printer.Println(";")

	case nodetype.CAST_EXPRESSION:

		//fmt.Println(node.Content())
		//node.PrintAST()
		//return
		//todo skip cast right now
		children := node.Children()
		for i, child := range children {
			if child.Type() == nodetype.RIGHT_PAREN {
				printer.Print(children[i+1].Content())
				break
			}
		}
	default:
		printer.VisitChildrenOf(node)
		//printer.Print(node.Content())
	}
}
