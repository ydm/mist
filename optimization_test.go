package mist_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ydm/mist"
)

func TestOptimizeIf(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(if 1 2 3)",
		"(if t 2 3)",
		"(if 0 2 3)",
		"(if nil 2 3)",
	}

	want := []string{
		"6002",
		"6002",
		"6003",
		"6003",
	}

	for i, c := range cases {
		have, err := mist.Compile(c, fmt.Sprintf("case%d", i), false, 0)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(want[i], have); diff != "" {
			t.Logf("Case #%d: %s", i, c)

			t.Logf("want:\n%s", mist.Decompile(want[i]))
			t.Logf("have:\n%s", mist.Decompile(have))

			t.Fatalf(diff)
		}
	}
}
