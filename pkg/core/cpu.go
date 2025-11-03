// Package core provides shared components for 6502-family CPU emulators.
//
// This package implements the common architecture shared by NMOS 6502 and WDC65C02
// processors. It provides:
//   - BaseCPU: Core CPU state and operations
//   - Bus: Memory interface abstraction
//   - Variant: CPU variant identification and behavior
//   - Common addressing modes
//   - Flag manipulation and stack operations
package core

// BaseCPU represents the core state and operations shared by all 6502-family processors.
// This includes registers, flags, memory interface, and common operations.
//
// Register layout:
//   - PC (Program Counter): 16-bit address of the next instruction
//   - SP (Stack Pointer): 8-bit offset into page 0x01 (0x0100-0x01FF)
//   - A (Accumulator): 8-bit general purpose register
//   - X, Y (Index registers): 8-bit registers for indexed addressing
//   - Status: 8-bit processor status register (flags)
type BaseCPU struct {
	PC     uint16 // Program Counter: current instruction address
	SP     byte   // Stack Pointer: offset into 0x0100-0x01FF (grows downward)
	A      byte   // Accumulator: primary 8-bit register
	X      byte   // X Index Register: 8-bit index/counter
	Y      byte   // Y Index Register: 8-bit index/counter
	Status byte   // Processor Status Register: NVUBDIZC flags

	Bus Bus // Memory and I/O interface

	Cycles byte // Remaining cycles for current instruction
	Halted bool // CPU halted (e.g., STP instruction on WDC65C02)

	// Interrupt pending flags
	NMIPending   bool // Non-Maskable Interrupt pending
	IRQPending   bool // Interrupt Request pending
	ResetPending bool // Reset pending

	Variant Variant // CPU variant (NMOS vs WDC65C02)
}

// Processor Status Register flags (8 bits: NV-BDIZC)
// Bit 5 is always 1 (unused flag)
const (
	FlagCarry            = 0x01 // C: Carry flag (bit 0)
	FlagZero             = 0x02 // Z: Zero flag (bit 1)
	FlagInterruptDisable = 0x04 // I: Interrupt disable flag (bit 2)
	FlagDecimal          = 0x08 // D: Decimal mode flag (bit 3)
	FlagBreak            = 0x10 // B: Break flag (bit 4, set by BRK/PHP)
	FlagUnused           = 0x20 // -: Unused, always 1 (bit 5)
	FlagOverflow         = 0x40 // V: Overflow flag (bit 6)
	FlagNegative         = 0x80 // N: Negative/sign flag (bit 7)
)

// NewBaseCPU creates and initializes a new BaseCPU with the specified bus and variant.
// Initial state:
//   - SP = 0xFD (stack pointer starts 3 bytes below top)
//   - Status = 0x34 (Interrupt Disable and Unused flags set)
//   - All other registers = 0
func NewBaseCPU(bus Bus, variant Variant) *BaseCPU {
	return &BaseCPU{
		SP:      0xFD,
		Status:  0x34, // I flag set, unused bit set
		Bus:     bus,
		Variant: variant,
	}
}

// GetFlag returns true if the specified flag bit is set in the Status register.
func (c *BaseCPU) GetFlag(flag byte) bool {
	return c.Status&flag != 0
}

// SetFlag sets or clears the specified flag bit in the Status register.
func (c *BaseCPU) SetFlag(flag byte, value bool) {
	if value {
		c.Status |= flag
	} else {
		c.Status &= ^flag
	}
}

// Push writes a byte to the stack and decrements the stack pointer.
// The stack is located at 0x0100-0x01FF and grows downward.
func (c *BaseCPU) Push(data byte) {
	c.Bus.Write(0x0100+uint16(c.SP), data)
	c.SP--
}

// Pull increments the stack pointer and reads a byte from the stack.
func (c *BaseCPU) Pull() byte {
	c.SP++
	return c.Bus.Read(0x0100 + uint16(c.SP))
}

// SetZN sets the Zero and Negative flags based on the given value.
// This is a common operation after most arithmetic and data movement instructions.
//   - Zero flag is set if value == 0
//   - Negative flag is set if bit 7 of value is set
func (c *BaseCPU) SetZN(value byte) {
	c.SetFlag(FlagZero, value == 0)
	c.SetFlag(FlagNegative, value&0x80 != 0)
}

// GetCycles returns the number of cycles remaining for the current instruction.
// This is primarily used for testing and debugging.
func (c *BaseCPU) GetCycles() byte {
	return c.Cycles
}

// Reset initializes the CPU to its power-on state.
// Registers are cleared, status is set to 0x34, and the PC is loaded from
// the reset vector at 0xFFFC-0xFFFD.
// Cycle count is set based on variant (6 for NMOS, 7 for WDC65C02).
func (c *BaseCPU) Reset() {
	c.A = 0
	c.X = 0
	c.Y = 0
	c.SP = 0xFD
	c.Status = 0x34 | FlagUnused // Set I flag and unused bit

	low := uint16(c.Bus.Read(0xFFFC))
	high := uint16(c.Bus.Read(0xFFFD))
	c.PC = (high << 8) | low

	c.Cycles = c.Variant.ResetCycles()
	c.Halted = false
}

// HandleNMI processes a Non-Maskable Interrupt.
// The NMI cannot be disabled and takes priority over IRQ.
// Saves PC and Status to stack, sets Interrupt Disable flag,
// and loads PC from the NMI vector at 0xFFFA-0xFFFB.
// Takes 7 cycles.
func (c *BaseCPU) HandleNMI() {
	c.Push(byte(c.PC >> 8))
	c.Push(byte(c.PC))
	c.Push(c.Status)
	c.SetFlag(FlagInterruptDisable, true)

	if c.Variant.ClearsDecimalOnInterrupt() {
		c.SetFlag(FlagDecimal, false)
	}

	low := uint16(c.Bus.Read(0xFFFA))
	high := uint16(c.Bus.Read(0xFFFB))
	c.PC = (high << 8) | low
	c.Cycles = 7
	c.NMIPending = false
}

// HandleIRQ processes an Interrupt Request.
// Only executed if the Interrupt Disable flag is clear.
// Saves PC and Status to stack, sets Interrupt Disable flag,
// and loads PC from the IRQ vector at 0xFFFE-0xFFFF.
// Takes 7 cycles.
func (c *BaseCPU) HandleIRQ() {
	c.Push(byte(c.PC >> 8))
	c.Push(byte(c.PC))
	c.Push(c.Status)
	c.SetFlag(FlagInterruptDisable, true)

	if c.Variant.ClearsDecimalOnInterrupt() {
		c.SetFlag(FlagDecimal, false)
	}

	low := uint16(c.Bus.Read(0xFFFE))
	high := uint16(c.Bus.Read(0xFFFF))
	c.PC = (high << 8) | low
	c.Cycles = 7
	c.IRQPending = false
}

// HandleReset processes a pending reset request.
// Calls Reset() and clears the ResetPending flag.
func (c *BaseCPU) HandleReset() {
	c.Reset()
	c.ResetPending = false
}
