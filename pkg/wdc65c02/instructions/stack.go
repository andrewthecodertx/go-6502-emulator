package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Stack instructions for WDC65C02

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
	c.Push(c.Status | core.FlagBreak | core.FlagUnused)
}

// PLP pulls the status register from the stack.
func PLP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Status = c.Pull()
	c.SetFlag(core.FlagUnused, true)  // Unused flag always set
	c.SetFlag(core.FlagBreak, false) // B flag not actually stored
}

// PHX pushes the X register onto the stack.
// NEW instruction in WDC65C02.
func PHX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Push(c.X)
}

// PHY pushes the Y register onto the stack.
// NEW instruction in WDC65C02.
func PHY(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Push(c.Y)
}

// PLX pulls the X register from the stack.
// NEW instruction in WDC65C02.
func PLX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.X = c.Pull()
	c.SetZN(c.X)
}

// PLY pulls the Y register from the stack.
// NEW instruction in WDC65C02.
func PLY(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Y = c.Pull()
	c.SetZN(c.Y)
}
