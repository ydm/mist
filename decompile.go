package mist

import (
	"fmt"
	"strconv"
	"strings"
)

var nargs = map[OpCode]int{
	PUSH1: 1,
	PUSH2: 2,
	PUSH3: 3,
	PUSH4: 4,
	PUSH5: 5,
	PUSH6: 6,
	PUSH7: 7,
	PUSH8: 8,
	PUSH9: 9,
	PUSH10: 10,
	PUSH11: 11,
	PUSH12: 12,
	PUSH13: 13,
	PUSH14: 14,
	PUSH15: 15,
	PUSH16: 16,
	PUSH17: 17,
	PUSH18: 18,
	PUSH19: 19,
	PUSH20: 20,
	PUSH21: 21,
	PUSH22: 22,
	PUSH23: 23,
	PUSH24: 24,
	PUSH25: 25,
	PUSH26: 26,
	PUSH27: 27,
	PUSH28: 28,
	PUSH29: 29,
	PUSH30: 30,
	PUSH31: 31,
	PUSH32: 32,
}

func Decompile(program string) {
	var (
		line strings.Builder
		ops  = parseOps(program)
	)

	for i := 0; i < len(ops); i++ {
		op := ops[i]
		n, ok := nargs[op]

		if ok {
			line.Reset()
			line.WriteString(op.String())

			for j := 1; j <= n && (i+j) < len(ops); j++ {
				line.WriteString(fmt.Sprintf(" %x", byte(ops[i+j])))
			}
			i += n

			fmt.Println(line.String())
		} else {
			fmt.Println(op)
		}
	}
}

func parseOps(program string) []OpCode {
	if strings.HasPrefix(program, "0x") {
		program = program[2:]
	}

	words := make([]OpCode, 0, 1024)
	for i := 0; i < len(program)/2; i += 2 {
		excerpt := program[i : i+2]
		word, err := strconv.ParseUint(excerpt, 16, 32)
		if err != nil {
			panic(err)
		}
		if 255 < word {
			panic("TODO")
		}
		words = append(words, OpCode(word))
	}

	return words
}
