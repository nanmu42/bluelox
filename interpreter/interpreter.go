package interpreter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"sync"

	"github.com/nanmu42/bluelox/ast"
	"github.com/nanmu42/bluelox/token"
)

var _ ast.ExprVisitor = (*Interpreter)(nil)
var _ ast.StmtVisitor = (*Interpreter)(nil)

type Interpreter struct {
	environment *Environment
	globals     *Environment
	// keys are all pointers, so it's fine if we stick with one interpreter.
	locals map[ast.Expression]int

	// protect stdout
	stdoutMu sync.RWMutex
	stdout   io.Writer
}

func NewInterpreter(stdout io.Writer) *Interpreter {
	globals := NewGlobalEnvironment()

	return &Interpreter{
		environment: globals,
		globals:     globals,
		locals:      make(map[ast.Expression]int),
		stdoutMu:    sync.RWMutex{},
		stdout:      stdout,
	}
}

func (i *Interpreter) ChangeStdoutTo(w io.Writer) {
	i.stdoutMu.Lock()
	i.stdout = w
	i.stdoutMu.Unlock()
}

func (i *Interpreter) Interpret(ctx context.Context, stmts []ast.Statement) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &RuntimeError{
				Reason: fmt.Sprintf("runtime panicking: \n%v\n", r),
			}
		}
	}()

	done := ctx.Done()
	for _, stmt := range stmts {
		select {
		case <-done:
			err = ctx.Err()
		default:
			// relax
		}
		err = i.execute(stmt)
		if err != nil {
			return
		}
	}

	return
}

func (i *Interpreter) evaluate(expr ast.Expression) (result interface{}, err error) {
	return expr.Accept(i)
}

func (i *Interpreter) ensureNumber(operator *token.Token, operand ...interface{}) (err error) {
	for _, item := range operand {
		if _, ok := item.(float64); !ok {
			err = &RuntimeError{
				Reason: fmt.Sprintf("operand(s) must be number(s), got %q(type %T)", item, item),
				Token:  operator,
			}
			return
		}
	}

	return
}

func (i *Interpreter) VisitBinaryExpr(v *ast.BinaryExpr) (result interface{}, err error) {
	left, err := i.evaluate(v.Left)
	if err != nil {
		return
	}
	right, err := i.evaluate(v.Right)
	if err != nil {
		return
	}

	switch v.Operator.Type {
	case token.Greater:
		err = i.ensureNumber(v.Operator, left, right)
		if err != nil {
			return
		}
		result = left.(float64) > right.(float64)
		return
	case token.GreaterEqual:
		err = i.ensureNumber(v.Operator, left, right)
		if err != nil {
			return
		}
		result = left.(float64) >= right.(float64)
		return
	case token.Less:
		err = i.ensureNumber(v.Operator, left, right)
		if err != nil {
			return
		}
		result = left.(float64) < right.(float64)
		return
	case token.LessEqual:
		err = i.ensureNumber(v.Operator, left, right)
		if err != nil {
			return
		}
		result = left.(float64) <= right.(float64)
		return
	case token.Minus:
		err = i.ensureNumber(v.Operator, left, right)
		if err != nil {
			return
		}
		result = left.(float64) - right.(float64)
		return
	case token.Plus:
		numLeft, okLeft := left.(float64)
		numRight, okRight := right.(float64)
		if okLeft && okRight {
			result = numLeft + numRight
			return
		}

		strLeft, okLeft := left.(string)
		strRight, okRight := right.(string)
		if okLeft && okRight {
			result = strLeft + strRight
			return
		}

		err = &RuntimeError{
			Reason: fmt.Sprintf("operands must be both numbers or strings, got %q(%T) and %q(%T)", left, left, right, right),
			Token:  v.Operator,
		}
		return
	case token.Slash:
		err = i.ensureNumber(v.Operator, left, right)
		if err != nil {
			return
		}
		if right.(float64) == 0 {
			if left.(float64) == 0 {
				result = math.NaN()
				return
			}

			err = &RuntimeError{
				Reason: "division by zero",
				Token:  v.Operator,
			}
			return
		}
		result = left.(float64) / right.(float64)
		return
	case token.Star:
		err = i.ensureNumber(v.Operator, left, right)
		if err != nil {
			return
		}
		result = left.(float64) * right.(float64)
		return
	case token.BangEqual:
		result = !i.isEqual(left, right)
		return
	case token.EqualEqual:
		result = i.isEqual(left, right)
		return
	}

	err = errors.New("evaluating Binary: unreachable code, implementation error")
	return
}

func (i *Interpreter) VisitGroupingExpr(v *ast.GroupingExpr) (result interface{}, err error) {
	return i.evaluate(v.Expr)
}

func (i *Interpreter) VisitLiteralExpr(v *ast.LiteralExpr) (result interface{}, err error) {
	result = v.Value
	return
}

func (i *Interpreter) VisitUnaryExpr(v *ast.UnaryExpr) (result interface{}, err error) {
	right, err := i.evaluate(v.Right)
	if err != nil {
		return
	}

	switch v.Operator.Type {
	case token.Bang:
		result = !i.isTruthy(right)
		return
	case token.Minus:
		err = i.ensureNumber(v.Operator, right)
		if err != nil {
			return
		}
		result = -right.(float64)
		return
	}

	err = errors.New("evaluating Unary: unreachable code, implementation error")
	return
}

// isTruthy false and nil are falsy,
// and everything else is truthy
func (i *Interpreter) isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	boolean, ok := v.(bool)
	if ok && !boolean {
		return false
	}

	return true
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil {
		return b == nil
	}

	switch ta := a.(type) {
	case bool:
		tb, ok := b.(bool)
		if !ok {
			return false
		}
		return ta == tb
	case float64:
		tb, ok := b.(float64)
		if !ok {
			return false
		}
		// so that we are compatible with jlox
		if math.IsNaN(ta) && math.IsNaN(tb) {
			return true
		}
		return ta == tb
	case string:
		tb, ok := b.(string)
		if !ok {
			return false
		}
		return ta == tb
	}

	panic("isEqual: unreachable code, implementation error")
}

func (i *Interpreter) stringify(v interface{}) string {
	if v == nil {
		return "nil"
	}

	if numV, ok := v.(float64); ok {
		return strconv.FormatFloat(numV, 'f', -1, 64)
	}

	return fmt.Sprintf("%v", v)
}

func (i *Interpreter) VisitExprStmt(v *ast.ExprStmt) (err error) {
	_, err = i.evaluate(v.Expr)
	return
}

func (i *Interpreter) VisitPrintStmt(v *ast.PrintStmt) (err error) {
	result, err := i.evaluate(v.Expr)
	if err != nil {
		return
	}

	i.stdoutMu.RLock()
	defer i.stdoutMu.RUnlock()
	_, err = fmt.Fprintln(i.stdout, i.stringify(result))
	if err != nil {
		err = fmt.Errorf("printing to io.Writer: %w", err)
		return
	}
	return
}

func (i *Interpreter) execute(stmt ast.Statement) error {
	return stmt.Accept(i)
}

func (i *Interpreter) Resolve(v ast.Expression, depth int) {
	i.locals[v] = depth
}

func (i *Interpreter) VisitVarStmt(v *ast.VarStmt) (err error) {
	var value interface{}

	if v.Initializer != nil {
		value, err = i.evaluate(v.Initializer)
		if err != nil {
			return
		}
	}

	i.environment.Define(v.Name.Lexeme, value)

	return
}

func (i *Interpreter) VisitVariableExpr(v *ast.VariableExpr) (result interface{}, err error) {
	result, err = i.lookUpVariable(v.Name, v)
	return
}

func (i *Interpreter) VisitAssignExpr(v *ast.AssignExpr) (result interface{}, err error) {
	result, err = i.evaluate(v.Value)
	if err != nil {
		return
	}

	distance, ok := i.locals[v]
	if ok {
		err = i.environment.AssignAt(distance, v.Name, result)
	} else {
		err = i.globals.Assign(v.Name, result)
	}

	return
}

func (i *Interpreter) VisitBlockStmt(v *ast.BlockStmt) (err error) {
	err = i.executeBlock(v.Stmts, NewChildEnvironment(i.environment))
	return
}

func (i *Interpreter) executeBlock(stmts []ast.Statement, blockEnv *Environment) (err error) {
	var previousEnv = i.environment
	i.environment = blockEnv
	defer func() {
		i.environment = previousEnv
	}()

	for _, stmt := range stmts {
		err = i.execute(stmt)
		if err != nil {
			return
		}
	}

	return
}

func (i *Interpreter) VisitIfStmt(v *ast.IfStmt) (err error) {
	evalCondition, err := i.evaluate(v.Condition)
	if err != nil {
		return
	}

	if i.isTruthy(evalCondition) {
		err = i.execute(v.ThenBranch)
		return
	}

	if v.ElseBranch == nil {
		return
	}

	err = i.execute(v.ElseBranch)
	return
}

func (i *Interpreter) VisitWhileStmt(v *ast.WhileStmt) (err error) {
	var evalCondition interface{}
	for {
		evalCondition, err = i.evaluate(v.Condition)
		if err != nil {
			return
		}
		if !i.isTruthy(evalCondition) {
			break
		}

		err = i.execute(v.Body)
		if err != nil {
			return
		}
	}

	return
}

func (i *Interpreter) VisitLogicalExpr(v *ast.LogicalExpr) (result interface{}, err error) {
	result, err = i.evaluate(v.Left)
	if err != nil {
		return
	}

	// short circuit
	if i.isTruthy(result) {
		if v.Operator.Type == token.Or {
			return
		}
	} else {
		if v.Operator.Type == token.And {
			return
		}
	}

	result, err = i.evaluate(v.Right)
	return
}

func (i *Interpreter) VisitSetExpr(v *ast.SetExpr) (result interface{}, err error) {
	object, err := i.evaluate(v.Object)
	if err != nil {
		return
	}

	instance, ok := object.(*Instance)
	if !ok {
		err = fmt.Errorf("only instances have fields, got %T", object)
		return
	}

	result, err = i.evaluate(v.Value)
	if err != nil {
		return
	}
	instance.Set(v.Name, result)

	return
}

func (i *Interpreter) VisitThisExpr(v *ast.ThisExpr) (result interface{}, err error) {
	return i.lookUpVariable(v.Keyword, v)
}

func (i *Interpreter) VisitCallExpr(v *ast.CallExpr) (result interface{}, err error) {
	callee, err := i.evaluate(v.Callee)
	if err != nil {
		return
	}

	var arguments []interface{}
	for _, arg := range v.Arguments {
		var evaledArg interface{}
		evaledArg, err = i.evaluate(arg)
		if err != nil {
			return
		}

		arguments = append(arguments, evaledArg)
	}

	function, ok := callee.(Callable)
	if !ok {
		err = fmt.Errorf("can only call functuins and classes, got %T", callee)
		return
	}
	if want, got := function.Arity(), len(arguments); want != got {
		err = fmt.Errorf("expected %d arguments but got %d", want, got)
		return
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(v *ast.GetExpr) (result interface{}, err error) {
	object, err := i.evaluate(v.Object)
	if err != nil {
		return
	}

	instance, ok := object.(*Instance)
	if !ok {
		err = errors.New("only instances have properties")
		return
	}

	return instance.Get(v.Name)
}

func (i *Interpreter) VisitFunctionStmt(v *ast.FunctionStmt) (err error) {
	i.environment.Define(v.Name.Lexeme, &Function{
		Declaration:   v,
		Closure:       i.environment,
		IsInitializer: false,
	})

	return nil
}

func (i *Interpreter) VisitClassStmt(v *ast.ClassStmt) (err error) {
	var superclass *Class
	if v.SuperClass != nil {
		var rawSuperclass interface{}
		rawSuperclass, err = i.evaluate(v.SuperClass)
		if err != nil {
			return
		}

		var ok bool
		superclass, ok = rawSuperclass.(*Class)
		if !ok {
			err = fmt.Errorf("superclass must be a class, got %T", rawSuperclass)
			return
		}
	}

	i.environment.Define(v.Name.Lexeme, nil)

	if v.SuperClass != nil {
		i.environment = NewChildEnvironment(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*Function, len(v.Methods))
	for _, method := range v.Methods {
		methods[method.Name.Lexeme] = &Function{
			Declaration:   method,
			Closure:       i.environment,
			IsInitializer: method.Name.Lexeme == "init",
		}
	}

	class := &Class{
		Name:       v.Name.Lexeme,
		SuperClass: superclass,
		Methods:    methods,
	}

	if v.SuperClass != nil {
		i.environment = i.environment.parent
	}

	return i.environment.Assign(v.Name, class)
}

func (i *Interpreter) VisitSuperExpr(v *ast.SuperExpr) (result interface{}, err error) {
	distance, ok := i.locals[v]
	if !ok {
		panic("no super expression predefined")
	}
	rawSuperclass, err := i.environment.GetAt(distance, "super")
	if err != nil {
		panic("no superclass predefined")
	}
	superClass := rawSuperclass.(*Class)

	rawObject, err := i.environment.GetAt(distance-1, "this")
	if err != nil {
		panic("no 'this' preceding super")
	}
	instance := rawObject.(*Instance)
	method, ok := superClass.FindMethod(v.Method.Lexeme)
	if !ok {
		err = fmt.Errorf("super class does not have method %q: %s", v.Method.Lexeme, v.Method)
		return
	}

	result = method.Bind(instance)

	return
}

func (i *Interpreter) VisitReturnStmt(v *ast.ReturnStmt) (err error) {
	var value interface{}
	if v.Value != nil {
		value, err = i.evaluate(v.Value)
		if err != nil {
			return
		}
	}

	// unwind the stack
	panic(&returnPayload{Value: value})
}

func (i *Interpreter) lookUpVariable(name *token.Token, v ast.Expression) (result interface{}, err error) {
	distance, ok := i.locals[v]
	if ok {
		return i.environment.GetAt(distance, name.Lexeme)
	}

	return i.globals.Get(name)
}
