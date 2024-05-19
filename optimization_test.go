package mist_test

import (
	"fmt"
	"testing"

	"github.com/ydm/mist"
)

func TestOptimizeIf(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(if 1 2 3)",
	}

	want := []string{
		"6002",
	}

	have, err := mist.Compile(c, fmt.Sprintf("case%d", i), false, offopt)
	if err != nil {
		t.Fatal(err)
	}

	if have != want {
		t.Logf("want:\n%s", mist.Decompile(want[i]))
		t.Logf("have:\n%s", mist.Decompile(have))
	}
}
