package wdc65c02

// Addressing mode wrappers that delegate to BaseCPU.
// These exist to satisfy the instruction map signature.
// WDC65C02 includes all NMOS 6502 addressing modes plus two new ones.

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
	// WDC65C02 has the JMP indirect bug FIXED
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

// NEW WDC65C02 ADDRESSING MODES

// addrZeroPageIndirect - Zero Page Indirect addressing mode (($nn))
// This is a NEW addressing mode in WDC65C02
func (c *CPU) addrZeroPageIndirect() (uint16, bool) {
	return c.BaseCPU.AddrZeroPageIndirect()
}

// addrAbsoluteIndexedIndirect - Absolute Indexed Indirect addressing mode (($nnnn,X))
// This is a NEW addressing mode in WDC65C02
func (c *CPU) addrAbsoluteIndexedIndirect() (uint16, bool) {
	return c.BaseCPU.AddrAbsoluteIndexedIndirect()
}
