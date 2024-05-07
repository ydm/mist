package mist_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ydm/mist"
)

func TestAbiEncode(t *testing.T) {
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
