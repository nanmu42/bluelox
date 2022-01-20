package token

import "fmt"

type Token struct {
	Type    Type
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("lexme %q with literal %s, type %s at line %d", t.Lexeme, t.Literal, t.Type, t.Line)
}
