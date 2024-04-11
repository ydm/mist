package mist

import "fmt"

func assertNumArgsEq(fn string, args []Node, want int) {
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

func handleNativeFunc(v *BytecodeVisitor, fn string, args []Node) bool {
	var (
		op    OpCode
		nargs int
		dir   = 0
	)

	switch fn {
	case "stop":
		op, nargs, dir = STOP, 0, -1
	// ADD is variadic.
	// MUL is variadic.
	// SUB is variadic.
	// DIV is variadic.
	// SDIV is NOT implemented.
	case "%":
		op, nargs, dir = MOD, 2, -1
	// SMOD is NOT implemented.
	case "+%":
		op, nargs, dir = ADDMOD, 2, -1
	case "*%":
		op, nargs, dir = MULMOD, 2, -1
	case "**":
		fallthrough
	case "expt":
		op, nargs, dir = EXP, 2, -1
	// SIGNEXTEND is NOT implemented.
	// LT (<) is variadic.
	// GT (>) is variadic.
	// SLT is NOT implemented.
	// SGT is NOT implemented.
	// EQ (=) is variadic.
	case "not":
		fallthrough
	case "zerop":
		op, nargs, dir = ISZERO, 1, -1
	// AND (logand, &) is variadic.
	// OR (logior, |) is variadic.
	// XOR (logxor, ^) is variadic.
	case "~":
		fallthrough
	case "lognot":
		op, nargs, dir = NOT, 1, -1
	case "byte": // (byte byte-index word)
		op, nargs, dir = BYTE, 2, -1
	case "<<": // (<< value count)
		op, nargs, dir = SHL, 2, 1
	case ">>": // (>> value count)
		op, nargs, dir = SHR, 2, 1
	// SAR is NOT implemented.
	// KECCAK256 is NOT implemented.
	case "address":
		op, nargs, dir = ADDRESS, 0, -1
	case "balance":
		op, nargs, dir = BALANCE, 1, -1
	case "origin":
		op, nargs, dir = ORIGIN, 0, -1
	case "caller":
		op, nargs, dir = CALLER, 0, -1
	case "call-value":
		op, nargs, dir = CALLVALUE, 0, -1
	case "calldata-load": // (calldata-load word-index)
		op, nargs, dir = CALLDATALOAD, 1, -1
	case "calldata-size":
		op, nargs, dir = CALLDATASIZE, 0, -1
	// case "calldata-copy": // (calldata-copy mm-start id-offset length)
	// 	op, nargs, dir = CALLDATACOPY, 3, -1
	case "code-size":
		op, nargs, dir = CODESIZE, 0, -1
	// case "code-copy": // (code-copy mm-start ib-offset length)
	// 	op, nargs, dir = CODECOPY, 3, -1
	case "gas-price":
		op, nargs, dir = GASPRICE, 0, -1
	// EXTCODESIZE is NOT implemented.
	// EXTCODECOPY is NOT implemented.
	// RETURNDATASIZE is NOT implemented.
	// RETURNDATACOPY is NOT implemented.
	// EXTCODEHASH is NOT implemented.
	// BLOCKHASH is NOT implemented.
	case "coinbase":
		op, nargs, dir = COINBASE, 0, -1
	case "timestamp":
		op, nargs, dir = TIMESTAMP, 0, -1
	case "block-number":
		op, nargs, dir = NUMBER, 0, -1
	case "prev-randao":
		op, nargs, dir = PREVRANDAO, 0, -1
	case "gas-limit":
		op, nargs, dir = GASLIMIT, 0, -1
	case "chain-id":
		op, nargs, dir = CHAINID, 0, -1
	case "self-balance":
		op, nargs, dir = SELFBALANCE, 0, -1
	case "base-fee":
		op, nargs, dir = BASEFEE, 0, -1
	// case "pop"
	// case "mload"
	// case "mstore"
	// case "mstore8"
	// case "sload"
	// case "sstore"
	// case "jump"
	// case "jumpi"
	case "program-counter":
		op, nargs, dir = PC, 0, -1
	case "memory-size":
		op, nargs, dir = MSIZE, 0, -1
	case "available-gas":
		op, nargs, dir = GAS, 0, -1
	// case "jumpdest"
	// case "push1..16"
	// case "dup1..16"
	// case "swap1..16"
	// case "log0..4"
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
	default:
	}

	if dir != 0 {
		assertNumArgsEq(fn, args, nargs)
		VisitSequence(v, args, dir)
		v.addOp(op)
		return true
	}
	return false
}

func handleVariadicFunc(v *BytecodeVisitor, fn string, args []Node) bool {
	var (
		op    OpCode
		match = false
	)

	switch fn {
	case "+":
		op, match = ADD, true
	case "*":
		op, match = MUL, true
	case "-":
		op, match = SUB, true
	case "/":
		op, match = DIV, true
	case "<":
		op, match = LT, true
	case ">":
		op, match = GT, true
	case "=":
		op, match = EQ, true
	case "&":
		fallthrough
	case "logand":
		op, match = AND, true
	case "|":
		fallthrough
	case "logior":
		op, match = OR, true
	case "^":
		fallthrough
	case "logxor":
		op, match = XOR, true
	}

	if match {
		assertArgsGte(fn, args, 2)
		last := len(args) - 1
		args[last].Accept(v)
		for i := last - 1; i >= 0; i-- {
			args[i].Accept(v)
			v.addOp(op)
		}
		return true
	}
	return false
}

func handleInlineFunc(v *BytecodeVisitor, fn string, args []Node) bool {
	switch fn {
	case "%":
		// (% x y) returns x%y, the remainder of x divided by y
		fnMod(v, args)
		return true
	case "+%":
		// (+% x y m) returns (x+y)%m
		fnAddmod(v, args)
		return true
	case "*%":
		// (*% x y m) returns (x*y)%m
		fnMulmod(v, args)
		return true
	case "discard":
		fnDiscard(v, args)
		return true
	case "if":
		fnIf(v, args)
		return true
	case "progn":
		fnProgn(v, args)
		return true
	case "return":
		// (return value)
		fnReturn(v, args)
		return true
	case "revert":
		// (revert value)
		fnRevert(v, args)
		return true
	case "setq":
		return true
	case "unless":
		fnUnless(v, args)
		return true
	case "when":
		fnWhen(v, args)
		return true

	default:
		return false
	}
}

// +------------------+
// | Inline functions |
// +------------------+

func fnAddmod(v *BytecodeVisitor, args []Node) {
	assertNumArgsEq("+%", args, 3)

	x, y, m := args[0], args[1], args[2]
	m.Accept(v)
	y.Accept(v)
	x.Accept(v)
	v.addOp(ADDMOD)
}

func fnDiscard(v *BytecodeVisitor, args []Node) {
	assertNumArgsEq("discard", args, 1)

	args[0].Accept(v)
	v.addOp(POP)
}

func fnMod(v *BytecodeVisitor, args []Node) {
	assertNumArgsEq("%", args, 2)

	x, y := args[0], args[1]
	y.Accept(v)
	x.Accept(v)
	v.addOp(MOD)
}

func fnMulmod(v *BytecodeVisitor, args []Node) {
	assertNumArgsEq("*%", args, 3)

	x, y, m := args[0], args[1], args[2]
	m.Accept(v)
	y.Accept(v)
	x.Accept(v)
	v.addOp(ADDMOD)
}

func fnProgn(v *BytecodeVisitor, args []Node) {
	for i := range args {
		last := i == len(args)-1
		args[i].Accept(v)
		if !last {
			v.addOp(POP)
		}
	}
}

func fnReturn(v *BytecodeVisitor, args []Node) {
	assertNumArgsEq("return", args, 1)

	v.pushU64(0x20)              // [20]
	v.pushU64(freeMemoryPointer) // [FP 20]
	v.addOp(MLOAD)               // [FM 20]
	args[0].Accept(v)            // [RV FM 20]
	v.addOp(DUP2)                // [FM RV FM 20]
	v.addOp(MSTORE)              // [FM 20]
	v.addOp(RETURN)              // []
}

func fnRevert(v *BytecodeVisitor, args []Node) {
	assertNumArgsEq("revert", args, 1)

	v.pushU64(0x20)              // [20]
	v.pushU64(freeMemoryPointer) // [FP 20]
	v.addOp(MLOAD)               // [FM 20]
	args[0].Accept(v)            // [RV FM 20]
	v.addOp(DUP2)                // [FM RV FM 20]
	v.addOp(MSTORE)              // [FM 20]
	v.addOp(REVERT)              // []
}

func fnUnless(v *BytecodeVisitor, args []Node) {
	assertArgsGte("unless", args, 1)

	// Prepare condition.
	cond := args[0]

	// Prepare the `then` branch.
	body := args[1:]
	yes := NewNodeList(NewOriginEmpty())
	if len(args) > 0 {
		yes.Origin = args[0].Origin
	}
	yes.AddChildren(body)

	// Prepare the `else` branch.
	no := NewNodeNil(NewOriginEmpty())

	fnIf(v, []Node{cond, yes, no})
}

func fnIf(v *BytecodeVisitor, args []Node) {
	assertNumArgsEq("if", args, 3)
	cond, yes, no := args[0], args[1], args[2]

	// Push the condition.
	cond.Accept(v)

	// Jump to the `then` branch if condition holds.
	dest := newSegmentJumpdest()
	v.addPointer(dest.id)
	v.addOp(JUMPI)

	// Otherwise, keep executing the `else` and jump after the `then`
	// at the end.
	no.Accept(v)
	after := newSegmentJumpdest()
	v.addPointer(after.id)
	v.addOp(JUMP)

	// Now add the `then`.
	v.addSegment(dest)
	yes.Accept(v)

	// Add the `after` label.
	v.addSegment(after)
}

func fnWhen(v *BytecodeVisitor, args []Node) {
	assertArgsGte("when", args, 1)

	// Prepare condition.
	cond := args[0]

	// Prepare the `then` branch.
	body := args[1:]
	yes := NewNodeNil(NewOriginEmpty())

	if len(body) > 0 {
		yes = NewNodeProgn(args[0].Origin)
		yes.AddChildren(body)
	}

	// Prepare the `else` branch.
	no := NewNodeNil(NewOriginEmpty())

	fnIf(v, []Node{cond, yes, no})
}

// +----------------------+
// | Contract constructor |
// +----------------------+

func MakeConstructor(deployedBytecode string) string {
	v := NewBytecodeVisitor()

	label := newEmptySegment()

	// (codecopy mm-offset@0 ib-offset@1 length@2)
	// has the following effect
	// M[mm-offset:+length] = Ib[ib-offset:+length]

	length := len(deployedBytecode) / 2
	v.pushU64(uint64(length)) // L
	v.addOp(DUP1)             // L L
	v.addPointer(label.id)    // P L L
	v.pushU64(0)              // 0 P L L
	v.addOp(CODECOPY)         // (codecopy 0 P L)
	v.pushU64(0)              // 0 L
	v.addOp(RETURN)           // return M[0:L]
	v.addOp(INVALID)
	v.addSegment(label)

	return v.String()
}
