package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Jump/Branch instructions for WDC65C02

// JMP jumps to a new location.
func JMP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.PC = addr
}

// JSR jumps to a subroutine.
func JSR(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	returnAddr := c.PC - 1
	c.Push(byte(returnAddr >> 8))
	c.Push(byte(returnAddr))
	c.PC = addr
}

// RTS returns from a subroutine.
func RTS(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	low := uint16(c.Pull())
	high := uint16(c.Pull())
	c.PC = ((high << 8) | low) + 1
}

// RTI returns from an interrupt.
func RTI(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Status = c.Pull()
	c.SetFlag(core.FlagUnused, true)  // Unused flag always set
	c.SetFlag(core.FlagBreak, false) // B flag not actually stored
	low := uint16(c.Pull())
	high := uint16(c.Pull())
	c.PC = (high << 8) | low
}

// BRK executes a software interrupt.
func BRK(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.PC++
	c.Push(byte(c.PC >> 8))
	c.Push(byte(c.PC))
	c.Push(c.Status | core.FlagBreak | core.FlagUnused)
	c.SetFlag(core.FlagInterruptDisable, true)

	// WDC65C02: Clear decimal mode on interrupt
	c.SetFlag(core.FlagDecimal, false)

	low := uint16(c.Bus.Read(0xFFFE))
	high := uint16(c.Bus.Read(0xFFFF))
	c.PC = (high << 8) | low
}

// Branch helper function
func branch(c *core.BaseCPU, condition bool, addr uint16) {
	if condition {
		oldPC := c.PC
		c.PC = addr
		c.Cycles++ // +1 cycle if branch taken
		if (oldPC & 0xFF00) != (c.PC & 0xFF00) {
			c.Cycles++ // +1 more cycle if page crossed
		}
	}
}

// BCC branches if carry clear.
func BCC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	branch(c, !c.GetFlag(core.FlagCarry), addr)
}

// BCS branches if carry set.
func BCS(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	branch(c, c.GetFlag(core.FlagCarry), addr)
}

// BEQ branches if equal (zero set).
func BEQ(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	branch(c, c.GetFlag(core.FlagZero), addr)
}

// BNE branches if not equal (zero clear).
func BNE(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	branch(c, !c.GetFlag(core.FlagZero), addr)
}

// BMI branches if minus (negative set).
func BMI(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	branch(c, c.GetFlag(core.FlagNegative), addr)
}

// BPL branches if plus (negative clear).
func BPL(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	branch(c, !c.GetFlag(core.FlagNegative), addr)
}

// BVC branches if overflow clear.
func BVC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	branch(c, !c.GetFlag(core.FlagOverflow), addr)
}

// BVS branches if overflow set.
func BVS(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	branch(c, c.GetFlag(core.FlagOverflow), addr)
}

// BRA branches always (unconditional relative branch).
// NEW instruction in WDC65C02 - always takes the branch.
func BRA(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	oldPC := c.PC
	c.PC = addr
	c.Cycles++ // +1 cycle for taken branch
	if (oldPC & 0xFF00) != (c.PC & 0xFF00) {
		c.Cycles++ // +1 more cycle if page crossed
	}
}
