package interpreter

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/nanmu42/bluelox/ast"
	"github.com/nanmu42/bluelox/token"
)

var _ ast.ExprVisitor = (*Interpreter)(nil)
var _ ast.StmtVisitor = (*Interpreter)(nil)

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewGlobalEnvironment(),
	}
}

func (i *Interpreter) Interpret(stmts []ast.Statement) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &RuntimeError{
				Reason: fmt.Sprintf("runtime panicking: \n%v\n", r),
			}
		}
	}()

	for _, stmt := range stmts {
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

	fmt.Println(i.stringify(result))
	return
}

func (i *Interpreter) execute(stmt ast.Statement) error {
	return stmt.Accept(i)
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
	result, err = i.environment.Get(v.Name)
	return
}

func (i *Interpreter) VisitAssignExpr(v *ast.AssignExpr) (result interface{}, err error) {
	result, err = i.evaluate(v.Value)
	if err != nil {
		return
	}

	err = i.environment.Assign(v.Name, result)

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

func (i *Interpreter) VisitFunctionStmt(v *ast.FunctionStmt) (err error) {
	i.environment.Define(v.Name.Lexeme, &function{
		Declaration: v,
		Closure:     i.environment,
	})
	return nil
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
