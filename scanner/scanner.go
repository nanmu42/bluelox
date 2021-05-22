package scanner

type Token struct {
}

type Scanner struct {
	source []byte
}

func NewScanner(source []byte) *Scanner {
	return &Scanner{
		source: source,
	}
}

func (s *Scanner) ScanTokens() (tokens []Token, err error) {
	panic("TODO")
}
