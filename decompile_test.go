package mist_test

import (
	"fmt"
	"testing"

	"github.com/ydm/mist"
)

func TestDecompile(t *testing.T) {
	t.Parallel()
	fmt.Println(mist.Decompile("0x60178061000c6000396000f340526000506020604051602060108181019150508152f3"))
}
