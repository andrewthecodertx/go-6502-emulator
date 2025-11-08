package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Shift/Rotate instructions for WDC65C02

// ASL shifts left.
func ASL(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.SetFlag(core.FlagCarry, data&0x80 != 0)
	data <<= 1
	c.Bus.Write(addr, data)
	c.SetZN(data)
}

// ASLAccumulator shifts the accumulator left.
func ASLAccumulator(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagCarry, c.A&0x80 != 0)
	c.A <<= 1
	c.SetZN(c.A)
}

// LSR shifts right.
func LSR(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.SetFlag(core.FlagCarry, data&0x01 != 0)
	data >>= 1
	c.Bus.Write(addr, data)
	c.SetZN(data)
}

// LSRAccumulator shifts the accumulator right.
func LSRAccumulator(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SetFlag(core.FlagCarry, c.A&0x01 != 0)
	c.A >>= 1
	c.SetZN(c.A)
}

// ROL rotates left.
func ROL(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	carry := byte(0)
	if c.GetFlag(core.FlagCarry) {
		carry = 1
	}
	c.SetFlag(core.FlagCarry, data&0x80 != 0)
	data = (data << 1) | carry
	c.Bus.Write(addr, data)
	c.SetZN(data)
}

// ROLAccumulator rotates the accumulator left.
func ROLAccumulator(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	carry := byte(0)
	if c.GetFlag(core.FlagCarry) {
		carry = 1
	}
	c.SetFlag(core.FlagCarry, c.A&0x80 != 0)
	c.A = (c.A << 1) | carry
	c.SetZN(c.A)
}

// ROR rotates right.
func ROR(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	carry := byte(0)
	if c.GetFlag(core.FlagCarry) {
		carry = 1
	}
	c.SetFlag(core.FlagCarry, data&0x01 != 0)
	data = (data >> 1) | (carry << 7)
	c.Bus.Write(addr, data)
	c.SetZN(data)
}

// RORAccumulator rotates the accumulator right.
func RORAccumulator(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	carry := byte(0)
	if c.GetFlag(core.FlagCarry) {
		carry = 1
	}
	c.SetFlag(core.FlagCarry, c.A&0x01 != 0)
	c.A = (c.A >> 1) | (carry << 7)
	c.SetZN(c.A)
}
