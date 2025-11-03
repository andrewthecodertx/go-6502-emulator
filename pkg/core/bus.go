package core

// Bus defines the memory and I/O interface for the 6502 CPU.
//
// The 6502 has a 16-bit address bus (0x0000-0xFFFF) and can access up to 64KB
// of address space. This interface abstracts memory access, allowing custom
// implementations for:
//   - Simple RAM/ROM
//   - Memory-mapped I/O devices
//   - Bank switching / memory mappers
//   - Debug/trace capabilities
//
// Example implementation:
//
//	type SimpleRAM struct {
//	    memory [0x10000]byte
//	}
//
//	func (r *SimpleRAM) Read(addr uint16) byte {
//	    return r.memory[addr]
//	}
//
//	func (r *SimpleRAM) Write(addr uint16, data byte) {
//	    r.memory[addr] = data
//	}
type Bus interface {
	// Read returns the byte at the specified address.
	// This may trigger side effects in memory-mapped I/O devices.
	Read(addr uint16) byte

	// Write stores a byte at the specified address.
	// This may trigger side effects in memory-mapped I/O devices.
	Write(addr uint16, data byte)
}
