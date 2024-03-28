package mist_test

import (
	"testing"

	"github.com/ydm/mist"
)

func TestCompileVariadic(t *testing.T) {
	t.Parallel()

	want := []string{
		"6001600201",
		"6001600201600301",

		"6003600203",
		"6003600203600103",

		"6003600201600103",

		"6003600101600460020402600603",
		"6003600101600460020402600603600214",

		"6001600316600716",
	}

	cases := []string{
		"(+ 1 2)",
		"(+ 1 2 3)",

		"(- 3 2)",
		"(- 3 2 1)",

		"(- (+ 3 2) 1)",

		"(- (* (+ 3 1) (/ 4 2)) 6)",
		"(= (- (* (+ 3 1) (/ 4 2)) 6) 2)",

		"(logand 1 3 7)",
	}

	for i, c := range cases {
		have := mist.Compile(c)
		if have != want[i] {
			t.Errorf("have=%s want=%s", have, want[i])
		}
	}
}

func TestCompileWhen(t *testing.T) {
	t.Parallel()

	want := []string{
		"600119610007575b",
		"600119610007575b00",

		"60011961000857005b",
		"60011961000857005b00",

		"60026002026004046001031961001157005b00",
		"6002600202600404600103196100185760016002016003015b00",

		"366004101961000e5760006000fd5b",
	}

	cases := []string{
		"(when 1)",
		"(when 1) (stop)",

		"(when 1 (stop))",
		"(when 1 (stop)) (stop)",

		"(when (- (/ (* 2 2) 4) 1) (stop)) (stop)",
		"(when (- (/ (* 2 2) 4) 1) (+ 1 2 3)) (stop)",

		"(when (< (call-data-size) 4) (revert))",
	}

	for i, c := range cases {
		have := mist.Compile(c)
		if have != want[i] {
			t.Errorf("have=%s want=%s", have, want[i])
		}
	}
}

func TestCompileNative(t *testing.T) {
	t.Parallel()

}
