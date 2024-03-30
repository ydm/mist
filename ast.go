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
	VisitNumber(value *uint256.Int)
	VisitSymbol(symbol string)

	VisitFunction(name string, args []Node)
}

func Visit(v Visitor, node Node) {
	node.Accept(v)
}

func VisitSequence(v Visitor, nodes []Node, dir int) {
	switch dir {
	case -1:
		for i := len(nodes) - 1; i >= 0; i-- {
			Visit(v, nodes[i])
		}
	case 1:
		for i := range nodes {
			Visit(v, nodes[i])
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
)

type Node struct {
	Type int

	ValueString string
	ValueNumber *uint256.Int

	Children []Node
}

func NewNodeU256(x *uint256.Int) Node {
	return Node{
		Type:        TypeNumber,
		ValueString: "",
		ValueNumber: x,
		Children:    nil,
	}
}

func NewNodeU64(x uint64) Node {
	return NewNodeU256(uint256.NewInt(x))
}

func NewNodeSymbol(symbol string) Node {
	return Node{
		Type:        TypeSymbol,
		ValueString: symbol,
		ValueNumber: nil,
		Children:    nil,
	}
}

func NewNodeList() Node {
	return Node{
		Type:        TypeList,
		ValueString: "",
		ValueNumber: nil,
		Children:    make([]Node, 0, 4),
	}
}

func (n *Node) AddChild(child Node) {
	if n.IsAtom() {
		panic("TODO")
	}

	n.Children = append(n.Children, child)
}

func (n *Node) IsAtom() bool {
	return n.Type != TypeList
}

// +---------+
// | Visitor |
// +---------+

func (n *Node) Accept(v Visitor) {
	switch n.Type {
	case TypeList:
		if len(n.Children) < 1 {
			// TODO: I should support arrays too...
			panic("TODO")
		}

		if n.Children[0].Type != TypeSymbol {
			panic("TODO")
		}

		v.VisitFunction(n.Children[0].ValueString, n.Children[1:])
	case TypeSymbol:
		v.VisitSymbol(n.ValueString)
	case TypeNumber:
		v.VisitNumber(n.ValueNumber)
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
