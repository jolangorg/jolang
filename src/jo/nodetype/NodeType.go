package nodetype

type NodeType string

const (
	PACKAGE_DECLARATION             = NodeType("package_declaration")
	CLASS_DECLARATION               = NodeType("class_declaration")
	ENUM_DECLARATION                = NodeType("enum_declaration")
	ENUM_CONSTANT                   = NodeType("enum_constant")
	ENUM_BODY                       = NodeType("enum_body")
	INTERFACE_DECLARATION           = NodeType("interface_declaration")
	INTERFACE_BODY                  = NodeType("interface_body")
	CONSTRUCTOR_DECLARATION         = NodeType("constructor_declaration")
	CONSTRUCTOR_BODY                = NodeType("constructor_body")
	SUPERCLASS                      = NodeType("superclass")
	EXPLICIT_CONSTRUCTOR_INVOCATION = NodeType("explicit_constructor_invocation")
	ARGUMENT_LIST                   = NodeType("argument_list")
	SCOPED_IDENTIFIER               = NodeType("scoped_identifier")
	ASTERISK                        = NodeType("asterisk")
	METHOD_DECLARATION              = NodeType("method_declaration")
	METHOD_INVOCATION               = NodeType("method_invocation")
	CLASS_BODY                      = NodeType("class_body")
	IMPORT_DECLARATION              = NodeType("import_declaration")
	IDENTIFIER                      = NodeType("identifier")

	TYPE_IDENTIFIER     = NodeType("type_identifier")
	FLOATING_POINT_TYPE = NodeType("floating_point_type")
	VOID_TYPE           = NodeType("void_type")
	BOOLEAN_TYPE        = NodeType("boolean_type")
	INTEGRAL_TYPE       = NodeType("integral_type")
	GENERIC_TYPE        = NodeType("generic_type")

	MODIFIERS                      = NodeType("modifiers")
	FIELD_DECLARATION              = NodeType("field_declaration")
	VARIABLE_DECLARATOR            = NodeType("variable_declarator")
	EXPRESSION_STATEMENT           = NodeType("expression_statement")
	FOR_STATEMENT                  = NodeType("for_statement")
	DECIMAL_INTEGER_LITERAL        = NodeType("decimal_integer_literal")
	DECIMAL_FLOATING_POINT_LITERAL = NodeType("decimal_floating_point_literal")
	LOCAL_VARIABLE_DECLARATION     = NodeType("local_variable_declaration")
	CAST_EXPRESSION                = NodeType("cast_expression")
	ASSIGNMENT_EXPRESSION          = NodeType("assignment_expression")
	FIELD_ACCESS                   = NodeType("field_access")
	ASSERT_STATEMENT               = NodeType("assert_statement")
	ASSERT                         = NodeType("assert")
	SUPER_INTERFACES               = NodeType("super_interfaces")
	TYPE_ARGUMENTS                 = NodeType("type_arguments")
	OBJECT_CREATION_EXPRESSION     = NodeType("object_creation_expression")

	//array stuff
	ARRAY_INITIALIZER         = NodeType("array_initializer")
	ARRAY_TYPE                = NodeType("array_type")
	ARRAY_CREATION_EXPRESSION = NodeType("array_creation_expression")
	DIMENSIONS                = NodeType("dimensions")
	DIMENSIONS_EXPR           = NodeType("dimensions_expr")

	BLOCK_COMMENT = NodeType("block_comment")
	LINE_COMMENT  = NodeType("line_comment")

	THIS  = NodeType("this")
	SUPER = NodeType("super")

	THROW = NodeType("throw")

	NEW    = NodeType("new")
	RETURN = NodeType("return")
	IF     = NodeType("if")
	ELSE   = NodeType("else")
	CASE   = NodeType("case")

	DOT                  = NodeType(".")
	EQUAL                = NodeType("=")
	DOUBLE_EQUAL         = NodeType("==")
	NOT_EQUAL            = NodeType("!=")
	LEFT_BRACE           = NodeType("{")
	RIGHT_BRACE          = NodeType("}")
	LEFT_SQUARE_BRACKET  = NodeType("[")
	RIGHT_SQUARE_BRACKET = NodeType("]")
	LEFT_PAREN           = NodeType("(")
	RIGHT_PAREN          = NodeType(")")
	SEMICOLON            = NodeType(";")
	COLON                = NodeType(":")

	MINUS = NodeType("-")
	PLUS  = NodeType("+")

	INC = NodeType("++")
	DEC = NodeType("--")

	STATIC            = NodeType("static")
	FINAL             = NodeType("final")
	FORMAL_PARAMETERS = NodeType("formal_parameters")
	FORMAL_PARAMETER  = NodeType("formal_parameter")
	BLOCK             = NodeType("block")
)

var DECLARATIONS = []NodeType{
	CLASS_DECLARATION,
	ENUM_DECLARATION,
	INTERFACE_DECLARATION,
	CONSTRUCTOR_DECLARATION,
	METHOD_DECLARATION,
	IMPORT_DECLARATION,
	FIELD_DECLARATION,
	LOCAL_VARIABLE_DECLARATION,
}

func (n NodeType) String() string {
	return string(n)
}
