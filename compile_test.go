package mist_test

import (
	"fmt"
	"testing"

	"github.com/ydm/mist"

	"github.com/google/go-cmp/cmp"
)

func compileAndCompare(t *testing.T, cases, want []string) {
	t.Helper()

	const offopt = mist.OffoptIf

	for i, c := range cases {
		have, err := mist.Compile(c, fmt.Sprintf("case%d", i), false, offopt)
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

func TestCompileAnd(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(and)",

		"(and 1)",
		"(and t)",
		
		"(and 0)",
		"(and nil)",

		"(and 0x20 0x30)",
		"(and 0x20 nil)",
	}

	want := []string{
		"6001",

		"6001",
		"6001",
		
		"6000",
		"6000",

		"6020801561000b575060305b", // "6001506020801561000e575060305b",
		"6020801561000b575060005b",
	}

	compileAndCompare(t, cases, want)	
}

func TestCompileCase(t *testing.T) {
	t.Parallel()

	cases := []string{
		"(case 1)",
		"(case 1 (1 0x10))",
		"(case 1 (1 0x10) (otherwise 0))",
		"(case 1 (1 0x10) (otherwise nil))",
		"(case 1 (1 0x10) (t nil))",
		"(case 2 (1 0x10) (2 0x20))",
		"(case 0x1234 (otherwise 0x10))",
		"(case 0x1234 (1 0x10) (otherwise 0x10))",
	}

	want := []string{
		"6000",
		"60018060011415610011576010610014565b60005b9050",
		"60018060011415610011576010610014565b60005b9050",
		"60018060011415610011576010610014565b60005b9050",
		"60018060011415610011576010610014565b60005b9050",
		"60028060011415610011576010610024565b8060021415610021576020610024565b60005b9050",
		"6010",
		"6112348060011415610012576010610015565b60105b9050",
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

		"(defun f (x) x) (f 2)",
		"(defun f (x) (+ x x)) (f 2)",
		"(defun f (x y) (- x y)) (f 0x20 0x10)",
		"(defun f (x y z) (/ (- x y) z)) (f 0x100 0x10 0x2)",

		"(defun f () 69) (+ (f) (f))",

		"(defun f () 1) (defun g () (+ (f) 2)) (g)",
		"(defun f () 1) (defun g () (+ (f) 2 (f))) (g)",

		// // "(defun f (x y) (- x y y)) (f 0x30 0x10)",
	}

	want := []string{
		"6000",
		"6100085b604590565b",

		"61000b60025b80905090565b",
		"61000d60025b808101905090565b",
		"610010601060205b81810391505090565b",
		"610016600260106101005b82828203049250505090565b",

		"6100085b604590565b610010610003565b01",

		"6100125b600261000e5b600190565b0190565b",
		"61001b5b61000c5b600190565b600201610017610007565b0190565b",
		
		// "60106020818103915050",
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

func TestCompileSelector(t *testing.T) {
	t.Parallel()

	cases := []string{
		`(selector "pause()")`,
		`(selector "something()")`,
	}

	want := []string{
		"638456cb59",
		"63a7a0d537",
	}

	compileAndCompare(t, cases, want)
}

func TestCompileString(t *testing.T) {
	t.Parallel()

	cases := []string {
		`"qwe"`,
		`"0123456789012345678901234567890"`,
	}

	want := []string{
		"6383717765",
		"7f9f30313233343536373839303132333435363738393031323334353637383930",
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
