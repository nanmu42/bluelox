package token

type Type int

const (
	SingleCharacterTokenStart Type = iota

	LeftParen
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	SingleCharacterTokenEnd

	OneOrTwoCharacterTokenStart

	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	OneOrTwoCharacterTokenEnd

	LiteralStart

	Identifier
	String
	Number

	LiteralEnd

	KeywordStart

	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	KeywordEnd

	EOF
)

var Str2Keyword = map[string]Type{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"for":    For,
	"fun":    Fun,
	"if":     If,
	"nil":    Nil,
	"or":     Nil,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

var Keyword2Str = make(map[Type]string, len(Str2Keyword))

func init() {
	for s, k := range Str2Keyword {
		Keyword2Str[k] = s
	}
}
