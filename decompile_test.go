package mist_test

import (
	"fmt"
	"testing"

	"github.com/ydm/mist"
)

func TestDecompile(t *testing.T) {
	t.Parallel()
	fmt.Println(mist.Decompile("0x60806040526020604051610011610016565b8152f3005b60459056"))
}
