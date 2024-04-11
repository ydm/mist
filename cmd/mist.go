package main

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/ydm/mist"
)

func main() {
	// source := "examples/something.mist"
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

	// Decorate with a contract constructor.
	code, err := mist.Compile(decoded, source, true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctor := mist.MakeConstructor(code)
	fmt.Print("0x" + ctor + code)

	// fmt.Println("bytecode:")
	// fmt.Println("0x" + ctor + code)

	// fmt.Println("deployedBytecode:")
	// fmt.Println("0x" + code)

	// mist.Decompile(code)
}
