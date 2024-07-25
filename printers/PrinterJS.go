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

var importJavaLang = map[string]bool{
	"RuntimeException": true,
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

		if _, ok := importJavaLang[s]; ok {
			_, _ = fmt.Fprintf(printer, "import {%s} from 'java/lang/%s.js';", s, s)
			typeIdentifiersReady[s] = true
			continue
		}

		// unknown typeIdentifiers
		//fmt.Println(s)
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
		decls := root.FindNodesByType(nodetype.INTERFACE_DECLARATION)
		for _, decl := range decls {
			printer.printInterface(decl, true)
		}
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
		unitName := absName
		if decls, ok := printer.Project.Declarations[absName]; ok {
			if len(decls) > 0 {
				decl := decls[0]
				unitName = decl.Unit.AbsName()
				//fmt.Println(name, absName, decl.Unit.AbsName())
			}
		}

		printer.importedNames[name] = absName
		path := printer.convertClassNameToPath(unitName)
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
				paramType := paramChild.Content()
				paramType = strings.ReplaceAll(paramType, "[]", "Arr")
				printer.Print(paramType)
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
			printer.Visit(paramChild)
			//t := paramChild.Content()
			//printer.Print(t)
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

func (printer *PrinterJS) printMethodHint(params jolang2.NodeList, resultType *jolang2.Node) {

	//example:

	/**
	 *
	 * @param fixtureA
	 * @param indexA
	 * @param fixtureB
	 * @param indexB
	 * @returns {null|*}
	 */

	printResultType := resultType != nil && resultType.Type() != nodetype.VOID_TYPE
	if len(params) == 0 && !printResultType {
		return
	}

	printer.Println("/**")
	for _, param := range params {
		t := printer.findType(param)
		if t == nil {
			_, _ = fmt.Fprintf(printer, "* @param %s", param.GetName())
		} else {
			s := printer.convertType(t)
			_, _ = fmt.Fprintf(printer, "* @param {%s} %s", s, param.GetName())
		}
		printer.Println()
	}

	if printResultType {
		_, _ = fmt.Fprintf(printer, "* @returns %s", resultType.Content())
		printer.Println()
	}

	printer.Println("*/")
}

func (printer *PrinterJS) convertType(t *jolang2.Node) string {
	s := t.Content()
	if v, ok := typeofTypes[s]; ok {
		return v
	}
	return s
}

func (printer *PrinterJS) findType(node *jolang2.Node) *jolang2.Node {
	result := node.FindNodeByType(
		nodetype.TYPE_IDENTIFIER,
		nodetype.FLOATING_POINT_TYPE,
		nodetype.VOID_TYPE,
		nodetype.BOOLEAN_TYPE,
		nodetype.INTEGRAL_TYPE,
		nodetype.GENERIC_TYPE,
	)

	if result == nil {
		return nil
	}

	if result.Type() == nodetype.GENERIC_TYPE {
		return printer.findType(result)
	}

	return result
}

func (printer *PrinterJS) printMethods(classBody *jolang2.Node) {
	methodDeclarations := classBody.FindNodesByType(nodetype.METHOD_DECLARATION)
	methodsByName := jolang2.NodeListMap{}

	for _, methodDeclaration := range methodDeclarations {
		name := methodDeclaration.GetName()
		methodsByName.AddNode(name, methodDeclaration)
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
		params := methodDeclaration.FindNodeByType(nodetype.FORMAL_PARAMETERS).FindNodesByType(nodetype.FORMAL_PARAMETER)
		resultType := printer.findType(methodDeclaration)

		printer.printMethodHint(params, resultType)

		if methodDeclaration.IsStatic() {
			printer.Print("static ")
		}

		if overloaded {
			printer.printOverloadName(methodDeclaration)
		} else {
			printer.Print(name)
		}

		block := methodDeclaration.FindNodeByType(nodetype.BLOCK)
		printer.printFormalParams(params)

		if block == nil {
			printer.Println("{}")
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
				printer.Println()
				printer.Println("/**")
				printer.Printf("* @var {%s}", printer.convertType(fieldType))
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

func (printer *PrinterJS) printInterface(interfaceDeclaration *jolang2.Node, shouldExport bool) {
	_, _ = fmt.Fprintf(printer, `export class %s extends jo.Interface {`, interfaceDeclaration.GetName())
	printer.Println()

	body := interfaceDeclaration.FindNodeByType(nodetype.INTERFACE_BODY)
	if body != nil {
		printer.Indent++
		printer.printMethods(body)
		printer.Indent--
	}
	printer.Println("}")
}

func (printer *PrinterJS) printEnum(enumDeclaration *jolang2.Node, shouldExport bool) {
	enumName := enumDeclaration.GetName()
	body := enumDeclaration.FindNodeByType(nodetype.ENUM_BODY)
	_, _ = fmt.Fprintf(printer, `export class %s extends jo.Enum {`, enumDeclaration.GetName())
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

func (printer *PrinterJS) printAnonClass(node *jolang2.Node) {
	//todo anon class here
	//todo print anon class content
	t := printer.findType(node)
	_, _ = fmt.Fprintf(printer, "($this => new class extends %s {", t.Content())
	printer.Println()
	printer.Indent++

	printer.Println("constructor(){")
	printer.Print("super")
	printer.Visit(node.FindNodeByType(nodetype.ARGUMENT_LIST))
	printer.Println(";")
	printer.Indent--
	printer.Println("}") //constructor end

	classBody := node.FindNodeByType(nodetype.CLASS_BODY)
	printer.printFields(classBody)
	printer.printMethods(classBody)

	printer.Indent--
	printer.Print("})(this)") //new class end
}

func (printer *PrinterJS) printClass(classDeclaration *jolang2.Node, shouldExport bool) {
	className := classDeclaration.FindNodeByType(nodetype.IDENTIFIER).Content()
	classBody := classDeclaration.FindNodeByType(nodetype.CLASS_BODY)

	subClassDeclarations := classBody.FindNodesByType(
		nodetype.CLASS_DECLARATION,
		nodetype.ENUM_DECLARATION,
		nodetype.INTERFACE_DECLARATION,
	)

	for _, subClassDeclaration := range subClassDeclarations {
		switch subClassDeclaration.Type() {
		case nodetype.CLASS_DECLARATION:
			printer.printClass(subClassDeclaration, false)

		case nodetype.ENUM_DECLARATION:
			printer.printEnum(subClassDeclaration, false)

		case nodetype.INTERFACE_DECLARATION:
			printer.printInterface(subClassDeclaration, false)
		}
		printer.Println()
		printer.Println()
	}

	printer.Print("export class", className, "")
	superclass := classDeclaration.FindNodeByType(nodetype.SUPERCLASS)
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
	if node == nil {
		return
	}

	switch node.Type() {
	case nodetype.NEW, nodetype.RETURN, nodetype.IF, nodetype.ELSE, nodetype.CASE:
		printer.VisitDefault(node)
		printer.Print(" ")

	case nodetype.FIELD_ACCESS:
		printer.VisitDefault(node)

	case nodetype.LINE_COMMENT:
		printer.Println(node.Content())

	case nodetype.THROW:
		printer.Print(node.Content())
		printer.Print(" ")

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
		if node.Parent().Type() == nodetype.ARRAY_INITIALIZER {
			printer.Print("[")
		} else {
			printer.Println(node.Content())
		}

	case nodetype.RIGHT_BRACE:
		if node.Parent().Type() == nodetype.ARRAY_INITIALIZER {
			printer.Print("]")
		} else {
			printer.Println(node.Content())
		}

	case nodetype.EQUAL:
		printer.Print(" ")
		printer.Print(node.Content())
		printer.Print(" ")

	case nodetype.DOUBLE_EQUAL:
		printer.Print(" === ")

	case nodetype.NOT_EQUAL:
		printer.Print(" !== ")

	case nodetype.MINUS, nodetype.PLUS:
		printer.Print("", node.Content(), "")

	case nodetype.SEMICOLON:
		printer.Print(node.Content())

	case nodetype.DECIMAL_INTEGER_LITERAL:
		printer.printIntegerLiteral(node)

	case nodetype.DECIMAL_FLOATING_POINT_LITERAL:
		printer.printFloatLiteral(node)

	case nodetype.LOCAL_VARIABLE_DECLARATION:
		//if node.IsFinal() { //todo too much problem with "const ". Just use "let " everywhere
		//	printer.Print("const ")
		//} else {
		//printer.Print("let ")
		//}

		printer.Print("let ")
		//varDecls := node.FindNodesByType(nodetype.VARIABLE_DECLARATOR)
		for i, decl := range node.FindNodesByType(nodetype.VARIABLE_DECLARATOR) {
			if i > 0 {
				printer.Print(", ")
			}
			printer.VisitChildrenOf(decl)
		}

		if node.Parent().Type() == nodetype.FOR_STATEMENT {
			printer.Print(";")
		} else {
			printer.Println(";")
		}

	case nodetype.DIMENSIONS:
		//just skip

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
		for _, child := range node.Children() {
			switch child.Type() {
			case nodetype.ASSERT:
				printer.Print("assert(")
			case nodetype.SEMICOLON:
				printer.Println(");")
			case nodetype.COLON:
				printer.Print(", ")
			default:
				printer.Visit(child)
			}
		}
	case nodetype.ENUM_DECLARATION:
		printer.printEnum(node, false)
		printer.Println()

	case nodetype.EXPLICIT_CONSTRUCTOR_INVOCATION:
		switch node.Child(0).Type() {
		case nodetype.THIS:
			printer.Print("$this")
		case nodetype.SUPER:
			printer.Print("super")
		default:
			fmt.Println("[WARN] strange EXPLICIT_CONSTRUCTOR_INVOCATION ...")
		}
		printer.VisitChildrenOf(node.FindNodeByType(nodetype.ARGUMENT_LIST))
		printer.Print(";")

	case nodetype.TYPE_ARGUMENTS:
		//todo skip generics now

	case nodetype.OBJECT_CREATION_EXPRESSION:
		if node.Contains(nodetype.CLASS_BODY) {
			printer.printAnonClass(node)
		} else {
			printer.VisitDefault(node)
		}

	case nodetype.ARRAY_TYPE: // used only for jo.suitable
		if node.Parent().Type() != nodetype.FORMAL_PARAMETER {
			break
		}
		dims_count := len(node.FindNodeByType(nodetype.DIMENSIONS).FindNodesByType(nodetype.LEFT_SQUARE_BRACKET))
		for i := 0; i < dims_count; i++ {
			printer.Print("[")
		}
		printer.Visit(node.Child(0))
		for i := 0; i < dims_count; i++ {
			printer.Print("]")
		}

	case nodetype.ARRAY_CREATION_EXPRESSION:
		printer.Print("jo.NewArray(")
		t := node.Child(1)
		//if t == nil {
		//	node.PrintAST()
		//	return
		//}
		printer.Print(t.Content())
		exprs := node.FindNodesByType(nodetype.DIMENSIONS_EXPR)
		for _, expr := range exprs {
			printer.Print(", ")
			printer.Visit(expr.Child(1))
		}
		printer.Print(")")

	default:
		printer.VisitDefault(node)
	}
}
