package mist

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/vm"
)

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
		op  vm.OpCode
		inp int // Number of input stack words.
		dir = 0
	)

	switch fn {
	case "stop":
		op, inp, dir = vm.STOP, 0, -1
	// ADD is variadic.
	// MUL is variadic.
	case "-":
		op, inp, dir = vm.SUB, 2, -1
	case "/":
		op, inp, dir = vm.DIV, 2, -1
	// SDIV is NOT implemented.
	case "%":
		op, inp, dir = vm.MOD, 2, -1
	// SMOD is NOT implemented.
	case "+%":
		op, inp, dir = vm.ADDMOD, 3, -1
	case "*%":
		op, inp, dir = vm.MULMOD, 3, -1
	case "**":
		fallthrough
	case "expt":
		op, inp, dir = vm.EXP, 2, -1
	// SIGNEXTEND is NOT implemented.
	case "<":
		op, inp, dir = vm.LT, 2, -1
	case ">":
		op, inp, dir = vm.GT, 2, -1
	// SLT is NOT implemented.
	// SGT is NOT implemented.
	case "=":
		op, inp, dir = vm.EQ, 2, -1
	case "not":
		fallthrough
	case "zerop":
		op, inp, dir = vm.ISZERO, 1, -1
	// AND is variadic.
	// OR is variadic.
	// XOR is variadic.
	case "~":
		fallthrough
	case "lognot":
		op, inp, dir = vm.NOT, 1, -1
	case "byte": // (byte byte-index word)
		op, inp, dir = vm.BYTE, 2, -1
	case "<<": // (<< value count)
		op, inp, dir = vm.SHL, 2, 1
	case ">>": // (>> value count)
		op, inp, dir = vm.SHR, 2, 1
	// SAR is NOT implemented.
	// KECCAK256 is NOT implemented.
	case "address":
		op, inp, dir = vm.ADDRESS, 0, -1
	case "balance":
		op, inp, dir = vm.BALANCE, 1, -1
	case "origin":
		op, inp, dir = vm.ORIGIN, 0, -1
	case "caller":
		op, inp, dir = vm.CALLER, 0, -1
	case "call-value":
		op, inp, dir = vm.CALLVALUE, 0, -1
	case "calldata-load": // (calldata-load word-index)
		op, inp, dir = vm.CALLDATALOAD, 1, -1
	case "calldata-size":
		op, inp, dir = vm.CALLDATASIZE, 0, -1
	// case "calldata-copy": // (calldata-copy mm-start id-offset length)
	// 	op, nargs, dir = CALLDATACOPY, 3, -1
	case "code-size":
		op, inp, dir = vm.CODESIZE, 0, -1
	// case "code-copy": // (code-copy mm-start ib-offset length)
	// 	op, nargs, dir = CODECOPY, 3, -1
	case "gas-price":
		op, inp, dir = vm.GASPRICE, 0, -1
	// EXTCODESIZE is NOT implemented.
	// EXTCODECOPY is NOT implemented.
	// RETURNDATASIZE is NOT implemented.
	// RETURNDATACOPY is NOT implemented.
	// EXTCODEHASH is NOT implemented.
	// BLOCKHASH is NOT implemented.
	case "coinbase":
		op, inp, dir = vm.COINBASE, 0, -1
	case "timestamp":
		op, inp, dir = vm.TIMESTAMP, 0, -1
	case "block-number":
		op, inp, dir = vm.NUMBER, 0, -1
	case "prev-randao":
		op, inp, dir = vm.PREVRANDAO, 0, -1
	case "gas-limit":
		op, inp, dir = vm.GASLIMIT, 0, -1
	case "chain-id":
		op, inp, dir = vm.CHAINID, 0, -1
	case "self-balance":
		op, inp, dir = vm.SELFBALANCE, 0, -1
	case "base-fee":
		op, inp, dir = vm.BASEFEE, 0, -1
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
		op, inp, dir = vm.GAS, 0, -1
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
		op vm.OpCode
		ok = false
	)

	fn := call.FunctionName()
	switch fn {
	case "+":
		op, ok = vm.ADD, true
	case "*":
		op, ok = vm.MUL, true
	case "&":
		fallthrough
	case "logand":
		op, ok = vm.AND, true
	case "|":
		fallthrough
	case "logior":
		op, ok = vm.OR, true
	case "^":
		fallthrough
	case "logxor":
		op, ok = vm.XOR, true
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
	case "and":
		fnAnd(v, s, esp, call)
		return true
	case "case":
		fnCase(v, s, esp, call)
		return true
	case "defconst":
		fnDefconst(v, s, esp, call)
		return true
	case "defun":
		fnDefun(v, s, esp, call)
		return true // TODO
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
	case "selector":
		fnSelector(v, s, esp, call)
		return true
	default:
		return false
	}
}

// Right now (defun) functions are inlined in the code.  Perhaps I
// should create a separate data segment for (defun) and introduce
// (definline) for inline functions?
func handleDefinedFunc(v *BytecodeVisitor, s *Scope, esp int, call Node) bool {
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
		v.addOp(vm.OpCode(vm.SWAP1 - 1 + len(fn.Args)))
		for range fn.Args {
			v.addOp(vm.POP)
		}
	}

	return true
}

// +--------------------+
// | Built-in functions |
// +--------------------+

func fnAnd(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsGte("and", call, 0)

	after := newSegmentJumpdest()

	v.VisitT()
	esp += 1

	last := len(args) - 1
	for i := range args {
		v.addOp(vm.POP)
		esp -= 1

		args[i].Accept(v, s, esp)
		esp += 1

		if i != last {
			v.addOp(vm.DUP1)
			esp += 1

			v.addOp(vm.ISZERO)
			esp += 0

			v.addPointer(after.id)
			esp += 1

			v.addOp(vm.JUMPI)
			esp -= 2
		}
	}

	if len(args) > 1 {
		v.addSegment(after)
	}
}

func fnCase(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsGte("case", call, 1)

	// Evaluate the given value.
	args[0].Accept(v, s, esp) // [XX]
	esp += 1

	// If this (case) has no clauses, push 0 and that's it.
	if len(args) <= 1 {
		v.pushU64(0)
		esp += 1
		return
	}

	n := len(args) - 2
	if n < 0 {
		n = 0
	}
	labels := make([]segment, n)
	for i := 0; i < n; i++ {
		labels[i] = newSegmentJumpdest()
	}

	after := newSegmentJumpdest()

	for i := 1; i < len(args); i++ {
		clause := args[i]

		if !clause.IsList() || clause.NumChildren() < 2 {
			panic("TODO")
		}

		// Push the label first.  Each case except the first one is
		// labeled.
		if i >= 2 {
			v.addSegment(labels[i-2])
		}

		// If values are not equal, jump to next clause.
		v.addOp(vm.DUP1)                     // [XX XX]
		esp += 1                             //
		clause.Children[0].Accept(v, s, esp) // [YY XX XX]
		esp += 1                             //
		v.addOp(vm.EQ)                       // [EQ XX]
		esp -= 1                             //
		v.addOp(vm.ISZERO)                   // [!E XX]
		esp += 0                             //
		v.addPointer(labels[i-1].id)         // [NX !E XX]
		esp += 1                             //
		v.addOp(vm.JUMPI)                    // [XX]
		esp -= 2                             //

		// If we got here, then this is the right clause to execute.
		// Evaluate the body.
		progn := NewNodeProgn()                //
		progn.AddChildren(clause.Children[1:]) //
		progn.Accept(v, s, esp)                // [AA XX]
		esp += 1                               //

		// Jump to the `after` label.
		v.addPointer(after.id) // [PP AA XX]
		esp += 1               //
		v.addOp(vm.JUMP)       // [AA XX]
		esp -= 2               //
	}

	v.addSegment(after)

	v.addOp(vm.SWAP1)			// [XX AA]
	esp += 0					//
	v.addOp(vm.POP)				// [AA]
	esp -= 1					// 
}

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

func fnIf(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsEq("if", call, 3)
	cond, yes, no := args[0], args[1], args[2]

	// Push the condition.
	cond.Accept(v, s, esp)
	esp += 1

	// Jump to the `then` branch if condition holds.
	dest := newSegmentJumpdest()
	v.addPointer(dest.id) // esp += 1
	v.addOp(vm.JUMPI)     // esp -= 2
	esp -= 1

	// Otherwise, keep executing the `else` and jump after the `then`
	// at the end.
	no.Accept(v, s, esp) // Pushing `no`, esp += 1
	after := newSegmentJumpdest()
	v.addPointer(after.id) // esp += 1
	v.addOp(vm.JUMP)       // esp -= 1

	// Now add the `then`.
	v.addSegment(dest)
	yes.Accept(v, s, esp) // Pushing `yes`, esp += 1

	// Add the `after` label.
	v.addSegment(after)

	// Either `yes` or `no` was evaluated, but not both.
}

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
	last := len(args) - 1
	for i := range args {
		args[i].Accept(v, s, esp)
		esp += 1

		if i != last {
			// This is not the last expression, so discard its result.
			v.addOp(vm.POP)
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
	v.addOp(vm.MLOAD)            // [FM 20]
	esp += 0                     //
	args[0].Accept(v, s, esp)    // [RV FM 20]
	esp += 1                     //
	v.addOp(vm.DUP2)             // [FM RV FM 20]
	esp += 1                     //
	v.addOp(vm.MSTORE)           // [FM 20]
	esp -= 2                     //
	v.addOp(vm.RETURN)           // []
	esp -= 2                     //
}

func fnRevert(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsEq("revert", call, 1)

	if args[0].Type == NodeString {
		encoded := EncodeWithSignature("Error(string)", args[0].ValueString)

		// Load free memory pointer.
		v.pushU64(freeMemoryPointer) // [FP]
		esp += 1                     //
		v.addOp(vm.MLOAD)            // [FM]
		esp += 0                     //

		// Push and store selector.
		v.addOp(vm.PUSH32)                 //
		v.addCode(padRight32(encoded[:8])) // [ER FM]
		esp += 1                           //
		v.addOp(vm.DUP2)                   // [FM ER FM]
		esp += 1                           //
		v.addOp(vm.MSTORE)                 // [FM], m[FM]=[ER]
		esp -= 2                           //

		// Push and store each word of 32 bytes (== 64 hex chars).
		n := uint64(len(encoded))
		for i := uint64(8); i < n; i += 64 {
			word := encoded[i : i+64]

			v.addOp(vm.PUSH32) //
			v.addCode(word)    // [WO FM]
			esp += 1           //
			v.addOp(vm.DUP2)   // [FM WO FM]
			esp += 1           //
			v.pushU64(i / 2)   // [OF FM WO FM], OF stands for offset
			esp += 1           //
			v.addOp(vm.ADD)    // [FO WO FM], FO=FM+OF
			esp -= 1           //
			v.addOp(vm.MSTORE) // [FM], m[FO]=WO
			esp -= 2           //
		}

		v.pushU64(n / 2)   //
		esp += 1           // [LE FM], esp=2
		v.addOp(vm.SWAP1)  // [FM LE], esp=2
		v.addOp(vm.REVERT) //
		esp -= 2           // [], esp=0
	} else {
		v.pushU64(0x20)              // [20]
		esp += 1                     //
		v.pushU64(freeMemoryPointer) // [FP 20]
		esp += 1                     //
		v.addOp(vm.MLOAD)            // [FM 20]
		esp += 0                     //
		args[0].Accept(v, s, esp)    // [RV FM 20]
		esp += 1                     //
		v.addOp(vm.DUP2)             // [FM RV FM 20]
		esp += 1                     //
		v.addOp(vm.MSTORE)           // [FM 20]
		esp -= 2                     //
		v.addOp(vm.REVERT)           // []
		esp -= 2                     //
	}
}

func fnSelector(v *BytecodeVisitor, _ *Scope, _ int, call Node) {
	args := assertNargsEq("selector", call, 1)

	if !args[0].IsString() {
		panic("TODO")
	}

	h := Keccak256Hash([]byte(args[0].ValueString))

	v.addOp(vm.PUSH4)
	v.addCode(fmt.Sprintf("%02x%02x%02x%02x", h[0], h[1], h[2], h[3]))
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
	v.addOp(vm.DUP1)          // L L
	v.addPointer(label.id)    // P L L
	v.pushU64(0)              // 0 P L L
	v.addOp(vm.CODECOPY)      // (codecopy 0 P L)
	v.pushU64(0)              // 0 L
	v.addOp(vm.RETURN)        // return M[0:L]
	v.addSegment(label)

	return v.String()
}
