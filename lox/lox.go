package lox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nanmu42/bluelox/scanner"
)

type Lox struct {
}

func NewLox() *Lox {
	return &Lox{}
}

func (l *Lox) RunFile(path string) (err error) {
	script, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("reading script file: %w", err)
		return
	}

	err = l.run(script)
	if err != nil {
		err = fmt.Errorf("running script: %w", err)
		return
	}

	return
}

func (l *Lox) RunPrompt() (err error) {
	lineReader := bufio.NewScanner(os.Stdin)

	fmt.Printf("> ")
	for lineReader.Scan() {
		fmt.Printf("> ")
		line := lineReader.Bytes()
		if len(line) == 0 {
			break
		}

		err = l.run(line)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = lineReader.Err()
	if err != nil {
		err = fmt.Errorf("reading input: %w", err)
		return
	}

	return
}

// run provided script.
//
// The provided script is read only, should not be modified.
func (l *Lox) run(script []byte) (err error) {
	s := scanner.NewScanner(script)
	tokens, err := s.ScanTokens()
	if err != nil {
		err = fmt.Errorf("scaning tokens: %w", err)
		return
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return
}
