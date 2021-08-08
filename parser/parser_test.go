package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nanmu42/bluelox/ast"
	"github.com/nanmu42/bluelox/token"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []*token.Token
		wantExpr ast.Expression
		wantErrs []error
	}{
		{
			name: "empty",
			tokens: []*token.Token{
				{
					Type:    token.EOF,
					Lexeme:  "",
					Literal: nil,
					Line:    1,
				},
			},
			wantExpr: nil,
			wantErrs: nil,
		},
		{
			name: "easy",
			tokens: []*token.Token{
				{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.Number,
					Lexeme:  "5",
					Literal: 5,
					Line:    1,
				},
				{
					Type:    token.EOF,
					Lexeme:  "",
					Literal: nil,
					Line:    1,
				},
			},
			wantExpr: &ast.UnaryExpr{
				Operator: &token.Token{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.LiteralExpr{Value: 5},
			},
			wantErrs: nil,
		},
		{
			name: "literal nil",
			tokens: []*token.Token{
				{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.Nil,
					Lexeme:  "nil",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.EOF,
					Lexeme:  "",
					Literal: nil,
					Line:    1,
				},
			},
			wantExpr: &ast.UnaryExpr{
				Operator: &token.Token{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.LiteralExpr{Value: nil},
			},
			wantErrs: nil,
		},
		{
			name: "literal nil",
			tokens: []*token.Token{
				{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.Nil,
					Lexeme:  "nil",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.EOF,
					Lexeme:  "",
					Literal: nil,
					Line:    1,
				},
			},
			wantExpr: &ast.UnaryExpr{
				Operator: &token.Token{
					Type:    token.Bang,
					Lexeme:  "!",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.LiteralExpr{Value: nil},
			},
			wantErrs: nil,
		},
		{
			name: "textbook",
			tokens: []*token.Token{
				{
					Type:    token.Minus,
					Lexeme:  "-",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.Number,
					Lexeme:  "123",
					Literal: 123,
					Line:    1,
				},
				{
					Type:    token.Star,
					Lexeme:  "*",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.LeftParen,
					Lexeme:  "(",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.Number,
					Lexeme:  "45.67",
					Literal: 45.67,
					Line:    1,
				},
				{
					Type:    token.RightParen,
					Lexeme:  ")",
					Literal: nil,
					Line:    1,
				},
				{
					Type:    token.EOF,
					Lexeme:  "",
					Literal: nil,
					Line:    1,
				},
			},
			wantExpr: &ast.BinaryExpr{
				Left: &ast.UnaryExpr{
					Operator: &token.Token{
						Type:    token.Minus,
						Lexeme:  "-",
						Literal: nil,
						Line:    1,
					},
					Right: &ast.LiteralExpr{Value: 123},
				},
				Operator: &token.Token{
					Type:    token.Star,
					Lexeme:  "*",
					Literal: nil,
					Line:    1,
				},
				Right: &ast.GroupingExpr{Expr: &ast.LiteralExpr{Value: 45.67}},
			},
			wantErrs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.tokens)
			gotExpr, gotErrs := p.Parse()
			if !reflect.DeepEqual(gotExpr, tt.wantExpr) {
				t.Errorf("Parse() gotExpr = \n%v\n, want \n%v\n", jsonify(gotExpr), jsonify(tt.wantExpr))
			}
			if !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("Parse() gotErrs = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func jsonify(v interface{}) string {
	marshaled, err := json.Marshal(v)
	if err != nil {
		err = fmt.Errorf("jsonify: %w", err)
		panic(err)
	}

	return string(marshaled)
}
