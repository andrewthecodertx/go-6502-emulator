// Package mos6502 provides a cycle-accurate emulator for the MOS 6502 CPU.
//
// The MOS 6502 is an 8-bit microprocessor that was widely used in home computers
// and gaming consoles during the 1970s and 1980s, including the Commodore 64,
// Apple II, and Nintendo Entertainment System.
//
// This implementation includes:
//   - All 56 legal 6502 instructions with proper flag handling
//   - 13 addressing modes including the JMP indirect page boundary bug
//   - Cycle-accurate execution with page-crossing penalties
//   - Interrupt support (NMI, IRQ, RESET)
//   - Configurable bus interface for flexible memory implementations
//
// The CPU communicates with memory and peripherals through the Bus interface,
// allowing custom memory mappers and I/O devices to be implemented.
//
// Example usage:
//
//	type SimpleRAM struct {
//	   memory [0x10000]byte
//	}
//
//	func (r *SimpleRAM) Read(addr uint16) byte  { return r.memory[addr] }
//	func (r *SimpleRAM) Write(addr uint16, data byte) { r.memory[addr] = data }
//
//	func main() {
//	   bus := &SimpleRAM{}
//	   cpu := mos6502.NewCPU(bus)
//	   cpu.Reset()
//	   cpu.Run()
//	}
package mos6502

import (
	"fmt"

	"github.com/andrewthecodertx/go-6502-emulator/pkg/core"
)

// CPU represents an NMOS 6502 processor
type CPU struct {
	*core.BaseCPU
}

// NewCPU creates a new NMOS 6502 CPU with the given bus
func NewCPU(bus core.Bus) *CPU {
	return &CPU{
		BaseCPU: core.NewBaseCPU(bus, core.VariantNMOS),
	}
}

// Reset is inherited from BaseCPU and uses VariantNMOS (6 cycles)

func (c *CPU) Run() {
	for !c.Halted {
		c.Step()
	}
}

func (c *CPU) Step() {
	if c.Cycles == 0 {
		if c.ResetPending {
			c.HandleReset()
			return
		}

		if c.NMIPending {
			c.HandleNMI()
			return
		}

		if c.IRQPending && !c.GetFlag(core.FlagInterruptDisable) {
			c.HandleIRQ()
			return
		}

		opcode := c.Bus.Read(c.PC)
		c.PC++

		instruction, ok := instructionMap[opcode]
		if !ok {
			// TODO: Handle illegal opcodes
			fmt.Printf("Unknown opcode: 0x%02X\n", opcode)
			c.Halted = true // Halt for now
			return
		}

		var addr uint16
		var pageCrossed bool

		if instruction.AddrMode != nil {
			addr, pageCrossed = instruction.AddrMode(c)
		}

		instruction.Operation(c, addr, pageCrossed)
		c.Cycles += instruction.Cycles
	}

	c.Cycles--
}

// Helper methods are inherited from BaseCPU:
// - GetFlag/SetFlag for status register manipulation
// - Push/Pull for stack operations
// - HandleReset/HandleNMI/HandleIRQ for interrupt handling
// - GetCycles for testing/debugging
