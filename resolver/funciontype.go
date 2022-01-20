//go:generate stringer -type FunctionType
package resolver

type FunctionType int

const (
	FuncTypeNone FunctionType = iota
	FuncTypeFunc
	FuncTypeInitializer
	FuncTypeMethod
)
