package interpreter

import (
	"fmt"
	"time"

	"github.com/nanmu42/bluelox/ast"
)

const nativeFuncStringForm = "<native fn>"

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []interface{}) (result interface{}, err error)
	String() string
}

type function struct {
	Declaration *ast.FunctionStmt
	Closure     *Environment
}

func (f *function) Arity() int {
	return len(f.Declaration.Params)
}

func (f *function) Call(interpreter *Interpreter, arguments []interface{}) (result interface{}, err error) {
	env := NewChildEnvironment(f.Closure)
	for i, param := range f.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		panicReason := recover()
		if panicReason == nil {
			return
		}
		returnValue, ok := panicReason.(*returnPayload)
		if !ok {
			panic(panicReason)
		}

		result = returnValue.Value
	}()

	err = interpreter.executeBlock(f.Declaration.Body, env)
	return
}

func (f *function) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}

type returnPayload struct {
	Value interface{}
}

type nativeFuncClock struct{}

func (n nativeFuncClock) Arity() int {
	return 0
}

func (n nativeFuncClock) Call(interpreter *Interpreter, arguments []interface{}) (result interface{}, err error) {
	return float64(time.Now().Unix()), nil
}

func (n nativeFuncClock) String() string {
	return nativeFuncStringForm
}
