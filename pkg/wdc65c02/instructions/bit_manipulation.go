package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Bit manipulation instructions - NEW in WDC65C02
// These instructions test, set, and reset specific bits in zero page memory

// ========== Reset Memory Bit (RMB) ==========

// RMB0 resets (clears) bit 0 in memory.
func RMB0(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data&^byte(1<<0))
}

// RMB1 resets bit 1 in memory.
func RMB1(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data&^byte(1<<1))
}

// RMB2 resets bit 2 in memory.
func RMB2(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data&^byte(1<<2))
}

// RMB3 resets bit 3 in memory.
func RMB3(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data&^byte(1<<3))
}

// RMB4 resets bit 4 in memory.
func RMB4(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data&^byte(1<<4))
}

// RMB5 resets bit 5 in memory.
func RMB5(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data&^byte(1<<5))
}

// RMB6 resets bit 6 in memory.
func RMB6(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data&^byte(1<<6))
}

// RMB7 resets bit 7 in memory.
func RMB7(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data&^byte(1<<7))
}

// ========== Set Memory Bit (SMB) ==========

// SMB0 sets bit 0 in memory.
func SMB0(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data|byte(1<<0))
}

// SMB1 sets bit 1 in memory.
func SMB1(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data|byte(1<<1))
}

// SMB2 sets bit 2 in memory.
func SMB2(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data|byte(1<<2))
}

// SMB3 sets bit 3 in memory.
func SMB3(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data|byte(1<<3))
}

// SMB4 sets bit 4 in memory.
func SMB4(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data|byte(1<<4))
}

// SMB5 sets bit 5 in memory.
func SMB5(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data|byte(1<<5))
}

// SMB6 sets bit 6 in memory.
func SMB6(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data|byte(1<<6))
}

// SMB7 sets bit 7 in memory.
func SMB7(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Bus.Write(addr, data|byte(1<<7))
}

// ========== Branch on Bit Reset (BBR) ==========

// BBR0 branches if bit 0 is reset.
func BBR0(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 0, false)
}

// BBR1 branches if bit 1 is reset.
func BBR1(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 1, false)
}

// BBR2 branches if bit 2 is reset.
func BBR2(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 2, false)
}

// BBR3 branches if bit 3 is reset.
func BBR3(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 3, false)
}

// BBR4 branches if bit 4 is reset.
func BBR4(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 4, false)
}

// BBR5 branches if bit 5 is reset.
func BBR5(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 5, false)
}

// BBR6 branches if bit 6 is reset.
func BBR6(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 6, false)
}

// BBR7 branches if bit 7 is reset.
func BBR7(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 7, false)
}

// ========== Branch on Bit Set (BBS) ==========

// BBS0 branches if bit 0 is set.
func BBS0(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 0, true)
}

// BBS1 branches if bit 1 is set.
func BBS1(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 1, true)
}

// BBS2 branches if bit 2 is set.
func BBS2(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 2, true)
}

// BBS3 branches if bit 3 is set.
func BBS3(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 3, true)
}

// BBS4 branches if bit 4 is set.
func BBS4(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 4, true)
}

// BBS5 branches if bit 5 is set.
func BBS5(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 5, true)
}

// BBS6 branches if bit 6 is set.
func BBS6(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 6, true)
}

// BBS7 branches if bit 7 is set.
func BBS7(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	branchOnBit(c, data, 7, true)
}

// branchOnBit is a helper function for BBR/BBS instructions.
// Note: BBR/BBS have a special 3-byte format: opcode, zero-page address, relative offset
func branchOnBit(c *core.BaseCPU, data byte, bit uint, testSet bool) {
	bitValue := (data >> bit) & 1
	shouldBranch := (bitValue == 1) == testSet

	if shouldBranch {
		// Read the relative offset
		offset := c.Bus.Read(c.PC)
		c.PC++

		// Calculate branch target
		var target uint16
		if offset < 0x80 {
			target = c.PC + uint16(offset)
		} else {
			target = c.PC + uint16(offset) - 0x100
		}

		oldPC := c.PC
		c.PC = target
		c.Cycles++ // +1 cycle for taken branch
		if (oldPC & 0xFF00) != (c.PC & 0xFF00) {
			c.Cycles++ // +1 more cycle if page crossed
		}
	} else {
		// Skip the offset byte
		c.PC++
	}
}
