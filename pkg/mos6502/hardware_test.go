package mos6502

import (
	"testing"
)

// ========== Addressing Mode Edge Cases ==========

// TestIndirectJMPBug tests the famous JMP indirect page boundary bug.
// When the indirect address is at a page boundary (e.g., $xxFF),
// the 6502 wraps within the page instead of crossing to the next page.
// Example: JMP ($10FF) reads low byte from $10FF and high byte from $1000 (not $1100)
func TestIndirectJMPBug(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// Set up memory for the bug:
	// At $10FF: 0x34 (low byte)
	// At $1000: 0x12 (high byte - wrong page!)
	// At $1100: 0x56 (what we would expect without the bug)
	bus.memory[0x10FF] = 0x34
	bus.memory[0x1000] = 0x12 // Bug causes this to be read
	bus.memory[0x1100] = 0x56 // This should NOT be read

	// JMP ($10FF) - opcode 0x6C
	program := []byte{
		0x6C, 0xFF, 0x10, // JMP ($10FF)
	}
	bus.LoadProgram(0x8000, program)

	cpu.Reset()

	// Consume reset cycles
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute JMP ($10FF)
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Due to the bug, PC should be $1234 (from $10FF=0x34, $1000=0x12)
	// NOT $5634 (which would be from $10FF=0x34, $1100=0x56)
	expectedPC := uint16(0x1234)
	if cpu.PC != expectedPC {
		t.Errorf("JMP indirect bug: expected PC to be 0x%04X (bug behavior), got 0x%04X", expectedPC, cpu.PC)
	}
}

// TestZeroPageWraparound tests that zero page addressing wraps within page 0.
func TestZeroPageWraparound(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// Put a value at $0010 (where $FF + $11 wraps to)
	bus.memory[0x0010] = 0x42

	// LDA $11,X with X=$FF should read from $0010, not $0110
	program := []byte{
		0xA2, 0xFF, // LDX #$FF (set X register to cause wraparound)
		0xB5, 0x11, // LDA $11,X
	}
	bus.LoadProgram(0x8000, program)

	cpu.Reset()

	// Consume reset cycles
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute LDX #$FF
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute LDA $11,X
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	if cpu.A != 0x42 {
		t.Errorf("Zero page wraparound: expected A to be 0x42, got 0x%02X", cpu.A)
	}
}

// TestPageBoundaryCrossingPenalty tests that crossing page boundaries adds cycles.
func TestPageBoundaryCrossingPenalty(t *testing.T) {
	tests := []struct {
		name           string
		setupX         byte
		baseAddr       uint16
		expectedCycles int // Total steps needed after fetch
	}{
		{"No page cross", 0x10, 0x2000, 4}, // LDA $2000,X with X=$10 = $2010 (same page), 4 base cycles
		{"Page cross", 0x20, 0x20F0, 5},    // LDA $20F0,X with X=$20 = $2110 (crosses from $20 to $21), 4+1 cycles
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			bus.SetResetVector(0x8000)

			// Put value in target memory
			targetAddr := tt.baseAddr + uint16(tt.setupX)
			bus.memory[targetAddr] = 0x99

			// Program: LDX #immediate, then LDA absolute,X
			program := []byte{
				0xA2, tt.setupX, // LDX #immediate
				0xBD, byte(tt.baseAddr & 0xFF), byte(tt.baseAddr >> 8), // LDA $2000,X
			}
			bus.LoadProgram(0x8000, program)

			cpu.Reset()

			// Consume reset cycles
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			// Execute LDX #immediate
			cpu.Step()
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			// Fetch LDA instruction
			cpu.Step()

			// Count steps needed to complete
			steps := 1 // Already did one step to fetch
			for cpu.Cycles > 0 {
				cpu.Step()
				steps++
			}

			if steps != tt.expectedCycles {
				t.Errorf("Expected %d total steps, got %d", tt.expectedCycles, steps)
			}

			if cpu.A != 0x99 {
				t.Errorf("Expected A to be 0x99, got 0x%02X", cpu.A)
			}
		})
	}
}

// ========== Flag Behavior Tests ==========

// TestADCOverflowFlag tests the overflow flag calculation for ADC.
func TestADCOverflowFlag(t *testing.T) {
	tests := []struct {
		name     string
		a        byte
		operand  byte
		carry    bool
		overflow bool
	}{
		{"Positive + Positive = Negative", 0x50, 0x50, false, true},  // 80 + 80 = 160 (overflow)
		{"Positive + Positive = Positive", 0x50, 0x10, false, false}, // 80 + 16 = 96 (no overflow)
		{"Negative + Negative = Positive", 0xD0, 0x90, false, true},  // -48 + -112 = -160 (overflow)
		{"Negative + Negative = Negative", 0xD0, 0x10, false, false}, // -48 + 16 = -32 (no overflow)
		{"Positive + Negative", 0x50, 0xD0, false, false},            // No overflow across signs
		{"Zero edge case", 0x00, 0x00, false, false},
		{"With carry set", 0x7F, 0x00, true, true}, // 127 + 0 + 1 = 128 (overflow)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			bus.SetResetVector(0x8000)

			// Set up A register and carry flag
			cpu.A = tt.a
			if tt.carry {
				cpu.SetFlag(0x01, true)
			}

			// ADC immediate
			program := []byte{
				0x69, tt.operand, // ADC #operand
			}
			bus.LoadProgram(0x8000, program)

			cpu.Reset()
			// Restore A (reset clears it)
			cpu.A = tt.a
			if tt.carry {
				cpu.SetFlag(0x01, true)
			}

			// Consume reset cycles
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			// Execute ADC
			cpu.Step()
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			// Check overflow flag (bit 6, 0x40)
			hasOverflow := cpu.GetFlag(0x40)
			if hasOverflow != tt.overflow {
				t.Errorf("Overflow flag: expected %v, got %v (A=0x%02X, operand=0x%02X, carry=%v, result=0x%02X)",
					tt.overflow, hasOverflow, tt.a, tt.operand, tt.carry, cpu.A)
			}
		})
	}
}

// TestSBCOverflowFlag tests the overflow flag calculation for SBC.
func TestSBCOverflowFlag(t *testing.T) {
	tests := []struct {
		name     string
		a        byte
		operand  byte
		carry    bool // In SBC, carry=1 means no borrow
		overflow bool
	}{
		{"Positive - Negative = Negative overflow", 0x50, 0xB0, true, true}, // 80 - (-80) = 160 (overflow)
		{"Negative - Positive = Positive overflow", 0xD0, 0x70, true, true}, // -48 - 112 = -160 (overflow)
		{"No overflow same signs", 0x50, 0x20, true, false},
		{"Zero", 0x00, 0x00, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			bus.SetResetVector(0x8000)

			cpu.A = tt.a
			cpu.SetFlag(0x01, tt.carry)

			// SBC immediate
			program := []byte{
				0xE9, tt.operand, // SBC #operand
			}
			bus.LoadProgram(0x8000, program)

			cpu.Reset()
			cpu.A = tt.a
			cpu.SetFlag(0x01, tt.carry)

			for cpu.Cycles > 0 {
				cpu.Step()
			}

			cpu.Step()
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			hasOverflow := cpu.GetFlag(0x40)
			if hasOverflow != tt.overflow {
				t.Errorf("Overflow flag: expected %v, got %v", tt.overflow, hasOverflow)
			}
		})
	}
}

// TestUnusedFlagAlwaysSet tests that bit 5 (unused flag) is always set.
func TestUnusedFlagAlwaysSet(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// Try to clear the unused flag by pulling 0x00 from stack
	program := []byte{
		0x68, // PLA - pull 0x00 into A
		0x08, // PHP - push status
		0x68, // PLA - pull status into A
	}
	bus.LoadProgram(0x8000, program)

	// Put 0x00 on the stack
	cpu.SP = 0xFD
	bus.memory[0x01FE] = 0x00

	cpu.Reset()

	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute PLA (load 0x00)
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute PHP (push status)
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute PLA (pull status into A)
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Check that bit 5 is set in A (which now contains the status)
	if cpu.A&0x20 == 0 {
		t.Error("Unused flag (bit 5) should always be set in status register")
	}
}

// TestBITFlagCopying tests that BIT copies bits 6 and 7 to V and N flags.
func TestBITFlagCopying(t *testing.T) {
	tests := []struct {
		name        string
		memValue    byte
		accumulator byte
		expectN     bool
		expectV     bool
		expectZ     bool
	}{
		{"Both flags set", 0xC0, 0xFF, true, true, false},   // bits 7 and 6 set
		{"Only N set", 0x80, 0xFF, true, false, false},      // only bit 7 set
		{"Only V set", 0x40, 0xFF, false, true, false},      // only bit 6 set
		{"Neither set", 0x00, 0xFF, false, false, true},     // no bits set, Z=1
		{"Zero result", 0xFF, 0x00, true, true, true},       // AND is zero
		{"Nonzero result", 0xFF, 0x01, true, true, false},   // bit 0 matches
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			bus.SetResetVector(0x8000)

			// Set memory location and accumulator
			bus.memory[0x0010] = tt.memValue
			cpu.A = tt.accumulator

			// BIT $10
			program := []byte{
				0x24, 0x10, // BIT $10
			}
			bus.LoadProgram(0x8000, program)

			cpu.Reset()
			cpu.A = tt.accumulator

			for cpu.Cycles > 0 {
				cpu.Step()
			}

			cpu.Step()
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			if cpu.GetFlag(0x80) != tt.expectN {
				t.Errorf("Negative flag: expected %v, got %v", tt.expectN, cpu.GetFlag(0x80))
			}
			if cpu.GetFlag(0x40) != tt.expectV {
				t.Errorf("Overflow flag: expected %v, got %v", tt.expectV, cpu.GetFlag(0x40))
			}
			if cpu.GetFlag(0x02) != tt.expectZ {
				t.Errorf("Zero flag: expected %v, got %v", tt.expectZ, cpu.GetFlag(0x02))
			}
		})
	}
}

// ========== Branch Timing Tests ==========

// TestBranchTiming tests cycle-accurate branch instruction timing.
func TestBranchTiming(t *testing.T) {
	tests := []struct {
		name           string
		setupFlag      bool
		targetOffset   int8
		baseCycles     int
		expectedPC     uint16
	}{
		{"Not taken", false, 5, 2, 0x8002},     // BEQ not taken (Z=0), 2 cycles, PC advances by 2
		{"Taken same page", true, 5, 3, 0x8007}, // BEQ taken (Z=1), forward 5, same page, 3 cycles
		{"Taken backward same page", true, -5, 4, 0x7FFD}, // BEQ taken, backward 5, crosses page (0x80 -> 0x7F), 4 cycles
		{"Taken cross page forward", true, 127, 3, 0x8081}, // BEQ taken, forward 127, same page (0x80 -> 0x80), 3 cycles
		{"Taken cross page backward", true, -127, 4, 0x7F83}, // BEQ taken, backward crosses page, 4 cycles
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			bus.SetResetVector(0x8000)

			// BEQ (0xF0) - branch if zero flag set
			program := []byte{
				0xF0, byte(tt.targetOffset), // BEQ offset
			}
			bus.LoadProgram(0x8000, program)

			cpu.Reset()

			// Set zero flag based on test case
			cpu.SetFlag(0x02, tt.setupFlag)

			for cpu.Cycles > 0 {
				cpu.Step()
			}

			// Execute BEQ and count cycles
			cpu.Step()
			steps := 1
			for cpu.Cycles > 0 {
				cpu.Step()
				steps++
			}

			if steps != tt.baseCycles {
				t.Errorf("Expected %d cycles, got %d", tt.baseCycles, steps)
			}

			if cpu.PC != tt.expectedPC {
				t.Errorf("Expected PC to be 0x%04X, got 0x%04X", tt.expectedPC, cpu.PC)
			}
		})
	}
}

// ========== Stack Wraparound Tests ==========

// TestStackWraparound tests that the stack wraps from $0100 to $01FF.
func TestStackWraparound(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// Set SP to 0x00 (bottom of stack)
	cpu.SP = 0x00

	// Push a value - should wrap to 0xFF
	program := []byte{
		0x48, // PHA
	}
	bus.LoadProgram(0x8000, program)

	cpu.Reset()
	cpu.A = 0x42
	cpu.SP = 0x00

	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute PHA
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// SP should wrap to 0xFF
	if cpu.SP != 0xFF {
		t.Errorf("Expected SP to wrap to 0xFF, got 0x%02X", cpu.SP)
	}

	// Value should be at $0100
	if bus.memory[0x0100] != 0x42 {
		t.Errorf("Expected value at $0100, got 0x%02X", bus.memory[0x0100])
	}
}

// ========== Complex Instruction Sequences ==========

// TestROLWithCarry tests rotate left through carry.
func TestROLWithCarry(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// ROL A with different carry states
	program := []byte{
		0x2A, // ROL A
	}
	bus.LoadProgram(0x8000, program)

	tests := []struct {
		name          string
		initialA      byte
		initialCarry  bool
		expectedA     byte
		expectedCarry bool
	}{
		{"No carry in, no carry out", 0x40, false, 0x80, false},
		{"Carry in, no carry out", 0x40, true, 0x81, false},
		{"No carry in, carry out", 0x80, false, 0x00, true},
		{"Carry in, carry out", 0x80, true, 0x01, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu.Reset()
			cpu.A = tt.initialA
			cpu.SetFlag(0x01, tt.initialCarry)

			for cpu.Cycles > 0 {
				cpu.Step()
			}

			cpu.Step()
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			if cpu.A != tt.expectedA {
				t.Errorf("Expected A to be 0x%02X, got 0x%02X", tt.expectedA, cpu.A)
			}

			if cpu.GetFlag(0x01) != tt.expectedCarry {
				t.Errorf("Expected carry flag to be %v, got %v", tt.expectedCarry, cpu.GetFlag(0x01))
			}
		})
	}
}

// TestJSRRTS tests JSR/RTS preserve return address correctly.
func TestJSRRTS(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// Program:
	// $8000: JSR $8010
	// $8003: NOP (return point)
	// $8010: RTS
	program := []byte{
		0x20, 0x10, 0x80, // JSR $8010
		0xEA,             // NOP (should execute after RTS)
	}
	bus.LoadProgram(0x8000, program)

	// Subroutine at $8010
	bus.memory[0x8010] = 0x60 // RTS

	cpu.Reset()

	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute JSR
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// PC should be at $8010 (subroutine)
	if cpu.PC != 0x8010 {
		t.Errorf("After JSR: expected PC to be 0x8010, got 0x%04X", cpu.PC)
	}

	// Execute RTS
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// PC should be at $8003 (after JSR instruction)
	if cpu.PC != 0x8003 {
		t.Errorf("After RTS: expected PC to be 0x8003, got 0x%04X", cpu.PC)
	}
}

// TestCompareInstructions tests CMP, CPX, CPY set flags correctly.
func TestCompareInstructions(t *testing.T) {
	tests := []struct {
		name   string
		reg    byte
		operand byte
		carry  bool // Result: reg >= operand
		zero   bool // Result: reg == operand
		neg    bool // Result: (reg - operand) is negative
	}{
		{"Equal", 0x42, 0x42, true, true, false},
		{"Greater", 0x50, 0x42, true, false, false},
		{"Less", 0x42, 0x50, false, false, true},
		{"Zero vs Zero", 0x00, 0x00, true, true, false},
		{"Max vs Max", 0xFF, 0xFF, true, true, false},
		{"Zero vs Max", 0x00, 0xFF, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			bus.SetResetVector(0x8000)

			// CMP immediate
			program := []byte{
				0xC9, tt.operand, // CMP #operand
			}
			bus.LoadProgram(0x8000, program)

			cpu.Reset()
			cpu.A = tt.reg

			for cpu.Cycles > 0 {
				cpu.Step()
			}

			cpu.Step()
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			if cpu.GetFlag(0x01) != tt.carry {
				t.Errorf("Carry flag: expected %v, got %v", tt.carry, cpu.GetFlag(0x01))
			}
			if cpu.GetFlag(0x02) != tt.zero {
				t.Errorf("Zero flag: expected %v, got %v", tt.zero, cpu.GetFlag(0x02))
			}
			if cpu.GetFlag(0x80) != tt.neg {
				t.Errorf("Negative flag: expected %v, got %v", tt.neg, cpu.GetFlag(0x80))
			}
		})
	}
}

// TestPHPBreakFlag tests that PHP pushes status with B flag set.
func TestPHPBreakFlag(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	program := []byte{
		0x08, // PHP
		0x68, // PLA
	}
	bus.LoadProgram(0x8000, program)

	cpu.Reset()

	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute PHP
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute PLA to get pushed status
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// A should contain status with B flag (0x10) and unused flag (0x20) set
	if cpu.A&0x10 == 0 {
		t.Error("PHP should push status with B flag set")
	}
	if cpu.A&0x20 == 0 {
		t.Error("PHP should push status with unused flag set")
	}
}
