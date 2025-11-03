# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a MOS 6502 CPU emulator written in Go. The project emulates the classic 8-bit MOS 6502 processor used in systems like the Commodore 64, Apple II, and NES.

## Building and Running

```bash
# Build the project
go build -o go-6502 ./cmd/go-6502

# Run the emulator
./go-6502

# Run all tests (when tests exist)
go test ./...

# Run tests for a specific package
go test ./pkg/mos6502
```

## Architecture

### Core Components

**CPU (`pkg/mos6502/cpu.go`)**
- The main CPU struct containing registers (A, X, Y, PC, SP, Status)
- Cycle-accurate execution with `Step()` and `Run()` methods
- Handles interrupts (NMI, IRQ, RESET)
- Stack pointer starts at 0xFD, status register at 0x34

**Bus Interface (`pkg/mos6502/bus.go`)**
- Simple interface for memory access: `Read(addr uint16) byte` and `Write(addr uint16, data byte)`
- Allows flexible memory implementations (RAM, ROM, memory-mapped I/O)

**Addressing Modes (`pkg/mos6502/addressing.go`)**
- All 13 6502 addressing modes implemented as CPU methods (e.g., `addrImmediate`, `addrZeroPageX`)
- Each returns the effective address and a boolean indicating page boundary crossing
- Includes hardware bug emulation (JMP indirect page boundary bug)

**Instruction System**

There are two parallel instruction systems in this codebase:

1. **Legacy System** (`pkg/mos6502/instructions/instructions.go`):
   - Map-based instruction dispatch (`instructionMap`)
   - Instructions defined inline with addressing mode and operation functions
   - Used by `cpu.Step()` for execution

2. **New System** (in development):
   - `pkg/core/opcode.go` defines abstract opcode structure
   - `pkg/mos6502/opcodes.go` loads opcodes from embedded JSON
   - `pkg/mos6502/interpreter.go` provides handler registration
   - `pkg/mos6502/instructions/*.go` contains modular instruction implementations

The new system aims to support JSON-driven opcode definitions with custom handlers.

### Status Register Flags

Flags are accessed using bit masks:
- `0x01`: Carry
- `0x02`: Zero
- `0x04`: Interrupt Disable
- `0x08`: Decimal Mode
- `0x10`: Break (B flag)
- `0x20`: Unused (always 1)
- `0x40`: Overflow
- `0x80`: Negative

### Instruction Implementation Pattern

Instructions follow a consistent pattern:
1. Read operand from memory using the computed address
2. Perform the operation
3. Update status flags (Z, N, C, V as appropriate)
4. Add extra cycle if page boundary crossed (for applicable instructions)

Branch instructions add 1 cycle if taken, 2 if crossing a page boundary.

## Important Implementation Details

- **Cycle Accuracy**: The emulator tracks cycles with the `cycles` field. Instructions consume cycles, and `Step()` decrements the counter each call.
- **Page Crossing**: Many instructions take an extra cycle when crossing a page boundary (when high byte of address changes). This is computed and returned by addressing mode functions.
- **Stack**: Located at 0x0100-0x01FF, grows downward. Push decrements SP, pull increments SP.
- **Reset Vector**: On reset, PC is loaded from 0xFFFC-0xFFFD
- **Interrupt Vectors**: BRK/IRQ at 0xFFFE-0xFFFF, NMI at 0xFFFA-0xFFFB

## Module Path

The module is `github.com/andrewthecodertx/go-65c02-emulator` - keep this in mind when adding imports.

## Development Notes

- No tests currently exist in the codebase
- The main.go currently only prints "go-6502" - the emulator needs a Bus implementation to run
- When adding new instructions, add to both the legacy instruction map and the new handler system for consistency
- Hardware quirks (like the JMP indirect bug) are intentionally emulated for accuracy
