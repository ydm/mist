package mist

import "github.com/holiman/uint256"

type CounterVisitor struct {
	functions map[string]int
}

func (c *CounterVisitor) VisitNil() {}

func (c *CounterVisitor) VisitT() {}

func (c *CounterVisitor) VisitNumber(value *uint256.Int) {}

func (c *CounterVisitor) VisitString(node Node) {}

func (c *CounterVisitor) VisitSymbol(s *Scope, esp int, symbol Node) {}

func (c *CounterVisitor) VisitFunction(s *Scope, esp int, call Node) {
	
}
