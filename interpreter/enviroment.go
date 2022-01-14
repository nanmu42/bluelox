package interpreter

import (
	"fmt"

	"github.com/nanmu42/bluelox/token"
)

type Environment struct {
	values map[string]interface{}
	parent *Environment
}

func NewGlobalEnvironment() (env *Environment) {
	env = &Environment{
		values: make(map[string]interface{}),
	}

	// native functions
	env.Define("clock", nativeFuncClock{})

	return
}

func NewChildEnvironment(parent *Environment) *Environment {
	return &Environment{
		values: make(map[string]interface{}),
		parent: parent,
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name *token.Token) (value interface{}, err error) {
	value, ok := e.values[name.Lexeme]
	if ok {
		return
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	err = fmt.Errorf("undefined variable %q", name.Lexeme)
	return
}

func (e *Environment) Assign(name *token.Token, value interface{}) (err error) {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return
	}

	if e.parent != nil {
		return e.parent.Assign(name, value)
	}

	err = fmt.Errorf("can not assign undecleared variable %q", name.Lexeme)
	return
}
