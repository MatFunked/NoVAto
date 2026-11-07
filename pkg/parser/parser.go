package parser

import (
	"novato/pkg/ast"
	"novato/pkg/lexer"
)

type Parser struct {
	l       *lexer.Lexer
	curTok  lexer.Token
	peekTok lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Leer dos tokens para establecer curTok y peekTok
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curTok.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curTok.Type {
	case "VARIABLE", "CONSTANTE":
		return p.parseVariableStatement()
	default:
		return nil
	}
}

func (p *Parser) parseVariableStatement() ast.Statement {
	stmt := &ast.VariableStatement{Token: p.curTok}

	// Siguiente token debe ser el identificador (nombre de la variable)
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}

	// Siguiente token debe ser el tipo (ej: ent, text, bool)
	if !p.expectPeek("T_ENT") && !p.expectPeek("T_TEXT") && !p.expectPeek("T_BOOL") && !p.expectPeek("T_DECI") {
		return nil
	}

	stmt.Type = p.curTok.Literal

	// Siguiente token debe ser '='
	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	// TODO: Aquí analizaremos la expresión del valor próximamente.
	// Por ahora avanzamos hasta el final de la sentencia o línea.
	for p.curTok.Type != lexer.EOF && p.curTok.Literal != "\n" {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekTok.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}
