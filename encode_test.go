package mist_test

import (
	"testing"

	"github.com/ydm/mist"

	"github.com/google/go-cmp/cmp"
)

func TestEncode(t *testing.T) {
	t.Parallel()

	have := mist.Encode("asd")
	want := "83617364"

	if diff := cmp.Diff(have, want); diff != "" {
		t.Error(diff)
	}
}
