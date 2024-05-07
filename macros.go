package mist

// Mist still doesn't have macros, but once it does, all of the
// functions in this file should be rewritten.

func handleMacro(v *BytecodeVisitor, s *Scope, esp int, call Node) bool {
	fn := call.FunctionName()
	switch fn {
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
