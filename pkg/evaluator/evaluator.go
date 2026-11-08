package evaluator

import (
	"novato/pkg/ast"
	"novato/pkg/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Evaluamos el programa raíz
	case *ast.Program:
		return evalProgram(node)

	// Sentencia de expresión individual
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Literales enteros
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	// Literales booleanos
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
	}
	return result
}
