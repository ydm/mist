package mist

type Visitor interface {
	VisitList(tokens []string)
	VisitAtom()
}

type BytecodeVisitor struct{}

func (v BytecodeVisitor) VisitFunction(tokens []string) {
	
}
