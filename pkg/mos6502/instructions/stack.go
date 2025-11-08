package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Stack instructions for NMOS 6502

// PHA pushes the accumulator onto the stack.
func PHA(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Push(c.A)
}

// PLA pulls the accumulator from the stack.
func PLA(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.A = c.Pull()
	c.SetZN(c.A)
}

// PHP pushes the status register onto the stack.
func PHP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	// Set B flag and unused flag
	c.Push(c.Status | 0x10 | 0x20)
}

// PLP pulls the status register from the stack.
func PLP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Status = c.Pull()
	c.SetFlag(core.FlagUnused, true) // Unused flag always set
	c.SetFlag(core.FlagBreak, false) // B flag not actually stored
}
