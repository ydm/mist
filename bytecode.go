package mist

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	freeMemoryPtr = 0x40
	freeMemory = 0x80
)

// +-----------------+
// | BytecodeVisitor |
// +-----------------+

type BytecodeVisitor struct {
	labels   []int // Label (index) to position.
	segments []string
}

func NewBytecodeVisitor() BytecodeVisitor {
	v := BytecodeVisitor{
		segments: make([]string, 0, 256),
	}

	// Initialize the free memory pointer.  Mist follows the same
	// memory layout as Solidity.
	v.pushUnsigned(freeMemory)
	v.pushUnsigned(freeMemoryPtr)
	v.pushOp(MSTORE)

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

func (v *BytecodeVisitor) pushUnsigned(x uint64) {
	hex := strconv.FormatUint(x, 16)

	length := len(hex)/2 + len(hex)%2
	if length < 1 || 32 < length {
		panic("TODO")
	}

	op := OpCode(byte(PUSH0) + byte(length))
	v.pushOp(op)

	padding := ""
	if len(hex) % 2 == 1 {
		padding = "0"
	}
	v.pushSegment(fmt.Sprintf(
		"%s%s",
		padding,
		hex,
	))
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
	//nolint:gocritic,nestif
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

	// TODO: KECCAK256

}

func (v *BytecodeVisitor) VisitSymbol(_ string) {
}

func (v *BytecodeVisitor) VisitUint256(literal uint64) {
	v.pushUnsigned(literal)
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

func (v *BytecodeVisitor) visitAlpha(fn string) {
	var name string
	switch fn {
	case "address":
		name = "ADDRESS"
	case "origin":
		name = "ORIGIN"
	case "caller":
		name = "CALLER"
	case "call-value":
		name = "CALLVALUE"
	case "calldata-load":
		name = "CALLDATALOAD"
	case "calldata-size":
		name = "CALLDATASIZE"
	case "code-size":
		name = "CODESIZE"
	case "gas-price":
		name = "GASPRICE"
	case "return-data-size":
		name = "RETURNDATASIZE"
	case "coinbase":
		name = "COINBASE"
	case "timestamp":
		name = "TIMESTAMP"
	case "block-number":
		name = "NUMBER"
	case "prev-randao":
		name = "PREVRANDAO"
	case "gas-limit":
		name = "GASLIMIT"
	case "chain-id":
		name = "CHAINID"
	case "self-balance":
		name = "SELFBALANCE"
	case "base-fee":
		name = "BASEFEE"
	default:
		panic("unrecognized alpha op: " + fn)
	}

	op, ok := stringToOp[name]
	if !ok {
		panic("TODO")
	}

	v.pushSegment(fmt.Sprintf("%02x", byte(op)))
}

func (v *BytecodeVisitor) visitOp(fn string, args []Node) {

}

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
	case "logand":
		op = AND
	case "logior":
		op = OR
	case "logxor":
		op = XOR
	default:
		panic("unrecognized arithmetic op: " + fn)
	}

	last := len(args)-1
	args[last].Accept(v)
	for i := last - 1; i >= 0; i-- {
		args[i].Accept(v)
		v.pushOp(op)
	}
}
