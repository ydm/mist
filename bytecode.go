package mist

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/holiman/uint256"
)

const (
	freeMemoryPtr   = 0x40
	freeMemoryStart = 0x80
)

// +---------+
// | segment |
// +---------+

var _segmentID int32 = 0

func makeSegmentID() int32 {
	return atomic.AddInt32(&_segmentID, 1)
}

type segment struct {
	id int32

	code    string
	pointer int32
}

func newSegmentCode(code string) segment {
	return segment{makeSegmentID(), code, 0}
}

func newEmptySegment() segment {
	return segment{makeSegmentID(), "", 0}
}

func newSegmentJumpdest() segment {
	return newSegmentOpCode(JUMPDEST)
}

func newSegmentOpCode(op OpCode) segment {
	return newSegmentCode(fmt.Sprintf("%02x", byte(op)))
}

func newSegmentPointer(jumpdest int32) segment {
	return segment{makeSegmentID(), "", jumpdest}
}

func (s *segment) isPointer() bool {
	return s.pointer != 0
}

func (s *segment) pointTo(pos int) {
	if !s.isPointer() {
		panic("not a pointer")
	}

	code := fmt.Sprintf("%02x%04x", byte(PUSH2), pos)
	if len(code) != 6 {
		panic("broken invariant")
	}

	if s.code != "" && s.code != code {
		panic(fmt.Sprintf(
			"trying to reassign a different position: old=%s new=%s",
			s.code,
			code,
		))
	}
	s.code = code
}

func (s *segment) getCode() string {
	if s.isPointer() && s.code == "" {
		panic("pointer not initialized")
	}

	return s.code
}

func (s *segment) len() int {
	if s.isPointer() {
		// The length of this segment is 3 bytes: (PUSH2 AA BB).
		return 3
	} else {
		// Each byte needs 2 hexadecimal characters.
		return len(s.code) / 2
	}
}

// +-----------------+
// | BytecodeVisitor |
// +-----------------+

type BytecodeVisitor struct {
	segments []segment
}

func NewBytecodeVisitor() BytecodeVisitor {
	v := BytecodeVisitor{
		segments: make([]segment, 0, 1024),
	}

	// Initialize the free memory pointer.  Mist follows the same
	// memory layout as Solidity.
	// v.pushU64(freeMemory)
	// v.pushU64(freeMemoryPtr)
	// v.pushOp(MSTORE)

	return v
}

// codeLength return the number of bytes.
func (v *BytecodeVisitor) codeLength() int {
	ans := 0
	for i := range v.segments {
		ans += v.segments[i].len()
	}
	return ans
}

// +---------+
// | Add fns |
// +---------+

func (v *BytecodeVisitor) addSegment(s segment) *segment {
	index := len(v.segments)
	v.segments = append(v.segments, s)
	return &v.segments[index]
}

func (v *BytecodeVisitor) addCode(code string) {
	v.addSegment(newSegmentCode(code))
}

func (v *BytecodeVisitor) addJumpdest() int32 {
	return v.addSegment(newSegmentJumpdest()).id
}

func (v *BytecodeVisitor) addOp(op OpCode) {
	v.addSegment(newSegmentOpCode(op))
}

func (v *BytecodeVisitor) addPointer(dest int32) *segment {
	return v.addSegment(newSegmentPointer(dest))
}

func (v *BytecodeVisitor) addU256(x *uint256.Int) {
	hex := x.Hex()

	padding := ""
	if len(hex)%2 == 1 {
		padding = "0"
	}

	code := fmt.Sprintf("%s%s", padding, hex[2:])
	v.addCode(code)
}

func (v *BytecodeVisitor) addU64(x uint64) {
	v.addU256(uint256.NewInt(x))
}

// +----------+
// | Push fns |
// +----------+

func (v *BytecodeVisitor) pushU256(x *uint256.Int) {
	hex := x.Hex()

	length := len(hex)/2 - 1 + len(hex)%2
	if length < 1 || 32 < length {
		panic("TODO")
	}

	op := OpCode(byte(PUSH0) + byte(length))
	v.addOp(op)
	v.addU256(x)
}

func (v *BytecodeVisitor) pushU64(x uint64) {
	v.pushU256(uint256.NewInt(x))
}

func (v *BytecodeVisitor) VisitList() {
}

func (v *BytecodeVisitor) VisitFunction(fn string, args []Node) {
	if handleNativeFunc(v, fn, args) {
		// noop
	} else if handleVariadicFunc(v, fn, args) {
		// noop
	} else if handleInlineFunc(v, fn, args) {
		// noop
	} else {
		panic("unrecognized function: " + fn)
	}
}

func (v *BytecodeVisitor) VisitSymbol(_ string) {
}

func (v *BytecodeVisitor) VisitNumber(x *uint256.Int) {
	v.pushU256(x)
}

func (v *BytecodeVisitor) getPosition(id int32) int {
	pos := 0
	for i := range v.segments {
		if v.segments[i].id == id {
			return pos
		}
		pos += v.segments[i].len()
	}
	panic("broken invariant")
}

func (v *BytecodeVisitor) populatePointers() {
	for i := range v.segments {
		if v.segments[i].isPointer() {
			pos := v.getPosition(v.segments[i].pointer)
			v.segments[i].pointTo(pos)
		}
	}
}

func (v *BytecodeVisitor) String() string {
	v.populatePointers()

	var b strings.Builder
	for i := range v.segments {
		b.WriteString(v.segments[i].getCode())
	}

	return b.String()
}

// +---------+
// | Private |
// +---------+

func (v *BytecodeVisitor) visitVariadicOp(fn string, args []Node) {
	assertArgsGte(fn, args, 2)

	var op OpCode
	switch fn {
	case "+":
		op = ADD
	case "*":
		op = MUL
	case "-":
		op = SUB
	case "/":
		op = DIV
	case "<":
		op = LT
	case ">":
		op = GT
	case "=":
		op = EQ
	case "&":
		fallthrough
	case "logand":
		op = AND
	case "|":
		fallthrough
	case "logior":
		op = OR
	case "^":
		fallthrough
	case "logxor":
		op = XOR
	default:
		panic("unrecognized arithmetic op: " + fn)
	}

	last := len(args) - 1
	args[last].Accept(v)
	for i := last - 1; i >= 0; i-- {
		args[i].Accept(v)
		v.addOp(op)
	}
}
