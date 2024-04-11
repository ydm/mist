package mist

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

// +---------+
// | Visitor |
// +---------+

type Visitor interface {
	VisitNil()
	VisitT()

	VisitNumber(value *uint256.Int)
	VisitSymbol(s *Scope, symbol string)

	VisitFunction(s *Scope, name string, args []Node)
}

func VisitSequence(v Visitor, s *Scope, nodes []Node, dir int) {
	switch dir {
	case -1:
		for i := len(nodes) - 1; i >= 0; i-- {
			nodes[i].Accept(v, s)
		}
	case 1:
		for i := range nodes {
			nodes[i].Accept(v, s)
		}
	default:
		panic("invalid direction")
	}
}

// +------+
// | Node |
// +------+

const (
	TypeList = iota

	TypeSymbol
	TypeNumber

	// string, function, primitive, macro
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
		Type:        TypeNumber,
		ValueString: "",
		ValueNumber: x,
		Children:    nil,
		Origin:      origin,
	}
}

func NewNodeU64(x uint64, origin Origin) Node {
	return NewNodeU256(uint256.NewInt(x), origin)
}

func NewNodeSymbol(symbol string, origin Origin) Node {
	return Node{
		Type:        TypeSymbol,
		ValueString: symbol,
		ValueNumber: nil,
		Children:    nil,
		Origin:      origin,
	}
}

func NewNodeList(origin Origin) Node {
	return Node{
		Type:        TypeList,
		ValueString: "",
		ValueNumber: nil,
		Children:    make([]Node, 0, 4),
		Origin:      origin,
	}
}

func NewNodeNil(origin Origin) Node {
	return NewNodeSymbol("nil", origin)
}

func NewNodeProgn(origin Origin) Node {
	progn := NewNodeList(NewOriginEmpty())
	progn.AddChild(NewNodeSymbol("progn", NewOriginEmpty()))
	return progn
}

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

func (n *Node) IsAtom() bool {
	return !n.IsList()
}

func (n *Node) IsConstant() bool {
	return n.Type == TypeNumber || n.Type == TypeSymbol
}

func (n *Node) IsEmptyList() bool {
	return n.IsList() && n.NumChildren() == 0
}

func (n *Node) IsFunction(name string) bool {
	return (n.IsList() &&
		n.NumChildren() > 1 &&
		n.Children[0].IsSymbol() &&
		n.Children[0].ValueString == name)
}

func (n *Node) IsList() bool {
	return n.Type == TypeList
}

func (n *Node) IsNil() bool {
	switch n.Type {
	case TypeList:
		return n.IsEmptyList() || (n.NumChildren() == 1 && n.Children[0].IsQuote())
	case TypeSymbol:
		return n.ValueString == "nil"
	case TypeNumber:
		return false
	default:
		return false
	}
}

func (n *Node) IsQuote() bool {
	return n.IsSymbol() && n.ValueString == "quote"
}

func (n *Node) IsSymbol() bool {
	return n.Type == TypeSymbol
}

func (n *Node) IsT() bool {
	return n.IsSymbol() && n.ValueString == "t"
}

func (n *Node) NumChildren() int {
	return len(n.Children)
}

// +---------+
// | Visitor |
// +---------+

func (n *Node) Accept(v Visitor, s *Scope) {
	if n.IsNil() {
		v.VisitNil()
		return
	} else if n.IsT() {
		v.VisitT()
		return
	}

	switch n.Type {
	case TypeNumber:
		v.VisitNumber(n.ValueNumber)
	case TypeSymbol:
		v.VisitSymbol(s, n.ValueString)
	case TypeList:
		if n.NumChildren() < 1 {
			// TODO: I should support (empty) arrays too!
			panic("TODO")
		} else if !n.Children[0].IsSymbol() {
			panic(fmt.Sprintf("%v: %s is not a symbol", n.Children[0].Origin, n.Children[0].String()))
		} else {
			fn, args := n.Children[0].ValueString, n.Children[1:]
			v.VisitFunction(s, fn, args)
		}
	}
}

func (n *Node) String() string {
	switch n.Type {
	case TypeList:
		inner := make([]string, 0, len(n.Children))
		for _, child := range n.Children {
			inner = append(inner, child.String())
		}
		return fmt.Sprintf("(%s)", strings.Join(inner, " "))
	case TypeSymbol:
		return n.ValueString
	case TypeNumber:
		return n.ValueNumber.Dec()
	default:
		panic("TODO")
	}
}
