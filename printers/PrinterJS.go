package printers

import (
	"fmt"
	"jolang2"
	"jolang2/nodetype"
	"strings"
)

type PrinterJS struct {
	*BasePrinter
	importedNames map[string]string //shortName -> fullName
	leftPart      bool
}

func NewPrinterJS(project *jolang2.Project) Printer {
	return &PrinterJS{
		BasePrinter:   NewBasePrinter(),
		importedNames: map[string]string{},
	}
}

func (printer *PrinterJS) PrintUnit(unit *jolang2.Unit) string {
	root := unit.Root

	importDeclarations := root.FindNodesByType(nodetype.IMPORT_DECLARATION)
	for _, importDeclaration := range importDeclarations {
		printer.printImport(importDeclaration)
	}

	//import assert
	if root.FindNodeByTypeRecursive(nodetype.ASSERT_STATEMENT) != nil {
		printer.Println("import {assert} from 'jo';")
	}

	siblingUnits := unit.GetSiblingUnits()

	typeIdentifiers := root.FindNodesByTypeRecursive(nodetype.TYPE_IDENTIFIER)
	typeIdentifiersReady := map[string]bool{}
	for _, typeIdentifier := range typeIdentifiers {
		parents := typeIdentifier.Parents()
		if len(parents) > 1 && parents[1].Type() == nodetype.SUPER_INTERFACES {
			continue
		}

		s := typeIdentifier.Content()
		if _, ok := typeIdentifiersReady[s]; ok {
			continue
		}
		if s == unit.Name {
			continue
		}
		if _, ok := printer.importedNames[s]; ok {
			continue
		}

		if siblingUnit, ok := siblingUnits[s]; ok {
			jsPath := printer.convertClassNameToPath(siblingUnit.AbsName())
			_, _ = fmt.Fprintf(printer, "import {%s} from '%s';", s, jsPath)
			printer.Println()
			typeIdentifiersReady[s] = true
			continue
		}

		fmt.Println(s)
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

func (printer *PrinterJS) convertClassNameToPath(name string) string {
	return strings.ReplaceAll(name, ".", "/") + ".js"
}

func (printer *PrinterJS) printImport(importDeclaration *jolang2.Node) {
	ids := importDeclaration.FindNodesByTypeRecursive(nodetype.IDENTIFIER)
	if len(ids) < 1 {
		return
	}
	name := ids[len(ids)-1].Content()
	absName := importDeclaration.Child(1).Content()
	printer.importedNames[name] = absName
	path := printer.convertClassNameToPath(absName)
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
	methodsByName := map[string]jolang2.NodeList{}

	for _, methodDeclaration := range methodDeclarations {
		name := methodDeclaration.GetName()
		if _, ok := methodsByName[name]; ok {
			methodsByName[name] = methodsByName[name].Add(methodDeclaration)
		} else {
			methodsByName[name] = jolang2.NodeList{methodDeclaration}
		}
	}

	//show overloaded methods
	for name, list := range methodsByName {
		if len(list) > 1 {
			fmt.Println(name)
		}
	}

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

		if block == nil {
			printer.Println(";")
		} else {
			printer.Indent++
			printer.Visit(block)
			printer.Indent--
		}

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

func (printer *PrinterJS) printFloatLiteral(node *jolang2.Node) {
	content := node.Content()
	content = strings.ReplaceAll(content, "f", "")
	printer.Print(content)
}

func (printer *PrinterJS) VisitDefault(node *jolang2.Node) {
	if node.ChildCount() > 0 {
		printer.VisitChildrenOf(node)
	} else {
		printer.Print(node.Content())
		//printer.Print(" " + node.Content() + " ")
	}
}

func (printer *PrinterJS) Visit(node *jolang2.Node) {
	if node == nil {
		fmt.Println(node)
	}
	switch node.Type() {
	case nodetype.NEW, nodetype.RETURN, nodetype.IF, nodetype.ELSE:
		printer.VisitDefault(node)
		printer.Print(" ")

	case nodetype.FIELD_ACCESS:
		printer.VisitDefault(node)

	case nodetype.LINE_COMMENT:
		printer.Println(node.Content())

	case nodetype.IDENTIFIER:
		prev := node.PrevSibling()
		firstIdentifier := prev == nil || prev.Type() != nodetype.DOT
		parent := node.Parent()

		if firstIdentifier && parent.Type() != nodetype.VARIABLE_DECLARATOR {
			decl := node.FindDeclaration()
			if decl != nil {
				if decl.Type() == nodetype.VARIABLE_DECLARATOR {
					parent := decl.Parent()
					if parent != nil && parent.Type() == nodetype.FIELD_DECLARATION {
						printer.Print("this.")
					}
				} else if decl.Type() == nodetype.METHOD_DECLARATION {
					printer.Print("this.")
				}
			}
			//printer.firstIdentifier = false
		}

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

	case nodetype.EQUAL:
		printer.Print(" ")
		printer.Print(node.Content())
		printer.Print(" ")

	case nodetype.RIGHT_BRACE, nodetype.SEMICOLON:
		printer.Print(node.Content())

	case nodetype.DECIMAL_INTEGER_LITERAL:
		printer.printIntegerLiteral(node)

	case nodetype.DECIMAL_FLOATING_POINT_LITERAL:
		printer.printFloatLiteral(node)

	case nodetype.LOCAL_VARIABLE_DECLARATION:
		printer.Print("let ")
		printer.VisitChildrenOf(node.FindNodeByType(nodetype.VARIABLE_DECLARATOR))
		printer.Println(";")

	case nodetype.METHOD_INVOCATION:
		printer.VisitDefault(node)

	case nodetype.CAST_EXPRESSION:
		//todo skip cast right now
		children := node.Children()
		for i, child := range children {
			if child.Type() == nodetype.RIGHT_PAREN {
				printer.Print(children[i+1].Content())
				break
			}
		}

	case nodetype.ASSERT_STATEMENT:
		printer.VisitDefault(node)
		printer.Println()

	default:
		printer.VisitDefault(node)
	}
}
