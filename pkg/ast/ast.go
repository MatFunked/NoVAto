package ast

import (
	"bytes"
	"novato/pkg/lexer"
	"strings"
)

// Nodo base del AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Nodo principal que engloba todo el programa (lista de sentencias)
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out string
	for _, s := range p.Statements {
		out += s.String()
	}
	return out
}

// Tipos de sentencias
type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Nodo para declaración de variables (ej: variable intentos_maximos ent = 3)
type VariableStatement struct {
	Token lexer.Token // El token 'variable'
	Name  *Identifier
	Type  string
	Value Expression
}

func (vs *VariableStatement) statementNode()       {}
func (vs *VariableStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VariableStatement) String() string {
	var out string
	out += vs.Token.Literal + " " + vs.Name.String()
	if vs.Type != "" {
		out += " " + vs.Type
	}
	out += " = "
	if vs.Value != nil {
		out += vs.Value.String()
	}
	return out
}

// Identificador (nombres de variables, funciones)
type Identifier struct {
	Token lexer.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Literal entero (ej: 3, 18)
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// Literal booleano (Verdadero, Falso)
type Boolean struct {
	Token lexer.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// Expresión infija (ej: edad >= 18, a + b)
type InfixExpression struct {
	Token    lexer.Token // El operador (ej: >=, +, and)
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// Expresión de asignación (ej: contador = contador + 1)
type AssignmentExpression struct {
	Token lexer.Token // El token '='
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	return ae.Name.String() + " = " + ae.Value.String()
}

// Bloque de código (agrupación de sentencias entre llaves)
type BlockStatement struct {
	Token      lexer.Token // El token '{'
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out string
	for _, s := range bs.Statements {
		out += s.String()
	}
	return out
}

// Expresión Condicional (si / osi / sino)
type IfExpression struct {
	Token       lexer.Token // El token 'si' u 'osi'
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // Bloque 'sino' o bloque sintético que contiene un IfExpression 'osi'
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.Token.Literal)
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		if ie.Alternative.Token.Type == "OSI" {
			out.WriteString(" ")
		} else {
			out.WriteString(" sino ")
		}
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// Expresión de bucle mientras (while)
type WhileExpression struct {
	Token     lexer.Token // El token 'mientras'
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString(we.Token.Literal)
	out.WriteString(" ")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Body.String())

	return out.String()
}

// Expresión de bucle para (for) con inicialización, condición e incremento
type ForExpression struct {
	Token     lexer.Token // El token 'para'
	Init      Statement   // Inicialización (ej: variable i ent = 0)
	Condition Expression  // Condición de parada (ej: i < 3)
	Increment Expression  // Incremento por iteración (ej: i = i + 1)
	Body      *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString(fe.Token.Literal)
	out.WriteString(" ")
	if fe.Init != nil {
		out.WriteString(fe.Init.String())
	}
	out.WriteString("; ")
	if fe.Condition != nil {
		out.WriteString(fe.Condition.String())
	}
	out.WriteString("; ")
	if fe.Increment != nil {
		out.WriteString(fe.Increment.String())
	}
	out.WriteString(" ")
	out.WriteString(fe.Body.String())

	return out.String()
}

// Sentencia romper (break): sale inmediatamente del bucle actual
type BreakStatement struct {
	Token lexer.Token // El token 'romper'
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return bs.Token.Literal }

// Sentencia continuar (continue): salta a la siguiente iteración del bucle actual
type ContinueStatement struct {
	Token lexer.Token // El token 'continuar'
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return cs.Token.Literal }

// Sentencia ceder (yield): cede/retorna un valor opcional de forma anticipada
type YieldStatement struct {
	Token lexer.Token // El token 'ceder'
	Value Expression  // Expresión opcional
}

func (ys *YieldStatement) statementNode()       {}
func (ys *YieldStatement) TokenLiteral() string { return ys.Token.Literal }
func (ys *YieldStatement) String() string {
	if ys.Value != nil {
		return ys.Token.Literal + " " + ys.Value.String()
	}
	return ys.Token.Literal
}

// Sentencia retornar (return): devuelve un valor opcional desde una función
type ReturnStatement struct {
	Token lexer.Token // El token 'retornar'
	Value Expression  // Expresión opcional
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	if rs.Value != nil {
		return rs.Token.Literal + " " + rs.Value.String()
	}
	return rs.Token.Literal
}

// Literal de función (ej: func sumar(a ent, b ent) ent { ... } o func(x) { ... })
type FunctionLiteral struct {
	Token          lexer.Token // El token 'func'
	Name           *Identifier // Nombre opcional
	Parameters     []*Identifier
	ParameterTypes []string // Tipos opcionales, alineados con Parameters
	ReturnType     string   // Tipo de retorno opcional
	Body           *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(fl.Token.Literal)
	if fl.Name != nil {
		out.WriteString(" ")
		out.WriteString(fl.Name.String())
	}

	params := make([]string, 0, len(fl.Parameters))
	for i, param := range fl.Parameters {
		paramStr := param.String()
		if i < len(fl.ParameterTypes) && fl.ParameterTypes[i] != "" {
			paramStr += " " + fl.ParameterTypes[i]
		}
		params = append(params, paramStr)
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	if fl.ReturnType != "" {
		out.WriteString(" ")
		out.WriteString(fl.ReturnType)
	}

	out.WriteString(" ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// Expresión de llamada a función (ej: sumar(2, 3))
type CallExpression struct {
	Token     lexer.Token // El token '('
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	args := make([]string, 0, len(ce.Arguments))
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}

// Sentencia de expresión (permite que una expresión viva como sentencia independiente)
type ExpressionStatement struct {
	Token      lexer.Token // El primer token de la expresión
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
