package mist

import "fmt"

// +-------+
// | Scope |
// +-------+

type LispFunction struct {
	Origin Origin
	Name   string
	Args   []Node
	Body   Node // Wrapped in (progn).
}

func NewLispFunction(n Node) (LispFunction, error) {
	var empty LispFunction

	// [0] defun
	// [1] name
	// [2] args
	// [3:] body

	if !n.IsList() {
		panic("TODO")
	}

	// (defun name args body...), length should be >= 3
	if n.NumChildren() < 3 {
		return empty, NewCompilationError(
			n.Origin,
			fmt.Sprintf("invalid function definition: %v", n),
		)
	}

	// [1] name
	identifier := n.Children[1]
	if !identifier.IsSymbol() {
		return empty, NewCompilationError(
			identifier.Origin,
			fmt.Sprintf("invalid function identifier: %v", &identifier),
		)
	}

	// [2] args
	if !n.Children[2].IsList() {
		return empty, NewCompilationError(
			n.Children[2].Origin,
			fmt.Sprintf("fn arguments are not a list: %v", &n.Children[2]),
		)
	}

	args := n.Children[2].Children

	// Each "argument" should be a symbol.
	for i := range args {
		if !args[i].IsSymbol() {
			return empty, NewCompilationError(
				args[i].Origin,
				fmt.Sprintf("fn argument is not a symbol: %v", &args[i]),
			)
		}
	}

	// [3:] body...
	body := NewNodeNil(n.Origin)
	if n.NumChildren() > 3 {
		body = NewNodeProgn()
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
	Position   int		// Absolute position in stack.
}

type Scope struct {
	Constants     map[string]Node
	Functions     map[string]LispFunction
	CallAddresses map[string]int32

	StackVariables   map[string]StackVariable
	StorageVariables map[string]int32

	Parent *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Constants:     make(map[string]Node),
		Functions:     make(map[string]LispFunction),
		CallAddresses: make(map[string]int32),

		StackVariables:   make(map[string]StackVariable),
		StorageVariables: make(map[string]int32),

		Parent: parent,
	}
}

func NewGlobalScope() *Scope {
	return NewScope(nil)
}

func (s *Scope) NewChildScope() *Scope {
	return NewScope(s)
}

func (s *Scope) IsGlobal() bool {
	return s.Parent == nil
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

func (s *Scope) GetCallAddress(identifier string) (int32, bool) {
	ptr, ok := s.CallAddresses[identifier]
	if !ok && s.Parent != nil {
		return s.Parent.GetCallAddress(identifier)
	}
	return ptr, ok
}

func (s *Scope) GetStackVariable(identifier string) (StackVariable, bool) {
	variable, ok := s.StackVariables[identifier]
	if !ok && s.Parent != nil {
		return s.Parent.GetStackVariable(identifier)
	}
	return variable, ok
}

func (s *Scope) GetStorageVariable(identifier string) (int32, bool) {
	pos, ok := s.StorageVariables[identifier]
	if !ok && s.Parent != nil {
		return s.Parent.GetStorageVariable(identifier)
	}
	return pos, ok
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

func (s *Scope) SetStorageVariable(name string, position int32) {
	if position < 0 {
		panic("TODO")
	}
	s.StorageVariables[name] = position
}

func (s *Scope) SetCallAddress(identifier string, segmentID int32) {
	if segmentID <= 0 {
		panic("TODO")
	}

	// CallAddresses match Functions one-to-one.
	if _, ok := s.Functions[identifier]; ok {
		s.CallAddresses[identifier] = segmentID
	} else {
		s.Parent.SetCallAddress(identifier, segmentID)
	}
}
