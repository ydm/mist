package mist

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/holiman/uint256"
)

// +--------+
// | Errors |
// +--------+

type LexicalError struct {
	Origin  Origin
	Message string
	Token   string
}

func NewLexicalError(filename string, line, column int, message, token string) error {
	return &LexicalError{
		NewOrigin(filename, line, column),
		message,
		token,
	}
}

func (e *LexicalError) Error() string {
	end := ""
	if e.Token != "" {
		end = ": " + e.Token
	}

	return fmt.Sprintf("%v: error: %s%s", e.Origin, e.Message, end)
}

// +-------+
// | Token |
// +-------+

const (
	TokenLeftParen  = iota // (
	TokenRightParen        // )
	TokenQuote             // '
	TokenNumber
	TokenString
	TokenSymbol
)

type Token struct {
	Type        int
	ValueString string       // Symbolic or string value.
	ValueNumber *uint256.Int // Numeric value.
	Origin      Origin
}

func (t Token) Short() string {
	switch t.Type {
	case TokenLeftParen:
		return "("
	case TokenRightParen:
		return ")"
	case TokenQuote:
		return "'"
	case TokenNumber:
		return fmt.Sprintf("number(%v)", t.ValueNumber)
	case TokenString:
		return fmt.Sprintf("string(\"%s\")", t.ValueString)
	case TokenSymbol:
		return fmt.Sprintf("symbol(%s)", t.ValueString)
	default:
		panic("TODO")
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%v %s", t.Origin, t.Short())
}

// +---------------+
// | TokenIterator |
// +---------------+

type TokenIterator struct {
	tokens []Token
	index  int
}

func NewTokenIterator() TokenIterator {
	return TokenIterator{
		tokens: make([]Token, 0, 512),
		index:  0,
	}
}

func (i *TokenIterator) HasNext() bool {
	return i.index < len(i.tokens)
}

func (i *TokenIterator) Next() Token {
	ans := i.tokens[i.index]
	i.index++
	return ans
}

func (i *TokenIterator) Peek() Token {
	return i.tokens[i.index]
}

func (i *TokenIterator) push(token Token) {
	i.tokens = append(i.tokens, token)
}

func Scan(code string, filename string) (TokenIterator, error) {
	var (
		// Multi-character tokens are built character by character.
		// Line and column of origin are also stored.
		builder       strings.Builder
		builderLine   int
		builderColumn int

		tokens = NewTokenIterator()

		push = func(tokenType int, s string, n *uint256.Int, line, col int) {
			tokens.push(Token{
				Type:        tokenType,
				ValueString: s,
				ValueNumber: n,
				Origin:      NewOrigin(filename, line, col),
			})
		}
	)

	// If there are characters already collected, we might be able to
	// build a token.
	maybeBuild := func() error {
		if builder.Len() <= 0 {
			return nil
		}

		built := strings.TrimSpace(builder.String())
		if len(built) <= 0 {
			panic("broken invariant")
		}

		e := func(msg string) error {
			return NewLexicalError(filename, builderLine, builderColumn, msg, built)
		}

		if strings.HasPrefix(built, "\"") {
			// Token starts with a double quote, treat it as string.
			panic("TODO") // TODO: Handle strings!
		} else if strings.HasPrefix(built, "0x") {
			// Token starts with a 0x prefix, treat it as number.
			if parsed, err := uint256.FromHex(built); err == nil {
				push(TokenNumber, "", parsed, builderLine, builderColumn)
			} else {
				return e("invalid hex literal")
			}
		} else if parsed, err := uint256.FromDecimal(built); err == nil {
			// If token can be parsed into a number, treat it as such.
			push(TokenNumber, "", parsed, builderLine, builderColumn)
		} else {
			// TODO: Check if that's a proper symbol, contains no
			// forbidden characters like quotes, etc.
			push(TokenSymbol, built, nil, builderLine, builderColumn)
		}

		builder.Reset()
		return nil
	}

	// Context variables.
	comment := false
	line := 1 // Lines start from 1.start
	offset := 0

	for index, r := range code {
		// Everything between a ';' and '\n' is considered a comment.
		// This would have to change once strings are supported.
		switch r {
		case ';':
			if err := maybeBuild(); err != nil {
				return tokens, err
			}
			comment = true
		case '\n':
			comment = false
			line++
			offset = index + 1
		}

		// If inside a comment, skip this character.
		if comment {
			continue
		}

		col := index - offset // Columns start from 0.

		// Character is part of the programming code.
		switch r {
		case '(':
			if err := maybeBuild(); err != nil {
				return tokens, err
			}
			push(TokenLeftParen, "", nil, line, col)
		case ')':
			if err := maybeBuild(); err != nil {
				return tokens, err
			}
			push(TokenRightParen, "", nil, line, col)
		case '\'':
			if err := maybeBuild(); err != nil {
				return tokens, err
			}
			push(TokenQuote, "", nil, line, col)
		default:
			if unicode.IsSpace(r) {
				if err := maybeBuild(); err != nil {
					return tokens, err
				}
			} else {
				// If this is the first character from a new token,
				// mark the line and column of origin.
				if builder.Len() <= 0 {
					builderLine = line
					builderColumn = col
				}
				builder.WriteRune(r)
			}
		}
	}

	err := maybeBuild()
	return tokens, err
}
