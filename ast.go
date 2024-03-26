package mist

import (
	"fmt"
	"strings"
)

// +-------+
// | Nodle |
// +-------+

const (
	TypeList = iota

	TypeSymbol
	TypeUint256
)


type Node struct {
	Type int

	ValueString string
	ValueUint256 uint64			// TODO

	Children []Node
}

func NewNodeUint256(literal uint64) Node {
	return Node{
		Type: TypeUint256,
		ValueString: "",
		ValueUint256: literal,
		Children: nil,
	}
}

func NewNodeSymbol(symbol string) Node {
	return Node {
		Type: TypeSymbol,
		ValueString: symbol,
		ValueUint256: 0,
		Children: nil,
	}
}

func NewNodeList() Node {
	return Node {
		Type: TypeList,
		ValueString: "",
		ValueUint256: 0,
		Children: make([]Node, 0, 4),
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
		return fmt.Sprintf("%d", n.ValueUint256)
	default:
		panic("TODO")
	}
}
