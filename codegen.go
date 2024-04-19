package mist

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/holiman/uint256"
)

const (
	freeMemoryPointer = 0x40
	freeMemoryInitial = 0x80
)

// +---------+
// | segment |
// +---------+

var _segmentID int32 //nolint:gochecknoglobals

func makeSegmentID() int32 {
	return atomic.AddInt32(&_segmentID, 1)
}

type segment struct {
	id int32

	opcode  int					// valid if >=0
	data    string				// 
	pointer int32				// valid if >0
}

func newSegmentData(data string) segment {
	return segment{makeSegmentID(), -1, data, 0}
}

func newEmptySegment() segment {
	return segment{makeSegmentID(), -1, "", 0}
}

func newSegmentJumpdest() segment {
	return newSegmentOpCode(JUMPDEST)
}

func newSegmentOpCode(op OpCode) segment {
	return segment{makeSegmentID(), int(op), "", 0}
}

func newSegmentPointer(jumpdest int32) segment {
	return segment{makeSegmentID(), -1, "", jumpdest}
}

func (s *segment) isData() bool {
	return !s.isOpcode() && !s.isPointer() && len(s.data) >= 2
}

func (s *segment) isOpcode() bool {
	return s.opcode >= 0
}

func (s *segment) isPointer() bool {
	return s.pointer != 0
}

func (s *segment) isPush() bool {
	return s.isOpcode() && OpCode(s.opcode).IsPush()
}

func (s *segment) isPop() bool {
	return s.isOpcode() && OpCode(s.opcode) == POP
}

func (s *segment) pointTo(pos int) {
	if !s.isPointer() {
		panic("not a pointer")
	}

	code := fmt.Sprintf("%02x%04x", byte(PUSH2), pos)
	if len(code) != 6 {
		panic("broken invariant")
	}

	if s.data != "" && s.data != code {
		panic(fmt.Sprintf(
			"trying to reassign a different position: old=%s new=%s",
			s.data,
			code,
		))
	}
	s.data = code
}

func (s *segment) getCode() string {
	if s.isPointer() && s.data == "" {
		panic("pointer not initialized")
	}

	if s.opcode >= 0 {
		return fmt.Sprintf("%02x", byte(s.opcode))
	}

	return s.data
}

func (s *segment) len() int {
	if s.isPointer() {
		// The length of this segment is 3 bytes: (PUSH2 AA BB).
		return 3
	}

	if s.opcode >= 0 {
		// Each opcode is exactly 1 byte.
		return 1
	}

	// Each byte needs 2 hexadecimal characters.
	return len(s.data) / 2
}

// +-----------------+
// | BytecodeVisitor |
// +-----------------+

type BytecodeVisitor struct {
	segments []segment
}

func NewBytecodeVisitor(init bool) *BytecodeVisitor {
	v := &BytecodeVisitor{
		segments: make([]segment, 0, 1024),
	}

	if init {
		// Initialize the free memory pointer.  Mist follows the same
		// memory layout as Solidity.
		v.pushU64(freeMemoryInitial)
		v.pushU64(freeMemoryPointer)
		v.addOp(MSTORE)
	}

	return v
}

// +---------+
// | Add fns |
// +---------+

func (v *BytecodeVisitor) addSegment(s segment) {
	v.segments = append(v.segments, s)
}

func (v *BytecodeVisitor) addCode(code string) {
	v.addSegment(newSegmentData(code))
}

func (v *BytecodeVisitor) addOp(op OpCode) {
	v.addSegment(newSegmentOpCode(op))
}

func (v *BytecodeVisitor) addPointer(dest int32) {
	v.addSegment(newSegmentPointer(dest))
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

//nolint:unused
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

func (v *BytecodeVisitor) VisitNil()  {
	v.pushU64(0)
}

func (v *BytecodeVisitor) VisitT()  {
	v.pushU64(1)
}

func (v *BytecodeVisitor) VisitNumber(x *uint256.Int)  {
	v.pushU256(x)
}

func (v *BytecodeVisitor) VisitSymbol(s *Scope, esp int, symbol Node) {
	node, ok := s.GetConstant(symbol.ValueString)
	if ok {
		node.Accept(v, s, esp)
		return
	}

	variable, ok := s.GetStackVariable(symbol.ValueString)
	if ok {
		delta := esp - variable.Position
		v.addOp(OpCode(DUP1 + delta - 1))
		return
	}

	panic(fmt.Sprintf("%v: void variable %s", symbol.Origin, symbol.ValueString))
}

func (v *BytecodeVisitor) VisitFunction(s *Scope, esp int, call Node) {
	handlers := []func(*BytecodeVisitor, *Scope, int, Node) bool{
		handleNativeFunc,
		handleVariadicFunc,
		handleInlineFunc,
		handleDefun,
	}

	for _, handler := range handlers {
		ok := handler(v, s, esp, call)
		if ok {
			return	
		}
	}

	panic(fmt.Sprintf("%v: void function: %s", call.Origin, call.FunctionName()))
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

func (v *BytecodeVisitor) Optimize() {
	v.segments = OptimizePushPop(v.segments)
	v.populatePointers()
}

func (v *BytecodeVisitor) String() string {
	v.populatePointers()

	var b strings.Builder
	for i := range v.segments {
		b.WriteString(v.segments[i].getCode())
	}

	return b.String()
}
