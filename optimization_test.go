package mist_test

import (
	"testing"

	"github.com/ydm/mist"
)

func TestOptimizeIfs(t *testing.T) {
	t.Parallel()

	const (
		program = "(if 1 2 3)"
		want = "6002"
	)

	have, err := mist.Compile(program, "test", false, 0)
	if err != nil {
		t.Fatal(err)
	}
	if have != want {
		t.Errorf("have %s, want %s", have, want)
	}
}
