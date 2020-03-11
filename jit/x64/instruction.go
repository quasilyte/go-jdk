package x64

// instruction is a ready to be encoded machine instruction.
//
// We need this intermediate representation mostly for jump linkage.
type instruction struct {
	imm    int64 // Immediate operand
	disp   int32 // Memory operand displacement (offset)
	offset int32

	buf [16]byte

	opcode uint8  // 8-bit instruction opcode
	size   uint8  // Calculated encoded instruction width
	reg1   uint8  // 3 bits for ModRM.reg
	reg2   uint8  // 3 bits for ModRM.rm
	index  uint8  // 3 bits for SIB.index
	flags  uint16 // Hints on how to encode this instruction
}

func (i instruction) Bytes() []byte {
	return i.buf[:i.size]
}

const (
	flagMemory uint16 = 1 << iota
	flagReg2Op
	flagModRM
	flagModSIB
	flagImm8
	flagImm32
	flagImm64
	flagScale4
	flagScale8
	flagPseudo
	flagRexW
	flag0F
)
