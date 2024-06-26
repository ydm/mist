package mist_test

import (
	"testing"

	"github.com/ydm/mist"
)

func expectToken(
	t *testing.T,
	tok mist.Token,
	tokType int,
	tokValueString string,
	// tokValueNumber uint64,
) {
	t.Helper()

	if tok.Type != tokType {
		t.Errorf("have %d, want %d", tok.Type, tokType)
	}

	if tok.ValueString != tokValueString {
		t.Errorf("have %s, want %s", tok.ValueString, tokValueString)
	}

	// if tok.Type == mist.TokenNumber && !uint256.NewInt(tokValueNumber).Eq(tok.ValueNumber) {
	// 	t.Errorf("have %v, want %v", tok.ValueNumber, tokValueNumber)
	// }
}

func TestScan(t *testing.T) {
	t.Parallel()

	// (selector "pause")
	tokens, err := mist.Scan(`(selector "pause")`, "test")
	if err != nil {
		t.Fatal(err)
	}
	expectToken(t, tokens.Next(), mist.TokenLeftParen, "")
	expectToken(t, tokens.Next(), mist.TokenSymbol, "selector")
	expectToken(t, tokens.Next(), mist.TokenString, "pause")
	expectToken(t, tokens.Next(), mist.TokenRightParen, "")
	for tokens.HasNext() {
		t.Fail()
	}

	// (selector "pause()")
	tokens, err = mist.Scan(`(selector "pause()")`, "test")
	if err != nil {
		t.Fatal(err)
	}
	expectToken(t, tokens.Next(), mist.TokenLeftParen, "")
	expectToken(t, tokens.Next(), mist.TokenSymbol, "selector")
	expectToken(t, tokens.Next(), mist.TokenString, "pause()")
	expectToken(t, tokens.Next(), mist.TokenRightParen, "")
	for tokens.HasNext() {
		t.Fail()
	}

	// (selector "pause();")
	tokens, err = mist.Scan(`(selector "pause();")`, "test")
	if err != nil {
		t.Fatal(err)
	}
	expectToken(t, tokens.Next(), mist.TokenLeftParen, "")
	expectToken(t, tokens.Next(), mist.TokenSymbol, "selector")
	expectToken(t, tokens.Next(), mist.TokenString, "pause();")
	expectToken(t, tokens.Next(), mist.TokenRightParen, "")
	for tokens.HasNext() {
		t.Fail()
	}

	// (selector "")
	tokens, err = mist.Scan(`(selector "")`, "test")
	if err != nil {
		t.Fatal(err)
	}
	expectToken(t, tokens.Next(), mist.TokenLeftParen, "")
	expectToken(t, tokens.Next(), mist.TokenSymbol, "selector")
	expectToken(t, tokens.Next(), mist.TokenString, "")
	expectToken(t, tokens.Next(), mist.TokenRightParen, "")
	for tokens.HasNext() {
		t.Fail()
	}
}
