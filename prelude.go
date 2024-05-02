package mist

import "fmt"

func assertNargsEq(fn string, call Node, want int) []Node {
	if call.NumChildren() != (want + 1) {
		panic(fmt.Sprintf(
			"%v: have %d arguments, want %d: %v",
			call.Origin,
			call.NumChildren(),
			(want + 1),
			&call,
		))
	}
	name := call.FunctionName()
	args := call.Children[1:]
	if fn != name {
		panic(fmt.Sprintf("%v: have %s, want %s", call.Origin, name, fn))
	}
	if have := len(args); have != want {
		panic(fmt.Sprintf(
			"wrong number of arguments for (%s): have %d, want %d",
			fn,
			have,
			want,
		))
	}
	return args
}

func assertNargsGte(fn string, call Node, want int) []Node {
	if call.NumChildren() < (want + 1) {
		panic(fmt.Sprintf(
			"%v: have %d arguments, want at least %d",
			call.Origin,
			call.NumChildren(),
			(want + 1),
		))
	}
	name := call.FunctionName()
	args := call.Children[1:]
	if fn != name {
		panic(fmt.Sprintf("%v: have %s, want %s", call.Origin, name, fn))
	}
	if have := len(args); have < want {
		panic(fmt.Sprintf(
			"wrong number of arguments for (%s): have %d, want at least %d",
			fn,
			have,
			want,
		))
	}
	return args
}

func handleNativeFunc(v *BytecodeVisitor, s *Scope, esp int, call Node) bool {
	ebp := esp
	fn := call.FunctionName()

	var (
		op  OpCode
		inp int // Number of input stack words.
		dir = 0
	)

	switch fn {
	case "stop":
		op, inp, dir = STOP, 0, -1
	// ADD is variadic.
	// MUL is variadic.
	case "-":
		op, inp, dir = SUB, 2, -1
	case "/":
		op, inp, dir = DIV, 2, -1
	// SDIV is NOT implemented.
	case "%":
		op, inp, dir = MOD, 2, -1
	// SMOD is NOT implemented.
	case "+%":
		op, inp, dir = ADDMOD, 3, -1
	case "*%":
		op, inp, dir = MULMOD, 3, -1
	case "**":
		fallthrough
	case "expt":
		op, inp, dir = EXP, 2, -1
	// SIGNEXTEND is NOT implemented.
	case "<":
		op, inp, dir = LT, 2, -1
	case ">":
		op, inp, dir = GT, 2, -1
	// SLT is NOT implemented.
	// SGT is NOT implemented.
	case "=":
		op, inp, dir = EQ, 2, -1
	case "not":
		fallthrough
	case "zerop":
		op, inp, dir = ISZERO, 1, -1
	// AND is variadic.
	// OR is variadic.
	// XOR is variadic.
	case "~":
		fallthrough
	case "lognot":
		op, inp, dir = NOT, 1, -1
	case "byte": // (byte byte-index word)
		op, inp, dir = BYTE, 2, -1
	case "<<": // (<< value count)
		op, inp, dir = SHL, 2, 1
	case ">>": // (>> value count)
		op, inp, dir = SHR, 2, 1
	// SAR is NOT implemented.
	// KECCAK256 is NOT implemented.
	case "address":
		op, inp, dir = ADDRESS, 0, -1
	case "balance":
		op, inp, dir = BALANCE, 1, -1
	case "origin":
		op, inp, dir = ORIGIN, 0, -1
	case "caller":
		op, inp, dir = CALLER, 0, -1
	case "call-value":
		op, inp, dir = CALLVALUE, 0, -1
	case "calldata-load": // (calldata-load word-index)
		op, inp, dir = CALLDATALOAD, 1, -1
	case "calldata-size":
		op, inp, dir = CALLDATASIZE, 0, -1
	// case "calldata-copy": // (calldata-copy mm-start id-offset length)
	// 	op, nargs, dir = CALLDATACOPY, 3, -1
	case "code-size":
		op, inp, dir = CODESIZE, 0, -1
	// case "code-copy": // (code-copy mm-start ib-offset length)
	// 	op, nargs, dir = CODECOPY, 3, -1
	case "gas-price":
		op, inp, dir = GASPRICE, 0, -1
	// EXTCODESIZE is NOT implemented.
	// EXTCODECOPY is NOT implemented.
	// RETURNDATASIZE is NOT implemented.
	// RETURNDATACOPY is NOT implemented.
	// EXTCODEHASH is NOT implemented.
	// BLOCKHASH is NOT implemented.
	case "coinbase":
		op, inp, dir = COINBASE, 0, -1
	case "timestamp":
		op, inp, dir = TIMESTAMP, 0, -1
	case "block-number":
		op, inp, dir = NUMBER, 0, -1
	case "prev-randao":
		op, inp, dir = PREVRANDAO, 0, -1
	case "gas-limit":
		op, inp, dir = GASLIMIT, 0, -1
	case "chain-id":
		op, inp, dir = CHAINID, 0, -1
	case "self-balance":
		op, inp, dir = SELFBALANCE, 0, -1
	case "base-fee":
		op, inp, dir = BASEFEE, 0, -1
	// case "pop"
	// case "mload"
	// case "mstore"
	// case "mstore8"
	// case "sload"
	// case "sstore"
	// case "jump"
	// case "jumpi"
	// case "program-counter": op, inp, dir = PC, 0, -1
	// case "memory-size": op, inp, dir = MSIZE, 0, -1
	case "available-gas":
		op, inp, dir = GAS, 0, -1
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
		args := assertNargsEq(fn, call, inp)
		esp += VisitSequence(v, s, esp, args, dir)
		delta := esp - ebp
		if delta != inp {
			panic("broken invariant")
		}
		v.addOp(op)
		return true
	}
	return false
}

func handleVariadicFunc(v *BytecodeVisitor, s *Scope, esp int, call Node) bool {
	var (
		op OpCode
		ok = false
	)

	fn := call.FunctionName()
	switch fn {
	case "+":
		op, ok = ADD, true
	case "*":
		op, ok = MUL, true
	case "&":
		fallthrough
	case "logand":
		op, ok = AND, true
	case "|":
		fallthrough
	case "logior":
		op, ok = OR, true
	case "^":
		fallthrough
	case "logxor":
		op, ok = XOR, true
	}

	if !ok {
		return false
	}

	args := assertNargsGte(fn, call, 2)
	last := len(args) - 1

	args[last].Accept(v, s, esp)
	esp += 1

	for i := last - 1; i >= 0; i-- {
		args[i].Accept(v, s, esp)
		esp += 1

		v.addOp(op)
		esp -= 1
	}

	return true
}

func handleBuiltinFunc(v *BytecodeVisitor, s *Scope, esp int, call Node) bool {
	fn := call.FunctionName()
	switch fn {
	case "defconst":
		fnDefconst(v, s, esp, call)
		return true
	case "defun":
		fnDefun(v, s, esp, call)
		return true // TODO
	// case "discard":
	// 	return fnDiscard(v, s, esp, call), true
	case "if":
		fnIf(v, s, esp, call)
		return true
	case "progn":
		fnProgn(v, s, esp, call)
		return true
	case "return": // (return value)
		fnReturn(v, s, esp, call)
		return true
	case "revert": // (revert value)
		fnRevert(v, s, esp, call)
		return true

	//
	// TODO: If I ever implement macros, these should be reimplemented!
	//

	case "unless":
		fnUnless(v, s, esp, call)
		return true
	case "when":
		fnWhen(v, s, esp, call)
		return true

	default:
		return false
	}
}

// Right now (defun) functions are inlined in the code.  Perhaps I
// should create a separate data segment for (defun) and introduce
// (definline) for inline functions?
func handleDefined(v *BytecodeVisitor, s *Scope, esp int, call Node) bool {
	ebp := esp

	name := call.FunctionName()
	args := assertNargsGte(name, call, 0)

	fn, ok := s.GetFunction(name)
	if !ok {
		return false
	}

	// Check the number of arguments.
	if len(args) != len(fn.Args) {
		panic(NewCompilationError(
			call.Origin,
			fmt.Sprintf(
				"wrong number of arguments for (%s): have %d, want %d",
				fn.Name,
				len(args),
				len(fn.Args),
			),
		))
	}

	// Create a child scope and evaluate all arguments.
	childScope := s.NewChildScope()
	esp += VisitSequence(v, s, esp, args, -1)

	for i := range fn.Args {
		identifier := fn.Args[i].ValueString
		position := ebp + len(fn.Args) - 1 - i
		childScope.SetStackVariable(identifier, StackVariable{
			Origin:     fn.Args[i].Origin,
			Identifier: identifier,
			Position:   position,
		})
	}

	fn.Body.Accept(v, childScope, esp)
	esp += 1

	if len(fn.Args) > 0 {
		v.addOp(OpCode(SWAP1 - 1 + len(fn.Args)))
		for range fn.Args {
			v.addOp(POP)
		}
	}

	return true
}

// +--------------------+
// | Built-in functions |
// +--------------------+

func fnDefconst(v *BytecodeVisitor, s *Scope, _ int, call Node) {
	args := assertNargsEq("defconst", call, 2)
	name, value := args[0], args[1]

	if !name.IsSymbol() {
		panic(fmt.Sprintf("%v: %v is not a symbol", value.Origin, value))
	}

	// Store into scope.
	s.Defconst(args[0].ValueString, args[1])

	// All expressions have a value.
	v.VisitNil()
}

func fnDefun(v *BytecodeVisitor, s *Scope, _ int, node Node) {
	fn, err := NewLispFunction(node)
	if err != nil {
		panic(err)
	}

	s.Defun(fn)

	// All expressions have a value.
	v.VisitNil()
}

// func fnDiscard(v *BytecodeVisitor, s *Scope, esp int, call Node) int {
// 	args := assertNargsEq("discard", call, 1)

// 	CheckAccept(args[0], v, s, esp, 1)
// 	v.addOp(POP)

// 	// Evaluates 3, consumes 3, pushes 1.
// 	// Consumes 1, pushes 0.
// 	return -1
// }

func fnProgn(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	ebp := esp
	args := assertNargsGte("progn", call, 0)

	// Empty (progn) results in nil.
	if len(args) <= 0 {
		v.VisitNil()
		return
	}

	// For each expression of progn's body, if it's not the last one,
	// discard (i.e. pop) it.  Otherwise, if it's the last, push it
	// onto the stack.
	for i := range args {
		args[i].Accept(v, s, esp)
		esp += 1

		if last := (i == len(args)-1); !last {
			// This is not the last expression, so discard its result.
			v.addOp(POP)
			esp -= 1
		} else {
			// This is the last expression, keep it on the stack.
		}
	}

	if esp != (ebp + 1) {
		panic("broken invariant")
	}
}

func fnReturn(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsEq("return", call, 1)

	v.pushU64(0x20)              // [20]
	esp += 1                     //
	v.pushU64(freeMemoryPointer) // [FP 20]
	esp += 1                     //
	v.addOp(MLOAD)               // [FM 20]
	esp += 0                     //
	args[0].Accept(v, s, esp)    // [RV FM 20]
	esp += 1                     //
	v.addOp(DUP2)                // [FM RV FM 20]
	esp += 1                     //
	v.addOp(MSTORE)              // [FM 20]
	esp -= 2                     //
	v.addOp(RETURN)              // []
	esp -= 2                     //
}

func fnRevert(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsEq("revert", call, 1)

	v.pushU64(0x20)              // [20]
	esp += 1                     //
	v.pushU64(freeMemoryPointer) // [FP 20]
	esp += 1                     //
	v.addOp(MLOAD)               // [FM 20]
	esp += 0                     //
	args[0].Accept(v, s, esp)    // [RV FM 20]
	esp += 1                     //
	v.addOp(DUP2)                // [FM RV FM 20]
	esp += 1                     //
	v.addOp(MSTORE)              // [FM 20]
	esp -= 2                     //
	v.addOp(REVERT)              // []
	esp -= 2                     //
}

func fnStop(v *BytecodeVisitor, _ *Scope, _ int, call Node) {
	assertNargsEq("stop", call, 0)
	v.addOp(STOP)
}

func fnUnless(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsGte("unless", call, 1)

	// Prepare condition.
	cond := args[0]

	// Prepare the `then` branch by wrapping it in a `progn`.
	body := args[1:]
	then := NewNodeNil(NewOriginEmpty())

	if len(body) > 0 {
		then = NewNodeProgn()
		then.AddChildren(body)
	}

	// Prepare the `else` branch.
	noop := NewNodeNil(NewOriginEmpty())

	replacement := NewNodeList(call.Origin)
	replacement.AddChild(NewNodeSymbol("if", NewOriginEmpty()))
	replacement.AddChild(cond)
	replacement.AddChild(noop)
	replacement.AddChild(then)
	fnIf(v, s, esp, replacement)
}

func fnIf(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsEq("if", call, 3)
	cond, yes, no := args[0], args[1], args[2]

	// Push the condition.
	cond.Accept(v, s, esp)
	esp += 1

	// Jump to the `then` branch if condition holds.
	dest := newSegmentJumpdest()
	v.addPointer(dest.id) // esp += 1
	v.addOp(JUMPI)        // esp -= 2
	esp -= 1

	// Otherwise, keep executing the `else` and jump after the `then`
	// at the end.
	no.Accept(v, s, esp) // Pushing `no`, esp += 1
	after := newSegmentJumpdest()
	v.addPointer(after.id) // esp += 1
	v.addOp(JUMP)          // esp -= 1

	// Now add the `then`.
	v.addSegment(dest)
	yes.Accept(v, s, esp) // Pushing `yes`, esp += 1

	// Add the `after` label.
	v.addSegment(after)

	// Either `yes` or `no` was evaluated, but not both.
}

func fnWhen(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsGte("when", call, 1)

	// Prepare condition.
	cond := args[0]

	// Prepare the `then` branch by wrapping it in a `progn`.
	body := args[1:]
	then := NewNodeNil(NewOriginEmpty())

	if len(body) > 0 {
		then = NewNodeProgn()
		then.AddChildren(body)
	}

	// Prepare the `else` branch.
	noop := NewNodeNil(NewOriginEmpty())

	replacement := NewNodeList(call.Origin)
	replacement.AddChild(NewNodeSymbol("if", NewOriginEmpty()))
	replacement.AddChild(cond)
	replacement.AddChild(then)
	replacement.AddChild(noop)
	fnIf(v, s, esp, replacement)
}

// +----------------------+
// | Contract constructor |
// +----------------------+

func MakeConstructor(deployedBytecode string) string {
	v := NewBytecodeVisitor(false)

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
	v.addSegment(label)

	return v.String()
}
