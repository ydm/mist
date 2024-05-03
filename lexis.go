package mist

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/holiman/uint256"
)

// A very simple state machine.

const (
	lexerStateCode = iota
	lexerStateComment
	lexerStateString
)

type lexerState struct {
	state  int
	lines  int
	offset int
}

func newLexerState() lexerState {
	return lexerState{lexerStateCode, 1, 0}
}

func (s *lexerState) transitionTo(state int) bool {
	// Allowed transitions:
	//
	// code -> comment
	// code -> string
	//
	// comment -> code
	//
	// string -> code
	switch s.state {
	case lexerStateCode:
		if state == lexerStateComment || state == lexerStateString {
			goto allowed
		}
	case lexerStateComment:
		fallthrough
	case lexerStateString:
		if state == lexerStateCode {
			goto allowed
		}
	}
	return false

allowed:
	s.state = state
	return true
}

func (s *lexerState) newLine(index int) {
	s.lines++
	s.offset = index + 1
}

func (s *lexerState) inCode() bool {
	return s.state == lexerStateCode
}

func (s *lexerState) inComment() bool {
	return s.state == lexerStateComment
}

func (s *lexerState) inString() bool {
	return s.state == lexerStateString
}

func (s *lexerState) getLine() int            { return s.lines }
func (s *lexerState) getColumn(index int) int { return index - s.offset }

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
		state  = newLexerState()
		// prev   = rune(0) // TODO

		pushRune = func(i int, r rune) {
			// If this is the first character from a new token,
			// mark the line and column of origin.
			if builder.Len() <= 0 {
				builderLine = state.getLine()
				builderColumn = state.getColumn(i)
			}
			builder.WriteRune(r)
		}

		pushToken = func(tokenType int, s string, n *uint256.Int, line, col int) {
			tokens.push(Token{
				Type:        tokenType,
				ValueString: s,
				ValueNumber: n,
				Origin:      NewOrigin(filename, line, col),
			})
		}

		// If there are characters already collected, we might be able to
		// build a token.
		maybeBuild = func() error {
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

			if strings.HasPrefix(built, "\"") && strings.HasSuffix(built, "\"") {
				end := len(built) - 1
				stripped := built[1:end]
				pushToken(TokenString, stripped, nil, builderLine, builderColumn)
			} else if strings.HasPrefix(built, "0x") {
				// Token starts with a 0x prefix, treat it as number.
				if parsed, err := uint256.FromHex(built); err == nil {
					pushToken(TokenNumber, "", parsed, builderLine, builderColumn)
				} else {
					return e("invalid hex literal")
				}
			} else if parsed, err := uint256.FromDecimal(built); err == nil {
				// If token can be parsed into a number, treat it as such.
				pushToken(TokenNumber, "", parsed, builderLine, builderColumn)
			} else {
				// TODO: Check if that's a proper symbol, contains no
				// forbidden characters like quotes, etc.
				pushToken(TokenSymbol, built, nil, builderLine, builderColumn)
			}

			builder.Reset()
			return nil
		}
	)

	for i, r := range code {
		if r == '"' && state.transitionTo(lexerStateString) {
			// Beginning of a string.
			pushRune(i, r)
			continue
		} else if r == ';' && state.transitionTo(lexerStateComment) {
			// Beginning of a comment.
			continue
		}

		if state.inComment() {
			// Inside a comment, ignore character.
		} else if state.inString() {
			pushRune(i, r)
			if r == '"' {
				if err := maybeBuild(); err != nil {
					// Error while completing the string.
					return tokens, err
				}
				state.transitionTo(lexerStateCode)
			}
		} else if state.inCode() {
			tokenType := -1
			switch r {
			case '(':
				tokenType = TokenLeftParen
			case ')':
				tokenType = TokenRightParen
			case '\'':
				tokenType = TokenQuote
			}
			if tokenType != -1 {
				if err := maybeBuild(); err != nil {
					return tokens, err
				}
				pushToken(tokenType, "", nil, state.getLine(), state.getColumn(i))
				continue
			}

			if unicode.IsSpace(r) {
				if err := maybeBuild(); err != nil {
					return tokens, err
				}
				continue
			}

			pushRune(i, r)
		} else {
			panic("broken invariant")
		}

		// if r == '"' && state.inString() {

		// 	continue
		// } else if r == '\n' {
		// 	// We doesn't support multi-line comments or strings, so
		// 	// get back to code.
		// 	state.transitionTo(lexerStateCode)
		// 	state.newLine(i)
		// 	continue
		// }

		// if true {
		// 	continue
		// }

		// // If inside a comment, skip this character.
		// if state.inComment() {
		// 	continue
		// } else if state.inString() {
		// 	pushRune(i, r)
		// }

		// // Character is part of the programming code.
		// fmt.Printf("before switch, %c\n", r)
		// switch r {
		// case '(':
		// 	if err := maybeBuild(); err != nil {
		// 		return tokens, err
		// 	}
		// 	pushToken(TokenLeftParen, "", nil, state.getLine(), state.getColumn(i))
		// case ')':
		// 	if err := maybeBuild(); err != nil {
		// 		return tokens, err
		// 	}
		// 	pushToken(TokenRightParen, "", nil, state.getLine(), state.getColumn(i))
		// case '\'':
		// 	if err := maybeBuild(); err != nil {
		// 		return tokens, err
		// 	}
		// 	pushToken(TokenQuote, "", nil, state.getLine(), state.getColumn(i))
		// default:
		// 	if unicode.IsSpace(r) {
		// 		if err := maybeBuild(); err != nil {
		// 			return tokens, err
		// 		}
		// 	} else {
		// 		pushRune(i, r)
		// 		fmt.Printf("[Y 3] %v >%s<\n", tokens.tokens, builder.String())
		// 	}
		// }
		// fmt.Printf("after switch, %c\n", r)

		// prev = r
	}

	err := maybeBuild()
	return tokens, err
}
