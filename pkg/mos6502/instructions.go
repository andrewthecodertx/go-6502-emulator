package mos6502

import "github.com/andrewthecodertx/go-65c02-emulator/pkg/mos6502/instructions"

// This file contains wrapper methods that delegate to the instruction implementations.
// These methods satisfy the instruction map signatures which expect *CPU methods.

// Load/Store instructions
func (c *CPU) lda(addr uint16, pageCrossed bool) { instructions.LDA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ldx(addr uint16, pageCrossed bool) { instructions.LDX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ldy(addr uint16, pageCrossed bool) { instructions.LDY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sta(addr uint16, pageCrossed bool) { instructions.STA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) stx(addr uint16, pageCrossed bool) { instructions.STX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sty(addr uint16, pageCrossed bool) { instructions.STY(c.BaseCPU, addr, pageCrossed) }

// Transfer instructions
func (c *CPU) tax(addr uint16, pageCrossed bool) { instructions.TAX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) tay(addr uint16, pageCrossed bool) { instructions.TAY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) txa(addr uint16, pageCrossed bool) { instructions.TXA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) tya(addr uint16, pageCrossed bool) { instructions.TYA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) tsx(addr uint16, pageCrossed bool) { instructions.TSX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) txs(addr uint16, pageCrossed bool) { instructions.TXS(c.BaseCPU, addr, pageCrossed) }

// Stack instructions
func (c *CPU) pha(addr uint16, pageCrossed bool) { instructions.PHA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) pla(addr uint16, pageCrossed bool) { instructions.PLA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) php(addr uint16, pageCrossed bool) { instructions.PHP(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) plp(addr uint16, pageCrossed bool) { instructions.PLP(c.BaseCPU, addr, pageCrossed) }

// Arithmetic instructions
func (c *CPU) adc(addr uint16, pageCrossed bool) { instructions.ADC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sbc(addr uint16, pageCrossed bool) { instructions.SBC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cmp(addr uint16, pageCrossed bool) { instructions.CMP(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cpx(addr uint16, pageCrossed bool) { instructions.CPX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cpy(addr uint16, pageCrossed bool) { instructions.CPY(c.BaseCPU, addr, pageCrossed) }

// Increment/Decrement instructions
func (c *CPU) inc(addr uint16, pageCrossed bool) { instructions.INC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) dec(addr uint16, pageCrossed bool) { instructions.DEC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) inx(addr uint16, pageCrossed bool) { instructions.INX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) dex(addr uint16, pageCrossed bool) { instructions.DEX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) iny(addr uint16, pageCrossed bool) { instructions.INY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) dey(addr uint16, pageCrossed bool) { instructions.DEY(c.BaseCPU, addr, pageCrossed) }

// Shift/Rotate instructions
func (c *CPU) asl(addr uint16, pageCrossed bool) { instructions.ASL(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) lsr(addr uint16, pageCrossed bool) { instructions.LSR(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rol(addr uint16, pageCrossed bool) { instructions.ROL(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ror(addr uint16, pageCrossed bool) { instructions.ROR(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) aslAccumulator(addr uint16, pageCrossed bool) { instructions.ASLAccumulator(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) lsrAccumulator(addr uint16, pageCrossed bool) { instructions.LSRAccumulator(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rolAccumulator(addr uint16, pageCrossed bool) { instructions.ROLAccumulator(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rorAccumulator(addr uint16, pageCrossed bool) { instructions.RORAccumulator(c.BaseCPU, addr, pageCrossed) }

// Logic instructions
func (c *CPU) and(addr uint16, pageCrossed bool) { instructions.AND(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ora(addr uint16, pageCrossed bool) { instructions.ORA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) eor(addr uint16, pageCrossed bool) { instructions.EOR(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bit(addr uint16, pageCrossed bool) { instructions.BIT(c.BaseCPU, addr, pageCrossed) }

// Jump/Branch instructions
func (c *CPU) jmp(addr uint16, pageCrossed bool) { instructions.JMP(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) jsr(addr uint16, pageCrossed bool) { instructions.JSR(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rts(addr uint16, pageCrossed bool) { instructions.RTS(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rti(addr uint16, pageCrossed bool) { instructions.RTI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) brk(addr uint16, pageCrossed bool) { instructions.BRK(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bcc(addr uint16, pageCrossed bool) { instructions.BCC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bcs(addr uint16, pageCrossed bool) { instructions.BCS(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) beq(addr uint16, pageCrossed bool) { instructions.BEQ(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bmi(addr uint16, pageCrossed bool) { instructions.BMI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bne(addr uint16, pageCrossed bool) { instructions.BNE(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bpl(addr uint16, pageCrossed bool) { instructions.BPL(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bvc(addr uint16, pageCrossed bool) { instructions.BVC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bvs(addr uint16, pageCrossed bool) { instructions.BVS(c.BaseCPU, addr, pageCrossed) }

// Flag instructions
func (c *CPU) clc(addr uint16, pageCrossed bool) { instructions.CLC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cld(addr uint16, pageCrossed bool) { instructions.CLD(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cli(addr uint16, pageCrossed bool) { instructions.CLI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) clv(addr uint16, pageCrossed bool) { instructions.CLV(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sec(addr uint16, pageCrossed bool) { instructions.SEC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sed(addr uint16, pageCrossed bool) { instructions.SED(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sei(addr uint16, pageCrossed bool) { instructions.SEI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) nop(addr uint16, pageCrossed bool) { instructions.NOP(c.BaseCPU, addr, pageCrossed) }
