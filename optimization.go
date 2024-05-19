package mist

// +-------------------+
// | AST optimizations |
// +-------------------+

// Flags to turn off optimizations.
const (
	OFFOPT_ARITHMETIC = 1 << iota
	OFFOPT_IF = 1 << iota
)

// TODO: If an arithmetic expression is made up of constants, replace
// it with the result instead.
func optimizeArithmetic(node Node) Node {
	return node
}

// If the condition of an (if) expression is constant, replace the
// whole (if) expression with the corresponding branch.
func optimizeIf(node Node) Node {
	if node.Type != NodeList || node.NumChildren() < 1 {
		return node
	}

	if !node.Children[0].IsThisSymbol("if") {
		ans := NewNodeList(node.Origin)
		for i := range node.Children {
			ans.AddChild(optimizeIf(node.Children[i]))
		}
		return ans
	}

	// (if cond yes no)
	//   0    1   2  3
	if node.Children[1].IsT() {
		return node.Children[2]
	} else if node.Children[1].IsNil() {
		return node.Children[3]
	}
	return node
}

func OptimizeAST(node Node, offs uint32) Node {
	type t func(node Node) Node
	fns := []t{
		optimizeArithmetic,
		optimizeIf,
	}
	for i, fn := range fns {
		bit := uint32(1) << i
		if offs&bit == 0 {
			node = fn(node)
		}
	}
	return node
}

// +------------------------+
// | Bytecode optimizations |
// +------------------------+

// optimizePushPop deletes sequences of the following form:
//
// PUSH[1-16]
// DATA
// POP
func optimizePushPop(segments []segment) []segment {
	marked := make([]int, 0, 16)
	n := len(segments) - 2
	for i := 0; i < n; i++ {
		if segments[i].isPush() &&
			segments[i+1].isData() &&
			segments[i+2].isPop() {
			//
			marked = append(marked, i)
			i += 2
		}
	}

	optimized := make([]segment, len(segments))
	copy(optimized, segments)

	for i := len(marked) - 1; i >= 0; i-- {
		start := marked[i]
		end := start + 3
		optimized = append(optimized[:start], optimized[end:]...)
	}

	return optimized
}

func OptimizeBytecode(segments []segment) []segment {
	type t func(segments []segment) []segment
	fns := []t{
		optimizePushPop,
	}
	for _, fn := range fns {
		segments = fn(segments)
	}
	return segments
}
