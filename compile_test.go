package mist_test

import (
	"testing"

	"github.com/ydm/mist"
)

func TestCompileVariadic(t *testing.T) {
	t.Parallel()

	want := []string{
		"6002600101",
		"6003600201600101",
		"6002600303",
		"6001600203600303",
		"6001600260030103",
		"6002600404",
		"6006600260040460016003010203",
		"6007600317600117",
		"6007600317600117",
	}

	cases := []string{
		"(+ 1 2)",
		"(+ 1 2 3)",
		"(- 3 2)",
		"(- 3 2 1)",
		"(- (+ 3 2) 1)",
		"(/ 4 2)",
		"(- (* (+ 3 1) (/ 4 2)) 6)",
		"(logand 1 3 7)",
		"(& 1 3 7)",
		"(logior 1 3 7)",
		"(| 1 3 7)",
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

		"60016004600260020204031961001157005b00",
		"6001600460026002020403196100185760036002016001015b00",

		"600436101961000e5760006000fd5b",
	}

	cases := []string{
		"(when 1)",
		"(when 1) (stop)",

		"(when 1 (stop))",
		"(when 1 (stop)) (stop)",

		"(when (- (/ (* 2 2) 4) 1) (stop)) (stop)",
		"(when (- (/ (* 2 2) 4) 1) (+ 1 2 3)) (stop)",

		"(when (< (calldata-size) 4) (revert))",
	}

	for i, c := range cases {
		have := mist.Compile(c)
		if have != want[i] {
			t.Errorf("have=%s want=%s", have, want[i])
		}
	}
}

func TestCompileComplex(t *testing.T) {
	t.Parallel()

	want := []string{
		// "366004101961000e5760006000fd5b"
		"60003560e01c63a7a0d53714",
		"602060405160458152f3fe",
		"60003560e01c63a7a0d537141961001c57602060405160458152f3fe5b00",
	}
	cases := []string{
		// "(when (< (calldata-size) 4) (revert))",
		"(= (>> (calldata-load 0) 0xe0) 0xa7a0d537)",
		"(return 69)",
		"(when (= (>> (calldata-load 0) 0xe0) 0xa7a0d537) (return 69)) (stop)",
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
