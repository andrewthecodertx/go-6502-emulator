package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Flag instructions for WDC65C02

// CLC clears the carry flag.
func CLC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagCarry, false)
}

// SEC sets the carry flag.
func SEC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagCarry, true)
}

// CLI clears the interrupt disable flag.
func CLI(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagInterruptDisable, false)
}

// SEI sets the interrupt disable flag.
func SEI(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagInterruptDisable, true)
}

// CLD clears the decimal mode flag.
func CLD(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagDecimal, false)
}

// SED sets the decimal mode flag.
func SED(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagDecimal, true)
}

// CLV clears the overflow flag.
func CLV(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagOverflow, false)
}

// NOP does nothing.
func NOP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	// Do nothing
}
