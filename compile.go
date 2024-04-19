package mist

func Compile(program, source string, init bool) (string, error) {
	tokens, err := Scan(program, source)
	if err != nil {
		return "", err
	}

	visitor := NewBytecodeVisitor(init)
	global := NewGlobalScope()

	progn := Parse(&tokens)
	progn.Accept(visitor, global, 0)

	visitor.Optimize()
	return visitor.String(), nil
}
