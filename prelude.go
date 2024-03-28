package mist

import "fmt"

func assertArgsEq(fn string, args []Node, want int) {
	if have := len(args); have != want {
		panic(fmt.Sprintf(
			"wront number of arguments for (%s): have %d, want %d",
			fn,
			have,
			want,
		))
	}
}

func assertArgsGte(fn string, args []Node, want int) {
	if have := len(args); have < want {
		panic(fmt.Sprintf(
			"wront number of arguments for (%s): have %d, want at least %d",
			fn,
			have,
			want,
		))
	}
}

func FnWhen(v *BytecodeVisitor, args []Node) {
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

func FnRevert(v *BytecodeVisitor, args []Node) {
	assertArgsEq("revert", args, 0)

	zero := NewNodeUint256(0)
	zero.Accept(v)
	zero.Accept(v)
	v.pushOp(REVERT)
}

// Alpha is what I call those opcodes that have δ=0 and α>=0,
// i.e. STOP, ADDRESS, ORIGIN, CALLER, etc.
func IsAlpha(tok string) (OpCode, bool) {
	switch tok {
	// α == 0
	case "stop":
		return STOP, true

	// α == 1
	case "address":
		return ADDRESS, true
	case "origin":
		return ORIGIN, true
	case "caller":
		return CALLER, true
	case "call-value":
		return CALLVALUE, true
	case "call-data-load":
		return CALLDATALOAD, true
	case "call-data-size":
		return CALLDATASIZE, true
	case "code-size":
		return CODESIZE, true
	case "gas-price":
		return GASPRICE, true
	case "return-data-size":
		return RETURNDATASIZE, true
	case "coinbase":
		return COINBASE, true
	case "timestamp":
		return TIMESTAMP, true
	case "block-number":
		return NUMBER, true
	case "prev-randao":
		return PREVRANDAO, true
	case "gas-limit":
		return GASLIMIT, true
	case "chain-id":
		return CHAINID, true
	case "self-balance":
		return SELFBALANCE, true
	case "base-fee":
		return BASEFEE, true
	}

	return 0, false
}

func IsVariadic(tok string) bool {
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

type Environment struct {
	// constants map[string]string
}
