package nodetype

type NodeType string

const (
	PACKAGE_DECLARATION            = NodeType("package_declaration")
	CLASS_DECLARATION              = NodeType("class_declaration")
	CONSTRUCTOR_DECLARATION        = NodeType("constructor_declaration")
	METHOD_DECLARATION             = NodeType("method_declaration")
	METHOD_INVOCATION              = NodeType("method_invocation")
	CLASS_BODY                     = NodeType("class_body")
	IMPORT_DECLARATION             = NodeType("import_declaration")
	IDENTIFIER                     = NodeType("identifier")
	TYPE_IDENTIFIER                = NodeType("type_identifier")
	MODIFIERS                      = NodeType("modifiers")
	FIELD_DECLARATION              = NodeType("field_declaration")
	VARIABLE_DECLARATOR            = NodeType("variable_declarator")
	EXPRESSION_STATEMENT           = NodeType("expression_statement")
	DECIMAL_INTEGER_LITERAL        = NodeType("decimal_integer_literal")
	DECIMAL_FLOATING_POINT_LITERAL = NodeType("decimal_floating_point_literal")
	LOCAL_VARIABLE_DECLARATION     = NodeType("local_variable_declaration")
	CAST_EXPRESSION                = NodeType("cast_expression")
	ASSIGNMENT_EXPRESSION          = NodeType("assignment_expression")
	FIELD_ACCESS                   = NodeType("field_access")
	ASSERT_STATEMENT               = NodeType("assert_statement")
	SUPER_INTERFACES               = NodeType("super_interfaces")

	NEW    = NodeType("new")
	RETURN = NodeType("return")
	IF     = NodeType("if")
	ELSE   = NodeType("else")

	DOT         = NodeType(".")
	EQUAL       = NodeType("=")
	LEFT_BRACE  = NodeType("{")
	RIGHT_BRACE = NodeType("}")
	LEFT_PAREN  = NodeType("(")
	RIGHT_PAREN = NodeType(")")
	SEMICOLON   = NodeType(";")

	STATIC            = NodeType("static")
	FORMAL_PARAMETERS = NodeType("formal_parameters")
	FORMAL_PARAMETER  = NodeType("formal_parameter")
	BLOCK             = NodeType("block")
)

func (n NodeType) String() string {
	return string(n)
}
