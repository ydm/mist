package mist_test

import (
	"fmt"
	"testing"

	"github.com/ydm/mist"
)

func TestDecompile(t *testing.T) {
	t.Parallel()
	fmt.Println(mist.Decompile(""))
}
