// Package core provides shared components for 6502-family CPU emulators.
//
// This package contains the common architecture shared by NMOS 6502 and WDC65C02
// processors, including:
//   - CPU variant identification
//   - Common bus interface
//   - Shared addressing modes
//   - Common register and flag operations
//
// Variant-specific implementations are in pkg/mos6502 (NMOS 6502) and
// pkg/wdc65c02 (WDC65C02).
package core

type Variant int

const (
	// VariantNMOS represents the original NMOS 6502 processor
	// Features:
	//   - 56 documented instructions
	//   - JMP indirect page boundary bug
	//   - 6 cycle reset
	//   - Decimal mode does NOT clear on interrupts
	//   - Undefined behavior for illegal opcodes
	VariantNMOS Variant = iota

	// VariantWDC65C02 represents the Western Design Center 65C02 processor
	// Features:
	//   - 70 instructions (adds 27 new instructions)
	//   - JMP indirect bug is FIXED
	//   - 7 cycle reset
	//   - Decimal mode CLEARS on interrupts
	//   - All illegal opcodes become NOPs
	//   - New addressing modes: Zero Page Indirect, Absolute Indexed Indirect
	VariantWDC65C02
)

func (v Variant) String() string {
	switch v {
	case VariantNMOS:
		return "NMOS6502"
	case VariantWDC65C02:
		return "WDC65C02"
	default:
		return "Unknown"
	}
}

func (v Variant) ResetCycles() byte {
	switch v {
	case VariantNMOS:
		return 6 // Per NMOS 6502 datasheet page 8
	case VariantWDC65C02:
		return 7 // Per W65C02S datasheet page 10
	default:
		return 6
	}
}

// HasJMPIndirectBug returns true if this variant has the JMP indirect page boundary bug
func (v Variant) HasJMPIndirectBug() bool {
	return v == VariantNMOS
}

func (v Variant) ClearsDecimalOnInterrupt() bool {
	return v == VariantWDC65C02
}
