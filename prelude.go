package mist

import "fmt"

// +-----------+
// | Utilities |
// +-----------+

func assertArgsEq(fn string, args []Node, want int) {
	if have := len(args); have != want {
		panic(fmt.Sprintf(
			"wrong number of arguments for (%s): have %d, want %d",
			fn,
			have,
			want,
		))
	}
}

func assertArgsGte(fn string, args []Node, want int) {
	if have := len(args); have < want {
		panic(fmt.Sprintf(
			"wrong number of arguments for (%s): have %d, want at least %d",
			fn,
			have,
			want,
		))
	}
}

func isNative(tok string) (OpCode, int, bool) {
	switch tok {
	case "stop":
		return STOP, 0, true

	// case SIGNEXTEND

	case "zerop":
		return ISZERO, 1, true
	case "not":
		return NOT, 1, true
	case "byte": // (byte word which)
		return BYTE, 2, true

	case "<<": // (<< value count)
		return SHL, 2, true
	case ">>": // (>> value count)
		return SHR, 2, true // TODO: SAR

	// α == 1
	case "address":
		return ADDRESS, 0, true
	case "balance":
		return BALANCE, 1, true
	case "origin":
		return ORIGIN, 0, true
	case "caller":
		return CALLER, 0, true
	case "call-value":
		return CALLVALUE, 0, true
	case "call-data-load": // (call-data-load start)
		return CALLDATALOAD, 1, true
	case "call-data-size":
		return CALLDATASIZE, 0, true
	case "call-data-copy": // (call-data-copy length id-offset mm-start)
		return CALLDATACOPY, 3, true
	case "code-size":
		return CODESIZE, 0, true
	case "code-copy":
		return CODECOPY, 3, true // (code-copy length ib-offset mm-start)
	case "gas-price":
		return GASPRICE, 0, true
	case "ext-code-size":
		return EXTCODESIZE, 1, true
	// case "ext-code-copy":
	// case "return-data-size":
	// case "return-data-copy":
	// case "ext-code-hash":
	// case "block-hash":
	case "return-data-size":
		return RETURNDATASIZE, 0, true
	case "coinbase":
		return COINBASE, 0, true
	case "timestamp":
		return TIMESTAMP, 0, true
	case "block-number":
		return NUMBER, 0, true
	case "prev-randao":
		return PREVRANDAO, 0, true
	case "gas-limit":
		return GASLIMIT, 0, true
	case "chain-id":
		return CHAINID, 0, true
	case "self-balance":
		return SELFBALANCE, 0, true
	case "base-fee":
		return BASEFEE, 0, true

	// case "pop":
	case "mload": // (mload start)
		return MLOAD, 1, true
	case "mstore": // (mstore value start)
		return MSTORE, 2, true
	case "mstore8":
		return MSTORE8, 2, true
	case "sload": // (sload word-index)
		return SLOAD, 1, true
	case "sstore": // (sstore value word-index)
		return SSTORE, 2, true

	// case "jump":
	// case "jumpi":
	// case "program-counter": return PC, 0, true
	case "memory-size":
		return MSIZE, 0, true
	case "gas":
		return GAS, 0, true
	}

	// case "jumpdest"
	// case "push1..16"
	// case "dup1..16"
	// case "swap1..16"

	// case CREATE
	// case CALL
	// case CALLCODE
	// case RETURN
	// case DELEGATECALL
	// case CREATE2
	// case STATICCALL
	// case REVERT
	// case INVALID
	// case SELFDESTRUCT

	return 0, 0, false
}

func isVariadic(tok string) bool {
	variadic := []string{
		"*", "+", "-", "/",
		"<", ">", "=",
		"logand", "logior", "logxor",
	}
	for _, x := range variadic {
		if tok == x {
			return true
		}
	}
	return false
}

// +-------------------+
// | Prelude functions |
// +-------------------+

func fnAddmod(v *BytecodeVisitor, args []Node) {
	assertArgsEq("+%", args, 3)

	x, y, m := args[0], args[1], args[2]
	m.Accept(v)
	y.Accept(v)
	x.Accept(v)
	v.pushOp(ADDMOD)
}

func fnExpt(v *BytecodeVisitor, args []Node) {
	assertArgsEq("expt", args, 2)

	x, y := args[0], args[1]
	y.Accept(v)
	x.Accept(v)
	v.pushOp(EXP)
}

func fnMod(v *BytecodeVisitor, args []Node) {
	assertArgsEq("%", args, 2)

	x, y := args[0], args[1]
	y.Accept(v)
	x.Accept(v)
	v.pushOp(MOD)
}

func fnMulmod(v *BytecodeVisitor, args []Node) {
	assertArgsEq("*%", args, 3)

	x, y, m := args[0], args[1], args[2]
	m.Accept(v)
	y.Accept(v)
	x.Accept(v)
	v.pushOp(ADDMOD)
}

func fnProgn(v *BytecodeVisitor, args []Node) {
	VisitSequence(v, args)
}

func fnReturn(v *BytecodeVisitor, args []Node) {
	// TODO
}

func fnRevert(v *BytecodeVisitor, args []Node) {
	assertArgsEq("revert", args, 0)

	zero := NewNodeUint256(0)
	zero.Accept(v)
	zero.Accept(v)
	v.pushOp(REVERT)
}

func fnWhen(v *BytecodeVisitor, args []Node) {
	assertArgsGte("when", args, 1)

	cond := args[0]
	body := args[1:]

	// Push condition onto stack.
	cond.Accept(v)

	// Invert the condition.
	v.pushOp(NOT)

	// Push a pointer and jump.
	pointer := v.pushPointer()
	v.pushOp(JUMPI)

	// Now push the body.
	VisitSequence(v, body)

	// Next, we're pushing a JUMPDEST that matches the JUMP
	// instruction, but not before we update the original pointer.
	v.output[pointer].code = fmt.Sprintf("61%04x", v.codeLength())
	v.pushOp(JUMPDEST)
}

type PreludeFunction func(v *BytecodeVisitor, args []Node)

func isPreludeFunc(tok string) (PreludeFunction, bool) {
	switch tok {
	case "%":
		// (% x y) returns remainder of x divided by y
		return fnMod, true
	case "+%":
		// (+% x y m) returns (x+y)%m
		return fnAddmod, true
	case "*%":
		// (*% x y m) returns (x*y)%m
		return fnMulmod, true
	case "expt":
		// (expt x y) returns x**y
		return fnExpt, true

	case "progn":
		return fnProgn, true
	case "revert":
		return fnRevert, true
	case "when":
		return fnWhen, true
	default:
		return nil, false
	}
}

// +-------------+
// | Environment |
// +-------------+

type Environment struct {
	// constants map[string]string
}
