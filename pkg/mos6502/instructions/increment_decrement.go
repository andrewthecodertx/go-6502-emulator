package instructions

import "github.com/andrewthecodertx/go-65c02-emulator/pkg/core"

// Increment/Decrement instructions for NMOS 6502

// INC increments a value in memory.
func INC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr) + 1
	c.Bus.Write(addr, data)
	c.SetZN(data)
}

// DEC decrements a value in memory.
func DEC(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	data := c.Bus.Read(addr) - 1
	c.Bus.Write(addr, data)
	c.SetZN(data)
}

// INX increments the X register.
func INX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.X++
	c.SetZN(c.X)
}

// DEX decrements the X register.
func DEX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.X--
	c.SetZN(c.X)
}

// INY increments the Y register.
func INY(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Y++
	c.SetZN(c.Y)
}

// DEY decrements the Y register.
func DEY(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Y--
	c.SetZN(c.Y)
}
