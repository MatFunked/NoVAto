package evaluator

import (
	"novato/pkg/ast"
	"novato/pkg/object"
)

var (
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	NULL     = &object.Null{}
	BREAK    = &object.BreakSignal{}
	CONTINUE = &object.ContinueSignal{}
)

func Eval(node ast.Node, env *Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.VariableStatement:
		val := Eval(node.Value, env)
		env.Define(node.Name.Value, val)
		return val

	case *ast.AssignmentExpression:
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val)
		return val

	case *ast.BreakStatement:
		return BREAK

	case *ast.ContinueStatement:
		return CONTINUE

	case *ast.YieldStatement:
		return evalYieldStatement(node, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)

	case *ast.FunctionLiteral:
		fn := &object.Function{
			Parameters:     node.Parameters,
			ParameterTypes: node.ParameterTypes,
			ReturnType:     node.ReturnType,
			Body:           node.Body,
			Env:            env,
		}
		if node.Name != nil {
			env.Define(node.Name.Value, fn)
		}
		return fn

	case *ast.CallExpression:
		return evalCallExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	}

	return NULL
}

func evalProgram(program *ast.Program, env *Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		if isControlSignal(result) {
			return result
		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return NULL
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	if left == nil || right == nil {
		return NULL
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	default:
		return NULL
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(ie *ast.IfExpression, env *Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func evalWhileExpression(we *ast.WhileExpression, env *Environment) object.Object {
	var result object.Object = NULL

	for isTruthy(Eval(we.Condition, env)) {
		result = Eval(we.Body, env)

		if result != nil && (result.Type() == object.YIELD_SIGNAL_OBJ || result.Type() == object.RETURN_SIGNAL_OBJ) {
			return result
		}

		if result == BREAK {
			return NULL
		}

		if result == CONTINUE {
			result = NULL
			continue
		}
	}

	return result
}

func evalForExpression(fe *ast.ForExpression, env *Environment) object.Object {
	loopEnv := NewEnclosedEnvironment(env)

	if fe.Init != nil {
		Eval(fe.Init, loopEnv)
	}

	var result object.Object = NULL
	for isTruthy(Eval(fe.Condition, loopEnv)) {
		result = Eval(fe.Body, loopEnv)

		if result != nil && (result.Type() == object.YIELD_SIGNAL_OBJ || result.Type() == object.RETURN_SIGNAL_OBJ) {
			return result
		}

		if result == BREAK {
			return NULL
		}

		if fe.Increment != nil {
			Eval(fe.Increment, loopEnv)
		}

		if result == CONTINUE {
			result = NULL
		}
	}

	return result
}

func evalYieldStatement(ys *ast.YieldStatement, env *Environment) object.Object {
	signal := &object.YieldSignal{}

	if ys.Value != nil {
		signal.Value = Eval(ys.Value, env)
	}

	return signal
}

func evalReturnStatement(rs *ast.ReturnStatement, env *Environment) object.Object {
	signal := &object.ReturnSignal{}

	if rs.Value != nil {
		signal.Value = Eval(rs.Value, env)
	}

	return signal
}

func evalCallExpression(ce *ast.CallExpression, env *Environment) object.Object {
	function := Eval(ce.Function, env)
	args := evalExpressions(ce.Arguments, env)

	return applyFunction(function, args)
}

func evalExpressions(exps []ast.Expression, env *Environment) []object.Object {
	result := make([]object.Object, 0, len(exps))

	for _, exp := range exps {
		result = append(result, Eval(exp, env))
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return NULL
	}

	if len(args) != len(function.Parameters) {
		return NULL
	}

	fnEnv, _ := function.Env.(*Environment)
	callEnv := NewEnclosedEnvironment(fnEnv)

	for i, param := range function.Parameters {
		callEnv.Define(param.Value, args[i])
	}

	return unwrapSignal(Eval(function.Body, callEnv))
}

func unwrapSignal(obj object.Object) object.Object {
	if obj == nil {
		return NULL
	}

	switch obj.Type() {
	case object.YIELD_SIGNAL_OBJ:
		if signal, ok := obj.(*object.YieldSignal); ok {
			if signal.Value == nil {
				return NULL
			}
			return signal.Value
		}
	case object.RETURN_SIGNAL_OBJ:
		if signal, ok := obj.(*object.ReturnSignal); ok {
			if signal.Value == nil {
				return NULL
			}
			return signal.Value
		}
	}

	return obj
}

func isTruthy(obj object.Object) bool {
	if obj == nil {
		return false
	}

	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		if boolObj, ok := obj.(*object.Boolean); ok {
			return boolObj.Value
		}
		return true
	}
}

func evalBlockStatement(block *ast.BlockStatement, env *Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if isControlSignal(result) {
			return result
		}
	}
	return result
}

func isControlSignal(obj object.Object) bool {
	if obj == nil {
		return false
	}

	switch obj.Type() {
	case object.BREAK_SIGNAL_OBJ, object.CONTINUE_SIGNAL_OBJ, object.YIELD_SIGNAL_OBJ, object.RETURN_SIGNAL_OBJ:
		return true
	}

	return false
}
