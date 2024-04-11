package main

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/ydm/mist"
)

func main() {
	// stream, err := os.Open("examples/something.mist")
	// if err != nil {
	// 	panic(err)
	// }
	// inp, err := io.ReadAll(stream)

	inp, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	decoded := string(inp)
	if !utf8.ValidString(decoded) {
		panic("TODO")
	}

	// Tokenize.
	tokens, err := mist.Scan(decoded, "stdin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// for tokens.HasNext() {
	// 	fmt.Println(tokens.Next())
	// }
	// return

	// Parse.
	progn := mist.Parse(&tokens)
	fmt.Println(&progn)

	return

	// Compile.
	v := mist.NewBytecodeVisitor()
	progn.Accept(&v)

	// Decorate with a contract constructor.
	code := v.String()
	ctor := mist.MakeConstructor(code)
	fmt.Print("0x" + ctor + code)

	// fmt.Println("bytecode:")
	// fmt.Println("0x" + ctor + code)

	// fmt.Println("deployedBytecode:")
	// fmt.Println("0x" + code)

	// mist.Decompile(code)
}
