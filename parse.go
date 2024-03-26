package mist

import "strconv"

func parseAtom(tokens *TokenIterator) Node {
	if !tokens.HasNext() {
		panic("TODO")
	}

	if next := tokens.Peek(); next == "(" || next == ")" {
		panic("TODO")
	}

	next := tokens.Next() // Consume

	parsed, err := strconv.ParseUint(next, 10, 64)
	if err == nil {
		return NewNodeUint256(parsed)
	}

	return NewNodeSymbol(next)
}

func parseList(tokens *TokenIterator) Node {
	if tokens.Peek() != "(" {
		panic("TODO")
	}

	// Consume (.
	tokens.Next()

	// Prepare stack.
	root := NewNodeList()

	for tokens.HasNext() {
		if tokens.Peek() == "(" {
			root.AddChild(parseList(tokens))
		} else if tokens.Peek() == ")" {
			// Consume ) and stop looping as this is the end of the
			// current list.
			tokens.Next()
			break
		} else {
			root.AddChild(parseAtom(tokens))
		}
	}

	return root
}

func parse(tokens *TokenIterator) Node {
	for tokens.HasNext() {
		if tokens.Peek() == "'" {
			tokens.Next()
			quote := NewNodeList()
			quote.AddChild(NewNodeSymbol("quote"))
			quote.AddChild(parse(tokens))
			return quote
		} else if tokens.Peek() == "(" {
			return parseList(tokens)
		} else {
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
