package parser

import (
	"novato/pkg/ast"
	"novato/pkg/lexer"
	"strconv"
)

// Precedencias de operadores
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > o <, >=, <=
	SUM         // + o -
	PRODUCT     // * o /
	PREFIX      // -X o !X
	CALL        // myFunction(X)
)

var precedences = map[lexer.TokenType]int{
	lexer.EQ:    EQUALS,
	lexer.LTE:   LESSGREATER,
	lexer.GTE:   LESSGREATER,
	lexer.LT:    LESSGREATER,
	lexer.RT:    LESSGREATER,
	lexer.PLUS:  SUM,
	lexer.MINUS: SUM,
	"and":       EQUALS, // Precedencia para operador lógico
}

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

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LTE, p.parseInfixExpression)
	p.registerInfix(lexer.GTE, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.RT, p.parseInfixExpression)
	p.registerInfix("and", p.parseInfixExpression)

	p.registerPrefix("si", p.parseIfExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curTok.Type != lexer.EOF {
		// Si es un salto de línea puro, simplemente avanzamos
		if p.curTok.Literal == "\n" || p.curTok.Type == "NL" { // Ajusta el tipo según tu token de salto de línea si aplica
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			// Evitar imprimir sentencias vacías en el string
			if stmt.String() != "" {
				program.Statements = append(program.Statements, stmt)
			}
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
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curTok}
	stmt.Expression = p.parseExpression(LOWEST)

	// Opcional: consumir saltos de línea si tu lexer los emite
	if p.peekTokenIs("\n") {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVariableStatement() ast.Statement {
	stmt := &ast.VariableStatement{Token: p.curTok}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}

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

	for !p.peekTokenIs(lexer.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekTok.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curTok,
		Operator: p.curTok.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curTok}

	value, err := strconv.ParseInt(p.curTok.Literal, 0, 64)
	if err != nil {
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curTok}

	// Avanzamos para evaluar la condición del 'si'
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	// Esperamos la apertura del bloque '{'
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// Verificamos si existe un bloque 'sino'
	if p.peekTokenIs("SINO") {
		p.nextToken() // Avanza a 'sino'

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curTok}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokIs(lexer.RBRACE) && !p.curTokIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) curTokIs(t lexer.TokenType) bool {
	return p.curTok.Type == t
}
