package mist

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

const (
	freeMemoryPtr = 0x40
	freeMemory    = 0x80
)

// +-----------------+
// | BytecodeVisitor |
// +-----------------+

type BytecodeVisitor struct {
	segments []string
}

func NewBytecodeVisitor() BytecodeVisitor {
	v := BytecodeVisitor{
		segments: make([]string, 0, 256),
	}

	// Initialize the free memory pointer.  Mist follows the same
	// memory layout as Solidity.
	// v.pushU64(freeMemory)
	// v.pushU64(freeMemoryPtr)
	// v.pushOp(MSTORE)

	return v
}

func (v *BytecodeVisitor) pushOp(op OpCode) {
	v.pushSegment(fmt.Sprintf("%02x", byte(op)))
}

func (v *BytecodeVisitor) pushPointer() int {
	index := len(v.segments)
	v.segments = append(v.segments, "POINTR") // PUSH2 + two bytes address
	return index
}

func (v *BytecodeVisitor) pushSegment(code string) {
	v.segments = append(v.segments, code)
}

func (v *BytecodeVisitor) pushU64(x uint64) {
	v.pushU256(uint256.NewInt(x))
}

func (v *BytecodeVisitor) pushU256(x *uint256.Int) {
	hex := x.Hex()

	length := len(hex)/2 - 1 + len(hex)%2
	if length < 1 || 32 < length {
		panic("TODO")
	}

	op := OpCode(byte(PUSH0) + byte(length))
	v.pushOp(op)

	padding := ""
	if len(hex)%2 == 1 {
		padding = "0"
	}
	v.pushSegment(fmt.Sprintf("%s%s", padding, hex[2:]))
}

func (v *BytecodeVisitor) codeLength() int {
	ans := 0
	for i := range v.segments {
		ans += len(v.segments[i]) / 2
	}
	return ans
}

func (v *BytecodeVisitor) VisitList() {
}

func (v *BytecodeVisitor) VisitFunction(fn string, args []Node) {
	if isVariadic(fn) {
		v.visitVariadicOp(fn, args)
	} else if op, nargs, ok := isNative(fn); ok {
		assertArgsEq(fn, args, nargs)
		VisitSequence(v, args)
		v.pushOp(op)
	} else if callable, ok := isPreludeFunc(fn); ok {
		callable(v, args)
	} else {
		panic("unrecognized function: " + fn)
	}
}

func (v *BytecodeVisitor) VisitSymbol(_ string) {
}

func (v *BytecodeVisitor) VisitNumber(x *uint256.Int) {
	v.pushU256(x)
}

func (v *BytecodeVisitor) String() string {
	var b strings.Builder
	for i := range v.segments {
		b.WriteString(v.segments[i])
	}
	return b.String()
}

// +---------+
// | Private |
// +---------+

func (v *BytecodeVisitor) visitVariadicOp(fn string, args []Node) {
	assertArgsGte(fn, args, 2)

	var op OpCode
	switch fn {
	case "+":
		op = ADD
	case "*":
		op = MUL
	case "-":
		op = SUB
	case "/":
		op = DIV // TODO: SDIV
	case "&":
		fallthrough
	case "logand":
		op = AND
	case "|":
		fallthrough
	case "logior":
		op = OR
	case "^":
		fallthrough
	case "logxor":
		op = XOR
	default:
		panic("unrecognized arithmetic op: " + fn)
	}

	last := len(args) - 1
	args[last].Accept(v)
	for i := last - 1; i >= 0; i-- {
		args[i].Accept(v)
		v.pushOp(op)
	}
}
