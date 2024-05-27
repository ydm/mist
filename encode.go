package mist

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
)

func EncodeRLP(x any) string {
	b := bytes.NewBuffer([]byte{})
	if err := rlp.Encode(b, x); err != nil {
		panic(err)
	}

	var s strings.Builder
	for {
		x, err := b.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&s, "%02x", x)
	}

	return s.String()
}

func EncodeString(s string) string {
	if len(s) > 32 {
		panic("TODO")
	}

	var b strings.Builder
	for i := range len(s) {
		fmt.Fprintf(&b, "%02x", s[i])
	}
	return padRight32(b.String())
}
