package mist_test

import (
	"testing"

	"github.com/ydm/mist"
)

func TestDecompile(t *testing.T) {
	t.Parallel()

	// ctor := "0x6080604052606b8060116000396000f3fe"
	code := "0x608060405260206040516000358152f3fe602060405160003560e01c8152f3fe600436101561003457602060405160208152f3fe5b602060405160308152f3fe"
	mist.Decompile(code)
}
