package resolver

type ClassType int

const (
	ClassTypeNone ClassType = iota
	ClassTypeClass
	ClassTypeSubclass
)
