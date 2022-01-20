package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nanmu42/bluelox/ast"
	"github.com/nanmu42/bluelox/token"
)

const maxFuncArgCounts = 255

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

// statement → exprStmt
// | forStmt
// | ifStmt
// | printStmt
// | returnStmt
// | whileStmt
// | block ;
func (p *Parser) statement() (stmt ast.Statement, err error) {
	if p.match(token.For) {
		return p.forStmt()
	}
	if p.match(token.If) {
		return p.ifStmt()
	}
	if p.match(token.Print) {
		return p.printStmt()
	}
	if p.match(token.Return) {
		return p.returnStmt()
	}
	if p.match(token.While) {
		return p.whileStmt()
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

// unary → ( "!" | "-" ) unary | call ;
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

	return p.call()
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
	if p.match(token.This) {
		expr = &ast.ThisExpr{Keyword: p.previous()}
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

// declaration → classDecl | funDecl | varDecl | statement ;
func (p *Parser) declaration() (stmt ast.Statement, err error) {
	if p.match(token.Class) {
		return p.classDecl()
	}

	if p.match(token.Fun) {
		return p.function("function")
	}

	if p.match(token.Var) {
		return p.varDecl()
	}

	return p.statement()
}

// function → IDENTIFIER "(" parameters? ")" block ;
func (p *Parser) function(kind string) (stmt ast.Statement, err error) {
	name, err := p.consume(token.Identifier)
	if err != nil {
		err = fmt.Errorf("expected %s name: %w", kind, err)
		return
	}

	_, err = p.consume(token.LeftParen)
	if err != nil {
		err = fmt.Errorf("expected '(' after %s name: %w", kind, err)
		return
	}

	var parameters []*token.Token
	if !p.check(token.RightParen) {
		var firstParam *token.Token
		firstParam, err = p.consume(token.Identifier)
		if err != nil {
			err = fmt.Errorf("expected parameter name: %w", err)
			return
		}
		parameters = append(parameters, firstParam)

		for p.match(token.Comma) {
			if len(parameters) >= maxFuncArgCounts {
				err = fmt.Errorf("can't have more than %d parameters", maxFuncArgCounts)
				return
			}

			var param *token.Token
			param, err = p.consume(token.Identifier)
			if err != nil {
				err = fmt.Errorf("expected parameter name: %w", err)
				return
			}
			parameters = append(parameters, param)
		}
	}

	_, err = p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expected ')' after parameters: %w", err)
		return
	}

	_, err = p.consume(token.LeftBrace)
	if err != nil {
		err = fmt.Errorf("expected '{' before %s body: %w", kind, err)
		return
	}

	body, err := p.block()
	if err != nil {
		return
	}

	stmt = &ast.FunctionStmt{
		Name:   name,
		Params: parameters,
		Body:   body,
	}
	return
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

// assignment → ( call "." )? IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() (expr ast.Expression, err error) {
	expr, err = p.or()
	if err != nil {
		return
	}

	if !p.match(token.Equal) {
		return
	}

	var value ast.Expression
	value, err = p.assignment()
	if err != nil {
		err = fmt.Errorf("parsing r-value: %w", err)
		return
	}

	if get, ok := expr.(*ast.GetExpr); ok {
		expr = &ast.SetExpr{
			Object: get.Object,
			Name:   get.Name,
			Value:  value,
		}
	} else {
		name, ok := expr.(*ast.VariableExpr)
		if !ok {
			err = errors.New("invalid assignment target")
			return
		}

		expr = &ast.AssignExpr{
			Name:  name.Name,
			Value: value,
		}
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

func (p *Parser) ifStmt() (stmt ast.Statement, err error) {
	_, err = p.consume(token.LeftParen)
	if err != nil {
		err = fmt.Errorf("expected '(' after 'if': %w", err)
		return
	}
	condition, err := p.expression()
	if err != nil {
		return
	}
	_, err = p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expected ')' after 'if' condition: %w", err)
		return
	}

	thenBranch, err := p.statement()
	if err != nil {
		return
	}
	var elseBranch ast.Statement
	if p.match(token.Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return
		}
	}

	stmt = &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}

	return
}

// or → and ( "or" and )* ;
func (p *Parser) or() (expr ast.Expression, err error) {
	expr, err = p.and()
	if err != nil {
		return
	}

	for p.match(token.Or) {
		operator := p.previous()

		var right ast.Expression
		right, err = p.and()
		if err != nil {
			return
		}

		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return
}

// and → equality ( "and" equality )* ;
func (p *Parser) and() (expr ast.Expression, err error) {
	expr, err = p.equality()
	if err != nil {
		return
	}

	for p.match(token.And) {
		operator := p.previous()

		var right ast.Expression
		right, err = p.equality()
		if err != nil {
			return
		}

		expr = &ast.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return
}

// whileStmt      → "while" "(" expression ")" statement ;
func (p *Parser) whileStmt() (stmt ast.Statement, err error) {
	_, err = p.consume(token.LeftParen)
	if err != nil {
		err = fmt.Errorf("expected '(' after 'while': %w", err)
		return
	}

	condition, err := p.expression()
	if err != nil {
		return
	}

	_, err = p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expected ')' after 'while': %w", err)
		return
	}

	body, err := p.statement()
	if err != nil {
		return
	}

	stmt = &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}
	return
}

// forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
//                 expression? ";"
//                 expression? ")" statement ;
func (p *Parser) forStmt() (stmt ast.Statement, err error) {
	_, err = p.consume(token.LeftParen)
	if err != nil {
		err = fmt.Errorf("expected '(' after 'for': %w", err)
		return
	}

	var initializer ast.Statement
	if p.match(token.Semicolon) {
		initializer = nil
	} else if p.match(token.Var) {
		initializer, err = p.varDecl()
		if err != nil {
			return
		}
	} else {
		initializer, err = p.exprStmt()
		if err != nil {
			return
		}
	}

	var condition ast.Expression
	if !p.check(token.Semicolon) {
		condition, err = p.expression()
		if err != nil {
			return
		}
	} else {
		condition = &ast.LiteralExpr{Value: true}
	}
	_, err = p.consume(token.Semicolon)
	if err != nil {
		err = fmt.Errorf("expected ';' after loop condition: %w", err)
		return
	}

	var increment ast.Expression
	if !p.check(token.RightParen) {
		increment, err = p.expression()
		if err != nil {
			return
		}
	}

	_, err = p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expected ')' after 'for' clause: %w", err)
		return
	}

	body, err := p.statement()
	if err != nil {
		return
	}
	if increment != nil {
		body = &ast.BlockStmt{
			Stmts: []ast.Statement{
				body,
				&ast.ExprStmt{Expr: increment},
			},
		}
	}

	stmt = &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		stmt = &ast.BlockStmt{
			Stmts: []ast.Statement{
				initializer,
				stmt,
			},
		}
	}

	return
}

// call  → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
func (p *Parser) call() (expr ast.Expression, err error) {
	expr, err = p.primary()
	if err != nil {
		return
	}

	for {
		if p.match(token.LeftParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return
			}
		} else if p.match(token.Dot) {
			var name *token.Token
			name, err = p.consume(token.Identifier)
			if err != nil {
				err = fmt.Errorf("expected property name after '.': %w", err)
				return
			}
			expr = &ast.GetExpr{
				Object: expr,
				Name:   name,
			}
		} else {
			break
		}
	}

	return
}

func (p *Parser) finishCall(callee ast.Expression) (expr ast.Expression, err error) {
	var arguments []ast.Expression
	if !p.check(token.RightParen) {
		var firstArg ast.Expression
		firstArg, err = p.expression()
		if err != nil {
			err = fmt.Errorf("parsing first arg: %w", err)
			return
		}
		arguments = append(arguments, firstArg)

		argSeq := 1
		for p.match(token.Comma) {
			argSeq++

			var arg ast.Expression
			arg, err = p.expression()
			if err != nil {
				err = fmt.Errorf("parsing arg #%d: %w", argSeq, err)
				return
			}
			if argSeq >= maxFuncArgCounts {
				err = fmt.Errorf("can't have more than %d arguments", maxFuncArgCounts)
				return
			}

			arguments = append(arguments, arg)
		}
	}

	paren, err := p.consume(token.RightParen)
	if err != nil {
		err = fmt.Errorf("expect ')' after arguments: %w", err)
		return
	}

	expr = &ast.CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
	return
}

func (p *Parser) returnStmt() (stmt ast.Statement, err error) {
	keyword := p.previous()
	var value ast.Expression
	if !p.check(token.Semicolon) {
		value, err = p.expression()
		if err != nil {
			return
		}
	}

	_, err = p.consume(token.Semicolon)
	if err != nil {
		err = fmt.Errorf("expected ';' after return value: %w", err)
		return
	}
	stmt = &ast.ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}

	return
}

// classDecl  → "class" IDENTIFIER "{" function* "}" ;
func (p *Parser) classDecl() (stmt ast.Statement, err error) {
	name, err := p.consume(token.Identifier)
	if err != nil {
		err = fmt.Errorf("expected class name: %w", err)
		return
	}

	_, err = p.consume(token.LeftBrace)
	if err != nil {
		err = fmt.Errorf("expected '{' before class body: %w", err)
		return
	}
	var methods []*ast.FunctionStmt
	for !p.isAtEnd() && !p.check(token.RightBrace) {
		var method ast.Statement
		method, err = p.function("method")
		if err != nil {
			return
		}
		methods = append(methods, method.(*ast.FunctionStmt))
	}

	_, err = p.consume(token.RightBrace)
	if err != nil {
		err = fmt.Errorf("expected '}' after class body: %w", err)
		return
	}

	stmt = &ast.ClassStmt{
		Name:    name,
		Methods: methods,
	}

	return
}
