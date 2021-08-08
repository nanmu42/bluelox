package parser

import (
	"fmt"

	"github.com/nanmu42/bluelox/ast"
	"github.com/nanmu42/bluelox/token"
)

type Parser struct {
	tokens []*token.Token

	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}

// Parse tokens into expressions.
//
// When everything is fine, errs is nil.
// Otherwise, more than one error may appear,
// and expr can not be considered valid.
func (p *Parser) Parse() (expr ast.Expression, errs []error) {
	for !p.isAtEnd() {
		var err error
		expr, err = p.expression()
		if err != nil {
			errs = append(errs, err)
			p.synchronize()
			continue
		}
	}

	return
}

// expression → equality
func (p *Parser) expression() (ast.Expression, error) {
	return p.equality()
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

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
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
