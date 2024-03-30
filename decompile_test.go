package mist_test

import (
	"testing"

	"github.com/ydm/mist"
)

func TestDecompile(t *testing.T) {
	t.Parallel()

	// ctor := "0x6080604052606b8060116000396000f3fe"
	code := "600115610007574a554d5044455354"
	mist.Decompile(code)
}
