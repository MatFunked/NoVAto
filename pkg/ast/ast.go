package ast

import (
	"novato/pkg/lexer"
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
	out += vs.Token.Literal + " " + vs.Name.String() + " " + vs.Type + " = "
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

// Expresión Condicional (si / sino)
type IfExpression struct {
	Token       lexer.Token // El token 'si'
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out string
	out += "si " + ie.Condition.String() + " { " + ie.Consequence.String() + " }"
	if ie.Alternative != nil {
		out += " sino { " + ie.Alternative.String() + " }"
	}
	return out
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
