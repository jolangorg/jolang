package printers

import (
	"fmt"
	"jolang2"
	"jolang2/nodetype"
	"strings"
)

var typeofTypes = map[string]string{ // java -> js
	"float":   "number",
	"int":     "number",
	"boolean": "boolean",
}

type PrinterJS struct {
	*BasePrinter
	importedNames map[string]string //shortName -> fullName
	leftPart      bool
}

func NewPrinterJS(project *jolang2.Project) Printer {
	return &PrinterJS{
		BasePrinter:   NewBasePrinter(project),
		importedNames: map[string]string{},
	}
}

func (printer *PrinterJS) PrintUnit(unit *jolang2.Unit) string {
	root := unit.Root

	//import core, types
	printer.Println(`import * as jo from "jo";`)
	printer.Println(`import {int, float, boolean} from "jo";`)

	//import assert
	if root.FindNodeByTypeRecursive(nodetype.ASSERT_STATEMENT) != nil {
		printer.Println(`import {assert} from "jo";`)
	}

	if len(root.FindNodesByTypeRecursive(nodetype.ENUM_DECLARATION)) > 0 {
		printer.Println(`import {Enum} from "jo";`)
	}

	importDeclarations := root.FindNodesByType(nodetype.IMPORT_DECLARATION)
	for _, importDeclaration := range importDeclarations {
		printer.printImport(importDeclaration)
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

		// unknown typeIdentifiers
		// fmt.Println(s)
	}

	printer.Println()

	//main package classes
	{
		classDeclarations := root.FindNodesByType(nodetype.CLASS_DECLARATION)
		for _, classDeclaration := range classDeclarations {
			printer.printClass(classDeclaration, true)
		}
	}

	//main package enums
	{
		enumDecls := root.FindNodesByType(nodetype.ENUM_DECLARATION)
		for _, decl := range enumDecls {
			printer.printEnum(decl, true)
		}
	}

	//main  package interfaces
	{
		//decls := root.FindNodesByType(nodetype.INTERFACE_DECLARATION)
		//for _, decl := range decls {
		//	printer.printInterface(decl, true)
		//}
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
	absName := importDeclaration.FindNodeByType(nodetype.SCOPED_IDENTIFIER).Content()
	//import * stuff
	if importDeclaration.Contains(nodetype.ASTERISK) {
		pkg := absName
		if units, ok := printer.Project.UnitsByPkg[pkg]; ok {
			for _, unit := range units {
				absName = unit.AbsName()
				name = unit.Name

				printer.importedNames[name] = absName
				path := printer.convertClassNameToPath(absName)
				_, _ = fmt.Fprintf(printer, `import {%s} from "%s"`, name, path)
				printer.Println(";")
			}
		}
	} else {
		printer.importedNames[name] = absName
		path := printer.convertClassNameToPath(absName)
		_, _ = fmt.Fprintf(printer, `import {%s} from "%s"`, name, path)
		printer.Println(";")
	}
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

func (printer *PrinterJS) printOverloadName(methodDeclaration *jolang2.Node) {
	name := methodDeclaration.GetName()
	if methodDeclaration.Type() == nodetype.CONSTRUCTOR_DECLARATION {
		name = "constructor"
	}
	params := methodDeclaration.FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)
	printer.Print(name)
	printer.Print("$")
	for _, param := range params {
		for _, paramChild := range param.Children() {
			if paramChild.Type() != nodetype.MODIFIERS {
				printer.Print(paramChild.Content())
				break
			}
		}
	}

}

func (printer *PrinterJS) printOverloadCheck(methodDeclaration *jolang2.Node) {
	printer.printOverloadCheckFull(methodDeclaration, methodDeclaration.Type() == nodetype.CONSTRUCTOR_DECLARATION)
}

func (printer *PrinterJS) printOverloadCheckFull(methodDeclaration *jolang2.Node, forConstructor bool) {
	params := methodDeclaration.FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)
	if forConstructor {
		printer.Print("case jo.suitable(arguments")
	} else {
		printer.Print("if (jo.suitable(arguments")
	}

	for _, param := range params {
		printer.Print(", ")
		for _, paramChild := range param.Children() {
			if paramChild.Type() == nodetype.MODIFIERS {
				continue
			}
			t := paramChild.Content()
			printer.Print(t)
			break
		}
	}

	if forConstructor {
		printer.Print("): ")
	} else {
		printer.Print(")) return this.")
		printer.printOverloadName(methodDeclaration)
		printer.Println("(...arguments);")
	}
}

func (printer *PrinterJS) printMethods(classBody *jolang2.Node) {
	methodDeclarations := classBody.FindNodesByType(nodetype.METHOD_DECLARATION)
	methodsByName := jolang2.NodeListMap{}

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
			methodDeclaration := list[0]
			printer.Println()
			if methodDeclaration.IsStatic() {
				printer.Print("static ")
			}
			printer.Println(name, "(){")
			printer.Indent++
			for _, methodDeclaration = range list {
				printer.printOverloadCheck(methodDeclaration)
			}
			printer.Println()
			printer.Println("return super." + name + "(...arguments);")
			printer.Indent--
			printer.Println("}")

			//fmt.Println(name)
		}
	}

	for _, methodDeclaration := range methodDeclarations {
		printer.Println()
		printer.Println()

		name := methodDeclaration.GetName()
		overloaded := len(methodsByName[name]) > 1

		if methodDeclaration.IsStatic() {
			printer.Print("static ")
		}

		params := methodDeclaration.FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)

		if overloaded {
			printer.printOverloadName(methodDeclaration)
		} else {
			printer.Print(name)
		}

		block := methodDeclaration.FindNodeByType(nodetype.BLOCK)
		printer.printFormalParams(params)

		if block == nil {
			printer.Println(";") //todo strange, need check on abstract methods
		} else {
			printer.Indent++
			printer.Visit(block)
			printer.Indent--
		}

		//constructorDeclaration.PrintAST()
	}
}

func (printer *PrinterJS) printConstructors(classBody *jolang2.Node, superclass *jolang2.Node) {
	constructorDeclarations := classBody.FindNodesByType(nodetype.CONSTRUCTOR_DECLARATION)
	count := len(constructorDeclarations)
	if count == 0 {
		return
	}

	printer.Println()

	if count == 1 {
		printer.Print("constructor")
		constructorDeclaration := constructorDeclarations[0]
		constructorBody := constructorDeclaration.FindNodeByType(nodetype.CONSTRUCTOR_BODY)
		params := constructorDeclaration.FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)
		printer.printFormalParams(params)
		printer.Indent++
		printer.Visit(constructorBody)
		printer.Indent--
	} else {
		printer.Println("constructor(){")
		printer.Indent++
		printer.Println("const $this = () => {")
		printer.Indent++
		printer.Println("switch (true){")
		printer.Indent++

		for _, constructorDeclaration := range constructorDeclarations {
			printer.printOverloadCheck(constructorDeclaration)
			printer.Println(" {") // start case:
			printer.Indent++
			params := constructorDeclaration.FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)
			if len(params) > 0 {
				printer.Print("let [")
				for i, param := range params {
					if i != 0 {
						printer.Print(", ")
					}
					printer.Print(param.GetName())
				}
				printer.Println("] = arguments;")
			}

			block := constructorDeclaration.FindNodeByType(nodetype.CONSTRUCTOR_BODY)
			if block == nil {
				printer.Println()
			} else {
				printer.Indent++
				printer.VisitChildrenOf(block)
				printer.Indent--
			}

			printer.Println("break;")
			printer.Indent--
			printer.Println("}") // finish case:
			printer.Println()
		}

		if superclass != nil {
			printer.Println()
			printer.Println("default: super(...arguments); break;")
		}

		printer.Indent--
		printer.Println("}") // end of switch

		printer.Indent--
		printer.Println("}") // end of $this helper func

		printer.Println()
		printer.Println("$this(...arguments)")

		printer.Indent--
		printer.Println("}") // end of constructor
	}
}

func (printer *PrinterJS) printSubClass(subClass *jolang2.Node) {
	name := subClass.GetName()
	_, _ = fmt.Fprintf(printer, "static %s = %s;", name, name)
	printer.Println()
}

func (printer *PrinterJS) printFields(classBody *jolang2.Node) {
	fieldDeclarations := classBody.FindNodesByType(nodetype.FIELD_DECLARATION)
	for _, fieldDeclaration := range fieldDeclarations {
		variableDeclarators := fieldDeclaration.FindNodesByType(nodetype.VARIABLE_DECLARATOR)
		for _, variableDeclarator := range variableDeclarators {
			fieldName := variableDeclarator.FindNodeByType(nodetype.IDENTIFIER).Content()

			fieldType := variableDeclarator.PrevSibling()
			if fieldType != nil {
				fieldTypeS := fieldType.Content()
				if s, ok := typeofTypes[fieldTypeS]; ok {
					fieldTypeS = s
				}
				printer.Println()
				printer.Println("/**")
				printer.Printf("* @var {%s}", fieldTypeS)
				printer.Println()
				printer.Println("*/")
			}

			/**
			 * @var {Shape}
			 */

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

func (printer *PrinterJS) printEnum(enumDeclaration *jolang2.Node, shouldExport bool) {
	enumName := enumDeclaration.GetName()
	body := enumDeclaration.FindNodeByType(nodetype.ENUM_BODY)
	if shouldExport {
		printer.Print("export ")
	}
	_, _ = fmt.Fprintf(printer, `class %s extends Enum {`, enumDeclaration.GetName())
	printer.Println()
	if body != nil {
		printer.Indent++
		decls := body.FindNodesByType(nodetype.ENUM_CONSTANT)
		for _, decl := range decls {
			name := decl.GetName()
			_, _ = fmt.Fprintf(printer, `static %s = new %s("%s");`, name, enumName, name)
			printer.Println()
		}
		printer.Indent--
	}
	printer.Println("}")
}

func (printer *PrinterJS) printClass(classDeclaration *jolang2.Node, shouldExport bool) {
	className := classDeclaration.FindNodeByType(nodetype.IDENTIFIER).Content()
	classBody := classDeclaration.FindNodeByType(nodetype.CLASS_BODY)
	superclass := classDeclaration.FindNodeByType(nodetype.SUPERCLASS)

	subClassDeclarations := classBody.FindNodesByTypeRecursive(nodetype.CLASS_DECLARATION)

	if len(subClassDeclarations) > 0 {
		//fmt.Println(subClassDeclarations)
	}

	for _, subClassDeclaration := range subClassDeclarations {
		printer.printClass(subClassDeclaration, false)
		printer.Println()
		printer.Println()
	}

	if shouldExport {
		printer.Print("export ")
	}

	printer.Print("class", className, "")
	if superclass != nil {
		printer.Print(superclass.Content())
	}
	printer.Println(" {")

	printer.Indent++

	for _, subClassDeclaration := range subClassDeclarations {
		printer.printSubClass(subClassDeclaration)
	}

	printer.printFields(classBody)
	printer.printConstructors(classBody, superclass)
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

// print "this." or "ClassName." if needed
func (printer *PrinterJS) printFullPath(node *jolang2.Node) {
	if p := node.FindParent(nodetype.METHOD_DECLARATION); p != nil {
		if p.GetName() == "initializeRegisters" && node.Content() == "pool" {
			fmt.Println(p.GetName(), node.Content())
		}
	}

	prev := node.PrevSibling()
	firstIdentifier := prev == nil || prev.Type() != nodetype.DOT
	parent := node.Parent()

	if !firstIdentifier {
		return
	}

	if parent.Type() == nodetype.VARIABLE_DECLARATOR {
		return
	}

	decl := node.FindDeclaration()
	if decl == nil {
		return
	}

	if decl.Type() == nodetype.VARIABLE_DECLARATOR {
		declParent := decl.Parent()
		if declParent == nil || declParent.Type() != nodetype.FIELD_DECLARATION {
			return
		}

		if declParent.IsStatic() {
			clsDecl := declParent.FindParent(nodetype.CLASS_DECLARATION)
			printer.Print(clsDecl.GetName() + ".")
		} else {
			printer.Print("this.")
		}

		return
	}

	if decl.Type() == nodetype.METHOD_DECLARATION {
		printer.Print("this.")
	}
}

func (printer *PrinterJS) Visit(node *jolang2.Node) {
	switch node.Type() {
	case nodetype.NEW, nodetype.RETURN, nodetype.IF, nodetype.ELSE, nodetype.CASE:
		printer.VisitDefault(node)
		printer.Print(" ")

	case nodetype.FIELD_ACCESS:
		printer.VisitDefault(node)

	case nodetype.LINE_COMMENT:
		printer.Println(node.Content())

	case nodetype.IDENTIFIER:
		printer.printFullPath(node)

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
		if node.IsFinal() {
			printer.Print("const ")
		} else {
			printer.Print("let ")
		}

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

	case nodetype.ENUM_DECLARATION:
		printer.printEnum(node, false)
		printer.Println()

	case nodetype.EXPLICIT_CONSTRUCTOR_INVOCATION:
		printer.Print("$this")
		printer.VisitChildrenOf(node.FindNodeByType(nodetype.ARGUMENT_LIST))
		printer.Print(";")

	default:
		printer.VisitDefault(node)
	}
}
