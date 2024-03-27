package mist

func Compile(program string) string {
	tokens := Tokenize(program)
	progn := Parse(&tokens)

	visitor := NewBytecodeVisitor()
	progn.Accept(&visitor)

	code := visitor.String()
	return code
}
