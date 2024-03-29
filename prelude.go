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
	// SIGNEXTEND
	case "=":
		return EQ, 2, true
	case "not":
		return ISZERO, 1, true
	case "zerop":
		return ISZERO, 1, true
	case "~":
		fallthrough
	case "lognot":
		return NOT, 1, true
	case "byte": // (byte word which)
		return BYTE, 2, true

	case "<<": // (<< value count)
		return SHL, 2, true
	case ">>": // (>> value count)
		return SHR, 2, true // TODO: SAR

	// TODO: case "keccak256": KECCAK256, 2, true

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
	case "calldata-load": // (calldata-load start)
		return CALLDATALOAD, 1, true
	case "calldata-size":
		return CALLDATASIZE, 0, true
	case "calldata-copy": // (calldata-copy length id-offset mm-start)
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
	// 	return POP, 0, true
	// case "mload": // (mload start)
	// 	return MLOAD, 1, true
	// case "mstore": // (mstore value start)
	// 	return MSTORE, 2, true
	// case "mstore8":
	// 	return MSTORE8, 2, true
	// case "sload": // (sload word-index)
	// 	return SLOAD, 1, true
	// case "sstore": // (sstore value word-index)
	// 	return SSTORE, 2, true
	// case "jump":
	//
	// case "jumpi":
	//
	// case "program-counter":
	// 	return PC, 0, true
	// case "memory-size":
	// 	return MSIZE, 0, true
	case "available-gas":
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
		"~", "lognot",
		"&", "logand",
		"|", "logior",
		"^", "logxor",
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

func fnCmpGT(v *BytecodeVisitor, args []Node) {
	assertArgsEq(">", args, 2)

	x, y := args[0], args[1]
	y.Accept(v)
	x.Accept(v)
	v.pushOp(GT)
}

func fnCmpLT(v *BytecodeVisitor, args []Node) {
	assertArgsEq("<", args, 2)

	x, y := args[0], args[1]
	y.Accept(v)
	x.Accept(v)
	v.pushOp(LT)
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
	assertArgsEq("return", args, 1)

	v.pushUnsigned(0x20)          // [20]
	v.pushUnsigned(freeMemoryPtr) // [FP 20]
	v.pushOp(MLOAD)               // [FM 20]
	args[0].Accept(v)             // [RV FM 20]
	v.pushOp(DUP2)                // [FM RV FM 20]
	v.pushOp(MSTORE)              // [FM 20]
	v.pushOp(RETURN)              // []
	v.pushOp(INVALID)
}

func fnRevert(v *BytecodeVisitor, args []Node) {
	assertArgsEq("revert", args, 0)

	zero := NewNodeUint256(0)
	zero.Accept(v)

	v.pushOp(DUP1)
	v.pushOp(REVERT)
}

func fnUnless(v *BytecodeVisitor, args []Node) {
	assertArgsGte("unless", args, 1)

	cond := args[0]
	body := args[1:]

	// Push condition onto stack.
	cond.Accept(v)

	// Push a pointer and jump.
	pointer := v.pushPointer()
	v.pushOp(JUMPI)

	// Now push the body.
	VisitSequence(v, body)

	// Next, we're pushing a JUMPDEST that matches the JUMP
	// instruction, but not before we update the original pointer to
	// point to the address of that JUMPDEST.
	dest := v.codeLength()
	code := fmt.Sprintf("%02x%04x", byte(PUSH2), dest)
	if len(code) != 6 {
		panic("TODO")
	}
	v.segments[pointer] = code
	v.pushOp(JUMPDEST)
}

func fnWhen(v *BytecodeVisitor, args []Node) {
	assertArgsGte("when", args, 1)

	cond := args[0]
	body := args[1:]

	// Push condition onto stack.
	cond.Accept(v)

	// Invert the condition.
	v.pushOp(ISZERO)

	// Push a pointer and jump.
	pointer := v.pushPointer()
	v.pushOp(JUMPI)

	// Now push the body.
	VisitSequence(v, body)

	// Next, we're pushing a JUMPDEST that matches the JUMP
	// instruction, but not before we update the original pointer to
	// point to the address of that JUMPDEST.
	dest := v.codeLength()
	code := fmt.Sprintf("%02x%04x", byte(PUSH2), dest)
	if len(code) != 6 {
		panic("TODO")
	}
	v.segments[pointer] = code
	v.pushOp(JUMPDEST)
}

type PreludeFunction func(v *BytecodeVisitor, args []Node)

func isPreludeFunc(tok string) (PreludeFunction, bool) {
	switch tok {

	case "<":
		return fnCmpLT, true // TODO: SLT
	case ">":
		return fnCmpGT, true // TODO: SGT

	case "%":
		// (% x y) returns x%y, the remainder of x divided by y
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
	case "return":
		return fnReturn, true
	case "revert":
		return fnRevert, true
	case "unless":
		return fnUnless, true
	case "when":
		return fnWhen, true
	default:
		return nil, false
	}
}

// +-------------+
// | Constructor |
// +-------------+

func MakeConstructor(deployedBytecode string) string {
	v := NewBytecodeVisitor()

	length := uint64(len(deployedBytecode) / 2)
	v.pushUnsigned(length)
	v.pushOp(DUP1)
	pointer := v.pushPointer()
	v.pushUnsigned(0)
	v.pushOp(CODECOPY)
	v.pushUnsigned(0)
	v.pushOp(RETURN)
	v.pushOp(INVALID)

	code := fmt.Sprintf("%02x%04x", byte(PUSH2), v.codeLength())
	if len(code) != 6 {
		panic("TODO")
	}
	v.segments[pointer] = code

	return v.String()
}

// +-------------+
// | Environment |
// +-------------+

type Environment struct {
	// constants map[string]string
}
