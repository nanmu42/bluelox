package scanner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nanmu42/bluelox/token"
)

// st stands for simple token
func st(tokenType token.Type, lexme string, line int) *token.Token {
	return &token.Token{
		Type:    tokenType,
		Lexeme:  lexme,
		Literal: nil,
		Line:    line,
	}
}

func TestScanner_ScanTokens(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		wantTokens []*token.Token
		wantErr    bool
	}{
		{
			name:   "empty",
			source: "",
			wantTokens: []*token.Token{
				st(token.EOF, "", 1),
			},
			wantErr: false,
		},
		{
			name: "easy",
			source: `// this is a comment
(( )){} // grouping stuff
!*+-/=<> <= == // operators`,
			wantTokens: []*token.Token{
				st(token.LeftParen, "(", 2),
				st(token.LeftParen, "(", 2),
				st(token.RightParen, ")", 2),
				st(token.RightParen, ")", 2),
				st(token.LeftBrace, "{", 2),
				st(token.RightBrace, "}", 2),

				st(token.Bang, "!", 3),
				st(token.Star, "*", 3),
				st(token.Plus, "+", 3),
				st(token.Minus, "-", 3),
				st(token.Slash, "/", 3),
				st(token.Equal, "=", 3),
				st(token.Less, "<", 3),
				st(token.Greater, ">", 3),
				st(token.LessEqual, "<=", 3),
				st(token.EqualEqual, "==", 3),

				st(token.EOF, "", 3),
			},
			wantErr: false,
		},
		{
			name: "some code",
			source: `var alibaba = blaster175 + 3.14
if 6 >= k { // hey!
print "hello world"
}
`,
			wantTokens: []*token.Token{
				st(token.Var, "var", 1),
				st(token.Identifier, "alibaba", 1),
				st(token.Equal, "=", 1),
				st(token.Identifier, "blaster175", 1),
				st(token.Plus, "+", 1),
				{
					Type:    token.Number,
					Lexeme:  "3.14",
					Literal: 3.14,
					Line:    1,
				},

				st(token.If, "if", 2),
				{
					Type:    token.Number,
					Lexeme:  "6",
					Literal: float64(6),
					Line:    2,
				},
				st(token.GreaterEqual, ">=", 2),
				st(token.Identifier, "k", 2),
				st(token.LeftBrace, "{", 2),

				st(token.Print, "print", 3),
				{
					Type:    token.String,
					Lexeme:  `"hello world"`,
					Literal: "hello world",
					Line:    3,
				},

				st(token.RightBrace, "}", 4),

				st(token.EOF, "", 5),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner([]byte(tt.source))
			gotTokens, err := s.ScanTokens()
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.wantTokens, gotTokens)
		})
	}
}
