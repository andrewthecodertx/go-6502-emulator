package mos6502

// Addressing mode wrappers that delegate to BaseCPU.
// These exist to satisfy the instruction map signature.

func (c *CPU) addrImmediate() (uint16, bool) {
	return c.BaseCPU.AddrImmediate()
}

func (c *CPU) addrZeroPage() (uint16, bool) {
	return c.BaseCPU.AddrZeroPage()
}

func (c *CPU) addrZeroPageX() (uint16, bool) {
	return c.BaseCPU.AddrZeroPageX()
}

func (c *CPU) addrZeroPageY() (uint16, bool) {
	return c.BaseCPU.AddrZeroPageY()
}

func (c *CPU) addrAbsolute() (uint16, bool) {
	return c.BaseCPU.AddrAbsolute()
}

func (c *CPU) addrAbsoluteX() (uint16, bool) {
	return c.BaseCPU.AddrAbsoluteX()
}

func (c *CPU) addrAbsoluteY() (uint16, bool) {
	return c.BaseCPU.AddrAbsoluteY()
}

func (c *CPU) addrIndirect() (uint16, bool) {
	return c.BaseCPU.AddrIndirect()
}

func (c *CPU) addrIndirectX() (uint16, bool) {
	return c.BaseCPU.AddrIndirectX()
}

func (c *CPU) addrIndirectY() (uint16, bool) {
	return c.BaseCPU.AddrIndirectY()
}

func (c *CPU) addrRelative() (uint16, bool) {
	return c.BaseCPU.AddrRelative()
}

