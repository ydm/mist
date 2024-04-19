package mist_test

import (
	"fmt"
	"testing"

	"github.com/ydm/mist"

	"github.com/google/go-cmp/cmp"
)

func compileAndCompare(t *testing.T, cases, want []string) {
	t.Helper()

	for i, c := range cases {
		have, err := mist.Compile(c, fmt.Sprintf("case%d", i), false)
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

func TestCompileDefconst(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(defconst +x+ 123)",
		"(defconst +x+ 123) +x+",
	}

	want := []string{
		// (defconst) results in nil and, if this is the last
		// expression to compile, it's left on the stack.
		"6000",
		"607b",
	}

	compileAndCompare(t, cases, want)
}

func TestCompileDefun(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(defun f () 69)",
		"(defun f () 69) (f)",

		"(defun f (x) (+ x 1)) (f 2)",
		"(defun f (x) (+ x x)) (f 2)",

		"(defun f (x y) (- x y)) (f 0x20 0x10)",
		// "(defun f (x y) (- x y y)) (f 0x30 0x10)",
	}

	want := []string{
		"6000",
		"6045",

		"6002600181019050",
		"60028081019050",

		"60106020818103915050",
	}

	compileAndCompare(t, cases, want)
}

func TestCompileVariadic(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(+ 1 2)",
		"(+ 1 2 3)",
		"(- 3 2)",
		// "(- 3 2 1)",
		"(- (+ 3 2) 1)",
		"(/ 4 2)",
		"(- (* (+ 3 1) (/ 4 2)) 6)",
		"(logand 1 3 7)",
		"(& 1 3 7)",
		"(logior 1 3 7)",
		"(| 1 3 7)",
	}

	want := []string{
		"6002600101",
		"6003600201600101",
		"6002600303",
		// "6001600203600303",
		"6001600260030103",
		"6002600404",
		"6006600260040460016003010203",
		"6007600316600116",
		"6007600316600116",
		"6007600317600117",
		"6007600317600117",
	}

	compileAndCompare(t, cases, want)
}

func TestCompileIf(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(if 1 2 3)",
		"(if 1 2 3) (stop)",

		"(when 1)",
		"(when 1) (stop)",

		"(when (- (/ (* 2 2) 4) 1) (stop))",
		"(when (- (/ (* 2 2) 4) 1) (+ 1 2 3)) (stop)",
	}

	want := []string{
		"600161000c57600361000f565b60025b",
		"600161000c57600361000f565b60025b5000",

		"600161000c57600061000f565b60005b",
		"600161000c57600061000f565b60005b5000",

		"6001600460026002020403610015576000610017565b005b",
		"600160046002600202040361001557600061001e565b60036002016001015b5000",
	}

	compileAndCompare(t, cases, want)
}

func TestCompileComplex(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(= (>> (calldata-load 0) 0xe0) 0xa7a0d537)",
		"(return 69)",
		"(calldata-load 0)",
		"(>> (calldata-load 0) 0xe0)",
	}

	want := []string{
		"63a7a0d53760003560e01c14",
		"602060405160458152f3",
		"600035",
		"60003560e01c",
	}

	compileAndCompare(t, cases, want)
}

func TestCompileNative(t *testing.T) {
	t.Parallel()
}
