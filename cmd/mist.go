package main

import (
	"fmt"

	"github.com/ydm/mist"
)

func main() {
	tokens := mist.Tokenize("(+ 1 2 3)")
	progn := mist.Parse(&tokens)
	fmt.Println(&progn)

	v := mist.NewBytecodeVisitor()
	progn.Accept(&v)

	code := v.String()

	fmt.Println(code)
	mist.Decompile(code)
}
