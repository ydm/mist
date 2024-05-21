package mist

import "fmt"

// +---------+
// | Lexical |
// +---------+

type LexicalError struct {
	Origin  Origin
	Message string
	Token   string
}

func NewLexicalError(filename string, line, column int, message, token string) error {
	return LexicalError{
		NewOrigin(filename, line, column),
		message,
		token,
	}
}

func (e LexicalError) Error() string {
	end := ""
	if e.Token != "" {
		end = ": " + e.Token
	}

	return fmt.Sprintf("%v: lexical error: %s%s", e.Origin, e.Message, end)
}

// +-------------+
// | Compilation |
// +-------------+

type CompilationError struct {
	Origin  Origin
	Message string
}

func NewCompilationError(origin Origin, message string) error {
	return CompilationError{
		origin,
		message,
	}
}

func (e CompilationError) Error() string {
	return fmt.Sprintf("%v: compilation error: %s", e.Origin, e.Message)
}
