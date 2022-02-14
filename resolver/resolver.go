package resolver

import (
	"errors"
	"fmt"

	"github.com/nanmu42/bluelox/ast"
	"github.com/nanmu42/bluelox/interpreter"
	"github.com/nanmu42/bluelox/token"
)

var (
	_ ast.ExprVisitor = (*Resolver)(nil)
	_ ast.StmtVisitor = (*Resolver)(nil)
)

// Resolver A variable usage refers to the preceding declaration
// with the same name in the innermost scope
// that encloses the expression where the variable is used.
type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          scopes // used as a stack
	currentFunction FunctionType
	currentClass    ClassType
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      newScopes(),
	}
}

func (r *Resolver) VisitBlockStmt(v *ast.BlockStmt) (err error) {
	r.beginScope()
	defer r.endScope()

	err = r.ResolveStmts(v.Stmts)
	return
}

func (r *Resolver) VisitExprStmt(v *ast.ExprStmt) (err error) {
	return r.resolveExpr(v.Expr)
}

func (r *Resolver) VisitFunctionStmt(v *ast.FunctionStmt) (err error) {
	err = r.declare(v.Name)
	if err != nil {
		return
	}
	r.define(v.Name)

	err = r.resolveFunction(v, FuncTypeFunc)
	return
}

func (r *Resolver) VisitIfStmt(v *ast.IfStmt) (err error) {
	err = r.resolveExpr(v.Condition)
	if err != nil {
		return
	}

	err = r.resolveStmt(v.ThenBranch)
	if err != nil {
		return
	}

	if v.ElseBranch != nil {
		err = r.resolveStmt(v.ElseBranch)
	}

	return
}

func (r *Resolver) VisitPrintStmt(v *ast.PrintStmt) (err error) {
	return r.resolveExpr(v.Expr)
}

func (r *Resolver) VisitReturnStmt(v *ast.ReturnStmt) (err error) {
	if r.currentFunction == FuncTypeNone {
		err = errors.New("can't return from top-level code")
		return
	}
	if v.Value != nil {
		if r.currentFunction == FuncTypeInitializer {
			err = fmt.Errorf("can't return a value from an initializer: %s", v.Keyword)
			return
		}
		return r.resolveExpr(v.Value)
	}

	return
}

func (r *Resolver) VisitVarStmt(v *ast.VarStmt) (err error) {
	err = r.declare(v.Name)
	if err != nil {
		return
	}
	if v.Initializer != nil {
		err = r.resolveExpr(v.Initializer)
		if err != nil {
			return
		}
	}
	r.define(v.Name)
	return
}

func (r *Resolver) VisitWhileStmt(v *ast.WhileStmt) (err error) {
	err = r.resolveExpr(v.Condition)
	if err != nil {
		return
	}

	return r.resolveStmt(v.Body)
}

func (r *Resolver) VisitAssignExpr(v *ast.AssignExpr) (result interface{}, err error) {
	err = r.resolveExpr(v.Value)
	if err != nil {
		return
	}

	err = r.resolveLocal(v, v.Name)
	return
}

func (r *Resolver) VisitBinaryExpr(v *ast.BinaryExpr) (result interface{}, err error) {
	err = r.resolveExpr(v.Left)
	if err != nil {
		return
	}

	err = r.resolveExpr(v.Right)
	return
}

func (r *Resolver) VisitCallExpr(v *ast.CallExpr) (result interface{}, err error) {
	err = r.resolveExpr(v.Callee)
	if err != nil {
		return
	}

	for _, arg := range v.Arguments {
		err = r.resolveExpr(arg)
		if err != nil {
			return
		}
	}

	return
}

func (r *Resolver) VisitGetExpr(v *ast.GetExpr) (result interface{}, err error) {
	err = r.resolveExpr(v.Object)
	return
}

func (r *Resolver) VisitGroupingExpr(v *ast.GroupingExpr) (result interface{}, err error) {
	err = r.resolveExpr(v.Expr)
	return
}

func (r *Resolver) VisitLiteralExpr(v *ast.LiteralExpr) (result interface{}, err error) {
	return
}

func (r *Resolver) VisitLogicalExpr(v *ast.LogicalExpr) (result interface{}, err error) {
	err = r.resolveExpr(v.Left)
	if err != nil {
		return
	}

	err = r.resolveExpr(v.Right)
	return
}

func (r *Resolver) VisitSetExpr(v *ast.SetExpr) (result interface{}, err error) {
	err = r.resolveExpr(v.Value)
	if err != nil {
		return
	}

	err = r.resolveExpr(v.Object)
	return
}

func (r *Resolver) VisitThisExpr(v *ast.ThisExpr) (result interface{}, err error) {
	if r.currentClass == ClassTypeNone {
		err = fmt.Errorf("can't use 'this' outside of a class: %s", v.Keyword)
		return
	}

	err = r.resolveLocal(v, v.Keyword)
	return
}

func (r *Resolver) VisitUnaryExpr(v *ast.UnaryExpr) (result interface{}, err error) {
	err = r.resolveExpr(v.Right)
	return
}

func (r *Resolver) VisitVariableExpr(v *ast.VariableExpr) (result interface{}, err error) {
	if !r.scopes.IsEmpty() {
		initialized, ok := r.scopes.Peek()[v.Name.Lexeme]
		if ok && !initialized {
			err = fmt.Errorf("can't read local variable %q in its own initializer", v.Name.Lexeme)
			return
		}
	}

	err = r.resolveLocal(v, v.Name)
	return
}

func (r *Resolver) VisitClassStmt(v *ast.ClassStmt) (err error) {
	enclosingClass := r.currentClass
	r.currentClass = ClassTypeClass
	defer func() {
		r.currentClass = enclosingClass
	}()

	err = r.declare(v.Name)
	if err != nil {
		return
	}

	r.define(v.Name)

	if v.SuperClass != nil {
		r.currentClass = ClassTypeSubclass

		if v.SuperClass.Name.Lexeme == v.Name.Lexeme {
			err = fmt.Errorf("the class %q can't inherit from itself, at line %d", v.Name.Lexeme, v.Name.Line)
			return
		}

		err = r.resolveExpr(v.SuperClass)
		if err != nil {
			return
		}

		r.beginScope()
		defer r.endScope()
		r.scopes.Peek()["super"] = true
	}

	r.beginScope()
	defer r.endScope()
	r.scopes.Peek()["this"] = true

	for _, method := range v.Methods {
		var funcType FunctionType
		if method.Name.Lexeme == "init" {
			funcType = FuncTypeInitializer
		} else {
			funcType = FuncTypeMethod
		}
		err = r.resolveFunction(method, funcType)
		if err != nil {
			return
		}
	}

	return
}

func (r *Resolver) VisitSuperExpr(v *ast.SuperExpr) (result interface{}, err error) {
	if r.currentClass == ClassTypeNone {
		err = fmt.Errorf("can't use 'super' outside of a class: %s", v)
		return
	} else if r.currentClass != ClassTypeSubclass {
		err = fmt.Errorf("can't use 'super' in a class with no superclass: %s", v)
		return
	}

	err = r.resolveLocal(v, v.Keyword)
	return
}

func (r *Resolver) ResolveStmts(stmts []ast.Statement) (err error) {
	for _, stmt := range stmts {
		err = r.resolveStmt(stmt)
		if err != nil {
			return
		}
	}

	return
}

func (r *Resolver) resolveStmt(stmt ast.Statement) error {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expression) (err error) {
	_, err = expr.Accept(r)
	return
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name *token.Token) (err error) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := r.scopes.Peek()
	if _, ok := scope[name.Lexeme]; ok {
		err = fmt.Errorf("variable %q at line %d already existed in this scope", name.Lexeme, name.Line)
		return
	}

	scope[name.Lexeme] = false
	return
}

func (r *Resolver) define(name *token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	r.scopes.Peek()[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(v ast.Expression, name *token.Token) error {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(v, len(r.scopes)-1-i)
			return nil
		}
	}

	return nil
}

func (r *Resolver) resolveFunction(v *ast.FunctionStmt, funcType FunctionType) (err error) {
	enclosingFuncType := r.currentFunction

	r.beginScope()
	r.currentFunction = funcType

	defer func() {
		r.endScope()
		r.currentFunction = enclosingFuncType
	}()

	for _, param := range v.Params {
		err = r.declare(param)
		if err != nil {
			return
		}
		r.define(param)
	}
	return r.ResolveStmts(v.Body)
}

type scopes []map[string]bool

func newScopes() scopes {
	return make(scopes, 0, 8)
}

func (s *scopes) Push(element ...map[string]bool) {
	*s = append(*s, element...)
}

func (s *scopes) Pop() (element map[string]bool) {
	element = s.Peek()
	*s = (*s)[:len(*s)-1]

	return
}

func (s *scopes) Peek() (element map[string]bool) {
	if s.IsEmpty() {
		panic("stack is empty")
	}
	element = (*s)[len(*s)-1]

	return
}

func (s *scopes) IsEmpty() bool {
	return len(*s) == 0
}
