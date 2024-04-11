package mist

import "fmt"

// +-------+
// | Scope |
// +-------+

type Scope struct {
	Constants map[string]Node
	Functions map[string]Node
	Variables map[string]Node

	Parent *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Constants: make(map[string]Node),
		Functions: make(map[string]Node),
		Variables: make(map[string]Node),

		Parent: parent,
	}
}

func NewGlobalScope() *Scope {
	return NewScope(nil)
}

func (e *Scope) NewChildScope() *Scope {
	return NewScope(e)
}

func (e *Scope) Defconst(name string, value Node) {
	if !value.IsConstant() {
		panic(fmt.Sprintf("%v: not constant: %v", value.Origin, value))
	}

	e.Constants[name] = value
}

func (e *Scope) Setq(name string, value Node) {
	e.Variables[name] = value
}
