package instructions

import "github.com/andrewthecodertx/go-65c02-emulator/pkg/core"

// Jump/Branch instructions for NMOS 6502

// JMP jumps to a new location.
func JMP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.PC = addr
}

// JSR jumps to a subroutine.
func JSR(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.PC--
	c.Push(byte(c.PC >> 8))
	c.Push(byte(c.PC))
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
	c.SetFlag(core.FlagBreak, false) // B flag not stored
	c.SetFlag(core.FlagUnused, true)  // Unused flag always set
	low := uint16(c.Pull())
	high := uint16(c.Pull())
	c.PC = (high << 8) | low
}

// BRK forces a break.
func BRK(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.PC++
	c.SetFlag(core.FlagInterruptDisable, true)
	c.Push(byte(c.PC >> 8))
	c.Push(byte(c.PC))
	c.SetFlag(core.FlagBreak, true)
	c.Push(c.Status)
	c.SetFlag(core.FlagBreak, false)
	low := uint16(c.Bus.Read(0xFFFE))
	high := uint16(c.Bus.Read(0xFFFF))
	c.PC = (high << 8) | low
}

// branch takes a branch if the condition is met.
func branch(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Cycles++
	if pageCrossed {
		c.Cycles++
	}
	c.PC = addr
}

// BCC branches if carry is clear.
func BCC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	if !c.GetFlag(core.FlagCarry) {
		branch(c, addr, pageCrossed)
	}
}

// BCS branches if carry is set.
func BCS(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	if c.GetFlag(core.FlagCarry) {
		branch(c, addr, pageCrossed)
	}
}

// BEQ branches if zero is set.
func BEQ(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	if c.GetFlag(core.FlagZero) {
		branch(c, addr, pageCrossed)
	}
}

// BMI branches if negative is set.
func BMI(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	if c.GetFlag(core.FlagNegative) {
		branch(c, addr, pageCrossed)
	}
}

// BNE branches if zero is clear.
func BNE(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	if !c.GetFlag(core.FlagZero) {
		branch(c, addr, pageCrossed)
	}
}

// BPL branches if negative is clear.
func BPL(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	if !c.GetFlag(core.FlagNegative) {
		branch(c, addr, pageCrossed)
	}
}

// BVC branches if overflow is clear.
func BVC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	if !c.GetFlag(core.FlagOverflow) {
		branch(c, addr, pageCrossed)
	}
}

// BVS branches if overflow is set.
func BVS(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	if c.GetFlag(core.FlagOverflow) {
		branch(c, addr, pageCrossed)
	}
}
