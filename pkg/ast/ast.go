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
