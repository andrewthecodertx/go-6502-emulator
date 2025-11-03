package instructions

import "github.com/andrewthecodertx/go-65c02-emulator/pkg/core"

// CPU control instructions - NEW in WDC65C02

// WAI waits for interrupt.
// NEW instruction in WDC65C02.
// Puts the CPU into a low-power state until an interrupt occurs.
func WAI(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	// WAI halts execution until an interrupt (NMI or IRQ) occurs.
	// The CPU remains in this state, consuming cycles but not executing instructions.
	// When an interrupt is triggered, the CPU wakes up and handles it normally.

	// For emulation purposes, we can implement this as:
	// - Set a "waiting" flag
	// - Keep consuming cycles until an interrupt is pending
	// - In a real implementation, this would be checked in the Step() loop

	// For now, we'll implement a simple version that just waits for the next interrupt
	// In a more sophisticated emulator, this would integrate with the Step() function
	// to check for interrupt flags before continuing execution.

	// Simple implementation: do nothing, let the normal interrupt handling wake us
	// The Step() function will handle interrupts on the next cycle
}

// STP stops the processor.
// NEW instruction in WDC65C02.
// Stops the processor completely until a hardware reset occurs.
func STP(c *core.BaseCPU, addr uint16, pageCrossed bool) {
	// STP halts the processor completely.
	// Only a hardware reset can restart execution.
	c.Halted = true
}
