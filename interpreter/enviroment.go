package interpreter

import (
	"fmt"

	"github.com/nanmu42/bluelox/token"
)

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) (err error) {
	e.values[name] = value
	return
}

func (e *Environment) Get(name *token.Token) (value interface{}, err error) {
	value, ok := e.values[name.Lexeme]
	if !ok {
		err = fmt.Errorf("undefined variable %q", name.Lexeme)
		return
	}

	return
}
