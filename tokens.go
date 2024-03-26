package mist

import (
	"strings"
	"unicode"
)

// +--------------+
// | Tokenization |
// +--------------+

type TokenIterator struct {
	tokens []string
	index  int
}

func NewTokenIterator() TokenIterator {
	return TokenIterator{
		tokens: make([]string, 0, 512),
		index:  0,
	}
}

func (i *TokenIterator) HasNext() bool {
	return i.index < len(i.tokens)
}

func (i *TokenIterator) Next() string {
	ans := i.tokens[i.index]
	i.index++
	return ans
}

func (i *TokenIterator) Peek() string {
	return i.tokens[i.index]
}

func (i *TokenIterator) push(token string) {
	i.tokens = append(i.tokens, token)
}

func Tokenize(code string) TokenIterator {
	var (
		builder strings.Builder
		tokens  = NewTokenIterator()
	)

	push := func(token string) {
		if builder.Len() > 0 {
			tokens.push(builder.String())
			builder.Reset()
		}

		if len(token) > 0 {
			tokens.push(token)
		}
	}

	for _, r := range code {
		switch r {
		case '\'':
			push("'")
		case '(':
			push("(")
		case ')':
			push(")")
		default:
			if unicode.IsSpace(r) {
				push("")
			} else {
				builder.WriteRune(r)
			}
		}
	}

	push("")
	return tokens
}
