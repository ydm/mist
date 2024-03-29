package mist

import (
	"fmt"
	"strconv"
	"strings"
)

// +---------+
// | Visitor |
// +---------+

type Visitor interface {
	VisitUint256(literal uint64)
	VisitSymbol(symbol string)

	VisitFunction(name string, args []Node)
}

func VisitSequence(v Visitor, nodes []Node) {
	for _, node := range nodes {
		node.Accept(v)
	}
}

// +------+
// | Node |
// +------+

const (
	TypeList = iota

	TypeSymbol
	TypeInt256
	TypeUint256
)

type Node struct {
	Type int

	ValueString  string
	ValueUint256 uint64 // TODO

	Children []Node
}

func NewNodeUint256(literal uint64) Node {
	return Node{
		Type:         TypeUint256,
		ValueString:  "",
		ValueUint256: literal,
		Children:     nil,
	}
}

func NewNodeSymbol(symbol string) Node {
	return Node{
		Type:         TypeSymbol,
		ValueString:  symbol,
		ValueUint256: 0,
		Children:     nil,
	}
}

func NewNodeList() Node {
	return Node{
		Type:         TypeList,
		ValueString:  "",
		ValueUint256: 0,
		Children:     make([]Node, 0, 4),
	}
}

func (n *Node) AddChild(child Node) {
	if n.IsAtom() {
		panic("TODO")
	}

	n.Children = append(n.Children, child)
}

// +------------+
// | Predicates |
// +------------+

func (n *Node) IsAtom() bool {
	return n.Type != TypeList
}

func (n *Node) IsNumber() bool {
	return n.Type == TypeInt256 || n.Type == TypeUint256
}

func (n *Node) IsSigned() bool {
	return n.Type == TypeInt256
}

func (n *Node) IsUnsigned() bool {
	return n.Type == TypeUint256
}

func AllAtoms(nodes []Node) bool {
	allAtoms := true
	for i := range nodes {
		allAtoms = allAtoms && nodes[i].IsAtom()
	}
	return allAtoms
}

func AllNumbers(nodes []Node) bool {
	allNumbers := true
	for i := range nodes {
		allNumbers = allNumbers && nodes[i].IsNumber()
	}
	return allNumbers
}

func AllSigned(nodes []Node) bool {
	allNumbers := true
	allSigned := true
	for i := range nodes {
		allNumbers = allNumbers && nodes[i].IsNumber()
		allSigned = allSigned && nodes[i].IsSigned()
	}
	return allNumbers && allSigned
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
	case TypeUint256:
		v.VisitUint256(n.ValueUint256)
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
	case TypeUint256:
		return strconv.FormatUint(n.ValueUint256, 10)
	default:
		panic("TODO")
	}
}
