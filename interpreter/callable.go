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

type Function struct {
	Declaration   *ast.FunctionStmt
	Closure       *Environment
	IsInitializer bool
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f *Function) Call(interpreter *Interpreter, arguments []interface{}) (result interface{}, err error) {
	env := NewChildEnvironment(f.Closure)
	for i, param := range f.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		if f.IsInitializer {
			var badThis error
			result, badThis = f.Closure.GetAt(0, "this")
			if badThis != nil {
				panic(badThis)
			}
		}
	}()
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

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}

func (f *Function) Bind(i *Instance) *Function {
	env := NewChildEnvironment(f.Closure)
	env.Define("this", i)
	return &Function{
		Declaration:   f.Declaration,
		Closure:       env,
		IsInitializer: f.IsInitializer,
	}
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
