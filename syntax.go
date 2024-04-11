package mist

import "fmt"

func consume(tokens *TokenIterator, types ...int) {
	if !tokens.HasNext() {
		panic("TODO")
	}

	next := tokens.Peek()

	for _, tokenType := range types {
		if next.Type == tokenType {
			tokens.Next()
			return
		}
	}

	panic("TODO")
}

func consumeExcept(tokens *TokenIterator, types ...int) Token {
	if !tokens.HasNext() {
		panic("TODO")
	}

	next := tokens.Peek()

	for _, tokenType := range types {
		if next.Type == tokenType {
			panic(fmt.Sprintf("%v: unexpected token %s",next.Origin, next.Short()))
		}
	}

	return tokens.Next()
}

func parseAtom(tokens *TokenIterator) Node {
	next := consumeExcept(tokens, TokenLeftParen, TokenRightParen, TokenQuote)
	switch next.Type {
	case TokenLeftParen:
		fallthrough
	case TokenRightParen:
		fallthrough
	case TokenQuote:
		panic("TODO")
	case TokenNumber:
		return NewNodeU256(next.ValueNumber)
	case TokenString:
		panic("TODO")
	case TokenSymbol:
		return NewNodeSymbol(next.ValueString)
	default:
		panic("TODO")
	}
}

func parseList(tokens *TokenIterator) Node {
	root := NewNodeList()
	consume(tokens, TokenLeftParen)
	
	for tokens.HasNext() {
		next := tokens.Peek()
		if next.Type == TokenLeftParen {
			// That's a nested list, go deeper.
			root.AddChild(parseList(tokens))
		} else if next.Type == TokenRightParen {
			break
		} else {
			root.AddChild(parse(tokens))
		}
	}

	consume(tokens, TokenRightParen)
	return root
}

func parse(tokens *TokenIterator) Node {
	for tokens.HasNext() {
		next := tokens.Peek()
		switch next.Type {
		case TokenLeftParen:
			return parseList(tokens)
		case TokenRightParen:
			panic("unbalanced parentheses")
		case TokenQuote:
			tokens.Next()
			quote := NewNodeList()
			quote.AddChild(NewNodeSymbol("quote"))
			quote.AddChild(parse(tokens))
			return quote
		case TokenNumber:
			fallthrough
		case TokenString:
			fallthrough
		case TokenSymbol:
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
