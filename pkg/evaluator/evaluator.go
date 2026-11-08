package evaluator

import (
	"fmt"
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

var builtins = map[string]*object.Builtin{
	"imprimir": {Fn: builtinImprimir},
	"mostrar":  {Fn: builtinImprimir},
}

func builtinImprimir(args ...object.Object) object.Object {
	values := make([]interface{}, 0, len(args))

	for _, arg := range args {
		if arg == nil {
			values = append(values, "nulo")
			continue
		}
		values = append(values, arg.Inspect())
	}

	fmt.Println(values...)
	return NULL
}

func Eval(node ast.Node, env *Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Decimal{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

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
		if isError(val) {
			return val
		}
		if err := validateType(node.Type, val); err != nil {
			return err
		}
		env.DefineTyped(node.Name.Value, val, node.Type)
		return val

	case *ast.AssignmentExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if typ, ok := env.GetType(node.Name.Value); ok && typ != "" {
			if err := validateType(typ, val); err != nil {
				return err
			}
		}
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
		if isControlSignal(result) || isError(result) {
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
	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}

	if left == nil || right == nil {
		return NULL
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.DECIMAL_OBJ && right.Type() == object.DECIMAL_OBJ:
		return evalDecimalInfixExpression(operator, left.(*object.Decimal).Value, right.(*object.Decimal).Value)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.DECIMAL_OBJ:
		leftVal := float64(left.(*object.Integer).Value)
		return evalDecimalInfixExpression(operator, leftVal, right.(*object.Decimal).Value)
	case left.Type() == object.DECIMAL_OBJ && right.Type() == object.INTEGER_OBJ:
		rightVal := float64(right.(*object.Integer).Value)
		return evalDecimalInfixExpression(operator, left.(*object.Decimal).Value, rightVal)
	default:
		return NULL
	}
}

func evalDecimalInfixExpression(operator string, leftVal, rightVal float64) object.Object {
	switch operator {
	case "+":
		return &object.Decimal{Value: leftVal + rightVal}
	case "-":
		return &object.Decimal{Value: leftVal - rightVal}
	case "*":
		return &object.Decimal{Value: leftVal * rightVal}
	case "/":
		return &object.Decimal{Value: leftVal / rightVal}
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
	if isError(condition) {
		return condition
	}

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

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}

		result = Eval(we.Body, env)

		if isError(result) {
			return result
		}

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
		if result := Eval(fe.Init, loopEnv); isError(result) {
			return result
		}
	}

	var result object.Object = NULL
	for {
		condition := Eval(fe.Condition, loopEnv)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}

		result = Eval(fe.Body, loopEnv)

		if isError(result) {
			return result
		}

		if result != nil && (result.Type() == object.YIELD_SIGNAL_OBJ || result.Type() == object.RETURN_SIGNAL_OBJ) {
			return result
		}

		if result == BREAK {
			return NULL
		}

		if fe.Increment != nil {
			if increased := Eval(fe.Increment, loopEnv); isError(increased) {
				return increased
			}
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
	if isError(function) {
		return function
	}

	args := evalExpressions(ce.Arguments, env)
	for _, arg := range args {
		if isError(arg) {
			return arg
		}
	}

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
	switch fn := fn.(type) {
	case *object.Builtin:
		return fn.Fn(args...)

	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return NULL
		}

		fnEnv, _ := fn.Env.(*Environment)
		callEnv := NewEnclosedEnvironment(fnEnv)

		for i, param := range fn.Parameters {
			callEnv.Define(param.Value, args[i])
		}

		return unwrapSignal(Eval(fn.Body, callEnv))
	}

	return NULL
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
		if isControlSignal(result) || isError(result) {
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

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func validateType(expectedType string, val object.Object) *object.Error {
	if expectedType == "" || val == nil {
		return nil
	}

	expectedObjType := object.ObjectType("")

	switch expectedType {
	case "ent":
		expectedObjType = object.INTEGER_OBJ
	case "deci":
		expectedObjType = object.DECIMAL_OBJ
	case "text":
		expectedObjType = object.STRING_OBJ
	case "bool":
		expectedObjType = object.BOOLEAN_OBJ
	default:
		return nil
	}

	if val.Type() != expectedObjType {
		return newError("tipo incompatible: se esperaba '%s' y se obtuvo '%s'", expectedType, objectTypeName(val))
	}

	return nil
}

func objectTypeName(obj object.Object) string {
	switch obj.Type() {
	case object.INTEGER_OBJ:
		return "ent"
	case object.DECIMAL_OBJ:
		return "deci"
	case object.STRING_OBJ:
		return "text"
	case object.BOOLEAN_OBJ:
		return "bool"
	case object.NULL_OBJ:
		return "nulo"
	case object.FUNCTION_OBJ:
		return "funcion"
	case object.BUILTIN_OBJ:
		return "nativa"
	case object.ERROR_OBJ:
		return "error"
	default:
		return string(obj.Type())
	}
}
