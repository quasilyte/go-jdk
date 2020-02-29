package x64

import (
	"fmt"
	"strings"
	"testing"
)

func TestAsm(t *testing.T) {
	type expected struct {
		line int
		asm  string
		enc  string
	}

	// Test cases are generated with a help of `testdata/gen.go`.
	tests := []struct {
		name string
		want []expected
		run  func(*Assembler)
	}{
		{
			name: "testJge1",
			want: []expected{
				{5, "JGE forward2", "7d00"},
				{7, "JGE forward1", "7d01"},
				{8, "NOP1", "90"},
				{10, "NOP1", "90"},
			},
			run: func(asm *Assembler) {
				asm.Jge(2)
				asm.Label(2)
				asm.Jge(1)
				asm.Nop(1)
				asm.Label(1)
				asm.Nop(1)
			},
		},

		{
			name: "testJge2",
			want: []expected{
				{15, "NOP1", "90"},
				{16, "JGE l1", "7d03"},
				{18, "NOP1", "90"},
				{19, "JGE l2", "7dfa"},
				{21, "NOP1", "90"},
				{22, "JGE l3", "7dfa"},
			},
			run: func(asm *Assembler) {
				asm.Label(2)
				asm.Nop(1)
				asm.Jge(1)
				asm.Label(3)
				asm.Nop(1)
				asm.Jge(2)
				asm.Label(1)
				asm.Nop(1)
				asm.Jge(3)
			},
		},

		{
			name: "testJmp1",
			want: []expected{
				{26, "JMP forward", "eb00"},
			},
			run: func(asm *Assembler) {
				asm.Jmp(7)
				asm.Label(7)
			},
		},

		{
			name: "testJmp2",
			want: []expected{
				{32, "JMP looping", "ebfe"},
			},
			run: func(asm *Assembler) {
				asm.Label(3)
				asm.Jmp(3)
			},
		},

		{
			name: "testJmp3",
			want: []expected{
				{37, "NOP1", "90"},
				{38, "JMP backward", "ebfd"},
			},
			run: func(asm *Assembler) {
				asm.Label(0)
				asm.Nop(1)
				asm.Jmp(0)
			},
		},

		{
			name: "testJmp4",
			want: []expected{
				{42, "JMP sharedlabel", "eb04"},
				{43, "NOP1", "90"},
				{44, "JMP sharedlabel", "eb01"},
				{45, "NOP1", "90"},
				{47, "NOP1", "90"},
			},
			run: func(asm *Assembler) {
				asm.Jmp(0)
				asm.Nop(1)
				asm.Jmp(0)
				asm.Nop(1)
				asm.Label(0)
				asm.Nop(1)
			},
		},

		{
			name: "testJmp5",
			want: []expected{
				{51, "NOP1", "90"},
				{52, "JMP l1", "eb03"},
				{54, "NOP1", "90"},
				{55, "JMP l2", "eb03"},
				{57, "NOP1", "90"},
				{58, "JMP l3", "ebfa"},
				{60, "NOP1", "90"},
			},
			run: func(asm *Assembler) {
				asm.Nop(1)
				asm.Jmp(1)
				asm.Label(3)
				asm.Nop(1)
				asm.Jmp(2)
				asm.Label(1)
				asm.Nop(1)
				asm.Jmp(3)
				asm.Label(2)
				asm.Nop(1)
			},
		},

		{
			name: "testJmp6",
			want: []expected{
				{65, "NOP1", "90"},
				{66, "JMP l1", "eb03"},
				{68, "NOP1", "90"},
				{69, "JMP l2", "ebfa"},
				{71, "NOP1", "90"},
				{72, "JMP l3", "ebfa"},
			},
			run: func(asm *Assembler) {
				asm.Label(2)
				asm.Nop(1)
				asm.Jmp(1)
				asm.Label(3)
				asm.Nop(1)
				asm.Jmp(2)
				asm.Label(1)
				asm.Nop(1)
				asm.Jmp(3)
			},
		},

		{
			name: "testJmpReg",
			want: []expected{
				{76, "JMP AX", "ffe0"},
				{77, "JMP DX", "ffe2"},
				{78, "JMP CX", "ffe1"},
			},
			run: func(asm *Assembler) {
				asm.JmpReg(RAX)
				asm.JmpReg(RDX)
				asm.JmpReg(RCX)
			},
		},

		{
			name: "testJmpMem",
			want: []expected{
				{82, "JMP (AX)", "ff20"},
				{83, "JMP -8(DI)", "ff67f8"},
				{84, "JMP 13935(CX)", "ffa16f360000"},
			},
			run: func(asm *Assembler) {
				asm.JmpMem(RAX, 0)
				asm.JmpMem(RDI, -8)
				asm.JmpMem(RCX, 13935)
			},
		},
		{
			name: "testAdd",
			want: []expected{
				{88, "ADDL (AX), DX", "0310"},
				{89, "ADDL 8(SI), AX", "034608"},
				{90, "ADDL $7, (AX)", "830007"},
				{91, "ADDL $-9, -8(DX)", "8342f8f7"},
				{92, "ADDL $9300, 16(SI)", "81461054240000"},
				{93, "ADDL $127, CX", "83c17f"},
				{94, "ADDL $-128, BX", "83c380"},
				{95, "ADDL $200, BP", "81c5c8000000"},
				{96, "ADDQ $0, 0*8(SI)", "48830600"},
				{97, "ADDQ $1, 0*8(SI)", "48830601"},
				{98, "ADDQ $1, 1*8(SI)", "4883460801"},
				{99, "ADDQ $-1, 3*8(SI)", "48834618ff"},
				{100, "ADDQ $14, 10*8(SI)", "488346500e"},
				{101, "ADDQ $14, 100*8(SI)", "488386200300000e"},
				{102, "ADDQ $0xff, 0*8(SI)", "488106ff000000"},
				{103, "ADDQ $0xff, 1*8(SI)", "48814608ff000000"},
				{104, "ADDQ $-129, 100*8(SI)", "488186200300007fffffff"},
				{105, "ADDQ $1, AX", "4883c001"},
				{106, "ADDQ $-1, DI", "4883c7ff"},
				{107, "ADDQ $5000, SP", "4881c488130000"},
			},
			run: func(asm *Assembler) {
				asm.AddlMemReg(RAX, RDX, 0)
				asm.AddlMemReg(RSI, RAX, 8)
				asm.AddlConst8Mem(7, RAX, 0)
				asm.AddlConst8Mem(-9, RDX, -8)
				asm.AddlConst32Mem(9300, RSI, 16)
				asm.AddlConst8Reg(127, RCX)
				asm.AddlConst8Reg(-128, RBX)
				asm.AddlConst32Reg(200, RBP)
				asm.AddqConst8Mem(0, RSI, 0*8)
				asm.AddqConst8Mem(1, RSI, 0*8)
				asm.AddqConst8Mem(1, RSI, 1*8)
				asm.AddqConst8Mem(-1, RSI, 3*8)
				asm.AddqConst8Mem(14, RSI, 10*8)
				asm.AddqConst8Mem(14, RSI, 100*8)
				asm.AddqConst32Mem(0xff, RSI, 0*8)
				asm.AddqConst32Mem(0xff, RSI, 1*8)
				asm.AddqConst32Mem(-129, RSI, 100*8)
				asm.AddqConst8Reg(1, RAX)
				asm.AddqConst8Reg(-1, RDI)
				asm.AddqConst32Reg(5000, RSP)
			},
		},

		{
			name: "testMov",
			want: []expected{
				{111, "MOVL $0, 0*8(SI)", "c70600000000"},
				{112, "MOVL $1, 0*8(DI)", "c70701000000"},
				{113, "MOVL $1, 1*8(AX)", "c7400801000000"},
				{114, "MOVL $-50000, 40*8(SI)", "c78640010000b03cffff"},
				{115, "MOVL (AX), AX", "8b00"},
				{116, "MOVL -16(CX), DX", "8b51f0"},
				{117, "MOVL AX, (AX)", "8900"},
				{118, "MOVL DX, -16(CX)", "8951f0"},
				{119, "MOVL $1355, AX", "b84b050000"},
				{120, "MOVL $-6643, DX", "ba0de6ffff"},
				{121, "MOVQ 0*8(AX), BX", "488b18"},
				{122, "MOVQ 16*8(BX), AX", "488b8380000000"},
				{123, "MOVQ AX, 0*8(DI)", "488907"},
				{124, "MOVQ DX, 3*8(DI)", "48895718"},
				{125, "MOVQ AX, 0*8(AX)", "488900"},
				{126, "MOVQ $140038723203072, AX", "48b800f0594e5d7f0000"},
				{127, "MOVQ $9223372036854775807, DX", "48baffffffffffffff7f"},
				{128, "MOVQ $-9223372036854775800, SI", "48be0800000000000080"},
				{129, "MOVQ $1423, AX", "48c7c08f050000"},
				{130, "MOVQ $-23, CX", "48c7c1e9ffffff"},
				{131, "MOVQ $1, DX", "48c7c201000000"},
				{132, "MOVQ 100(BP), DX", "488b5564"},
				{133, "MOVQ $1, 1(AX)", "48c7400101000000"},
				{134, "MOVQ $-1, 2(AX)", "48c74002ffffffff"},
				{135, "MOVQ $0, -96(BP)", "48c745a000000000"},
				{136, "MOVQ $100, -96(BP)", "48c745a064000000"},
			},
			run: func(asm *Assembler) {
				asm.MovlConst32Mem(0, RSI, 0*8)
				asm.MovlConst32Mem(1, RDI, 0*8)
				asm.MovlConst32Mem(1, RAX, 1*8)
				asm.MovlConst32Mem(-50000, RSI, 40*8)
				asm.MovlMemReg(RAX, RAX, 0)
				asm.MovlMemReg(RCX, RDX, -16)
				asm.MovlRegMem(RAX, RAX, 0)
				asm.MovlRegMem(RDX, RCX, -16)
				asm.MovlConst32Reg(1355, RAX)
				asm.MovlConst32Reg(-6643, RDX)
				asm.MovqMemReg(RAX, RBX, 0*8)
				asm.MovqMemReg(RBX, RAX, 16*8)
				asm.MovqRegMem(RAX, RDI, 0*8)
				asm.MovqRegMem(RDX, RDI, 3*8)
				asm.MovqRegMem(RAX, RAX, 0*8)
				asm.MovqConst64Reg(140038723203072, RAX)
				asm.MovqConst64Reg(9223372036854775807, RDX)
				asm.MovqConst64Reg(-9223372036854775800, RSI)
				asm.MovqConst32Reg(1423, RAX)
				asm.MovqConst32Reg(-23, RCX)
				asm.MovqConst32Reg(1, RDX)
				asm.MovqMemReg(RBP, RDX, 100)
				asm.MovqConst32Mem(1, RAX, 1)
				asm.MovqConst32Mem(-1, RAX, 2)
				asm.MovqConst32Mem(0, RBP, -96)
				asm.MovqConst32Mem(100, RBP, -96)
			},
		},

		{
			name: "testCmp",
			want: []expected{
				{140, "CMPL AX, 0*8(DI)", "3b07"},
				{141, "CMPL BX, 1*8(AX)", "3b5808"},
				{142, "CMPL 16(SI), $0", "837e1000"},
				{143, "CMPL (AX), $15", "83380f"},
				{144, "CMPL (DI), $242", "813ff2000000"},
				{145, "CMPL -8(BX), $-5343", "817bf821ebffff"},
				{146, "CMPQ 6*8(SI), $0", "48837e3000"},
				{147, "CMPQ (SI), $999", "48813ee7030000"},
				{148, "CMPQ 8(DI), $-999", "48817f0819fcffff"},
			},
			run: func(asm *Assembler) {
				asm.CmplRegMem(RAX, RDI, 0*8)
				asm.CmplRegMem(RBX, RAX, 1*8)
				asm.CmplConst8Mem(0, RSI, 16)
				asm.CmplConst8Mem(15, RAX, 0)
				asm.CmplConst32Mem(242, RDI, 0)
				asm.CmplConst32Mem(-5343, RBX, -8)
				asm.CmpqConst8Mem(0, RSI, 6*8)
				asm.CmpqConst32Mem(999, RSI, 0)
				asm.CmpqConst32Mem(-999, RDI, 8)
			},
		},

		{
			name: "testNeg",
			want: []expected{
				{152, "NEGQ 0*8(SI)", "48f71e"},
				{153, "NEGQ 5*8(AX)", "48f75828"},
				{154, "NEGL AX", "f7d8"},
				{155, "NEGL DX", "f7da"},
				{156, "NEGL (AX)", "f718"},
				{157, "NEGL 100(BX)", "f75b64"},
				{158, "NEGQ CX", "48f7d9"},
				{159, "NEGQ BX", "48f7db"},
			},
			run: func(asm *Assembler) {
				asm.NegqMem(RSI, 0*8)
				asm.NegqMem(RAX, 5*8)
				asm.NeglReg(RAX)
				asm.NeglReg(RDX)
				asm.NeglMem(RAX, 0)
				asm.NeglMem(RBX, 100)
				asm.NegqReg(RCX)
				asm.NegqReg(RBX)
			},
		},

		{
			name: "testRaw",
			want: []expected{
				{163, "MOVL -16(CX), DX", "8b51f0"},
				{164, "JMP AX", "ffe0"},
				{165, "CMPQ 6*8(SI), $0", "48837e3000"},
			},
			run: func(asm *Assembler) {
				asm.Raw(0x8b, 0x51, 0xf0)
				asm.Raw(0xff, 0xe0)
				asm.Raw(0x48, 0x83, 0x7e, 0x30, 0x00)
			},
		},

		{
			name: "testCall",
			want: []expected{
				{169, "CALL AX", "ffd0"},
				{170, "CALL BX", "ffd3"},
			},
			run: func(asm *Assembler) {
				asm.CallReg(RAX)
				asm.CallReg(RBX)
			},
		},

		{
			name: "testSub",
			want: []expected{
				{174, "SUBL (AX), DI", "2b38"},
				{175, "SUBL 16(SI), AX", "2b4610"},
				{176, "SUBL 640(BX), DX", "2b9380020000"},
			},
			run: func(asm *Assembler) {
				asm.SublMemReg(RAX, RDI, 0)
				asm.SublMemReg(RSI, RAX, 16)
				asm.SublMemReg(RBX, RDX, 640)
			},
		},

		{
			name: "testJgt1",
			want: []expected{
				{180, "JGT forward2", "7f00"},
				{182, "JGT forward1", "7f01"},
				{183, "NOP1", "90"},
				{185, "NOP1", "90"},
			},
			run: func(asm *Assembler) {
				asm.Jgt(2)
				asm.Label(2)
				asm.Jgt(1)
				asm.Nop(1)
				asm.Label(1)
				asm.Nop(1)
			},
		},

		{
			name: "testJgt2",
			want: []expected{
				{190, "NOP1", "90"},
				{191, "JGT l1", "7f03"},
				{193, "NOP1", "90"},
				{194, "JGT l2", "7ffa"},
				{196, "NOP1", "90"},
				{197, "JGT l3", "7ffa"},
			},
			run: func(asm *Assembler) {
				asm.Label(2)
				asm.Nop(1)
				asm.Jgt(1)
				asm.Label(3)
				asm.Nop(1)
				asm.Jgt(2)
				asm.Label(1)
				asm.Nop(1)
				asm.Jgt(3)
			},
		},

		{
			name: "testImul",
			want: []expected{
				{201, "IMULL (AX), CX", "0faf08"},
				{202, "IMULL 4(SI), CX", "0faf4e04"},
				{203, "IMULL -8(DX), AX", "0faf42f8"},
			},
			run: func(asm *Assembler) {
				asm.ImullMemReg(RAX, RCX, 0)
				asm.ImullMemReg(RSI, RCX, 4)
				asm.ImullMemReg(RDX, RAX, -8)
			},
		},

		{
			name: "testMovlqsx",
			want: []expected{
				{207, "MOVLQSX 4(AX), BX", "48635804"},
				{208, "MOVLQSX 8(AX), AX", "48634008"},
			},
			run: func(asm *Assembler) {
				asm.MovlqsxMemReg(RAX, RBX, 4)
				asm.MovlqsxMemReg(RAX, RAX, 8)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			asm := NewAssembler()
			test.run(asm)
			have := fmt.Sprintf("%x", linkToBytes(asm))
			head := have
			for _, v := range test.want {
				if !strings.HasPrefix(head, v.enc) {
					if len(head) >= len(v.enc) {
						head = head[:len(v.enc)]
					}
					t.Fatalf("asmtest.s:%d: %s:\nhave: %s\nwant: %s",
						v.line, v.asm,
						head, v.enc)
				}
				head = head[len(v.enc):]
			}
			if head != "" {
				t.Fatalf("extra trailing bytes: %s", head)
			}
		})
	}
}

func linkToBytes(asm *Assembler) []byte {
	length := asm.Link()
	buf := make([]byte, length)
	asm.Put(buf)
	return buf
}
