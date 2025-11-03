package instructions

import "github.com/andrewthecodertx/go-65c02-emulator/pkg/core"

// Logic instructions for NMOS 6502

// AND performs a bitwise AND.
func AND(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.A &= data
	c.SetZN(c.A)

	if pageCrossed {
		c.Cycles++
	}
}

// ORA performs a bitwise OR.
func ORA(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.A |= data
	c.SetZN(c.A)

	if pageCrossed {
		c.Cycles++
	}
}

// EOR performs a bitwise XOR.
func EOR(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.A ^= data
	c.SetZN(c.A)

	if pageCrossed {
		c.Cycles++
	}
}

// BIT tests bits.
func BIT(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.SetFlag(core.FlagNegative, data&0x80 != 0)  // Negative from bit 7
	c.SetFlag(core.FlagOverflow, data&0x40 != 0)  // Overflow from bit 6
	c.SetFlag(core.FlagZero, data&c.A == 0)   // Zero if AND is zero
}
