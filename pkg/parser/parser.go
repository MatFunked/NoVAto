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
	ASSIGNMENT  // =
	EQUALS      // ==
	LESSGREATER // > o <, >=, <=
	SUM         // + o -
	PRODUCT     // * o /
	PREFIX      // -X o !X
	CALL        // myFunction(X)
)

var precedences = map[lexer.TokenType]int{
	lexer.ASSIGN:   ASSIGNMENT,
	lexer.EQ:       EQUALS,
	lexer.NOT_EQ:   EQUALS,
	lexer.LTE:      LESSGREATER,
	lexer.GTE:      LESSGREATER,
	lexer.LT:       LESSGREATER,
	lexer.RT:       LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.ASTERISK: PRODUCT,
	lexer.SLASH:    PRODUCT,
	lexer.LPAREN:   CALL,
	"and":          EQUALS, // Precedencia para operador lógico
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
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix("VERDADERO", p.parseBoolean)
	p.registerPrefix("FALSO", p.parseBoolean)
	p.registerPrefix("FUNC", p.parseFunctionLiteral)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LTE, p.parseInfixExpression)
	p.registerInfix(lexer.GTE, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.RT, p.parseInfixExpression)
	p.registerInfix("and", p.parseInfixExpression)

	p.registerPrefix("SI", p.parseIfExpression)
	p.registerPrefix("OSI", p.parseIfExpression)
	p.registerPrefix("COMPARAR", p.parseCompararExpression)
	p.registerPrefix("MIENTRAS", p.parseWhileExpression)
	p.registerPrefix("PARA", p.parseForExpression)

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
		if p.curTok.Literal == "\n" || p.curTok.Type == "NL" {
			p.nextToken()
			continue
		}

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
	case "SI":
		stmt := &ast.ExpressionStatement{Token: p.curTok}
		stmt.Expression = p.parseExpression(LOWEST)
		return stmt
	case "ROMPER":
		return p.parseBreakStatement()
	case "CONTINUAR":
		return p.parseContinueStatement()
	case "CEDER":
		return p.parseYieldStatement()
	case "RETORNAR":
		return p.parseReturnStatement()
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

	if isTypeToken(p.peekTok.Type) {
		p.nextToken()
		stmt.Type = p.curTok.Literal
	}

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

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return nil
	}

	expression := &ast.AssignmentExpression{
		Token: p.curTok,
		Name:  ident,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Value = p.parseExpression(precedence - 1)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
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

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curTok}

	value, err := strconv.ParseFloat(p.curTok.Literal, 64)
	if err != nil {
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curTok, Value: p.curTok.Literal}
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

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs("OSI") {
		// 'osi' encadena otro IfExpression dentro del bloque alternativo (elif)
		p.nextToken()
		osiToken := p.curTok
		nested := p.parseIfExpression()
		if nested == nil {
			return nil
		}

		expression.Alternative = &ast.BlockStatement{
			Token: osiToken,
			Statements: []ast.Statement{
				&ast.ExpressionStatement{Token: osiToken, Expression: nested},
			},
		}
	} else if p.peekTokenIs("SINO") {
		p.nextToken()

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseCompararExpression() ast.Expression {
	expression := &ast.SwitchExpression{Token: p.curTok}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)
	if expression.Value == nil {
		return nil
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken()

	for !p.curTokIs(lexer.RBRACE) && !p.curTokIs(lexer.EOF) {
		switch p.curTok.Type {
		case "ENCASO":
			clause := p.parseCaseClause()
			if clause == nil {
				return nil
			}
			expression.Cases = append(expression.Cases, clause)

		case "SINO", "DEFECTO":
			if expression.Default != nil {
				return nil
			}
			if !p.expectPeek(lexer.LBRACE) {
				return nil
			}
			expression.Default = p.parseBlockStatement()

		default:
			return nil
		}

		p.nextToken()
	}

	return expression
}

func (p *Parser) parseCaseClause() *ast.CaseClause {
	clause := &ast.CaseClause{Token: p.curTok}

	p.nextToken()
	value := p.parseExpression(LOWEST)
	if value == nil {
		return nil
	}
	clause.Values = append(clause.Values, value)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume la coma
		p.nextToken() // avanza al siguiente valor
		value = p.parseExpression(LOWEST)
		if value == nil {
			return nil
		}
		clause.Values = append(clause.Values, value)
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	clause.Body = p.parseBlockStatement()

	return clause
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curTok}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curTok}

	p.nextToken()

	init, ok := p.parseVariableStatement().(*ast.VariableStatement)
	if !ok {
		return nil
	}
	expression.Init = init

	if !p.expectPeek("HASTA") {
		return nil
	}

	p.nextToken()
	limit := p.parseExpression(LOWEST)
	if limit == nil {
		return nil
	}

	expression.Condition = &ast.InfixExpression{
		Token:    lexer.Token{Type: lexer.LT, Literal: "<"},
		Left:     init.Name,
		Operator: "<",
		Right:    limit,
	}

	expression.Increment = &ast.AssignmentExpression{
		Token: lexer.Token{Type: lexer.ASSIGN, Literal: "="},
		Name:  init.Name,
		Value: &ast.InfixExpression{
			Token:    lexer.Token{Type: lexer.PLUS, Literal: "+"},
			Left:     init.Name,
			Operator: "+",
			Right:    &ast.IntegerLiteral{Token: lexer.Token{Type: lexer.INT, Literal: "1"}, Value: 1},
		},
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseBreakStatement() ast.Statement {
	return &ast.BreakStatement{Token: p.curTok}
}

func (p *Parser) parseContinueStatement() ast.Statement {
	return &ast.ContinueStatement{Token: p.curTok}
}

func (p *Parser) parseYieldStatement() ast.Statement {
	stmt := &ast.YieldStatement{Token: p.curTok}

	if p.peekStartsStatement() {
		return stmt
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curTok}

	if p.peekStartsStatement() {
		return stmt
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curTok}

	if p.peekTokenIs(lexer.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal}
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	lit.Parameters, lit.ParameterTypes = p.parseFunctionParameters()
	if lit.Parameters == nil {
		return nil
	}

	if isTypeToken(p.peekTok.Type) {
		p.nextToken()
		lit.ReturnType = p.curTok.Literal
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() ([]*ast.Identifier, []string) {
	identifiers := []*ast.Identifier{}
	types := []string{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return identifiers, types
	}

	p.nextToken()
	identifiers = append(identifiers, &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal})
	types = append(types, p.parseOptionalParameterType())

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, &ast.Identifier{Token: p.curTok, Value: p.curTok.Literal})
		types = append(types, p.parseOptionalParameterType())
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil, nil
	}

	return identifiers, types
}

func (p *Parser) parseOptionalParameterType() string {
	if isTypeToken(p.peekTok.Type) {
		p.nextToken()
		return p.curTok.Literal
	}
	return ""
}

func isTypeToken(t lexer.TokenType) bool {
	switch t {
	case "T_ENT", "T_TEXT", "T_BOOL", "T_DECI":
		return true
	}
	return false
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: p.curTok, Function: function}
	expression.Arguments = p.parseExpressionList(lexer.RPAREN)

	return expression
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) peekStartsStatement() bool {
	switch p.peekTok.Type {
	case lexer.EOF, lexer.RBRACE, "VARIABLE", "CONSTANTE", "SI", "OSI", "SINO", "MIENTRAS", "PARA", "ROMPER", "CONTINUAR", "CEDER", "RETORNAR":
		return true
	}
	return false
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
