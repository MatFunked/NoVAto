package parser

import (
	"fmt"
	"novato/pkg/ast"
	"novato/pkg/lexer"
	"strconv"
)

// Precedencias de operadores
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > o <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X o !X
	CALL        // myFunction(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l       *lexer.Lexer
	curTok  lexer.Token
	peekTok lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix("VERDADERO", p.parseBoolean)
	p.registerPrefix("FALSO", p.parseBoolean)

	// Leer dos tokens para establecer curTok y peekTok
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
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

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}

	// Siguiente token debe ser el tipo (ej: ent, text, bool)
	if !p.expectPeek("T_ENT") && !p.expectPeek("T_TEXT") && !p.expectPeek("T_BOOL") && !p.expectPeek("T_DECI") {
		return nil
	}

	stmt.Type = p.curTok.Literal

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curTok.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curTok}

	value, err := strconv.ParseInt(p.curTok.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("no se pudo parsear %q como entero", p.curTok.Literal)
		// Aquí podríamos registrar errores del parser en el futuro
		_ = msg
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curTok, Value: p.curTok.Type == "VERDADERO"}
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
