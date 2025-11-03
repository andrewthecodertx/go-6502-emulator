package core

// Addressing mode functions compute effective addresses for 6502 instructions.
//
// Each function returns:
//   - uint16: The effective memory address to operate on
//   - bool: Whether a page boundary was crossed (affects cycle count)
//
// Page crossing occurs when an indexed address moves to a different page
// (different high byte). Many instructions take an extra cycle when this happens.

// AddrImmediate handles immediate addressing mode (#$nn).
// The operand is the byte immediately following the opcode.
// Returns the address of the operand (PC) and advances PC.
func (c *BaseCPU) AddrImmediate() (uint16, bool) {
	addr := c.PC
	c.PC++
	return addr, false
}

// AddrZeroPage handles zero page addressing mode ($nn).
// Addresses the first 256 bytes of memory (0x0000-0x00FF).
// Uses single-byte address, faster than absolute addressing.
func (c *BaseCPU) AddrZeroPage() (uint16, bool) {
	addr := uint16(c.Bus.Read(c.PC))
	c.PC++
	return addr, false
}

// AddrZeroPageX handles zero page indexed with X addressing mode ($nn,X).
// Adds X register to zero page address with wraparound within zero page.
func (c *BaseCPU) AddrZeroPageX() (uint16, bool) {
	addr := uint16(c.Bus.Read(c.PC) + c.X)
	c.PC++
	return addr & 0x00FF, false // Wrap to stay in zero page
}

// AddrZeroPageY handles zero page indexed with Y addressing mode ($nn,Y).
// Adds Y register to zero page address with wraparound within zero page.
func (c *BaseCPU) AddrZeroPageY() (uint16, bool) {
	addr := uint16(c.Bus.Read(c.PC) + c.Y)
	c.PC++
	return addr & 0x00FF, false // Wrap to stay in zero page
}

// AddrAbsolute handles absolute addressing mode ($nnnn).
// Uses a full 16-bit address specified by two bytes (low, high).
func (c *BaseCPU) AddrAbsolute() (uint16, bool) {
	low := uint16(c.Bus.Read(c.PC))
	c.PC++
	high := uint16(c.Bus.Read(c.PC))
	c.PC++
	return (high << 8) | low, false
}

// AddrAbsoluteX handles absolute indexed with X addressing mode ($nnnn,X).
// Adds X register to the 16-bit address.
// Returns true if page boundary was crossed (high byte changed).
func (c *BaseCPU) AddrAbsoluteX() (uint16, bool) {
	low := uint16(c.Bus.Read(c.PC))
	c.PC++
	high := uint16(c.Bus.Read(c.PC))
	c.PC++
	addr := ((high << 8) | low) + uint16(c.X)
	return addr, (addr & 0xFF00) != (high << 8)
}

// AddrAbsoluteY handles absolute indexed with Y addressing mode ($nnnn,Y).
// Adds Y register to the 16-bit address.
// Returns true if page boundary was crossed (high byte changed).
func (c *BaseCPU) AddrAbsoluteY() (uint16, bool) {
	low := uint16(c.Bus.Read(c.PC))
	c.PC++
	high := uint16(c.Bus.Read(c.PC))
	c.PC++
	addr := ((high << 8) | low) + uint16(c.Y)
	return addr, (addr & 0xFF00) != (high << 8)
}

// AddrIndirect - indirect addressing mode (($nnnn))
// NMOS6502: Has page boundary bug where JMP ($10FF) reads high byte from $1000
// WDC65C02: Bug is fixed
func (c *BaseCPU) AddrIndirect() (uint16, bool) {
	low := uint16(c.Bus.Read(c.PC))
	c.PC++
	high := uint16(c.Bus.Read(c.PC))
	c.PC++
	ptr := (high << 8) | low

	// Check if this variant has the JMP indirect bug
	if c.Variant.HasJMPIndirectBug() && low == 0x00FF {
		// NMOS 6502: Read high byte from start of page (bug)
		return (uint16(c.Bus.Read(ptr&0xFF00)) << 8) | uint16(c.Bus.Read(ptr)), false
	}

	// Normal behavior (or WDC65C02 fixed behavior)
	return (uint16(c.Bus.Read(ptr+1)) << 8) | uint16(c.Bus.Read(ptr)), false
}

// AddrIndirectX handles indexed indirect addressing mode (($nn,X)).
// Adds X to the zero page address, then reads a 16-bit address from that location.
// Used pattern: LDA ($40,X) where X=0x05 reads address from $0045-$0046.
// Wraps within zero page.
func (c *BaseCPU) AddrIndirectX() (uint16, bool) {
	zeroPageAddr := uint16(c.Bus.Read(c.PC) + c.X)
	c.PC++
	low := uint16(c.Bus.Read(zeroPageAddr & 0x00FF))
	high := uint16(c.Bus.Read((zeroPageAddr + 1) & 0x00FF))
	return (high << 8) | low, false
}

// AddrIndirectY handles indirect indexed addressing mode (($nn),Y).
// Reads a 16-bit address from zero page, then adds Y to it.
// Used pattern: LDA ($40),Y where $0040-$0041 contains address, then add Y.
// Returns true if adding Y crossed a page boundary.
func (c *BaseCPU) AddrIndirectY() (uint16, bool) {
	zeroPageAddr := uint16(c.Bus.Read(c.PC))
	c.PC++
	low := uint16(c.Bus.Read(zeroPageAddr & 0x00FF))
	high := uint16(c.Bus.Read((zeroPageAddr + 1) & 0x00FF))
	addr := ((high << 8) | low) + uint16(c.Y)
	return addr, (addr & 0xFF00) != (high << 8)
}

// AddrRelative handles relative addressing mode (used by branch instructions).
// The operand is a signed 8-bit offset from the PC.
// Range: -128 to +127 bytes from the instruction following the branch.
// Returns true if the branch crosses a page boundary.
func (c *BaseCPU) AddrRelative() (uint16, bool) {
	offset := uint16(c.Bus.Read(c.PC))
	c.PC++
	if offset < 0x80 {
		return c.PC + offset, (c.PC & 0xFF00) != ((c.PC + offset) & 0xFF00)
	}
	return c.PC + offset - 0x100, (c.PC & 0xFF00) != ((c.PC + offset - 0x100) & 0xFF00)
}

// AddrZeroPageIndirect - zero page indirect addressing mode (($nn))
// This is a NEW addressing mode in WDC65C02
func (c *BaseCPU) AddrZeroPageIndirect() (uint16, bool) {
	zeroPageAddr := uint16(c.Bus.Read(c.PC))
	c.PC++
	low := uint16(c.Bus.Read(zeroPageAddr & 0x00FF))
	high := uint16(c.Bus.Read((zeroPageAddr + 1) & 0x00FF))
	return (high << 8) | low, false
}

// AddrAbsoluteIndexedIndirect - absolute indexed indirect addressing mode (($nnnn,X))
// This is a NEW addressing mode in WDC65C02
func (c *BaseCPU) AddrAbsoluteIndexedIndirect() (uint16, bool) {
	low := uint16(c.Bus.Read(c.PC))
	c.PC++
	high := uint16(c.Bus.Read(c.PC))
	c.PC++
	ptr := ((high << 8) | low) + uint16(c.X)

	// Read target address from pointer
	targetLow := uint16(c.Bus.Read(ptr))
	targetHigh := uint16(c.Bus.Read(ptr + 1))
	return (targetHigh << 8) | targetLow, false
}
