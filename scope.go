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

func (s *Scope) NewChildScope() *Scope {
	return NewScope(s)
}

// +---------+
// | Getters |
// +---------+

func (s *Scope) GetConstant(name string) (Node, bool) {
	node, ok := s.Constants[name]
	if !ok && s.Parent != nil {
		return s.Parent.GetConstant(name)
	}
	return node, ok
}

// +---------+
// | Setters |
// +---------+

func (s *Scope) Defconst(name string, value Node) {
	if !value.IsConstant() {
		panic(fmt.Sprintf("%v: %v is not constant", value.Origin, value))
	}

	if _, ok := s.GetConstant(name); ok {
		panic(fmt.Sprintf("%v: constant %s already defined", value.Origin, name))
	}

	s.Constants[name] = value
}

func (s *Scope) Setq(name string, value Node) {
	s.Variables[name] = value
}
