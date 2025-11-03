package mos6502

// This file contains the instruction map for NMOS 6502.
// Instruction implementations are in pkg/mos6502/instructions/.

// Instruction represents a 6502 instruction.
type Instruction struct {
	Name      string
	AddrMode  func(c *CPU) (uint16, bool)
	Operation func(c *CPU, addr uint16, pageCrossed bool)
	Cycles    byte
}

// instructionMap maps opcodes to their instruction implementations.
var instructionMap = map[byte]Instruction{
	0x00: {"BRK", nil, (*CPU).brk, 7},
	0x01: {"ORA", (*CPU).addrIndirectX, (*CPU).ora, 6},
	0x05: {"ORA", (*CPU).addrZeroPage, (*CPU).ora, 3},
	0x06: {"ASL", (*CPU).addrZeroPage, (*CPU).asl, 5},
	0x08: {"PHP", nil, (*CPU).php, 3},
	0x09: {"ORA", (*CPU).addrImmediate, (*CPU).ora, 2},
	0x0A: {"ASL", nil, (*CPU).aslAccumulator, 2},
	0x0D: {"ORA", (*CPU).addrAbsolute, (*CPU).ora, 4},
	0x0E: {"ASL", (*CPU).addrAbsolute, (*CPU).asl, 6},
	0x10: {"BPL", (*CPU).addrRelative, (*CPU).bpl, 2}, // +1 if branch succeeds, +2 if to new page
	0x11: {"ORA", (*CPU).addrIndirectY, (*CPU).ora, 5}, // +1 if page crossed
	0x15: {"ORA", (*CPU).addrZeroPageX, (*CPU).ora, 4},
	0x16: {"ASL", (*CPU).addrZeroPageX, (*CPU).asl, 6},
	0x18: {"CLC", nil, (*CPU).clc, 2},
	0x19: {"ORA", (*CPU).addrAbsoluteY, (*CPU).ora, 4}, // +1 if page crossed
	0x1D: {"ORA", (*CPU).addrAbsoluteX, (*CPU).ora, 4}, // +1 if page crossed
	0x1E: {"ASL", (*CPU).addrAbsoluteX, (*CPU).asl, 7},
	0x20: {"JSR", (*CPU).addrAbsolute, (*CPU).jsr, 6},
	0x21: {"AND", (*CPU).addrIndirectX, (*CPU).and, 6},
	0x24: {"BIT", (*CPU).addrZeroPage, (*CPU).bit, 3},
	0x25: {"AND", (*CPU).addrZeroPage, (*CPU).and, 3},
	0x26: {"ROL", (*CPU).addrZeroPage, (*CPU).rol, 5},
	0x28: {"PLP", nil, (*CPU).plp, 4},
	0x29: {"AND", (*CPU).addrImmediate, (*CPU).and, 2},
	0x2A: {"ROL", nil, (*CPU).rolAccumulator, 2},
	0x2C: {"BIT", (*CPU).addrAbsolute, (*CPU).bit, 4},
	0x2D: {"AND", (*CPU).addrAbsolute, (*CPU).and, 4},
	0x2E: {"ROL", (*CPU).addrAbsolute, (*CPU).rol, 7},
	0x30: {"BMI", (*CPU).addrRelative, (*CPU).bmi, 2}, // +1 if branch succeeds, +2 if to new page
	0x31: {"AND", (*CPU).addrIndirectY, (*CPU).and, 5}, // +1 if page crossed
	0x35: {"AND", (*CPU).addrZeroPageX, (*CPU).and, 4},
	0x36: {"ROL", (*CPU).addrZeroPageX, (*CPU).rol, 6},
	0x38: {"SEC", nil, (*CPU).sec, 2},
	0x39: {"AND", (*CPU).addrAbsoluteY, (*CPU).and, 4}, // +1 if page crossed
	0x3D: {"AND", (*CPU).addrAbsoluteX, (*CPU).and, 4}, // +1 if page crossed
	0x3E: {"ROL", (*CPU).addrAbsoluteX, (*CPU).rol, 7},
	0x40: {"RTI", nil, (*CPU).rti, 6},
	0x41: {"EOR", (*CPU).addrIndirectX, (*CPU).eor, 6},
	0x45: {"EOR", (*CPU).addrZeroPage, (*CPU).eor, 3},
	0x46: {"LSR", (*CPU).addrZeroPage, (*CPU).lsr, 5},
	0x48: {"PHA", nil, (*CPU).pha, 3},
	0x49: {"EOR", (*CPU).addrImmediate, (*CPU).eor, 2},
	0x4A: {"LSR", nil, (*CPU).lsrAccumulator, 2},
	0x4C: {"JMP", (*CPU).addrAbsolute, (*CPU).jmp, 3},
	0x4D: {"EOR", (*CPU).addrAbsolute, (*CPU).eor, 4},
	0x4E: {"LSR", (*CPU).addrAbsolute, (*CPU).lsr, 7},
	0x50: {"BVC", (*CPU).addrRelative, (*CPU).bvc, 2}, // +1 if branch succeeds, +2 if to new page
	0x51: {"EOR", (*CPU).addrIndirectY, (*CPU).eor, 5}, // +1 if page crossed
	0x55: {"EOR", (*CPU).addrZeroPageX, (*CPU).eor, 4},
	0x56: {"LSR", (*CPU).addrZeroPageX, (*CPU).lsr, 6},
	0x58: {"CLI", nil, (*CPU).cli, 2},
	0x59: {"EOR", (*CPU).addrAbsoluteY, (*CPU).eor, 4}, // +1 if page crossed
	0x5D: {"EOR", (*CPU).addrAbsoluteX, (*CPU).eor, 4}, // +1 if page crossed
	0x5E: {"LSR", (*CPU).addrAbsoluteX, (*CPU).lsr, 7},
	0x60: {"RTS", nil, (*CPU).rts, 6},
	0x61: {"ADC", (*CPU).addrIndirectX, (*CPU).adc, 6},
	0x65: {"ADC", (*CPU).addrZeroPage, (*CPU).adc, 3},
	0x66: {"ROR", (*CPU).addrZeroPage, (*CPU).ror, 5},
	0x68: {"PLA", nil, (*CPU).pla, 4},
	0x69: {"ADC", (*CPU).addrImmediate, (*CPU).adc, 2},
	0x6A: {"ROR", nil, (*CPU).rorAccumulator, 2},
	0x6C: {"JMP", (*CPU).addrIndirect, (*CPU).jmp, 5},
	0x6D: {"ADC", (*CPU).addrAbsolute, (*CPU).adc, 4},
	0x6E: {"ROR", (*CPU).addrAbsolute, (*CPU).ror, 7},
	0x70: {"BVS", (*CPU).addrRelative, (*CPU).bvs, 2}, // +1 if branch succeeds, +2 if to new page
	0x71: {"ADC", (*CPU).addrIndirectY, (*CPU).adc, 5}, // +1 if page crossed
	0x75: {"ADC", (*CPU).addrZeroPageX, (*CPU).adc, 4},
	0x76: {"ROR", (*CPU).addrZeroPageX, (*CPU).ror, 6},
	0x78: {"SEI", nil, (*CPU).sei, 2},
	0x79: {"ADC", (*CPU).addrAbsoluteY, (*CPU).adc, 4}, // +1 if page crossed
	0x7D: {"ADC", (*CPU).addrAbsoluteX, (*CPU).adc, 4}, // +1 if page crossed
	0x7E: {"ROR", (*CPU).addrAbsoluteX, (*CPU).ror, 7},
	0x81: {"STA", (*CPU).addrIndirectX, (*CPU).sta, 6},
	0x84: {"STY", (*CPU).addrZeroPage, (*CPU).sty, 3},
	0x85: {"STA", (*CPU).addrZeroPage, (*CPU).sta, 3},
	0x86: {"STX", (*CPU).addrZeroPage, (*CPU).stx, 3},
	0x88: {"DEY", nil, (*CPU).dey, 2},
	0x8A: {"TXA", nil, (*CPU).txa, 2},
	0x8C: {"STY", (*CPU).addrAbsolute, (*CPU).sty, 4},
	0x8D: {"STA", (*CPU).addrAbsolute, (*CPU).sta, 4},
	0x8E: {"STX", (*CPU).addrAbsolute, (*CPU).stx, 4},
	0x90: {"BCC", (*CPU).addrRelative, (*CPU).bcc, 2}, // +1 if branch succeeds, +2 if to new page
	0x91: {"STA", (*CPU).addrIndirectY, (*CPU).sta, 6},
	0x94: {"STY", (*CPU).addrZeroPageX, (*CPU).sty, 3},
	0x95: {"STA", (*CPU).addrZeroPageX, (*CPU).sta, 3},
	0x96: {"STX", (*CPU).addrZeroPageY, (*CPU).stx, 4},
	0x98: {"TYA", nil, (*CPU).tya, 2},
	0x99: {"STA", (*CPU).addrAbsoluteY, (*CPU).sta, 5},
	0x9A: {"TXS", nil, (*CPU).txs, 2},
	0x9D: {"STA", (*CPU).addrAbsoluteX, (*CPU).sta, 5},
	0xA0: {"LDY", (*CPU).addrImmediate, (*CPU).ldy, 2},
	0xA1: {"LDA", (*CPU).addrIndirectX, (*CPU).lda, 6},
	0xA2: {"LDX", (*CPU).addrImmediate, (*CPU).ldx, 2},
	0xA4: {"LDY", (*CPU).addrZeroPage, (*CPU).ldy, 3},
	0xA5: {"LDA", (*CPU).addrZeroPage, (*CPU).lda, 3},
	0xA6: {"LDX", (*CPU).addrZeroPage, (*CPU).ldx, 3},
	0xA8: {"TAY", nil, (*CPU).tay, 2},
	0xA9: {"LDA", (*CPU).addrImmediate, (*CPU).lda, 2},
	0xAA: {"TAX", nil, (*CPU).tax, 2},
	0xAC: {"LDY", (*CPU).addrAbsolute, (*CPU).ldy, 4},
	0xAD: {"LDA", (*CPU).addrAbsolute, (*CPU).lda, 4},
	0xAE: {"LDX", (*CPU).addrAbsolute, (*CPU).ldx, 4},
	0xB0: {"BCS", (*CPU).addrRelative, (*CPU).bcs, 2}, // +1 if branch succeeds, +2 if to new page
	0xB1: {"LDA", (*CPU).addrIndirectY, (*CPU).lda, 5}, // +1 if page crossed
	0xB4: {"LDY", (*CPU).addrZeroPageX, (*CPU).ldy, 3},
	0xB5: {"LDA", (*CPU).addrZeroPageX, (*CPU).lda, 3},
	0xB6: {"LDX", (*CPU).addrZeroPageY, (*CPU).ldx, 3},
	0xB8: {"CLV", nil, (*CPU).clv, 2},
	0xB9: {"LDA", (*CPU).addrAbsoluteY, (*CPU).lda, 4}, // +1 if page crossed
	0xBA: {"TSX", nil, (*CPU).tsx, 2},
	0xBC: {"LDY", (*CPU).addrAbsoluteX, (*CPU).ldy, 4}, // +1 if page crossed
	0xBD: {"LDA", (*CPU).addrAbsoluteX, (*CPU).lda, 4}, // +1 if page crossed
	0xBE: {"LDX", (*CPU).addrAbsoluteY, (*CPU).ldx, 4}, // +1 if page crossed
	0xC0: {"CPY", (*CPU).addrImmediate, (*CPU).cpy, 2},
	0xC1: {"CMP", (*CPU).addrIndirectX, (*CPU).cmp, 6},
	0xC4: {"CPY", (*CPU).addrZeroPage, (*CPU).cpy, 3},
	0xC5: {"CMP", (*CPU).addrZeroPage, (*CPU).cmp, 3},
	0xC6: {"DEC", (*CPU).addrZeroPage, (*CPU).dec, 5},
	0xC8: {"INY", nil, (*CPU).iny, 2},
	0xC9: {"CMP", (*CPU).addrImmediate, (*CPU).cmp, 2},
	0xCA: {"DEX", nil, (*CPU).dex, 2},
	0xCC: {"CPY", (*CPU).addrAbsolute, (*CPU).cpy, 4},
	0xCD: {"CMP", (*CPU).addrAbsolute, (*CPU).cmp, 4},
	0xCE: {"DEC", (*CPU).addrAbsolute, (*CPU).dec, 6},
	0xD0: {"BNE", (*CPU).addrRelative, (*CPU).bne, 2}, // +1 if branch succeeds, +2 if to new page
	0xD1: {"CMP", (*CPU).addrIndirectY, (*CPU).cmp, 5}, // +1 if page crossed
	0xD5: {"CMP", (*CPU).addrZeroPageX, (*CPU).cmp, 3},
	0xD6: {"DEC", (*CPU).addrZeroPageX, (*CPU).dec, 5},
	0xD8: {"CLD", nil, (*CPU).cld, 2},
	0xD9: {"CMP", (*CPU).addrAbsoluteY, (*CPU).cmp, 4}, // +1 if page crossed
	0xDD: {"CMP", (*CPU).addrAbsoluteX, (*CPU).cmp, 4}, // +1 if page crossed
	0xDE: {"DEC", (*CPU).addrAbsoluteX, (*CPU).dec, 7},
	0xE0: {"CPX", (*CPU).addrImmediate, (*CPU).cpx, 2},
	0xE1: {"SBC", (*CPU).addrIndirectX, (*CPU).sbc, 6},
	0xE4: {"CPX", (*CPU).addrZeroPage, (*CPU).cpx, 3},
	0xE5: {"SBC", (*CPU).addrZeroPage, (*CPU).sbc, 3},
	0xE6: {"INC", (*CPU).addrZeroPage, (*CPU).inc, 5},
	0xE8: {"INX", nil, (*CPU).inx, 2},
	0xE9: {"SBC", (*CPU).addrImmediate, (*CPU).sbc, 2},
	0xEA: {"NOP", nil, (*CPU).nop, 2},
	0xEC: {"CPX", (*CPU).addrAbsolute, (*CPU).cpx, 4},
	0xED: {"SBC", (*CPU).addrAbsolute, (*CPU).sbc, 4},
	0xEE: {"INC", (*CPU).addrAbsolute, (*CPU).inc, 7},
	0xF0: {"BEQ", (*CPU).addrRelative, (*CPU).beq, 2}, // +1 if branch succeeds, +2 if to new page
	0xF1: {"SBC", (*CPU).addrIndirectY, (*CPU).sbc, 5}, // +1 if page crossed
	0xF5: {"SBC", (*CPU).addrZeroPageX, (*CPU).sbc, 3},
	0xF6: {"INC", (*CPU).addrZeroPageX, (*CPU).inc, 5},
	0xF8: {"SED", nil, (*CPU).sed, 2},
	0xF9: {"SBC", (*CPU).addrAbsoluteY, (*CPU).sbc, 4}, // +1 if page crossed
	0xFD: {"SBC", (*CPU).addrAbsoluteX, (*CPU).sbc, 4}, // +1 if page crossed
	0xFE: {"INC", (*CPU).addrAbsoluteX, (*CPU).inc, 7},
}
