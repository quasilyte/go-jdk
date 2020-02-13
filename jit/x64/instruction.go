package x64

// instruction is a ready to be encoded machine instruction.
//
// We need this intermediate representation mostly for jump linkage.
type instruction struct {
	imm    int64 // Immediate operand
	disp   int32 // Memory operand displacement (offset)
	offset int32

	buf [16]byte

	opcode uint8 // 8-bit instruction opcode
	size   uint8 // Calculated encoded instruction width
	prefix uint8 // 8-bit instruction prefix (usually REX or none)
	reg1   uint8 // 3 bits for ModRM.reg
	reg2   uint8 // 3 bits for ModRM.rm
	flags  uint8 // Hints on how to encode this instruction
}

func (i instruction) Bytes() []byte {
	return i.buf[:i.size]
}

const (
	flagMemory uint8 = 1 << iota
	flagModRM
	flagImm8
	flagImm32
	flagImm64
	flagPseudo
	flag0F
)
