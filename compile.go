package mist

// import "fmt"

func Compile(program, source string, init bool) (string, error) {
	tokens, err := Scan(program, source)
	if err != nil {
		return "", err
	}

	visitor := NewBytecodeVisitor(init)
	global := NewGlobalScope()

	progn := Parse(&tokens)
	progn.Accept(visitor, global)

	return visitor.String(), nil
}
