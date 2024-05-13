package mist_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ydm/mist"
)

func TestEncodeWithSignature(t *testing.T) {
	t.Parallel()

	const want = ("08c379a0" +
		"0000000000000000000000000000000000000000000000000000000000000020" +
		"0000000000000000000000000000000000000000000000000000000000000003" +
		"6173640000000000000000000000000000000000000000000000000000000000")

	have := mist.EncodeWithSignature("Error(string)", "asd")

	if diff := cmp.Diff(have, want); diff != "" {
		t.Error(diff)
	}
}

func TestNumArguments(t *testing.T) {
	t.Parallel()

	signatures := []string{
		"totalSupply()",
		"balanceOf(address)",
		"transfer(address,uint256)",
		"allowance(address,address)",
		"transferFrom(address,address,uint256)",
	}

	want := []int{0, 1, 2, 2, 3}

	for i := range signatures {
		have := mist.NumArguments(signatures[i])
		if have != want[i] {
			t.Errorf("%s: have %d, want %d", signatures[i], have, want[i])
		}
	}
}
