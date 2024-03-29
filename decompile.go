package mist

import (
	"fmt"
	"strconv"
	"strings"
)

func nargs(op OpCode) (int, bool) {
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

func Decompile(program string) {
	var (
		line strings.Builder
		ops  = parseOps(program)
	)

	for i := 0; i < len(ops); i++ {
		op := ops[i]
		n, ok := nargs(op)

		if ok {
			line.Reset()
			line.WriteString(fmt.Sprintf("| %02x | %s", i, op.String()))

			for j := 1; j <= n && (i+j) < len(ops); j++ {
				line.WriteString(fmt.Sprintf(" %02x", byte(ops[i+j])))
			}
			i += n

			fmt.Println(line.String())
		} else {
			fmt.Printf("| %02x | %s\n", i, op.String())
		}
	}
}

func parseOps(program string) []OpCode {
	program = strings.TrimPrefix(program, "0x")

	words := make([]OpCode, 0, 1024)
	for i := 0; i < len(program); i += 2 {
		excerpt := program[i : i+2]

		word, err := strconv.ParseUint(excerpt, 16, 32)
		if err != nil {
			panic(err)
		}
		if word > 255 {
			panic("TODO")
		}
		words = append(words, OpCode(word))
	}

	return words
}
