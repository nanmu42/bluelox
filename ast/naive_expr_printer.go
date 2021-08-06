package ast

import (
	"fmt"
	"strings"
)

// NaiveExprPrinter prints expr in polish notation.
type NaiveExprPrinter struct {
	StubExprVisitor
}

func (p *NaiveExprPrinter) VisitBinaryExpr(v *BinaryExpr) (result interface{}, err error) {
	return p.parenthesize(v.Operator.Lexeme, v.Left, v.Right), nil
}

func (p *NaiveExprPrinter) VisitGroupingExpr(v *GroupingExpr) (result interface{}, err error) {
	return p.parenthesize("group", v.Expr), nil
}

func (p *NaiveExprPrinter) VisitLiteralExpr(v *LiteralExpr) (result interface{}, err error) {
	if v.Value == nil {
		return "nil", err
	}
	return fmt.Sprintf("%v", v.Value), nil
}

func (p *NaiveExprPrinter) VisitUnaryExpr(v *UnaryExpr) (result interface{}, err error) {
	return p.parenthesize(v.Operator.Lexeme, v.Right), nil
}

func (p *NaiveExprPrinter) parenthesize(operator string, expr ...Expression) string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(operator)
	for _, item := range expr {
		b.WriteString(" ")
		b.WriteString(noErrStringResult(item.Accept(p)))
	}
	b.WriteString(")")

	return b.String()
}
