package interpreter

import (
	"fmt"

	"github.com/nanmu42/bluelox/token"
)

type Class struct {
	Name       string
	SuperClass *Class
	Methods    map[string]*Function
}

func (c *Class) String() string {
	return c.Name
}

func (c *Class) Arity() int {
	initializer, ok := c.FindMethod("init")
	if !ok {
		return 0
	}
	return initializer.Arity()
}

func (c *Class) Call(interpreter *Interpreter, arguments []interface{}) (result interface{}, err error) {
	instance := NewInstance(c)

	initializer, ok := c.FindMethod("init")
	if ok {
		_, err = initializer.Bind(instance).Call(interpreter, arguments)
		if err != nil {
			return
		}
	}

	return instance, nil
}

func (c *Class) FindMethod(name string) (method *Function, ok bool) {
	method, ok = c.Methods[name]
	if !ok && c.SuperClass != nil {
		method, ok = c.SuperClass.FindMethod(name)
	}
	return
}

type Instance struct {
	class  *Class
	fields map[string]interface{}
}

func (i Instance) String() string {
	return i.class.Name + " instance"
}

func (i *Instance) Get(name *token.Token) (property interface{}, err error) {
	var ok bool
	property, ok = i.fields[name.Lexeme]
	if ok {
		return
	}
	property, ok = i.class.FindMethod(name.Lexeme)
	if ok {
		method := property.(*Function)
		return method.Bind(i), nil
	}

	err = fmt.Errorf("undefiend property %q", name.Lexeme)
	return
}

func (i *Instance) Set(name *token.Token, result interface{}) {
	i.fields[name.Lexeme] = result
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}
