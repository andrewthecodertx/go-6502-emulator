package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Flag instructions for NMOS 6502

// CLC clears the carry flag.
func CLC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagCarry, false)
}

// CLD clears the decimal mode flag.
func CLD(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagDecimal, false)
}

// CLI clears the interrupt disable flag.
func CLI(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagInterruptDisable, false)
}

// CLV clears the overflow flag.
func CLV(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagOverflow, false)
}

// SEC sets the carry flag.
func SEC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagCarry, true)
}

// SED sets the decimal mode flag.
func SED(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagDecimal, true)
}

// SEI sets the interrupt disable flag.
func SEI(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagInterruptDisable, true)
}

// NOP does nothing.
func NOP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	// Do nothing
}
