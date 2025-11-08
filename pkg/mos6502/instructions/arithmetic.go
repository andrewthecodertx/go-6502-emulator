package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Arithmetic instructions for NMOS 6502

// ADC adds with carry.
func ADC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	carry := byte(0)
	if c.GetFlag(core.FlagCarry) {
		carry = 1
	}

	result := uint16(c.A) + uint16(data) + uint16(carry)
	c.SetFlag(core.FlagCarry, result > 0xFF) // Carry
	c.SetFlag(core.FlagOverflow, ((uint16(c.A)^result)&(uint16(data)^result))&0x80 != 0) // Overflow
	c.A = byte(result)
	c.SetZN(c.A)

	if pageCrossed {
		c.Cycles++
	}
}

// SBC subtracts with carry.
func SBC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	carry := byte(0)
	if c.GetFlag(core.FlagCarry) {
		carry = 1
	}

	result := uint16(c.A) - uint16(data) - (1 - uint16(carry))
	c.SetFlag(core.FlagCarry, result < 0x100) // Carry (borrow)
	c.SetFlag(core.FlagOverflow, ((uint16(c.A)^result)&(^uint16(data)^result))&0x80 != 0) // Overflow
	c.A = byte(result)
	c.SetZN(c.A)

	if pageCrossed {
		c.Cycles++
	}
}

// compare is a helper function for comparison operations.
func compare(c *core.BaseCPU, a, b byte) {
	result := a - b
	c.SetFlag(core.FlagCarry, a >= b)        // Carry
	c.SetFlag(core.FlagNegative, result&0x80 != 0) // Negative
	c.SetFlag(core.FlagZero, result == 0)      // Zero
}

// CMP compares the accumulator.
func CMP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	compare(c, c.A, data)

	if pageCrossed {
		c.Cycles++
	}
}

// CPX compares the X register.
func CPX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	compare(c, c.X, data)
}

// CPY compares the Y register.
func CPY(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	compare(c, c.Y, data)
}
