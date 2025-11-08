package main

import (
	"fmt"

	"github.com/andrewthecodertx/go-6502-emulator/pkg/mos6502"
)

type SimpleRAM struct {
	memory [0x10000]byte
}

func (r *SimpleRAM) Read(addr uint16) byte {
	return r.memory[addr]
}

func (r *SimpleRAM) Write(addr uint16, data byte) {
	r.memory[addr] = data
}

func main() {
	bus := &SimpleRAM{}
	cpu := mos6502.NewCPU(bus)

	// Set reset vector
	bus.memory[0xFFFC] = 0x00
	bus.memory[0xFFFD] = 0x80

	// Set X register to cause wraparound
	cpu.X = 0xFF

	// Put a value at $0010 (where $FF + $11 wraps to)
	bus.memory[0x0010] = 0x42

	// LDA $11,X with X=$FF should read from $0010, not $0110
	bus.memory[0x8000] = 0xB5
	bus.memory[0x8001] = 0x11

	cpu.Reset()

	fmt.Printf("After reset: PC=0x%04X, X=0x%02X, A=0x%02X\n", cpu.PC, cpu.X, cpu.A)

	// Consume reset cycles
	for cpu.PC == 0x8000 && cpu.GetCycles() > 0 {
		cpu.Step()
	}

	fmt.Printf("After reset cycles: PC=0x%04X, cycles=%d\n", cpu.PC, cpu.GetCycles())

	// Execute LDA $11,X
	fmt.Printf("About to execute opcode 0x%02X at 0x%04X\n", bus.memory[cpu.PC], cpu.PC)
	fmt.Printf("X=0x%02X, $11+X=0x%02X\n", cpu.X, 0x11+cpu.X)
	fmt.Printf("Value at $0010: 0x%02X\n", bus.memory[0x0010])
	fmt.Printf("Value at $0110: 0x%02X\n", bus.memory[0x0110])

	cpu.Step()
	fmt.Printf("After fetch: PC=0x%04X, cycles=%d, A=0x%02X\n", cpu.PC, cpu.GetCycles(), cpu.A)

	for cpu.GetCycles() > 0 {
		cpu.Step()
	}

	fmt.Printf("After execution: PC=0x%04X, A=0x%02X\n", cpu.PC, cpu.A)
}
