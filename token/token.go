package token

type Token struct {
	Type    Type
	Lexeme  string
	Literal interface{}
	Line    int
}
