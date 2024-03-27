package mist

func Nargs(op OpCode) (int, bool) {
	switch op { //nolint:exhaustive
	case STOP:
		return 0, true
	case PUSH1:
		return 1, true
	case PUSH2:
		return 2, true
	case PUSH3:
		return 3, true
	case PUSH4:
		return 4, true
	case PUSH5:
		return 5, true
	case PUSH6:
		return 6, true
	case PUSH7:
		return 7, true
	case PUSH8:
		return 8, true
	case PUSH9:
		return 9, true
	case PUSH10:
		return 10, true
	case PUSH11:
		return 11, true
	case PUSH12:
		return 12, true
	case PUSH13:
		return 13, true
	case PUSH14:
		return 14, true
	case PUSH15:
		return 15, true
	case PUSH16:
		return 16, true
	case PUSH17:
		return 17, true
	case PUSH18:
		return 18, true
	case PUSH19:
		return 19, true
	case PUSH20:
		return 20, true
	case PUSH21:
		return 21, true
	case PUSH22:
		return 22, true
	case PUSH23:
		return 23, true
	case PUSH24:
		return 24, true
	case PUSH25:
		return 25, true
	case PUSH26:
		return 26, true
	case PUSH27:
		return 27, true
	case PUSH28:
		return 28, true
	case PUSH29:
		return 29, true
	case PUSH30:
		return 30, true
	case PUSH31:
		return 31, true
	case PUSH32:
		return 32, true
	}

	return 0, false
}
