package mist

import (
	"fmt"
	"strings"
	"unicode"
)

func Tokenize(code string) {
	var (
		token strings.Builder
		tokens = make([]string, 0, 128)
	)

	push := func (x string) {
		if token.Len() > 0 {
			tokens = append(tokens, token.String())
			token.Reset()
		}

		if len(x) > 0 {
			tokens = append(tokens, x)
		}
	}

	for _, r := range code {
		switch r {
		case '(':
			push("(")
		case ')':
			push(")")
		default:
			if unicode.IsSpace(r) {
				push("")
			} else {
				token.WriteRune(r)
			}
		}
	}

	push("")
	fmt.Println(strings.Join(tokens, ", "))
	fmt.Println(len(tokens))
}
