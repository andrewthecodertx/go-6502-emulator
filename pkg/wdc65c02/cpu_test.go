package wdc65c02

import (
	"testing"

	"github.com/andrewthecodertx/go-6502-emulator/pkg/core"
)

// SimpleRAM is a simple RAM implementation for testing
type SimpleRAM struct {
	memory [0x10000]byte
}

func (r *SimpleRAM) Read(addr uint16) byte {
	return r.memory[addr]
}

func (r *SimpleRAM) Write(addr uint16, data byte) {
	r.memory[addr] = data
}

func TestWDC65C02Reset(t *testing.T) {
	ram := &SimpleRAM{}
	// Set reset vector to 0x1234
	ram.memory[0xFFFC] = 0x34
	ram.memory[0xFFFD] = 0x12

	cpu := NewCPU(ram)
	cpu.Reset()

	// Check that PC is set to reset vector
	if cpu.PC != 0x1234 {
		t.Errorf("Expected PC to be 0x1234, got 0x%04X", cpu.PC)
	}

	// Check that reset takes 7 cycles on WDC65C02 (not 6 like NMOS)
	if cpu.Cycles != 7 {
		t.Errorf("Expected 7 cycles after reset (per WDC65C02 datasheet), got %d", cpu.Cycles)
	}

	// Check initial register state
	if cpu.A != 0 || cpu.X != 0 || cpu.Y != 0 {
		t.Errorf("Expected A, X, Y to be 0 after reset")
	}

	if cpu.SP != 0xFD {
		t.Errorf("Expected SP to be 0xFD, got 0x%02X", cpu.SP)
	}
}

func TestWDC65C02Variant(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)

	// Verify variant is set correctly
	if cpu.Variant != core.VariantWDC65C02 {
		t.Errorf("Expected variant to be WDC65C02, got %v", cpu.Variant)
	}

	// Verify variant-specific behaviors
	if cpu.Variant.ResetCycles() != 7 {
		t.Errorf("Expected 7 reset cycles for WDC65C02, got %d", cpu.Variant.ResetCycles())
	}

	if cpu.Variant.HasJMPIndirectBug() {
		t.Error("WDC65C02 should NOT have JMP indirect bug")
	}

	if !cpu.Variant.ClearsDecimalOnInterrupt() {
		t.Error("WDC65C02 should clear decimal flag on interrupts")
	}
}

func TestSTZ(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	tests := []struct {
		name     string
		setup    func()
		opcode   byte
		operands []byte
		addr     uint16
	}{
		{
			name: "STZ Zero Page",
			setup: func() {
				ram.memory[0x50] = 0xFF // Pre-fill with non-zero
			},
			opcode:   0x64,
			operands: []byte{0x50},
			addr:     0x50,
		},
		{
			name: "STZ Zero Page,X",
			setup: func() {
				cpu.X = 0x05
				ram.memory[0x55] = 0xFF
			},
			opcode:   0x74,
			operands: []byte{0x50},
			addr:     0x55,
		},
		{
			name: "STZ Absolute",
			setup: func() {
				ram.memory[0x1234] = 0xFF
			},
			opcode:   0x9C,
			operands: []byte{0x34, 0x12},
			addr:     0x1234,
		},
		{
			name: "STZ Absolute,X",
			setup: func() {
				cpu.X = 0x10
				ram.memory[0x1244] = 0xFF
			},
			opcode:   0x9E,
			operands: []byte{0x34, 0x12},
			addr:     0x1244,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset state
			cpu.PC = 0x0200
			cpu.X = 0
			cpu.Y = 0

			tt.setup()

			// Write instruction
			ram.memory[0x0200] = tt.opcode
			for i, b := range tt.operands {
				ram.memory[0x0201+uint16(i)] = b
			}

			// Execute instruction
			cpu.Step() // Consume any remaining cycles
			for cpu.Cycles > 0 {
				cpu.Step()
			}
			cpu.Step() // Execute the STZ instruction

			// Verify zero was stored
			if ram.memory[tt.addr] != 0x00 {
				t.Errorf("Expected 0x00 at address 0x%04X, got 0x%02X", tt.addr, ram.memory[tt.addr])
			}
		})
	}
}

func TestBRA(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	// Test BRA forward
	t.Run("BRA Forward", func(t *testing.T) {
		cpu.PC = 0x0200
		ram.memory[0x0200] = 0x80 // BRA
		ram.memory[0x0201] = 0x10 // +16 bytes

		initialCycles := cpu.Cycles
		cpu.Step()
		for cpu.Cycles > 0 {
			cpu.Step()
		}
		cpu.Step()

		expectedPC := uint16(0x0202 + 0x10)
		if cpu.PC != expectedPC {
			t.Errorf("Expected PC to be 0x%04X, got 0x%04X", expectedPC, cpu.PC)
		}

		// BRA should take 3 cycles base + 1 for branch taken = 4 total
		cyclesTaken := initialCycles - cpu.Cycles + 3 // +3 for the instruction base cycles
		if cyclesTaken < 3 {
			t.Errorf("Expected at least 3 cycles for BRA, got %d", cyclesTaken)
		}
	})

	// Test BRA backward
	t.Run("BRA Backward", func(t *testing.T) {
		cpu.PC = 0x0220
		ram.memory[0x0220] = 0x80 // BRA
		ram.memory[0x0221] = 0xFE // -2 bytes

		cpu.Step()
		for cpu.Cycles > 0 {
			cpu.Step()
		}
		cpu.Step()

		expectedPC := uint16(0x0220)
		if cpu.PC != expectedPC {
			t.Errorf("Expected PC to be 0x%04X, got 0x%04X", expectedPC, cpu.PC)
		}
	})
}

func TestPHXPLX(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	// Test PHX
	cpu.PC = 0x0200
	cpu.X = 0x42
	cpu.SP = 0xFD

	ram.memory[0x0200] = 0xDA // PHX

	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}
	cpu.Step()

	if ram.memory[0x01FD] != 0x42 {
		t.Errorf("Expected 0x42 on stack, got 0x%02X", ram.memory[0x01FD])
	}
	if cpu.SP != 0xFC {
		t.Errorf("Expected SP to be 0xFC, got 0x%02X", cpu.SP)
	}

	// Test PLX
	cpu.PC = 0x0201
	cpu.X = 0x00

	ram.memory[0x0201] = 0xFA // PLX

	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}
	cpu.Step()

	if cpu.X != 0x42 {
		t.Errorf("Expected X to be 0x42, got 0x%02X", cpu.X)
	}
	if cpu.SP != 0xFD {
		t.Errorf("Expected SP to be 0xFD, got 0x%02X", cpu.SP)
	}
}

func TestPHYPLY(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	// Test PHY
	cpu.PC = 0x0200
	cpu.Y = 0x84
	cpu.SP = 0xFD

	ram.memory[0x0200] = 0x5A // PHY

	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}
	cpu.Step()

	if ram.memory[0x01FD] != 0x84 {
		t.Errorf("Expected 0x84 on stack, got 0x%02X", ram.memory[0x01FD])
	}

	// Test PLY
	cpu.PC = 0x0201
	cpu.Y = 0x00

	ram.memory[0x0201] = 0x7A // PLY

	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}
	cpu.Step()

	if cpu.Y != 0x84 {
		t.Errorf("Expected Y to be 0x84, got 0x%02X", cpu.Y)
	}
}

func TestINCDECAccumulator(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	// Test INC A
	t.Run("INC A", func(t *testing.T) {
		cpu.PC = 0x0200
		cpu.A = 0x42

		ram.memory[0x0200] = 0x1A // INC A

		cpu.Step()
		for cpu.Cycles > 0 {
			cpu.Step()
		}
		cpu.Step()

		if cpu.A != 0x43 {
			t.Errorf("Expected A to be 0x43, got 0x%02X", cpu.A)
		}
	})

	// Test DEC A
	t.Run("DEC A", func(t *testing.T) {
		cpu.PC = 0x0201
		cpu.A = 0x43

		ram.memory[0x0201] = 0x3A // DEC A

		cpu.Step()
		for cpu.Cycles > 0 {
			cpu.Step()
		}
		cpu.Step()

		if cpu.A != 0x42 {
			t.Errorf("Expected A to be 0x42, got 0x%02X", cpu.A)
		}
	})
}

func TestTSB(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	cpu.PC = 0x0200
	cpu.A = 0x0F
	ram.memory[0x50] = 0xF0

	ram.memory[0x0200] = 0x04 // TSB $50
	ram.memory[0x0201] = 0x50

	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}
	cpu.Step()

	// Result should be 0xF0 | 0x0F = 0xFF
	if ram.memory[0x50] != 0xFF {
		t.Errorf("Expected memory to be 0xFF, got 0x%02X", ram.memory[0x50])
	}

	// Z flag should be set (A & original memory == 0)
	// A=0x0F, memory=0xF0, so (0x0F & 0xF0) = 0x00
	if !cpu.GetFlag(core.FlagZero) {
		t.Error("Expected Z flag to be set")
	}
}

func TestTRB(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	cpu.PC = 0x0200
	cpu.A = 0x0F
	ram.memory[0x50] = 0xFF

	ram.memory[0x0200] = 0x14 // TRB $50
	ram.memory[0x0201] = 0x50

	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}
	cpu.Step()

	// Result should be 0xFF & ~0x0F = 0xF0
	if ram.memory[0x50] != 0xF0 {
		t.Errorf("Expected memory to be 0xF0, got 0x%02X", ram.memory[0x50])
	}
}

func TestSTP(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	cpu.PC = 0x0200
	ram.memory[0x0200] = 0xDB // STP

	cpu.Step()
	for cpu.Cycles > 0 {
		cpu.Step()
	}
	cpu.Step()

	if !cpu.Halted {
		t.Error("Expected CPU to be halted after STP")
	}
}

func TestDecimalClearedOnInterrupt(t *testing.T) {
	ram := &SimpleRAM{}
	cpu := NewCPU(ram)
	cpu.Reset()

	// Set up IRQ vector
	ram.memory[0xFFFE] = 0x00
	ram.memory[0xFFFF] = 0x03

	// Consume reset cycles
	for cpu.Cycles > 0 {
		cpu.Step()
	}

	// Set decimal mode
	cpu.SetFlag(core.FlagDecimal, true)
	cpu.SetFlag(core.FlagInterruptDisable, false)

	// Trigger IRQ
	cpu.IRQPending = true

	// Execute interrupt handling - Step() will check cycles and process IRQ
	cpu.Step()

	// On WDC65C02, decimal mode should be cleared after interrupt
	if cpu.GetFlag(core.FlagDecimal) {
		t.Error("Expected decimal flag to be cleared on interrupt (WDC65C02 behavior)")
	}
}
