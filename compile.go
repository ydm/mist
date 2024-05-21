package mist

func Compile(program, source string, init bool, offopt uint32) (string, error) {
	tokens, err := Scan(program, source)
	if err != nil {
		return "", err
	}

	progn := Parse(&tokens)
	ast := OptimizeAST(progn, offopt)

	visitor := NewBytecodeVisitor(init)
	global := NewGlobalScope()
	ast.Accept(visitor, global, 0)

	visitor.OptimizeBytecode()
	code := visitor.String()

	return code, nil
}
