package mist

import (
	"fmt"
	"sync/atomic"
)

// Mist still doesn't support macros, but once it does, all of the
// functions in this file should be rewritten.  Manipulating the AST
// using Go is ugly.

// TODO: Once I implement macros and (cond) is implemented, (case)
// should be rewritten as a macro that translates to (cond).

func handleMacroFunc(v *BytecodeVisitor, s *Scope, esp int, call Node) bool {
	fn := call.FunctionName()
	switch fn {
	case "apply":
		fnApply(v, s, esp, call)
		return true
	case "let":
		fnLet(v, s, esp, call)
		return true
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

// Translate (apply 'fn args) to (fn args...).
func fnApply(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsGte("apply", call, 1)

	fst := args[0]
	snd := args[1]

	if !fst.IsList() ||
		fst.NumChildren() != 2 ||
		!fst.Children[0].IsQuote() ||
		!fst.Children[1].IsSymbol() {
		//
		panic(fmt.Sprintf(
			"%v: invalid function: want symbol, have %v",
			call.Origin,
			&fst,
		))
	}

	ans := NewNodeList(call.Origin)
	ans.AddChild(fst.Children[1])

	if !snd.IsNil() {
		if !snd.IsList() ||
			snd.NumChildren() < 1 ||
			!snd.Children[0].IsQuote() {
			//
			panic(fmt.Sprintf(
				"%v: wrong argument type: want quoted list, have %v",
				call.Origin,
				&snd,
			))
		}

		if snd.NumChildren() == 2 {
			if !snd.Children[1].IsList() {
				panic(fmt.Sprintf(
					"%v: wrong argument type: want list, have %v",
					call.Origin,
					&snd.Children[1],
				))
			}
			ans.AddChildren(snd.Children[1].Children)
		}
	}

	ans.Accept(v, s, esp)
}

var _lambdaCounter uint32 = 0 //nolint:gochecknoglobals

func makeUniqueLambdaName() string {
	return fmt.Sprintf("lambda%d", atomic.AddUint32(&_lambdaCounter, 1))
}

// Translate (let varlist body...), where varlist is
// ((key1 value1)
//
//	(key2 value2))
//
// to
//
// (progn (defun unique (keys...) body...)
//
//	(apply 'unique values))
func fnLet(v *BytecodeVisitor, s *Scope, esp int, call Node) {
	args := assertNargsGte("let", call, 0)

	// Split the varlist to two separate lists: keys and values.
	varlist := NewNodeNil(NewOriginEmpty())
	if len(args) > 0 {
		varlist = args[0]
	}

	keys := NewNodeList(NewOriginEmpty())
	values := NewNodeList(NewOriginEmpty())
	for i := 0; i < varlist.NumChildren(); i++ {
		pair := varlist.Children[0]
		if !pair.IsList() || pair.NumChildren() != 2 {
			panic("TODO")
		}

		key := pair.Children[0]
		if !key.IsSymbol() {
			panic("TODO")
		}
		keys.AddChild(key)

		value := pair.Children[1]
		values.AddChild(value)
	}

	// Wrap the body expressions into a single (progn).
	body := NewNodeNil(NewOriginEmpty())
	if len(args) > 1 {
		body = NewNodeProgn()
		body.AddChildren(args[1:])
	}

	unique := NewNodeSymbol(makeUniqueLambdaName(), NewOriginEmpty())

	defun := NewNodeList(NewOriginEmpty())
	defun.AddChild(NewNodeSymbol("defun", NewOriginEmpty()))
	defun.AddChild(unique)
	defun.AddChild(keys)
	defun.AddChildren(args[1:])

	apply := NewNodeList(NewOriginEmpty())
	apply.AddChild(NewNodeSymbol("apply", NewOriginEmpty()))
	apply.AddChild(NewNodeQuote(unique, NewOriginEmpty()))
	apply.AddChild(NewNodeQuote(values, NewOriginEmpty()))

	progn := NewNodeProgn()
	progn.AddChild(defun)
	progn.AddChild(apply)

	progn.Accept(v, s, esp)
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
