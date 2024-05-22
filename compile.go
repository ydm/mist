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

	// fmt.Println("BEFORE OPTIMIZATION:", visitor.getSegments())
	segments := visitor.GetOptimizedSegments()
	// fmt.Println("AFTER OPTIMIZATION:", segments)
	segments = SegmentsPopulatePointers(segments)
	// fmt.Println("AFTER POPULATION:", segments)
	code := SegmentsToString(segments)
	// fmt.Println("CODE:", code)

	return code, nil
}
