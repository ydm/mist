package mist

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
)

const (
	freeMemoryPointer = 0x40
	freeMemoryInitial = 0x80
)

// +---------+
// | segment |
// +---------+

var _segmentID int32 = 0 //nolint:gochecknoglobals

func makeSegmentID() int32 {
	return atomic.AddInt32(&_segmentID, 1)
}

type segment struct {
	id int32

	opcode  int    // valid if >=0
	data    string //
	pointer int32  // valid if >0
}

func (s segment) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[")

	if s.isData() {
		fmt.Fprintf(&b, "data %s", s.data)
	} else if s.isOpcode() {
		fmt.Fprintf(&b, "op %v", vm.OpCode(s.opcode))
	} else if s.isPointer() {
		fmt.Fprintf(&b, "ptr to %d", s.pointer)
	} else {
		panic("broken invariant")
	}

	fmt.Fprintf(&b, "]")
	return b.String()
}

func newSegmentData(data string) segment {
	return segment{makeSegmentID(), -1, data, 0}
}

func newEmptySegment() segment {
	return segment{makeSegmentID(), -1, "", 0}
}

func newSegmentJumpdest() segment {
	return newSegmentOpCode(vm.JUMPDEST)
}

func newSegmentOpCode(op vm.OpCode) segment {
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
	return s.isOpcode() && vm.OpCode(s.opcode).IsPush()
}

func (s *segment) isPop() bool {
	return s.isOpcode() && vm.OpCode(s.opcode) == vm.POP
}

func (s *segment) pointTo(pos int) {
	if !s.isPointer() {
		panic("not a pointer")
	}

	code := fmt.Sprintf("%02x%04x", byte(vm.PUSH2), pos)
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

	if s.isOpcode() {
		return fmt.Sprintf("%02x", byte(s.opcode))
	}

	if s.data == "" {
		panic("empty data segment")
	}

	return s.data
}

func (s *segment) len() int {
	if s.isPointer() {
		// The length of this segment is 3 bytes: (PUSH2 AA BB).
		return 3
	}

	if s.isOpcode() {
		// Each opcode is exactly 1 byte.
		return 1
	}

	// Each byte needs 2 hexadecimal characters.
	if len(s.data)%2 != 0 {
		panic("broken invariant")
	}
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
		v.addOp(vm.MSTORE)
	}

	return v
}

// +---------------+
// | Add functions |
// +---------------+

func (v *BytecodeVisitor) addSegment(s segment) {
	v.segments = append(v.segments, s)
}

func (v *BytecodeVisitor) addHex(code string) {
	v.addSegment(newSegmentData(code))
}

func (v *BytecodeVisitor) addOp(op vm.OpCode) {
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
	v.addHex(code)
}

func (v *BytecodeVisitor) addU64(x uint64) {
	v.addU256(uint256.NewInt(x))
}

// +----------------+
// | Push functions |
// +----------------+

func (v *BytecodeVisitor) pushU256(x *uint256.Int) {
	hex := x.Hex()

	length := len(hex)/2 - 1 + len(hex)%2
	if length < 1 || 32 < length {
		panic("TODO")
	}

	op := vm.OpCode(byte(vm.PUSH0) + byte(length))
	v.addOp(op)
	v.addU256(x)
}

func (v *BytecodeVisitor) pushU64(x uint64) {
	v.pushU256(uint256.NewInt(x))
}

// +-----------------+
// | Visit functions |
// +-----------------+

func (v *BytecodeVisitor) VisitNil() {
	v.pushU64(0)
}

func (v *BytecodeVisitor) VisitT() {
	v.pushU64(1)
}

func (v *BytecodeVisitor) VisitNumber(x *uint256.Int) {
	v.pushU256(x)
}

func (v *BytecodeVisitor) VisitString(n Node) {
	if !n.IsString() {
		panic("TODO")
	}

	encoded := EncodeRLP(n.ValueString)
	length := len(encoded) / 2
	if length > 32 {
		// Still not supporting strings bigger than a single word.
		panic("string literal is longer than 31 characters")
	}

	op := vm.OpCode(byte(vm.PUSH0) + byte(length))
	v.addOp(op)
	v.addHex(encoded)
}

func (v *BytecodeVisitor) VisitSymbol(s *Scope, esp int, symbol Node) {
	if variable, ok := s.GetStackVariable(symbol.ValueString); ok {
		delta := esp - variable.Position
		if delta <= 0 {
			panic("broken invariant")
		}
		opcode := vm.OpCode(vm.DUP1 + (delta - 1))
		// fmt.Printf(
		// 	"var=%s pos=%d esp=%d delta=%d opcode=%s\n",
		// 	symbol.ValueString,
		// 	variable.Position,
		// 	esp,
		// 	delta,
		// 	opcode,
		// )
		v.addOp(opcode)
		return
	}

	if pos, ok := s.GetStorageVariable(symbol.ValueString); ok {
		v.pushU64(uint64(pos))
		v.addOp(vm.SLOAD)
		return
	}

	if node, ok := s.GetConstant(symbol.ValueString); ok {
		node.Accept(v, s, esp)
		return
	}

	panic(fmt.Sprintf("%v: void variable %s", symbol.Origin, symbol.ValueString))
}

func (v *BytecodeVisitor) VisitFunction(s *Scope, esp int, call Node) {
	handlers := []func(*BytecodeVisitor, *Scope, int, Node) bool{
		// Custom functions have precedence over
		// native/builtin.
		handleDefinedFunc,

		handleNativeFunc,
		handleVariadicFunc,
		handleBuiltinFunc,
		handleMacroFunc,
	}

	for _, handler := range handlers {
		ok := handler(v, s, esp, call)
		if ok {
			return
		}
	}

	panic(fmt.Sprintf("%v: void function: %s", call.Origin, call.FunctionName()))
}

// +-------------------+
// | Segment functions |
// +-------------------+

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

func (v *BytecodeVisitor) OptimizeBytecode() {
	v.segments = OptimizeBytecode(v.segments)
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
