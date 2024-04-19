package mist

// OptimizePushPop deletes sequences of the following form:
//
// PUSH[1-16]
// DATA
// POP
func OptimizePushPop(segments []segment) []segment {

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
		start := i
		end := start + 3
		optimized = append(optimized[:start], optimized[end:]...)
	}

	return optimized
}
