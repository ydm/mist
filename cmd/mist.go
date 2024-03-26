package main

import (
	"fmt"

	"github.com/ydm/mist"
)

func main() {
	// mist.Tokenize("(+ 1 2 3)")
	tokens := mist.Tokenize("(+1 2 3) 123 123 222 '(1 2 3) ''''(1 2 3) '1 ''''2")
	// for tokens.HasNext() {
	// 	fmt.Println(tokens.Next())
	// }
	ast := mist.Parse(&tokens)
	fmt.Println(&ast)
}
