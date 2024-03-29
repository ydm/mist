package mist

import (
	"strconv"
	"strings"
)

func consume(iterator *TokenIterator, tokens ...string) {
	if !iterator.HasNext() {
		panic("TODO")
	}

	for _, token := range tokens {
		if iterator.Peek() == token {
			iterator.Next()
			return
		}
	}

	panic("TODO")
}

func consumeNot(iterator *TokenIterator, tokens ...string) string {
	if !iterator.HasNext() {
		panic("TODO")
	}

	for _, token := range tokens {
		if iterator.Peek() == token {
			panic("TODO")
		}
	}

	return iterator.Next()
}

func parseAtom(tokens *TokenIterator) Node {
	next := consumeNot(tokens, "(", ")")

	if strings.HasPrefix(next, "0x") {
		// TODO: Parse uint256
		parsed, err := strconv.ParseUint(next[2:], 16, 64)
		if err != nil {
			panic(err) // TODO
		}
		return NewNodeUint256(parsed)
	}

	parsed, err := strconv.ParseUint(next, 10, 64)
	if err == nil {
		return NewNodeUint256(parsed)
	}

	return NewNodeSymbol(next)
}

func parseList(tokens *TokenIterator) Node {
	consume(tokens, "(")

	root := NewNodeList()
	for tokens.HasNext() {
		if tokens.Peek() == "(" { //nolint:gocritic
			// That's a nested list, go deeper.
			root.AddChild(parseList(tokens))
		} else if tokens.Peek() == ")" {
			break
		} else {
			root.AddChild(parseAtom(tokens))
		}
	}

	consume(tokens, ")")

	return root
}

func parse(tokens *TokenIterator) Node {
	for tokens.HasNext() {
		switch tokens.Peek() {
		case "'":
			tokens.Next()
			quote := NewNodeList()
			quote.AddChild(NewNodeSymbol("quote"))
			quote.AddChild(parse(tokens))
			return quote
		case "(":
			return parseList(tokens)
		default:
			return parseAtom(tokens)
		}
	}

	panic("unreachable")
}

func Parse(tokens *TokenIterator) Node {
	progn := NewNodeList()
	progn.AddChild(NewNodeSymbol("progn"))

	for tokens.HasNext() {
		progn.AddChild(parse(tokens))
	}

	return progn
}
