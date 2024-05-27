package main

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/ydm/mist"
)

func main() {
	// source := argv[1]
	// stream, err := os.Open(source)
	// if err != nil {
	// 	panic(err)
	// }
	// inp, err := io.ReadAll(stream)

	source := "stdin"
	inp, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	decoded := string(inp)
	if !utf8.ValidString(decoded) {
		panic("TODO")
	}

	const (
		// TODO: Turn into cli args.
		init    = true
		verbose = true
		decompile = false
	)

	// Decorate with a contract constructor.
	code, err := mist.Compile(decoded, source, init, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if verbose {
		ctor := mist.MakeConstructor(code)

		fmt.Println("combined:")
		fmt.Println("0x" + ctor + code)
		fmt.Println()

		fmt.Println("constructor:")
		fmt.Println("0x" + ctor)
		fmt.Println()

		fmt.Println("deployedBytecode:")
		fmt.Println("0x" + code)
		fmt.Println()

		if decompile {
			fmt.Print(mist.Decompile(code))
		}
	} else {
		ctor := mist.MakeConstructor(code)
		fmt.Print("0x" + ctor + code)
	}
}
