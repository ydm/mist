package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ydm/mist"
)

func main() {
	inp, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	// Strip comments.
	stripped := mist.StripComments(string(inp))

	// Tokenize.
	tokens := mist.Tokenize(stripped)

	// Parse.
	progn := mist.Parse(&tokens)
	// fmt.Println(&progn)

	// Compile.
	v := mist.NewBytecodeVisitor()
	progn.Accept(&v)

	// Decorate with a contract constructor.
	code := v.String()
	ctor := mist.MakeConstructor(code)

	// fmt.Println("bytecode:")
	fmt.Print("0x" + ctor + code)

	// fmt.Println("deployedBytecode:")
	// fmt.Println("0x" + code)

	// fmt.Println()
	// mist.Decompile(code)
}
