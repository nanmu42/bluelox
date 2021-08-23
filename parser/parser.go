package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nanmu42/bluelox/ast"
	"github.com/nanmu42/bluelox/token"
)

type Parser struct {
	tokens []*token.Token

	current int
}

type ParsingErr struct {
	errs []error
}

func (p *ParsingErr) Error() string {
	length := len(p.errs)

	if length == 0 {
		return "parsing error with 0 detail, it's likely that there is a problem in implementation"
	}

	if length == 1 {
		return fmt.Sprintf("parsing error: %s", p.errs[0])
	}

	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "got %d parsing error(s):\n", length)

	for index, err := range p.errs {
		b.WriteString(strconv.Itoa(index+1) + ". ")
		b.WriteString(err.Error())
		b.WriteRune('\n')

		if index >= 9 {
			b.WriteString("too many errors, more contents omitted...")
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}

// Parse tokens into expressions.
//
// When everything is fine, errs is nil.
// Otherwise, more than one error may appear,
// and expr can not be considered valid.
func (p *Parser) Parse() (stmts []ast.Statement, err error) {
	var errs []error

	for !p.isAtEnd() {
		stmt, stmtErr := p.declaration()
		if stmtErr != nil {
			errs = append(errs, stmtErr)
			p.synchronize()
			continue
		}
		stmts = append(stmts, stmt)
	}

	if len(errs) > 0 {
		err = &ParsingErr{errs}
		return
	}

	return
}

// statement → exprStmt | printStmt | block ;
func (p *Parser) statement() (stmt ast.Statement, err error) {
	if p.match(token.Print) {
		return p.printStmt()
	}
	if p.match(token.LeftBrace) {
		var innerStmts []ast.Statement

		innerStmts, err = p.block()
		if err != nil {
			return
		}

		stmt = &ast.BlockStmt{Stmts: innerStmts}
		return
	}

	return p.exprStmt()
}

// expression → assignment
func (p *Parser) expression() (ast.Expression, error) {
	return p.assignment()
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (expr ast.Expression, err error) {
	expr, err = p.comparison()
	if err != nil {
		return
	}

	for p.match(token.BangEqual, token.EqualEqual) {
		var (
			operator = p.previous()
			right    ast.Expression
		)
		right, err = p.comparison()
		if err != nil {
			return
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (expr ast.Expression, err error) {
	expr, err = p.term()
	if err != nil {
		return
	}

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		var (
			operator = p.previous()
			right    ast.Expression
		)
		right, err = p.term()
		if err != nil {
			return
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return
}

func (p *Parser) match(tokenType ...token.Type) bool {
	for _, oneType := range tokenType {
		if p.check(oneType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) check(tokenType token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (expr ast.Expression, err error) {
	expr, err = p.factor()
	if err != nil {
		return
	}

	for p.match(token.Minus, token.Plus) {
		var (
			operator = p.previous()
			right    ast.Expression
		)
		right, err = p.factor()
		if err != nil {
			return
		}
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (expr ast.Expression, err error) {
	expr, err = p.unary()
	if err != nil {
		return
	}

	for p.match(token.Slash, token.Star) {
		var (
			operator = p.previous()
			right    ast.Expression
		)
		right, err = p.unary()
		if err != nil {
			return
		}
		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return
}

// unary → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() (expr ast.Expression, err error) {
	if p.match(token.Bang, token.Minus) {
		var (
			operator = p.previous()
			right    ast.Expression
		)
		right, err = p.unary()
		if err != nil {
			return
		}
		expr = &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}
		return
	}

	return p.primary()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | IDENTIFIER | "(" expression ")" ;
func (p *Parser) primary() (expr ast.Expression, err error) {
	if p.match(token.Number, token.String) {
		expr = &ast.LiteralExpr{Value: p.previous().Literal}
		return
	}

	if p.match(token.True) {
		expr = &ast.LiteralExpr{Value: true}
		return
	}
	if p.match(token.False) {
		expr = &ast.LiteralExpr{Value: false}
		return
	}
	if p.match(token.Nil) {
		expr = &ast.LiteralExpr{Value: nil}
		return
	}
	if p.match(token.Identifier) {
		expr = &ast.VariableExpr{Name: p.previous()}
		return
	}

	if p.match(token.LeftParen) {
		expr, err = p.expression()
		if err != nil {
			return
		}
		_, err = p.consume(token.RightParen)
		if err != nil {
			err = fmt.Errorf("expected ')' after expression: %w", err)
			return
		}

		expr = &ast.GroupingExpr{Expr: expr}
		return
	}

	unexpected := p.peek()
	err = fmt.Errorf("parsing primary: unexpected token %s %q at line %d", unexpected.Type, unexpected.Lexeme, unexpected.Line)
	return
}

func (p *Parser) consume(wantType token.Type) (consumed *token.Token, err error) {
	if p.check(wantType) {
		consumed = p.advance()
		return
	}

	err = fmt.Errorf("want token type %s, got %s", wantType, p.peek().Type)
	return
}

// synchronize discards tokens until it thinks it has found a statement boundary.
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.Semicolon {
			return
		}

		switch p.peek().Type {
		case token.Class, token.Fun, token.Var, token.For,
			token.If, token.While, token.Print, token.Return:
			return
		}

		p.advance()
	}
}

// printStmt → "print" expression ";" ;
func (p *Parser) printStmt() (stmt ast.Statement, err error) {
	value, err := p.expression()
	if err != nil {
		return
	}

	_, err = p.consume(token.Semicolon)
	if err != nil {
		err = fmt.Errorf("expected ';' after value: %w", err)
		return
	}

	stmt = &ast.PrintStmt{Expr: value}
	return
}

// exprStmt → expression ";"
func (p *Parser) exprStmt() (stmt ast.Statement, err error) {
	value, err := p.expression()
	if err != nil {
		return
	}

	_, err = p.consume(token.Semicolon)
	if err != nil {
		err = fmt.Errorf("expected ';' after value: %w", err)
		return
	}

	stmt = &ast.ExprStmt{Expr: value}
	return
}

// declaration → varDecl | statement ;
func (p *Parser) declaration() (stmt ast.Statement, err error) {
	if p.match(token.Var) {
		return p.varDecl()
	}

	return p.statement()
}

// varDecl → "var" IDENTIFIER ( "=" expression )? ";" ;
func (p *Parser) varDecl() (stmt ast.Statement, err error) {
	name, err := p.consume(token.Identifier)
	if err != nil {
		err = fmt.Errorf("expected a variable name: %w", err)
		return
	}

	var initializer ast.Expression
	if p.match(token.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return
		}
	}

	_, err = p.consume(token.Semicolon)
	if err != nil {
		err = fmt.Errorf("expected ';' after variable declaration: %w", err)
		return
	}

	stmt = &ast.VarStmt{
		Name:        name,
		Initializer: initializer,
	}

	return
}

// assignment → IDENTIFIER "=" assignment | equality ;
func (p *Parser) assignment() (expr ast.Expression, err error) {
	expr, err = p.equality()
	if err != nil {
		return
	}

	if !p.match(token.Equal) {
		return
	}

	value, err := p.assignment()
	if err != nil {
		err = fmt.Errorf("parsing r-value: %w", err)
		return
	}
	name, ok := expr.(*ast.VariableExpr)
	if !ok {
		err = errors.New("invalid assignment target")
		return
	}

	expr = &ast.AssignExpr{
		Name:  name.Name,
		Value: value,
	}

	return
}

func (p *Parser) block() (stmts []ast.Statement, err error) {
	for !p.check(token.RightBrace) {
		var innerStmt ast.Statement
		innerStmt, err = p.declaration()
		if err != nil {
			return
		}

		stmts = append(stmts, innerStmt)
	}

	_, err = p.consume(token.RightBrace)
	if err != nil {
		err = fmt.Errorf("expected '}' after block: %w", err)
		return
	}

	return
}
