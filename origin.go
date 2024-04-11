package mist

import "fmt"

type Origin struct {
	Filename string
	Line     int
	Column   int
}

func NewOrigin(filename string, line, column int) Origin {
	return Origin{filename, line, column}
}

func (o Origin) String() string {
	return fmt.Sprintf("%s:%d:%d", o.Filename, o.Line, o.Column)
}

