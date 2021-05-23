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
