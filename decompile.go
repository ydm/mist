package mist

import (
	"fmt"
	"strconv"
	"strings"
)

func Decompile(program string) {
	var (
		line strings.Builder
		ops  = parseOps(program)
	)

	for i := 0; i < len(ops); i++ {
		op := ops[i]
		n, ok := Nargs(op)

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
