package mist

import "strconv"

const (
	TypeInt = iota
	TypeSymbol
)

type Atom struct {
	Type int

	Symbol string
	Uint   uint
}

func NewAtomUint(x uint) Atom {
	return Atom{
		Type: TypeInt,

		Symbol: "",
		Uint:   x,
	}
}

func NewAtomSymbol(x string) Atom {
	return Atom{
		Type: TypeSymbol,

		Symbol: x,
		Uint:   0,
	}
}

func NewAtomInferred(s string) Atom {
	x, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return NewAtomUint(uint(x))
	}

	return NewAtomSymbol(s)
}
