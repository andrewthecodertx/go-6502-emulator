package instructions

import "github.com/andrewthecodertx/go-6502-emulator/pkg/core"

// Load/Store instructions for NMOS 6502

// LDA loads a byte from memory into the accumulator.
func LDA(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.A = data
	c.SetZN(c.A)

	if pageCrossed {
		c.Cycles++
	}
}

// LDX loads a byte from memory into the X register.
func LDX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.X = data
	c.SetZN(c.X)

	if pageCrossed {
		c.Cycles++
	}
}

// LDY loads a byte from memory into the Y register.
func LDY(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr)
	c.Y = data
	c.SetZN(c.Y)

	if pageCrossed {
		c.Cycles++
	}
}

// STA stores the accumulator in memory.
func STA(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Bus.Write(addr, c.A)
}

// STX stores the X register in memory.
func STX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Bus.Write(addr, c.X)
}

// STY stores the Y register in memory.
func STY(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Bus.Write(addr, c.Y)
}
