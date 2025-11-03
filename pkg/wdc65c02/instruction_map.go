package wdc65c02

// Instruction represents a WDC65C02 instruction.
type Instruction struct {
	Name      string
	AddrMode  func(c *CPU) (uint16, bool)
	Operation func(c *CPU, addr uint16, pageCrossed bool)
	Cycles    byte
}

// instructionMap maps opcodes to their instruction implementations.
// This includes all NMOS 6502 instructions plus WDC65C02 enhancements.
var instructionMap = map[byte]Instruction{
	// ========== Original 6502 Instructions ==========
	0x00: {"BRK", nil, (*CPU).brk, 7},
	0x01: {"ORA", (*CPU).addrIndirectX, (*CPU).ora, 6},
	0x05: {"ORA", (*CPU).addrZeroPage, (*CPU).ora, 3},
	0x06: {"ASL", (*CPU).addrZeroPage, (*CPU).asl, 5},
	0x08: {"PHP", nil, (*CPU).php, 3},
	0x09: {"ORA", (*CPU).addrImmediate, (*CPU).ora, 2},
	0x0A: {"ASL", nil, (*CPU).aslAccumulator, 2},
	0x0D: {"ORA", (*CPU).addrAbsolute, (*CPU).ora, 4},
	0x0E: {"ASL", (*CPU).addrAbsolute, (*CPU).asl, 6},
	0x10: {"BPL", (*CPU).addrRelative, (*CPU).bpl, 2},
	0x11: {"ORA", (*CPU).addrIndirectY, (*CPU).ora, 5}, // +1 if page crossed
	0x12: {"ORA", (*CPU).addrZeroPageIndirect, (*CPU).ora, 5}, // NEW: 65C02
	0x15: {"ORA", (*CPU).addrZeroPageX, (*CPU).ora, 4},
	0x16: {"ASL", (*CPU).addrZeroPageX, (*CPU).asl, 6},
	0x18: {"CLC", nil, (*CPU).clc, 2},
	0x19: {"ORA", (*CPU).addrAbsoluteY, (*CPU).ora, 4}, // +1 if page crossed
	0x1A: {"INC", nil, (*CPU).inca, 2}, // NEW: 65C02 - INC A
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
	0x2E: {"ROL", (*CPU).addrAbsolute, (*CPU).rol, 6},
	0x30: {"BMI", (*CPU).addrRelative, (*CPU).bmi, 2},
	0x31: {"AND", (*CPU).addrIndirectY, (*CPU).and, 5}, // +1 if page crossed
	0x32: {"AND", (*CPU).addrZeroPageIndirect, (*CPU).and, 5}, // NEW: 65C02
	0x35: {"AND", (*CPU).addrZeroPageX, (*CPU).and, 4},
	0x36: {"ROL", (*CPU).addrZeroPageX, (*CPU).rol, 6},
	0x38: {"SEC", nil, (*CPU).sec, 2},
	0x39: {"AND", (*CPU).addrAbsoluteY, (*CPU).and, 4}, // +1 if page crossed
	0x3A: {"DEC", nil, (*CPU).deca, 2}, // NEW: 65C02 - DEC A
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
	0x4E: {"LSR", (*CPU).addrAbsolute, (*CPU).lsr, 6},
	0x50: {"BVC", (*CPU).addrRelative, (*CPU).bvc, 2},
	0x51: {"EOR", (*CPU).addrIndirectY, (*CPU).eor, 5}, // +1 if page crossed
	0x52: {"EOR", (*CPU).addrZeroPageIndirect, (*CPU).eor, 5}, // NEW: 65C02
	0x55: {"EOR", (*CPU).addrZeroPageX, (*CPU).eor, 4},
	0x56: {"LSR", (*CPU).addrZeroPageX, (*CPU).lsr, 6},
	0x58: {"CLI", nil, (*CPU).cli, 2},
	0x59: {"EOR", (*CPU).addrAbsoluteY, (*CPU).eor, 4}, // +1 if page crossed
	0x5A: {"PHY", nil, (*CPU).phy, 3}, // NEW: 65C02
	0x5D: {"EOR", (*CPU).addrAbsoluteX, (*CPU).eor, 4}, // +1 if page crossed
	0x5E: {"LSR", (*CPU).addrAbsoluteX, (*CPU).lsr, 6},
	0x60: {"RTS", nil, (*CPU).rts, 6},
	0x61: {"ADC", (*CPU).addrIndirectX, (*CPU).adc, 6},
	0x64: {"STZ", (*CPU).addrZeroPage, (*CPU).stz, 3}, // NEW: 65C02
	0x65: {"ADC", (*CPU).addrZeroPage, (*CPU).adc, 3},
	0x66: {"ROR", (*CPU).addrZeroPage, (*CPU).ror, 5},
	0x68: {"PLA", nil, (*CPU).pla, 4},
	0x69: {"ADC", (*CPU).addrImmediate, (*CPU).adc, 2},
	0x6A: {"ROR", nil, (*CPU).rorAccumulator, 2},
	0x6C: {"JMP", (*CPU).addrIndirect, (*CPU).jmp, 5}, // Bug FIXED on 65C02
	0x6D: {"ADC", (*CPU).addrAbsolute, (*CPU).adc, 4},
	0x6E: {"ROR", (*CPU).addrAbsolute, (*CPU).ror, 6},
	0x70: {"BVS", (*CPU).addrRelative, (*CPU).bvs, 2},
	0x71: {"ADC", (*CPU).addrIndirectY, (*CPU).adc, 5}, // +1 if page crossed
	0x72: {"ADC", (*CPU).addrZeroPageIndirect, (*CPU).adc, 5}, // NEW: 65C02
	0x74: {"STZ", (*CPU).addrZeroPageX, (*CPU).stz, 4}, // NEW: 65C02
	0x75: {"ADC", (*CPU).addrZeroPageX, (*CPU).adc, 4},
	0x76: {"ROR", (*CPU).addrZeroPageX, (*CPU).ror, 6},
	0x78: {"SEI", nil, (*CPU).sei, 2},
	0x79: {"ADC", (*CPU).addrAbsoluteY, (*CPU).adc, 4}, // +1 if page crossed
	0x7A: {"PLY", nil, (*CPU).ply, 4}, // NEW: 65C02
	0x7D: {"ADC", (*CPU).addrAbsoluteX, (*CPU).adc, 4}, // +1 if page crossed
	0x7E: {"ROR", (*CPU).addrAbsoluteX, (*CPU).ror, 7},
	0x80: {"BRA", (*CPU).addrRelative, (*CPU).bra, 3}, // NEW: 65C02
	0x81: {"STA", (*CPU).addrIndirectX, (*CPU).sta, 6},
	0x84: {"STY", (*CPU).addrZeroPage, (*CPU).sty, 3},
	0x85: {"STA", (*CPU).addrZeroPage, (*CPU).sta, 3},
	0x86: {"STX", (*CPU).addrZeroPage, (*CPU).stx, 3},
	0x88: {"DEY", nil, (*CPU).dey, 2},
	0x89: {"BIT", (*CPU).addrImmediate, (*CPU).bit, 2}, // NEW: 65C02 - BIT immediate
	0x8A: {"TXA", nil, (*CPU).txa, 2},
	0x8C: {"STY", (*CPU).addrAbsolute, (*CPU).sty, 4},
	0x8D: {"STA", (*CPU).addrAbsolute, (*CPU).sta, 4},
	0x8E: {"STX", (*CPU).addrAbsolute, (*CPU).stx, 4},
	0x90: {"BCC", (*CPU).addrRelative, (*CPU).bcc, 2},
	0x91: {"STA", (*CPU).addrIndirectY, (*CPU).sta, 6},
	0x92: {"STA", (*CPU).addrZeroPageIndirect, (*CPU).sta, 5}, // NEW: 65C02
	0x94: {"STY", (*CPU).addrZeroPageX, (*CPU).sty, 4},
	0x95: {"STA", (*CPU).addrZeroPageX, (*CPU).sta, 4},
	0x96: {"STX", (*CPU).addrZeroPageY, (*CPU).stx, 4},
	0x98: {"TYA", nil, (*CPU).tya, 2},
	0x99: {"STA", (*CPU).addrAbsoluteY, (*CPU).sta, 5},
	0x9A: {"TXS", nil, (*CPU).txs, 2},
	0x9C: {"STZ", (*CPU).addrAbsolute, (*CPU).stz, 4}, // NEW: 65C02
	0x9D: {"STA", (*CPU).addrAbsoluteX, (*CPU).sta, 5},
	0x9E: {"STZ", (*CPU).addrAbsoluteX, (*CPU).stz, 5}, // NEW: 65C02
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
	0xB0: {"BCS", (*CPU).addrRelative, (*CPU).bcs, 2},
	0xB1: {"LDA", (*CPU).addrIndirectY, (*CPU).lda, 5}, // +1 if page crossed
	0xB2: {"LDA", (*CPU).addrZeroPageIndirect, (*CPU).lda, 5}, // NEW: 65C02
	0xB4: {"LDY", (*CPU).addrZeroPageX, (*CPU).ldy, 4},
	0xB5: {"LDA", (*CPU).addrZeroPageX, (*CPU).lda, 4},
	0xB6: {"LDX", (*CPU).addrZeroPageY, (*CPU).ldx, 4},
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
	0xCB: {"WAI", nil, (*CPU).wai, 3}, // NEW: 65C02
	0xCC: {"CPY", (*CPU).addrAbsolute, (*CPU).cpy, 4},
	0xCD: {"CMP", (*CPU).addrAbsolute, (*CPU).cmp, 4},
	0xCE: {"DEC", (*CPU).addrAbsolute, (*CPU).dec, 6},
	0xD0: {"BNE", (*CPU).addrRelative, (*CPU).bne, 2},
	0xD1: {"CMP", (*CPU).addrIndirectY, (*CPU).cmp, 5}, // +1 if page crossed
	0xD2: {"CMP", (*CPU).addrZeroPageIndirect, (*CPU).cmp, 5}, // NEW: 65C02
	0xD5: {"CMP", (*CPU).addrZeroPageX, (*CPU).cmp, 4},
	0xD6: {"DEC", (*CPU).addrZeroPageX, (*CPU).dec, 6},
	0xD8: {"CLD", nil, (*CPU).cld, 2},
	0xD9: {"CMP", (*CPU).addrAbsoluteY, (*CPU).cmp, 4}, // +1 if page crossed
	0xDA: {"PHX", nil, (*CPU).phx, 3}, // NEW: 65C02
	0xDB: {"STP", nil, (*CPU).stp, 3}, // NEW: 65C02
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
	0xEE: {"INC", (*CPU).addrAbsolute, (*CPU).inc, 6},
	0xF0: {"BEQ", (*CPU).addrRelative, (*CPU).beq, 2},
	0xF1: {"SBC", (*CPU).addrIndirectY, (*CPU).sbc, 5}, // +1 if page crossed
	0xF2: {"SBC", (*CPU).addrZeroPageIndirect, (*CPU).sbc, 5}, // NEW: 65C02
	0xF5: {"SBC", (*CPU).addrZeroPageX, (*CPU).sbc, 4},
	0xF6: {"INC", (*CPU).addrZeroPageX, (*CPU).inc, 6},
	0xF8: {"SED", nil, (*CPU).sed, 2},
	0xF9: {"SBC", (*CPU).addrAbsoluteY, (*CPU).sbc, 4}, // +1 if page crossed
	0xFA: {"PLX", nil, (*CPU).plx, 4}, // NEW: 65C02
	0xFD: {"SBC", (*CPU).addrAbsoluteX, (*CPU).sbc, 4}, // +1 if page crossed
	0xFE: {"INC", (*CPU).addrAbsoluteX, (*CPU).inc, 7},

	// ========== NEW WDC65C02 Instructions: TSB/TRB ==========
	0x04: {"TSB", (*CPU).addrZeroPage, (*CPU).tsb, 5}, // NEW: 65C02
	0x0C: {"TSB", (*CPU).addrAbsolute, (*CPU).tsb, 6}, // NEW: 65C02
	0x14: {"TRB", (*CPU).addrZeroPage, (*CPU).trb, 5}, // NEW: 65C02
	0x1C: {"TRB", (*CPU).addrAbsolute, (*CPU).trb, 6}, // NEW: 65C02

	// ========== NEW WDC65C02 Instructions: RMB (Reset Memory Bit) ==========
	0x07: {"RMB0", (*CPU).addrZeroPage, (*CPU).rmb0, 5}, // NEW: 65C02
	0x17: {"RMB1", (*CPU).addrZeroPage, (*CPU).rmb1, 5}, // NEW: 65C02
	0x27: {"RMB2", (*CPU).addrZeroPage, (*CPU).rmb2, 5}, // NEW: 65C02
	0x37: {"RMB3", (*CPU).addrZeroPage, (*CPU).rmb3, 5}, // NEW: 65C02
	0x47: {"RMB4", (*CPU).addrZeroPage, (*CPU).rmb4, 5}, // NEW: 65C02
	0x57: {"RMB5", (*CPU).addrZeroPage, (*CPU).rmb5, 5}, // NEW: 65C02
	0x67: {"RMB6", (*CPU).addrZeroPage, (*CPU).rmb6, 5}, // NEW: 65C02
	0x77: {"RMB7", (*CPU).addrZeroPage, (*CPU).rmb7, 5}, // NEW: 65C02

	// ========== NEW WDC65C02 Instructions: SMB (Set Memory Bit) ==========
	0x87: {"SMB0", (*CPU).addrZeroPage, (*CPU).smb0, 5}, // NEW: 65C02
	0x97: {"SMB1", (*CPU).addrZeroPage, (*CPU).smb1, 5}, // NEW: 65C02
	0xA7: {"SMB2", (*CPU).addrZeroPage, (*CPU).smb2, 5}, // NEW: 65C02
	0xB7: {"SMB3", (*CPU).addrZeroPage, (*CPU).smb3, 5}, // NEW: 65C02
	0xC7: {"SMB4", (*CPU).addrZeroPage, (*CPU).smb4, 5}, // NEW: 65C02
	0xD7: {"SMB5", (*CPU).addrZeroPage, (*CPU).smb5, 5}, // NEW: 65C02
	0xE7: {"SMB6", (*CPU).addrZeroPage, (*CPU).smb6, 5}, // NEW: 65C02
	0xF7: {"SMB7", (*CPU).addrZeroPage, (*CPU).smb7, 5}, // NEW: 65C02

	// ========== NEW WDC65C02 Instructions: BBR (Branch on Bit Reset) ==========
	0x0F: {"BBR0", (*CPU).addrZeroPage, (*CPU).bbr0, 5}, // NEW: 65C02
	0x1F: {"BBR1", (*CPU).addrZeroPage, (*CPU).bbr1, 5}, // NEW: 65C02
	0x2F: {"BBR2", (*CPU).addrZeroPage, (*CPU).bbr2, 5}, // NEW: 65C02
	0x3F: {"BBR3", (*CPU).addrZeroPage, (*CPU).bbr3, 5}, // NEW: 65C02
	0x4F: {"BBR4", (*CPU).addrZeroPage, (*CPU).bbr4, 5}, // NEW: 65C02
	0x5F: {"BBR5", (*CPU).addrZeroPage, (*CPU).bbr5, 5}, // NEW: 65C02
	0x6F: {"BBR6", (*CPU).addrZeroPage, (*CPU).bbr6, 5}, // NEW: 65C02
	0x7F: {"BBR7", (*CPU).addrZeroPage, (*CPU).bbr7, 5}, // NEW: 65C02

	// ========== NEW WDC65C02 Instructions: BBS (Branch on Bit Set) ==========
	0x8F: {"BBS0", (*CPU).addrZeroPage, (*CPU).bbs0, 5}, // NEW: 65C02
	0x9F: {"BBS1", (*CPU).addrZeroPage, (*CPU).bbs1, 5}, // NEW: 65C02
	0xAF: {"BBS2", (*CPU).addrZeroPage, (*CPU).bbs2, 5}, // NEW: 65C02
	0xBF: {"BBS3", (*CPU).addrZeroPage, (*CPU).bbs3, 5}, // NEW: 65C02
	0xCF: {"BBS4", (*CPU).addrZeroPage, (*CPU).bbs4, 5}, // NEW: 65C02
	0xDF: {"BBS5", (*CPU).addrZeroPage, (*CPU).bbs5, 5}, // NEW: 65C02
	0xEF: {"BBS6", (*CPU).addrZeroPage, (*CPU).bbs6, 5}, // NEW: 65C02
	0xFF: {"BBS7", (*CPU).addrZeroPage, (*CPU).bbs7, 5}, // NEW: 65C02

	// ========== Additional NOP variants (65C02 fills illegal opcodes with NOPs) ==========
	0x02: {"NOP", (*CPU).addrImmediate, (*CPU).nop, 2},
	0x03: {"NOP", nil, (*CPU).nop, 1},
	0x0B: {"NOP", nil, (*CPU).nop, 1},
	0x13: {"NOP", nil, (*CPU).nop, 1},
	0x1B: {"NOP", nil, (*CPU).nop, 1},
	0x22: {"NOP", (*CPU).addrImmediate, (*CPU).nop, 2},
	0x23: {"NOP", nil, (*CPU).nop, 1},
	0x2B: {"NOP", nil, (*CPU).nop, 1},
	0x33: {"NOP", nil, (*CPU).nop, 1},
	0x3B: {"NOP", nil, (*CPU).nop, 1},
	0x42: {"NOP", (*CPU).addrImmediate, (*CPU).nop, 2},
	0x43: {"NOP", nil, (*CPU).nop, 1},
	0x44: {"NOP", (*CPU).addrZeroPage, (*CPU).nop, 3},
	0x4B: {"NOP", nil, (*CPU).nop, 1},
	0x53: {"NOP", nil, (*CPU).nop, 1},
	0x54: {"NOP", (*CPU).addrZeroPage, (*CPU).nop, 4},
	0x5B: {"NOP", nil, (*CPU).nop, 1},
	0x5C: {"NOP", (*CPU).addrAbsolute, (*CPU).nop, 8},
	0x62: {"NOP", (*CPU).addrImmediate, (*CPU).nop, 2},
	0x63: {"NOP", nil, (*CPU).nop, 1},
	0x6B: {"NOP", nil, (*CPU).nop, 1},
	0x73: {"NOP", nil, (*CPU).nop, 1},
	0x7B: {"NOP", nil, (*CPU).nop, 1},
	0x82: {"NOP", (*CPU).addrImmediate, (*CPU).nop, 2},
	0x83: {"NOP", nil, (*CPU).nop, 1},
	0x8B: {"NOP", nil, (*CPU).nop, 1},
	0x93: {"NOP", nil, (*CPU).nop, 1},
	0x9B: {"NOP", nil, (*CPU).nop, 1},
	0xA3: {"NOP", nil, (*CPU).nop, 1},
	0xAB: {"NOP", nil, (*CPU).nop, 1},
	0xB3: {"NOP", nil, (*CPU).nop, 1},
	0xBB: {"NOP", nil, (*CPU).nop, 1},
	0xC2: {"NOP", (*CPU).addrImmediate, (*CPU).nop, 2},
	0xC3: {"NOP", nil, (*CPU).nop, 1},
	0xD3: {"NOP", nil, (*CPU).nop, 1},
	0xD4: {"NOP", (*CPU).addrZeroPage, (*CPU).nop, 4},
	0xDC: {"NOP", (*CPU).addrAbsolute, (*CPU).nop, 4},
	0xE2: {"NOP", (*CPU).addrImmediate, (*CPU).nop, 2},
	0xE3: {"NOP", nil, (*CPU).nop, 1},
	0xEB: {"NOP", nil, (*CPU).nop, 1},
	0xF3: {"NOP", nil, (*CPU).nop, 1},
	0xF4: {"NOP", (*CPU).addrZeroPage, (*CPU).nop, 4},
	0xFB: {"NOP", nil, (*CPU).nop, 1},
	0xFC: {"NOP", (*CPU).addrAbsolute, (*CPU).nop, 4},

	// JMP (absolute,X) - NEW addressing mode on 65C02
	0x7C: {"JMP", (*CPU).addrAbsoluteIndexedIndirect, (*CPU).jmp, 6}, // NEW: 65C02
}
