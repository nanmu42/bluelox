package interpreter

import (
	"math"
	"reflect"
	"testing"

	"github.com/nanmu42/bluelox/token"

	"github.com/nanmu42/bluelox/ast"
)

// TODO: 6(3) now outputs 3, which is odd
func TestInterpreter_evaluate(t *testing.T) {
	tests := []struct {
		name       string
		expr       ast.Expression
		wantResult interface{}
		wantErr    bool
	}{
		{
			name: "!5",
			expr: &ast.UnaryExpr{
				Operator: &token.Token{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.LiteralExpr{Value: 5},
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name: "!nil",
			expr: &ast.UnaryExpr{
				Operator: &token.Token{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.LiteralExpr{Value: nil},
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "neglect a muffin",
			expr: &ast.UnaryExpr{
				Operator: &token.Token{
					Type:    token.Minus,
					Lexeme:  "-",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.LiteralExpr{Value: "muffin"},
			},
			wantResult: nil,
			wantErr:    true,
		},
		{
			name: "textbook",
			expr: &ast.BinaryExpr{
				Left: &ast.UnaryExpr{
					Operator: &token.Token{
						Type:    token.Minus,
						Lexeme:  "-",
						Literal: nil,
						Line:    1,
					},
					Right: &ast.LiteralExpr{Value: float64(123)},
				},
				Operator: &token.Token{
					Type:    token.Star,
					Lexeme:  "*",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.GroupingExpr{Expr: &ast.LiteralExpr{Value: 45.67}},
			},
			wantResult: -5617.41,
			wantErr:    false,
		},
		{
			name: "java quirk",
			expr: &ast.BinaryExpr{
				Left: &ast.LiteralExpr{Value: math.NaN()},
				Operator: &token.Token{
					Type:    token.EqualEqual,
					Lexeme:  "==",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.BinaryExpr{
					Left: &ast.LiteralExpr{Value: float64(0)},
					Operator: &token.Token{
						Type:    token.Slash,
						Lexeme:  "/",
						Literal: nil,
						Line:    1,
					},
					Right: &ast.LiteralExpr{Value: float64(0)},
				},
			},
			wantResult: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interpreter{}
			gotResult, err := i.evaluate(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("evaluate() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
