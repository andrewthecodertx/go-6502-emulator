package mos6502

import (
	"testing"
)

// SimpleRAM implements a basic 64KB RAM for testing.
type SimpleRAM struct {
	memory [0x10000]byte
}

func (r *SimpleRAM) Read(addr uint16) byte {
	return r.memory[addr]
}

func (r *SimpleRAM) Write(addr uint16, data byte) {
	r.memory[addr] = data
}

// NewSimpleRAM creates a new SimpleRAM instance.
func NewSimpleRAM() *SimpleRAM {
	return &SimpleRAM{}
}

// LoadProgram loads a program into memory at the specified address.
func (r *SimpleRAM) LoadProgram(addr uint16, program []byte) {
	for i, b := range program {
		r.memory[addr+uint16(i)] = b
	}
}

// SetResetVector sets the reset vector to point to the given address.
func (r *SimpleRAM) SetResetVector(addr uint16) {
	r.memory[0xFFFC] = byte(addr & 0xFF)
	r.memory[0xFFFD] = byte(addr >> 8)
}

// TestCPUInitialization tests that a new CPU is created with correct initial state.
func TestCPUInitialization(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	if cpu.SP != 0xFD {
		t.Errorf("Expected SP to be 0xFD, got 0x%02X", cpu.SP)
	}

	if cpu.Status != 0x34 {
		t.Errorf("Expected Status to be 0x34, got 0x%02X", cpu.Status)
	}

	if cpu.Bus == nil {
		t.Error("Expected Bus to be set")
	}
}

// TestCPUReset tests the reset behavior according to the datasheet.
// According to the MOS 6502 datasheet:
// - Reset takes 7 cycles
// - A, X, Y are set to 0
// - SP is set to 0xFD
// - Status is set to 0x34 (with unused bit set)
// - PC is loaded from reset vector at $FFFC-$FFFD
func TestCPUReset(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	// Set reset vector to point to 0x8000
	bus.SetResetVector(0x8000)

	// Reset the CPU
	cpu.Reset()

	// Check register initialization
	if cpu.A != 0 {
		t.Errorf("Expected A to be 0 after reset, got 0x%02X", cpu.A)
	}

	if cpu.X != 0 {
		t.Errorf("Expected X to be 0 after reset, got 0x%02X", cpu.X)
	}

	if cpu.Y != 0 {
		t.Errorf("Expected Y to be 0 after reset, got 0x%02X", cpu.Y)
	}

	if cpu.SP != 0xFD {
		t.Errorf("Expected SP to be 0xFD after reset, got 0x%02X", cpu.SP)
	}

	// Check status register (0x34 = 0b00110100)
	// Bit 5 (unused) should be 1
	// Bit 4 (B flag) should be 1
	// Bit 2 (I flag - interrupt disable) should be 1
	expectedStatus := byte(0x34 | 0x20) // Ensure unused bit is set
	if cpu.Status != expectedStatus {
		t.Errorf("Expected Status to be 0x%02X after reset, got 0x%02X", expectedStatus, cpu.Status)
	}

	// Check PC is loaded from reset vector
	if cpu.PC != 0x8000 {
		t.Errorf("Expected PC to be 0x8000 (from reset vector), got 0x%04X", cpu.PC)
	}

	// Check that reset takes 6 cycles (per NMOS 6502 datasheet page 8)
	if cpu.Cycles != 6 {
		t.Errorf("Expected 6 cycles after reset, got %d", cpu.Cycles)
	}
}

// TestResetVectorReading tests that the reset vector is read correctly.
func TestResetVectorReading(t *testing.T) {
	tests := []struct {
		name     string
		vectorLo byte
		vectorHi byte
		expected uint16
	}{
		{"Vector at 0x0000", 0x00, 0x00, 0x0000},
		{"Vector at 0x8000", 0x00, 0x80, 0x8000},
		{"Vector at 0xFFFF", 0xFF, 0xFF, 0xFFFF},
		{"Vector at 0x1234", 0x34, 0x12, 0x1234},
		{"Vector at 0xC000", 0x00, 0xC0, 0xC000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			// Set reset vector
			bus.memory[0xFFFC] = tt.vectorLo
			bus.memory[0xFFFD] = tt.vectorHi

			cpu.Reset()

			if cpu.PC != tt.expected {
				t.Errorf("Expected PC to be 0x%04X, got 0x%04X", tt.expected, cpu.PC)
			}
		})
	}
}

// TestClockCycling tests that the CPU cycles through instructions correctly.
func TestClockCycling(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	// Set reset vector
	bus.SetResetVector(0x8000)

	// Load a NOP instruction at 0x8000 (NOP = 0xEA, takes 2 cycles)
	bus.LoadProgram(0x8000, []byte{0xEA}) // NOP

	cpu.Reset()

	// After reset, cycles should be 6
	initialCycles := cpu.Cycles
	if initialCycles != 6 {
		t.Errorf("Expected 6 cycles after reset (per NMOS 6502 datasheet page 8), got %d", initialCycles)
	}

	// Execute cycles until the reset cycles are consumed (6 cycles)
	for i := 0; i < 6; i++ {
		cpu.Step()
	}

	// After 6 steps, the reset cycles should be consumed
	// Now the CPU should fetch the next instruction
	if cpu.Cycles != 0 {
		t.Errorf("Expected 0 cycles after consuming reset cycles, got %d", cpu.Cycles)
	}

	// Step once to fetch and execute NOP
	// Note: Step() fetches the instruction, sets cycles to instruction.Cycles,
	// then decrements by 1, so after one Step() with a 2-cycle instruction,
	// cycles will be 1
	cpu.Step()

	// NOP takes 2 cycles, so after one Step() that fetches it, cycles should be 1
	if cpu.Cycles != 1 {
		t.Errorf("Expected 1 cycle remaining after fetching NOP, got %d", cpu.Cycles)
	}

	// PC should have advanced by 1 (NOP is 1 byte)
	expectedPC := uint16(0x8001)
	if cpu.PC != expectedPC {
		t.Errorf("Expected PC to be 0x%04X after NOP, got 0x%04X", expectedPC, cpu.PC)
	}
}

// TestStepExecution tests that Step() correctly decrements cycles.
func TestStepExecution(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)
	cpu.Reset()

	// cycles should be 6 after reset (per NMOS 6502 datasheet page 8)
	if cpu.Cycles != 6 {
		t.Fatalf("Expected 6 cycles after reset, got %d", cpu.Cycles)
	}

	// Step through the reset cycles
	for cpu.Cycles > 0 {
		before := cpu.Cycles
		cpu.Step()
		after := cpu.Cycles

		// Each Step() should decrement cycles by 1
		if before-after != 1 {
			t.Errorf("Expected cycles to decrement by 1, went from %d to %d", before, after)
		}
	}

	// After all cycles are consumed, cycles should be 0
	if cpu.Cycles != 0 {
		t.Errorf("Expected 0 cycles after stepping through reset, got %d", cpu.Cycles)
	}
}

// TestBasicInstructionExecution tests that the CPU can execute basic instructions.
func TestBasicInstructionExecution(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	// Set reset vector to 0x8000
	bus.SetResetVector(0x8000)

	// Load a simple program:
	// LDA #$42  (0xA9 0x42) - Load 0x42 into accumulator
	// NOP       (0xEA)      - No operation
	program := []byte{
		0xA9, 0x42, // LDA #$42
		0xEA,       // NOP
	}
	bus.LoadProgram(0x8000, program)

	cpu.Reset()

	// Consume reset cycles
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute LDA #$42 (takes 2 cycles)
	cpu.Step() // Fetch and queue instruction
	for cpu.Cycles > 0 {
		cpu.Step() // Execute cycles
	}

	// Check that A register contains 0x42
	if cpu.A != 0x42 {
		t.Errorf("Expected A to be 0x42 after LDA #$42, got 0x%02X", cpu.A)
	}

	// Check that PC advanced by 2 bytes
	if cpu.PC != 0x8002 {
		t.Errorf("Expected PC to be 0x8002, got 0x%04X", cpu.PC)
	}

	// Check flags: 0x42 is neither zero nor negative
	if cpu.GetFlag(0x02) { // Zero flag
		t.Error("Expected Zero flag to be clear")
	}
	if cpu.GetFlag(0x80) { // Negative flag
		t.Error("Expected Negative flag to be clear")
	}
}

// TestLDAFlagBehavior tests that LDA sets flags correctly.
func TestLDAFlagBehavior(t *testing.T) {
	tests := []struct {
		name     string
		value    byte
		zeroFlag bool
		negFlag  bool
	}{
		{"Load zero", 0x00, true, false},
		{"Load positive", 0x42, false, false},
		{"Load negative", 0x80, false, true},
		{"Load 0xFF", 0xFF, false, true},
		{"Load 0x7F", 0x7F, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			bus.SetResetVector(0x8000)

			// LDA immediate
			program := []byte{0xA9, tt.value}
			bus.LoadProgram(0x8000, program)

			cpu.Reset()

			// Consume reset cycles
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			// Execute LDA
			cpu.Step()
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			// Check accumulator value
			if cpu.A != tt.value {
				t.Errorf("Expected A to be 0x%02X, got 0x%02X", tt.value, cpu.A)
			}

			// Check zero flag
			if cpu.GetFlag(0x02) != tt.zeroFlag {
				t.Errorf("Expected Zero flag to be %v, got %v", tt.zeroFlag, cpu.GetFlag(0x02))
			}

			// Check negative flag
			if cpu.GetFlag(0x80) != tt.negFlag {
				t.Errorf("Expected Negative flag to be %v, got %v", tt.negFlag, cpu.GetFlag(0x80))
			}
		})
	}
}

// TestMemoryReadWrite tests that the CPU correctly reads and writes to memory.
func TestMemoryReadWrite(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// Program:
	// LDA #$55      (0xA9 0x55)
	// STA $0200     (0x8D 0x00 0x02)
	program := []byte{
		0xA9, 0x55,       // LDA #$55
		0x8D, 0x00, 0x02, // STA $0200
	}
	bus.LoadProgram(0x8000, program)

	cpu.Reset()

	// Consume reset cycles
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute LDA #$55
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Check A register
	if cpu.A != 0x55 {
		t.Errorf("Expected A to be 0x55, got 0x%02X", cpu.A)
	}

	// Execute STA $0200
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Check that memory at $0200 contains 0x55
	if bus.memory[0x0200] != 0x55 {
		t.Errorf("Expected memory at $0200 to be 0x55, got 0x%02X", bus.memory[0x0200])
	}
}

// TestTransferInstructions tests register transfer instructions.
func TestTransferInstructions(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// Program:
	// LDA #$42      (0xA9 0x42)
	// TAX           (0xAA)
	// TAY           (0xA8)
	program := []byte{
		0xA9, 0x42, // LDA #$42
		0xAA,       // TAX
		0xA8,       // TAY
	}
	bus.LoadProgram(0x8000, program)

	cpu.Reset()

	// Consume reset cycles
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute LDA #$42
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	if cpu.A != 0x42 {
		t.Fatalf("Expected A to be 0x42, got 0x%02X", cpu.A)
	}

	// Execute TAX
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	if cpu.X != 0x42 {
		t.Errorf("Expected X to be 0x42 after TAX, got 0x%02X", cpu.X)
	}

	// Execute TAY
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	if cpu.Y != 0x42 {
		t.Errorf("Expected Y to be 0x42 after TAY, got 0x%02X", cpu.Y)
	}
}

// TestStackOperations tests push and pull operations.
func TestStackOperations(t *testing.T) {
	bus := NewSimpleRAM()
	cpu := NewCPU(bus)

	bus.SetResetVector(0x8000)

	// Program:
	// LDA #$AA      (0xA9 0xAA)
	// PHA           (0x48)
	// LDA #$00      (0xA9 0x00)
	// PLA           (0x68)
	program := []byte{
		0xA9, 0xAA, // LDA #$AA
		0x48,       // PHA
		0xA9, 0x00, // LDA #$00
		0x68,       // PLA
	}
	bus.LoadProgram(0x8000, program)

	cpu.Reset()

	initialSP := cpu.SP

	// Consume reset cycles
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute LDA #$AA
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Execute PHA
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Stack pointer should have decremented
	if cpu.SP != initialSP-1 {
		t.Errorf("Expected SP to be 0x%02X after PHA, got 0x%02X", initialSP-1, cpu.SP)
	}

	// Check that 0xAA is on the stack
	stackAddr := uint16(0x0100) + uint16(cpu.SP) + 1
	if bus.memory[stackAddr] != 0xAA {
		t.Errorf("Expected 0xAA on stack at 0x%04X, got 0x%02X", stackAddr, bus.memory[stackAddr])
	}

	// Execute LDA #$00
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	if cpu.A != 0x00 {
		t.Fatalf("Expected A to be 0x00, got 0x%02X", cpu.A)
	}

	// Execute PLA
	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// A should be restored to 0xAA
	if cpu.A != 0xAA {
		t.Errorf("Expected A to be 0xAA after PLA, got 0x%02X", cpu.A)
	}

	// Stack pointer should be back to initial value
	if cpu.SP != initialSP {
		t.Errorf("Expected SP to be 0x%02X after PLA, got 0x%02X", initialSP, cpu.SP)
	}
}

// TestCycleAccuracy tests that instructions consume the correct number of cycles.
func TestCycleAccuracy(t *testing.T) {
	tests := []struct {
		name          string
		instruction   []byte
		expectedCycles byte
	}{
		{"NOP", []byte{0xEA}, 2},
		{"LDA Immediate", []byte{0xA9, 0x42}, 2},
		{"LDA Zero Page", []byte{0xA5, 0x10}, 3},
		{"LDA Absolute", []byte{0xAD, 0x00, 0x02}, 4},
		{"STA Zero Page", []byte{0x85, 0x10}, 3},
		{"STA Absolute", []byte{0x8D, 0x00, 0x02}, 4},
		{"TAX", []byte{0xAA}, 2},
		{"INX", []byte{0xE8}, 2},
		{"DEX", []byte{0xCA}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bus := NewSimpleRAM()
			cpu := NewCPU(bus)

			bus.SetResetVector(0x8000)
			bus.LoadProgram(0x8000, tt.instruction)

			cpu.Reset()

			// Consume reset cycles
			for cpu.Cycles > 0 {
				cpu.Step()
			}

			// Execute instruction - Step() fetches it and decrements cycles by 1
			// So after one Step(), cycles will be instruction.Cycles - 1
			cpu.Step()

			expectedRemaining := tt.expectedCycles - 1
			if cpu.Cycles != expectedRemaining {
				t.Errorf("Expected %d cycles remaining (instruction takes %d cycles, Step decrements by 1), got %d",
					expectedRemaining, tt.expectedCycles, cpu.Cycles)
			}

			// Now step through the remaining cycles to verify the total
			stepsNeeded := int(cpu.Cycles)
			for i := 0; i < stepsNeeded; i++ {
				cpu.Step()
			}

			// After all cycles are consumed, the CPU should be ready to fetch the next instruction
			if cpu.Cycles != 0 {
				t.Errorf("Expected 0 cycles after consuming all instruction cycles, got %d", cpu.Cycles)
			}
		})
	}
}
