package mist

import "fmt"

// +-------+
// | Scope |
// +-------+

type LispFunction struct {
	Origin Origin
	Name   string
	Args   []Node
	Body   Node
}

func NewLispFunction(n Node) (empty LispFunction, _ error) {
	// [0] defun
	// [1] function-name
	// [2] function-args
	// [3:] function-body

	if !n.IsList() {
		panic("TODO")
	}

	if n.NumChildren() < 3 {
		return empty, NewCompilationError(
			n.Origin,
			fmt.Sprintf("invalid function definition: %v", n),
		)
	}

	identifier := n.Children[1]
	if !identifier.IsSymbol() {
		return empty, NewCompilationError(
			identifier.Origin,
			fmt.Sprintf("invalid function identifier: %v", &identifier),
		)
	}

	if !n.Children[2].IsList() {
		return empty, NewCompilationError(
			n.Children[2].Origin,
			fmt.Sprintf("invalid function arguments: %v", &n.Children[2]),
		)
	}

	args := n.Children[2].Children

	body := NewNodeNil(n.Origin)
	if n.NumChildren() > 3 {
		body = NewNodeProgn(n.Children[3].Origin)
		body.AddChildren(n.Children[3:])
	}

	return LispFunction{
		Origin: n.Origin,
		Name:   identifier.ValueString,
		Args:   args,
		Body:   body,
	}, nil
}

type StackVariable struct {
	Origin     Origin
	Identifier string
	Position   int
}

type Scope struct {
	Constants map[string]Node
	Functions map[string]LispFunction
	StackVariables map[string]StackVariable

	Parent *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Constants:      make(map[string]Node),
		Functions:      make(map[string]LispFunction),
		StackVariables: make(map[string]StackVariable),

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

func (s *Scope) GetConstant(identifier string) (Node, bool) {
	node, ok := s.Constants[identifier]
	if !ok && s.Parent != nil {
		return s.Parent.GetConstant(identifier)
	}
	return node, ok
}

func (s *Scope) GetFunction(identifier string) (LispFunction, bool) {
	fn, ok := s.Functions[identifier]
	if !ok && s.Parent != nil {
		return s.Parent.GetFunction(identifier)
	}
	return fn, ok
}

func (s *Scope) GetStackVariable(identifier string) (StackVariable, bool) {
	variable, ok := s.StackVariables[identifier]
	if !ok && s.Parent != nil {
		return s.Parent.GetStackVariable(identifier)
	}
	return variable, ok
}

// +---------+
// | Setters |
// +---------+

func (s *Scope) Defconst(identifier string, value Node) {
	if !value.IsConstant() {
		panic(fmt.Sprintf("%v: %v is not constant", value.Origin, value))
	}

	if _, ok := s.GetConstant(identifier); ok {
		panic(fmt.Sprintf("%v: constant %s is already defined", value.Origin, identifier))
	}

	s.Constants[identifier] = value
}

func (s *Scope) Defun(fn LispFunction) {
	s.Functions[fn.Name] = fn
}

func (s *Scope) SetStackVariable(identifier string, variable StackVariable) {
	s.StackVariables[identifier] = variable
}

// func (s *Scope) Setq(name string, value Node) {
// 	s.Variables[name] = value
// }
