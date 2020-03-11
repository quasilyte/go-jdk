package x64

// Various encoding-related constants taken from the x86-64 manual.

const (
	// Rex prefix delivers 4 main bit fields: REX.W, REX.R, REX.X and REX.B.
	//
	//       0100 WRXB
	rexW = 0b0100_1000
	rexR = 0b0100_0100
	rexX = 0b0100_0010
	rexB = 0b0100_0001

	// First three modes are memory addressing with 0, 8 or 32-bit displacement.
	// Last mode is for reg-reg instructions (with no explicit memory operand).
	disp0  = 0b00
	disp8  = 0b01
	disp32 = 0b10
	regreg = 0b11

	scale1 = 0b00
	scale2 = 0b01
	scale4 = 0b10
	scale8 = 0b11

	// Opbytes occupy ModRM/reg field and expressed via /<digit> notation
	// in Intel manual. For example, /2 maps to the op2 constant.
	op0 = 0b000
	op1 = 0b001
	op2 = 0b010
	op3 = 0b011
	op4 = 0b100
	op5 = 0b101
	op6 = 0b110
	op7 = 0b111

	// Registers occupy ModRM/reg and ModRM/rm fields.
	regBitMask = 0b111
	RAX        = 0b000
	RCX        = 0b001
	RDX        = 0b010
	RBX        = 0b011
	RSP        = 0b100
	RBP        = 0b101
	RSI        = 0b110
	RDI        = 0b111
	R8         = 0b1000
	R9         = 0b1001
	R10        = 0b1010
	R11        = 0b1011
	R12        = 0b1100
	R13        = 0b1101
	R14        = 0b1110
	R15        = 0b1111
)

const (
	jmp8op  = 0xEB
	jmp32op = 0xE9
	jge8op  = 0x7D
	jge32op = 0x8D
	jgt8op  = 0x7F
	jgt32op = 0x8F
	jlt8op  = 0x7C
	jlt32op = 0x8C
)

var jumpRel8ToRel32 = [256]byte{
	jmp8op: jmp32op,
	jge8op: jge32op,
	jgt8op: jgt32op,
}
