package mist

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/core/vm"
)

func nargs(op vm.OpCode) int {
	if op.IsPush() {
		return int(op - vm.PUSH0)
	}

	return 0
}

func Decompile(program string) string {
	var (
		out  strings.Builder
		ops  = parseOps(program)
	)

	for i := 0; i < len(ops); i++ {
		op := ops[i]

		n := 0
		if op.IsPush() {
			n = int(op - vm.PUSH0)
		}

		if n > 0 {
			var line strings.Builder
			fmt.Fprintf(&line, "| %02x | %s", i, op.String())

			for j := 1; j <= n && (i+j) < len(ops); j++ {
				fmt.Fprintf(&line, " %02x", byte(ops[i+j]))
			}
			i += n

			fmt.Fprintln(&out, line.String())
		} else {
			fmt.Fprintf(&out, "| %02x | %s\n", i, op.String())
		}
	}

	return out.String()
}

func parseOps(program string) []vm.OpCode {
	program = strings.TrimPrefix(program, "0x")

	words := make([]vm.OpCode, 0, 1024)
	for i := 0; i < len(program); i += 2 {
		excerpt := program[i : i+2]

		word, err := strconv.ParseUint(excerpt, 16, 32)
		if err != nil {
			panic(err)
		}
		if word > 255 {
			panic("TODO")
		}
		words = append(words, vm.OpCode(word))
	}

	return words
}
