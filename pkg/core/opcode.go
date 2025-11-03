package core

type (
	AddressingMode     uint8
	InstructionHandler func(cpu interface{}, opcode *Opcode) uint8
)

const (
	Implied AddressingMode = iota
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
	Relative
	Accumulator

	// 65C02 SPECIFIC MODES
	ZeroPageIndirect
	AbsoluteXIndirect
)

type Opcode struct {
	Mnemonic       string
	Code           uint8
	Mode           AddressingMode
	Cycles         uint8
	PageCrossCycle bool
	BranchCycle    bool
	Length         uint8
	Handler        InstructionHandler
	Illegal        bool
}
