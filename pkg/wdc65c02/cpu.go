// Package wdc65c02 provides a cycle-accurate emulator for the WDC65C02 CPU.
//
// The WDC65C02 is an enhanced CMOS version of the NMOS 6502, produced by the
// Western Design Center. It includes bug fixes and new features compared to
// the original NMOS 6502:
//
//   - 27 new instructions including STZ, BRA, PHX/PHY/PLX/PLY, TSB/TRB, WAI/STP,
//     and bit manipulation instructions (BBR/BBS, RMB/SMB)
//   - 2 new addressing modes: Zero Page Indirect and Absolute Indexed Indirect
//   - Fixed JMP indirect page boundary bug
//   - Decimal mode automatically cleared on interrupts
//   - All illegal opcodes become NOPs
//   - 7-cycle reset sequence (vs 6 cycles on NMOS 6502)
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
//	   cpu := wdc65c02.NewCPU(bus)
//	   cpu.Reset()
//	   cpu.Run()
//	}
package wdc65c02

import (
	"github.com/andrewthecodertx/go-6502-emulator/pkg/core"
)

// CPU represents a WDC65C02 processor
type CPU struct {
	*core.BaseCPU
}

// NewCPU creates a new WDC65C02 CPU with the given bus
func NewCPU(bus core.Bus) *CPU {
	return &CPU{
		BaseCPU: core.NewBaseCPU(bus, core.VariantWDC65C02),
	}
}

// Run executes instructions until the CPU is halted
func (c *CPU) Run() {
	for !c.Halted {
		c.Step()
	}
}

// Step executes a single CPU cycle
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
			// On WDC65C02, all illegal opcodes are NOP
			// Most are 1-byte, 1-cycle NOPs, but some vary
			// For now, treat as 1-byte, 1-cycle NOP
			c.Cycles = 1
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
