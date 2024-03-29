package mist

import (
	"fmt"
	"strings"
)

// +---------+
// | Segment |
// +---------+

type segment struct {
	position int    // Position in global code.
	code     string // Actual code in hex.
}

func (s segment) Len() int {
	return len(s.code)
}

// +-----------------+
// | BytecodeVisitor |
// +-----------------+

type BytecodeVisitor struct {
	labels []int // Label (index) to position.
	output []segment
}

func NewBytecodeVisitor() BytecodeVisitor {
	return BytecodeVisitor{
		output: make([]segment, 0, 256),
	}
}

func (v *BytecodeVisitor) pushOp(op OpCode) {
	v.pushSegment(fmt.Sprintf("%02x", byte(op)))
}

func (v *BytecodeVisitor) pushPointer() int {
	index := len(v.output)

	position := v.codeLength()
	v.output = append(v.output, segment{
		position: position,
		code:     "POINTR", // PUSH2 + two more bytes
	})

	return index
}

func (v *BytecodeVisitor) pushSegment(code string) {
	position := v.codeLength()
	v.output = append(v.output, segment{position, code})
}

func (v *BytecodeVisitor) codeLength() int {
	length := 0
	if len(v.output) > 0 {
		last := v.output[len(v.output)-1]
		length = last.position + len(last.code)/2
	}
	return length
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
	// TODO: Support PUSH[2-32]
	v.pushSegment(fmt.Sprintf("60%02x", literal))
}

func (v *BytecodeVisitor) String() string {
	var b strings.Builder
	for i := range v.output {
		b.WriteString(v.output[i].code)
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
	case "call-data-load":
		name = "CALLDATALOAD"
	case "call-data-size":
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
	case "<":
		op = LT // TODO: SLT
	case ">":
		op = GT // TODO: SGT
	case "=":
		op = EQ
	case "logand":
		op = AND
	case "logior":
		op = OR
	case "logxor":
		op = XOR
	default:
		panic("unrecognized arithmetic op: " + fn)
	}

	args[0].Accept(v)
	for _, arg := range args[1:] {
		arg.Accept(v)
		v.pushOp(op)
	}
}

// TODO: MOD SMOD EXP NOT ISZERO SIGNEXTEND BYTE SHL SHR SAR ADDMOD MULMOD
// TODO: BALANCE

// TODO:
// BALANCE, CALLDATALOAD, CALLDATACOPY, CODECOPY, EXTCODESIZE, EXTCODECOPY,
// RETURNDATACOPY, EXTCODEHASH, BLOCKHASH,

/*
 (when (< (call-data-size) 4) (revert 00 00))
 (shr (call-data-load) e0)
 (dup1)
*/
