package mist

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
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
			"wrong number of arguments for (%s): want %d, have %d",
			fn,
			want,
			have,
		))
	}
	return args
}

func assertNargsGte(fn string, call Node, want int) []Node {
	if call.NumChildren() < (want + 1) {
		panic(fmt.Sprintf(
			"%v: %s: have %d arguments, want at least %d: %v",
			call.Origin,
			fn,
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
	if have := len(args); have < want {
		panic(fmt.Sprintf(
			"wrong number of arguments for (%s): want at least %d, have %d",
			fn,
			want,
			have,
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
	case "current-address":
		op, inp, dir = vm.ADDRESS, 0, -1
	case "balance":
		op, inp, dir = vm.BALANCE, 1, -1
	case "origin":
		op, inp, dir = vm.ORIGIN, 0, -1
	case "caller":
		op, inp, dir = vm.CALLER, 0, -1
	case "call-value":
		op, inp, dir = vm.CALLVALUE, 0, -1
	case "calldata-load": // (calldata-load byte-index)
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
		return true
	case "defvar":
		fnDefvar(v, s, esp, call)
		return true
	case "ether":
		fnEther(v, s, esp, call)
		return true
	case "gethash": // (gethash table keys...)
		fnGethash(v, s, esp, call)
		return true
	case "if":
		fnIf(v, s, esp, call)
		return true
	case "progn":
		fnProgn(v, s, esp, call)
		return true
	case "puthash": // (puthash table value keys...)
		fnPuthash(v, s, esp, call)
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
	case "setq":
		fnSetq(v, s, esp, call)
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
				"wrong number of arguments for (%s): want %d, have %d",
				fn.Name,
				len(fn.Args),
				len(args),
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

	// If the clause list doesn't end up with (otherwise expression),
	// add (otherwise nil) manually.
	hasOtherwise := false
	for i := 1; i < len(args); i++ {
		if !args[i].IsList() || args[i].NumChildren() < 2 {
			panic(fmt.Sprintf(
				"%v: wrong argument type for (case): want (key body...), have: %v",
				call.Origin,
				&args[i],
			))
		}
		if args[i].Children[0].IsThisSymbol("otherwise") || args[i].Children[0].IsThisSymbol("t") {
			hasOtherwise = true
			if i != len(args)-1 {
				// We do have an otherwise clause, but
				// it's not last in the list.
				panic("misplaced otherwise or t clause")
			}
		}
	}

	head := args[0]
	tail := args[1:]
	if !hasOtherwise {
		otherwise := NewNodeList(NewOriginEmpty())
		otherwise.AddChild(NewNodeSymbol("otherwise", NewOriginEmpty()))
		otherwise.AddChild(NewNodeNil(NewOriginEmpty()))
		tail = append(tail, otherwise)
	}

	// If this (case) has no clauses or just a single otherwise,
	// push it without evaluating the switch value.
	if len(tail) == 1 {
		clause := tail[0]
		body := clause.Children[1]
		if clause.NumChildren() > 2 {
			body = NewNodeProgn()
			body.AddChildren(clause.Children[1:])
		}
		body.Accept(v, s, esp)
		esp += 1
		return
	}

	// Evaluate the given switch value.
	head.Accept(v, s, esp) // [XX]
	esp += 1

	// For each clause after the first one there should be a
	// label.  Indices correspond to clauses.  The very last one
	// -- `after` -- is placed after the whole (case).
	after := len(tail)
	labels := make([]segment, after+1)
	for i := 1; i < after+1; i++ {
		labels[i] = newSegmentJumpdest()
	}

	// For all the clauses except the last one (which is always
	// `otherwise`), compare the key and eventually execute the
	// body.
	last := len(tail) - 1
	for i := 0; i < last; i++ {
		// Starting stack is always [XX].

		clause := tail[i]
		if !clause.IsList() || clause.NumChildren() < 2 {
			panic("TODO")
		}

		// Extract the key.
		key := clause.Children[0]

		// Extract the body.  If it's more than a single
		// expression, wrap it all in a (progn).
		body := clause.Children[1]
		if clause.NumChildren() > 2 {
			body = NewNodeProgn()
			body.AddChildren(clause.Children[1:])
		}

		// Push the label first.  Each case except the first
		// one is labeled.
		if i >= 1 {
			v.addSegment(labels[i])
		}

		// If values are not equal, jump to the next clause.
		v.addOp(vm.DUP1)             // [XX XX]
		esp += 1                     //
		key.Accept(v, s, esp)        // [KY XX XX]
		esp += 1                     //
		v.addOp(vm.EQ)               // [EQ XX]
		esp -= 1                     //
		v.addOp(vm.ISZERO)           // [!E XX]
		esp += 0                     //
		v.addPointer(labels[i+1].id) // [NX !E XX]
		esp += 1                     //
		v.addOp(vm.JUMPI)            // [XX]
		esp -= 2                     //

		// If we reach this, values ARE equal and the clause
		// body will be executed.
		//
		// NB: We do not tweak the esp from this point onward,
		// otherwise we'd mess the jumps.
		body.Accept(v, s, esp) // [RR XX]

		// Jump to the `after` label.
		v.addPointer(labels[after].id) // [PP RR XX]
		v.addOp(vm.JUMP)               // [RR XX]
	}

	// Handle the `otherwise` clause manually.  Stack is [XX].
	if last > 0 {
		v.addSegment(labels[last])
	}
	otherwise := tail[last]
	body := otherwise.Children[1]
	if otherwise.NumChildren() > 2 {
		body = NewNodeProgn()
		body.AddChildren(otherwise.Children[1:])
	}
	body.Accept(v, s, esp) // [RR XX]

	// Now the `after`.  Stack is always [RR XX], where RR is the
	// result of the evaluated body and XX is the original switch
	// value.
	esp += 1                    // Only 1 body was executed.
	v.addSegment(labels[after]) //
	v.addOp(vm.SWAP1)           // [XX RR]
	esp += 0                    //
	v.addOp(vm.POP)             // [RR]
	esp -= 1                    //
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

var _storagePosition int32 = -1

func fnDefvar(v *BytecodeVisitor, s *Scope, _ int, call Node) {
	args := assertNargsEq("defvar", call, 2)

	if !s.IsGlobal() {
		panic("defvar can be used only globally")
	}

	if !args[0].IsSymbol() {
		panic("TODO")
	}
	identifier := args[0].ValueString

	s.SetStorageVariable(identifier, atomic.AddInt32(&_storagePosition, 1))

	v.VisitNil()
}

func fnEther(v *BytecodeVisitor, _ *Scope, esp int, call Node) {
	args := assertNargsEq("ether", call, 1)

	if !args[0].IsString() {
		panic("TODO")
	}

	inp := args[0].ValueString
	sep := strings.Index(inp, ".")
	if sep < 0 {
		sep = len(inp)
	}
	val := inp + "000000000000000000"
	rep := strings.Replace(val, ".", "", 1)
	cut := rep[:sep+18]
	ans, err := uint256.FromDecimal(cut)
	if err != nil {
		panic(err)
	}

	v.pushU256(ans)
	esp += 1
}

func fnGethash(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsGte("gethash", call, 2) // (gethash table keys...)

	if !args[0].IsSymbol() {
		panic("TODO")
	}
	table := args[0].ValueString
	pos, ok := s.GetStorageVariable(table)
	if !ok {
		panic("TODO, void variable")
	}

	v.pushU64(uint64(pos)) // [PP]
	esp += 1               //  --> esp=1

	for i := 1; i < len(args); i++ {
		key := args[i]

		v.pushU64(20)         // [20 PP]
		esp += 1              //  --> esp=2
		v.addOp(vm.MSTORE)    // [], m[20]=PP
		esp -= 2              //  --> esp=0
		key.Accept(v, s, esp) // [KK]
		esp += 1              //  --> esp=1
		v.pushU64(0)          // [00 KK]
		esp += 1              //  --> esp=2
		v.addOp(vm.MSTORE)    // [], m[00]=KK
		esp -= 2              //  --> esp=0
		v.pushU64(0x40)       // [40]
		esp += 1              //  --> esp=1
		v.pushU64(0x00)       // [00 40]
		esp += 1              //  --> esp=2
		v.addOp(vm.KECCAK256) // [HH]
		esp -= 1              //  --> esp=1
	}

	v.addOp(vm.SLOAD) // [VV]
	esp += 0          //  --> esp=1
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

func fnPuthash(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsGte("puthash", call, 3) // (puthash table value keys...)

	if !args[0].IsSymbol() {
		panic("TODO")
	}
	table := args[0].ValueString

	value := args[1]

	pos, ok := s.GetStorageVariable(table)
	if !ok {
		panic("TODO, void variable")
	}

	value.Accept(v, s, esp) // [VV]
	esp += 1                //  --> esp=1
	v.addOp(vm.DUP1)        // [VV VV]
	esp += 1                //  --> esp=2
	v.pushU64(uint64(pos))  // [PP VV VV]
	esp += 1                //  --> esp=3

	for i := 2; i < len(args); i++ {
		key := args[i]

		v.pushU64(20)         // [20 PP VV VV]
		esp += 1              //  --> esp=4
		v.addOp(vm.MSTORE)    // [VV VV], m[20]=PP
		esp -= 2              //  --> esp=2
		key.Accept(v, s, esp) // [KK VV VV]
		esp += 1              //  --> esp=3
		v.pushU64(0)          // [00 KK VV VV]
		esp += 1              //  --> esp=4
		v.addOp(vm.MSTORE)    // [VV VV], m[00]=KK
		esp -= 2              //  --> esp=2
		v.pushU64(0x40)       // [40 VV VV]
		esp += 1              //  --> esp=3
		v.pushU64(0x00)       // [00 40 VV VV]
		esp += 1              //  --> esp=4
		v.addOp(vm.KECCAK256) // [HH VV VV]
		esp -= 1              //  --> esp=3
	}

	v.addOp(vm.SSTORE)	// [VV], s[HH]=VV
	esp -= 2		//  --> esp=1
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

func fnSetq(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsEq("setq", call, 2)

	if !args[0].IsSymbol() {
		panic("TODO")
	}
	identifier := args[0].ValueString

	pos, ok := s.GetStorageVariable(identifier)
	if !ok {
		panic("TODO, void variable")
	}

	// Evaluate the expression and push to stack.
	args[1].Accept(v, s, esp) // [X]
	esp += 1

	v.addOp(vm.DUP1) // [X X]
	esp += 1

	// Push position.
	v.pushU64(uint64(pos)) // [P X X]
	esp += 1

	v.addOp(vm.SSTORE) // [X]
	esp -= 2
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
