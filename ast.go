package mist

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

// +---------+
// | Visitor |
// +---------+

// Each visit function returns the stack delta, i.e. how many stack
// elements were pushed or taken.
type Visitor interface {
	// Simple literals, do not need a scope.
	VisitNil()
	VisitT()
	VisitNumber(value *uint256.Int)
	VisitString(node Node)

	// Symbols currently serve only as const/variable identifiers.
	VisitSymbol(s *Scope, esp int, symbol Node)

	VisitFunction(s *Scope, esp int, call Node)
}

func VisitSequence(v Visitor, s *Scope, esp int, nodes []Node, dir int) int {
	switch dir {
	case -1:
		for i := len(nodes) - 1; i >= 0; i-- {
			nodes[i].Accept(v, s, esp)
			esp += 1
		}
		return len(nodes)
	case 1:
		for i := range nodes {
			nodes[i].Accept(v, s, esp)
			esp += 1
		}
		return len(nodes)
	default:
		panic("invalid direction")
	}
}

// +------+
// | Node |
// +------+

const (
	NodeList   = iota // 0
	NodeSymbol        // 1
	NodeNumber        // 2
	NodeString        // 3

	// function, primitive, macro
)

type Node struct {
	Type int

	ValueString string
	ValueNumber *uint256.Int

	Children []Node

	Origin Origin
}

func NewNodeU256(x *uint256.Int, origin Origin) Node {
	return Node{
		Type:        NodeNumber,
		ValueString: "",
		ValueNumber: x,
		Children:    nil,
		Origin:      origin,
	}
}

func NewNodeU64(x uint64, origin Origin) Node {
	return NewNodeU256(uint256.NewInt(x), origin)
}

func NewNodeString(value string, origin Origin) Node {
	return Node{
		Type:        NodeString,
		ValueString: value,
		ValueNumber: nil,
		Children:    nil,
		Origin:      origin,
	}
}

func NewNodeSymbol(symbol string, origin Origin) Node {
	return Node{
		Type:        NodeSymbol,
		ValueString: symbol,
		ValueNumber: nil,
		Children:    nil,
		Origin:      origin,
	}
}

func NewNodeList(origin Origin) Node {
	return Node{
		Type:        NodeList,
		ValueString: "",
		ValueNumber: nil,
		Children:    make([]Node, 0, 4),
		Origin:      origin,
	}
}

func NewNodeNil(origin Origin) Node {
	return NewNodeSymbol("nil", origin)
}

func NewNodeProgn() Node {
	progn := NewNodeList(NewOriginEmpty())
	progn.AddChild(NewNodeSymbol("progn", NewOriginEmpty()))
	return progn
}

// TODO: Maybe accepting (Node) is better?  And then:
//
// parent = parent.AddChil(child)
func (n *Node) AddChild(child Node) {
	if n.IsAtom() {
		panic(fmt.Sprintf("%v: atom %s cannot have children", n.Origin, n.String()))
	}

	n.Children = append(n.Children, child)
}

func (n *Node) AddChildren(children []Node) {
	for i := range children {
		n.AddChild(children[i])
	}
}

func (n *Node) FunctionName() string {
	if !n.IsList() || n.NumChildren() < 1 || !n.Children[0].IsSymbol() {
		panic("")
	}

	return n.Children[0].ValueString
}

func (n *Node) IsAtom() bool {
	return !n.IsList()
}

func (n *Node) IsConstant() bool {
	return n.Type == NodeNumber || n.Type == NodeSymbol
}

func (n *Node) IsEmptyList() bool {
	return n.IsList() && n.NumChildren() == 0
}

func (n *Node) IsFunctionCall(name string) bool {
	return (n.IsList() &&
		n.NumChildren() > 1 &&
		n.Children[0].IsSymbol() &&
		n.FunctionName() == name)
}

func (n *Node) IsList() bool {
	return n.Type == NodeList
}

func (n *Node) IsNil() bool {
	switch n.Type {
	case NodeList:
		if n.NumChildren() == 2 && n.Children[0].IsQuote() {
			return n.Children[1].IsNil()
		}
		return n.IsEmptyList()
	case NodeSymbol:
		return n.ValueString == "nil"
	case NodeNumber:
		return n.ValueNumber.IsZero()
	default:
		return false
	}
}

func (n *Node) IsQuote() bool {
	return n.IsThisSymbol("quote")
}

func (n *Node) IsString() bool {
	return n.Type == NodeString
}

func (n *Node) IsSymbol() bool {
	return n.Type == NodeSymbol && n.ValueString != ""
}

func (n *Node) IsT() bool {
	switch n.Type {
	case NodeList:
		if n.NumChildren() == 2 && n.Children[0].IsQuote() {
			return n.Children[1].IsT()
		}
		return false // Could be true, could be not.
	case NodeSymbol:
		return n.ValueString == "t"
	case NodeNumber:
		return !n.ValueNumber.IsZero()
	default:
		return false
	}
}

func (n *Node) IsThisSymbol(s string) bool {
	return n.IsSymbol() && n.ValueString == s
}

func (n *Node) NumChildren() int {
	return len(n.Children)
}

// +---------+
// | Visitor |
// +---------+

func (n *Node) Accept(v Visitor, s *Scope, esp int) {
	if n.IsNil() {
		v.VisitNil()
		return
	}

	switch n.Type {
	case NodeNumber:
		v.VisitNumber(n.ValueNumber)
		return
	case NodeString:
		v.VisitString(*n)
		return
	case NodeSymbol:
		if n.IsT() {
			v.VisitT()
		} else {
			v.VisitSymbol(s, esp, *n)
		}
		return
	case NodeList:
		if n.NumChildren() < 1 {
			// TODO: I should support (empty) arrays too!
			panic("TODO")
		} else if !n.Children[0].IsSymbol() {
			panic(fmt.Sprintf("%v: %s is not a symbol", n.Children[0].Origin, n.Children[0].String()))
		} else {
			v.VisitFunction(s, esp, *n)
			return
		}
	default:
		panic("broken invariant")
	}
}

func (n *Node) String() string {
	switch n.Type {
	case NodeList:
		inner := make([]string, 0, len(n.Children))
		for _, child := range n.Children {
			inner = append(inner, child.String())
		}
		return fmt.Sprintf("(%s)", strings.Join(inner, " "))
	case NodeSymbol:
		return n.ValueString
	case NodeNumber:
		return n.ValueNumber.Dec()
	default:
		panic("TODO")
	}
}
