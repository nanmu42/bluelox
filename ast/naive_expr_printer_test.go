package ast

import (
	"reflect"
	"testing"

	"github.com/nanmu42/bluelox/token"
)

func TestNaiveExprPrinter(t *testing.T) {
	tests := []struct {
		name       string
		expr       Expression
		wantResult string
		wantErr    bool
	}{
		{
			name: "easy",
			expr: &UnaryExpr{
				Operator: token.Token{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    0,
				},
				Right: &LiteralExpr{Value: 5},
			},
			wantResult: "(! 5)",
			wantErr:    false,
		},
		{
			name: "literal nil",
			expr: &UnaryExpr{
				Operator: token.Token{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    0,
				},
				Right: &LiteralExpr{Value: nil},
			},
			wantResult: "(! nil)",
			wantErr:    false,
		},
		{
			name: "textbook",
			expr: &BinaryExpr{
				Left: &UnaryExpr{
					Operator: token.Token{
						Type:    token.Minus,
						Lexeme:  "-",
						Literal: nil,
						Line:    1,
					},
					Right: &LiteralExpr{Value: 123},
				},
				Operator: token.Token{
					Type:    token.Star,
					Lexeme:  "*",
					Literal: nil,
					Line:    1,
				},
				Right: &GroupingExpr{Expr: &LiteralExpr{Value: 45.67}},
			},
			wantResult: "(* (- 123) (group 45.67))",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &NaiveExprPrinter{}
			gotResult, err := stringResult(tt.expr.Accept(p))
			if (err != nil) != tt.wantErr {
				t.Errorf("VisitBinaryExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("VisitBinaryExpr() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
