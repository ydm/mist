package mist_test

import (
	"testing"

	"github.com/ydm/mist"
)

func TestDecompile(t *testing.T) {
	t.Parallel()
	mist.Decompile("600161000a5761000d565b60005b")
}
