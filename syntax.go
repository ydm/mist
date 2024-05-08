package mist

import "fmt"

func consume(tokens *TokenIterator, types ...int) Token {
	if !tokens.HasNext() {
		panic("incomplete code")
	}

	next := tokens.Peek()

	for _, tokenType := range types {
		if next.Type == tokenType {
			return tokens.Next()
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
			panic(fmt.Sprintf("%v: unexpected token %s", next.Origin, next.Short()))
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
		return NewNodeU256(next.ValueNumber, next.Origin)
	case TokenString:
		return NewNodeString(next.ValueString, next.Origin)
	case TokenSymbol:
		return NewNodeSymbol(next.ValueString, next.Origin)
	default:
		panic("TODO")
	}
}

func parseList(tokens *TokenIterator) Node {
	left := consume(tokens, TokenLeftParen)
	defer consume(tokens, TokenRightParen)

	root := NewNodeList(left.Origin)

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
			tokens.Next() // Consume the quote token.
			child := parse(tokens)
			quote := NewNodeQuote(child, next.Origin)
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
	progn := NewNodeProgn()

	for tokens.HasNext() {
		progn.AddChild(parse(tokens))
	}

	return progn
}
