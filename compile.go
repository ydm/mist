package mist

// import "fmt"

func Compile(program, source string) (string, error) {
	tokens, err := Scan(program, source)
	if err != nil {
		return "", err
	}

	progn := Parse(&tokens)

	visitor := NewBytecodeVisitor(true)
	progn.Accept(&visitor)

	return visitor.String(), nil
}
