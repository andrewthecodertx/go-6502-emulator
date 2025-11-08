package wdc65c02

import (
	"github.com/andrewthecodertx/go-6502-emulator/pkg/wdc65c02/instructions"
)

// Wrapper methods that delegate to the instructions package
// These exist to satisfy the instruction map signature

// ========== Load/Store ==========
func (c *CPU) lda(addr uint16, pageCrossed bool) { instructions.LDA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ldx(addr uint16, pageCrossed bool) { instructions.LDX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ldy(addr uint16, pageCrossed bool) { instructions.LDY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sta(addr uint16, pageCrossed bool) { instructions.STA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) stx(addr uint16, pageCrossed bool) { instructions.STX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sty(addr uint16, pageCrossed bool) { instructions.STY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) stz(addr uint16, pageCrossed bool) { instructions.STZ(c.BaseCPU, addr, pageCrossed) }

// ========== Transfer ==========
func (c *CPU) tax(addr uint16, pageCrossed bool) { instructions.TAX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) tay(addr uint16, pageCrossed bool) { instructions.TAY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) txa(addr uint16, pageCrossed bool) { instructions.TXA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) tya(addr uint16, pageCrossed bool) { instructions.TYA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) tsx(addr uint16, pageCrossed bool) { instructions.TSX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) txs(addr uint16, pageCrossed bool) { instructions.TXS(c.BaseCPU, addr, pageCrossed) }

// ========== Stack ==========
func (c *CPU) pha(addr uint16, pageCrossed bool) { instructions.PHA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) pla(addr uint16, pageCrossed bool) { instructions.PLA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) php(addr uint16, pageCrossed bool) { instructions.PHP(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) plp(addr uint16, pageCrossed bool) { instructions.PLP(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) phx(addr uint16, pageCrossed bool) { instructions.PHX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) phy(addr uint16, pageCrossed bool) { instructions.PHY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) plx(addr uint16, pageCrossed bool) { instructions.PLX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ply(addr uint16, pageCrossed bool) { instructions.PLY(c.BaseCPU, addr, pageCrossed) }

// ========== Arithmetic ==========
func (c *CPU) adc(addr uint16, pageCrossed bool) { instructions.ADC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sbc(addr uint16, pageCrossed bool) { instructions.SBC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cmp(addr uint16, pageCrossed bool) { instructions.CMP(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cpx(addr uint16, pageCrossed bool) { instructions.CPX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cpy(addr uint16, pageCrossed bool) { instructions.CPY(c.BaseCPU, addr, pageCrossed) }

// ========== Increment/Decrement ==========
func (c *CPU) inc(addr uint16, pageCrossed bool)  { instructions.INC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) dec(addr uint16, pageCrossed bool)  { instructions.DEC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) inx(addr uint16, pageCrossed bool)  { instructions.INX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) dex(addr uint16, pageCrossed bool)  { instructions.DEX(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) iny(addr uint16, pageCrossed bool)  { instructions.INY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) dey(addr uint16, pageCrossed bool)  { instructions.DEY(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) inca(addr uint16, pageCrossed bool) { instructions.INCA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) deca(addr uint16, pageCrossed bool) { instructions.DECA(c.BaseCPU, addr, pageCrossed) }

// ========== Shift/Rotate ==========
func (c *CPU) asl(addr uint16, pageCrossed bool)            { instructions.ASL(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) aslAccumulator(addr uint16, pageCrossed bool) { instructions.ASLAccumulator(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) lsr(addr uint16, pageCrossed bool)            { instructions.LSR(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) lsrAccumulator(addr uint16, pageCrossed bool) { instructions.LSRAccumulator(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rol(addr uint16, pageCrossed bool)            { instructions.ROL(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rolAccumulator(addr uint16, pageCrossed bool) { instructions.ROLAccumulator(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ror(addr uint16, pageCrossed bool)            { instructions.ROR(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rorAccumulator(addr uint16, pageCrossed bool) { instructions.RORAccumulator(c.BaseCPU, addr, pageCrossed) }

// ========== Logic ==========
func (c *CPU) and(addr uint16, pageCrossed bool) { instructions.AND(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) ora(addr uint16, pageCrossed bool) { instructions.ORA(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) eor(addr uint16, pageCrossed bool) { instructions.EOR(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bit(addr uint16, pageCrossed bool) { instructions.BIT(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) tsb(addr uint16, pageCrossed bool) { instructions.TSB(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) trb(addr uint16, pageCrossed bool) { instructions.TRB(c.BaseCPU, addr, pageCrossed) }

// ========== Jump/Branch ==========
func (c *CPU) jmp(addr uint16, pageCrossed bool) { instructions.JMP(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) jsr(addr uint16, pageCrossed bool) { instructions.JSR(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rts(addr uint16, pageCrossed bool) { instructions.RTS(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rti(addr uint16, pageCrossed bool) { instructions.RTI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) brk(addr uint16, pageCrossed bool) { instructions.BRK(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bcc(addr uint16, pageCrossed bool) { instructions.BCC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bcs(addr uint16, pageCrossed bool) { instructions.BCS(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) beq(addr uint16, pageCrossed bool) { instructions.BEQ(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bne(addr uint16, pageCrossed bool) { instructions.BNE(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bmi(addr uint16, pageCrossed bool) { instructions.BMI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bpl(addr uint16, pageCrossed bool) { instructions.BPL(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bvc(addr uint16, pageCrossed bool) { instructions.BVC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bvs(addr uint16, pageCrossed bool) { instructions.BVS(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bra(addr uint16, pageCrossed bool) { instructions.BRA(c.BaseCPU, addr, pageCrossed) }

// ========== Flags ==========
func (c *CPU) clc(addr uint16, pageCrossed bool) { instructions.CLC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sec(addr uint16, pageCrossed bool) { instructions.SEC(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cli(addr uint16, pageCrossed bool) { instructions.CLI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sei(addr uint16, pageCrossed bool) { instructions.SEI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) cld(addr uint16, pageCrossed bool) { instructions.CLD(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) sed(addr uint16, pageCrossed bool) { instructions.SED(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) clv(addr uint16, pageCrossed bool) { instructions.CLV(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) nop(addr uint16, pageCrossed bool) { instructions.NOP(c.BaseCPU, addr, pageCrossed) }

// ========== Bit Manipulation ==========
func (c *CPU) rmb0(addr uint16, pageCrossed bool) { instructions.RMB0(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rmb1(addr uint16, pageCrossed bool) { instructions.RMB1(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rmb2(addr uint16, pageCrossed bool) { instructions.RMB2(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rmb3(addr uint16, pageCrossed bool) { instructions.RMB3(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rmb4(addr uint16, pageCrossed bool) { instructions.RMB4(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rmb5(addr uint16, pageCrossed bool) { instructions.RMB5(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rmb6(addr uint16, pageCrossed bool) { instructions.RMB6(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) rmb7(addr uint16, pageCrossed bool) { instructions.RMB7(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) smb0(addr uint16, pageCrossed bool) { instructions.SMB0(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) smb1(addr uint16, pageCrossed bool) { instructions.SMB1(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) smb2(addr uint16, pageCrossed bool) { instructions.SMB2(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) smb3(addr uint16, pageCrossed bool) { instructions.SMB3(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) smb4(addr uint16, pageCrossed bool) { instructions.SMB4(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) smb5(addr uint16, pageCrossed bool) { instructions.SMB5(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) smb6(addr uint16, pageCrossed bool) { instructions.SMB6(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) smb7(addr uint16, pageCrossed bool) { instructions.SMB7(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbr0(addr uint16, pageCrossed bool) { instructions.BBR0(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbr1(addr uint16, pageCrossed bool) { instructions.BBR1(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbr2(addr uint16, pageCrossed bool) { instructions.BBR2(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbr3(addr uint16, pageCrossed bool) { instructions.BBR3(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbr4(addr uint16, pageCrossed bool) { instructions.BBR4(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbr5(addr uint16, pageCrossed bool) { instructions.BBR5(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbr6(addr uint16, pageCrossed bool) { instructions.BBR6(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbr7(addr uint16, pageCrossed bool) { instructions.BBR7(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbs0(addr uint16, pageCrossed bool) { instructions.BBS0(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbs1(addr uint16, pageCrossed bool) { instructions.BBS1(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbs2(addr uint16, pageCrossed bool) { instructions.BBS2(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbs3(addr uint16, pageCrossed bool) { instructions.BBS3(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbs4(addr uint16, pageCrossed bool) { instructions.BBS4(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbs5(addr uint16, pageCrossed bool) { instructions.BBS5(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbs6(addr uint16, pageCrossed bool) { instructions.BBS6(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) bbs7(addr uint16, pageCrossed bool) { instructions.BBS7(c.BaseCPU, addr, pageCrossed) }

// ========== CPU Control ==========
func (c *CPU) wai(addr uint16, pageCrossed bool) { instructions.WAI(c.BaseCPU, addr, pageCrossed) }
func (c *CPU) stp(addr uint16, pageCrossed bool) { instructions.STP(c.BaseCPU, addr, pageCrossed) }
