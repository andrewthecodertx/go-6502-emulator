package instructions

import "github.com/andrewthecodertx/go-65c02-emulator/pkg/core"

// Transfer instructions for NMOS 6502

// TAX transfers the accumulator to the X register.
func TAX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.X = c.A
	c.SetZN(c.X)
}

// TAY transfers the accumulator to the Y register.
func TAY(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.Y = c.A
	c.SetZN(c.Y)
}

// TXA transfers the X register to the accumulator.
func TXA(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.A = c.X
	c.SetZN(c.A)
}

// TYA transfers the Y register to the accumulator.
func TYA(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.A = c.Y
	c.SetZN(c.A)
}

// TSX transfers the stack pointer to the X register.
func TSX(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.X = c.SP
	c.SetZN(c.X)
}

// TXS transfers the X register to the stack pointer.
func TXS(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	c.SP = c.X
}
