package mist_test

import (
	"fmt"
	"testing"

	"github.com/ydm/mist"
)

func TestDecompile(t *testing.T) {
	t.Parallel()
	fmt.Println(mist.Decompile("565bfe5b60459056"))
}
