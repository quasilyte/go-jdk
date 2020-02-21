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
				{92, "ADDL $127, CX", "83c17f"},
				{93, "ADDL $-128, BX", "83c380"},
				{94, "ADDQ $0, 0*8(SI)", "48830600"},
				{95, "ADDQ $1, 0*8(SI)", "48830601"},
				{96, "ADDQ $1, 1*8(SI)", "4883460801"},
				{97, "ADDQ $-1, 3*8(SI)", "48834618ff"},
				{98, "ADDQ $14, 10*8(SI)", "488346500e"},
				{99, "ADDQ $14, 100*8(SI)", "488386200300000e"},
				{100, "ADDQ $0xff, 0*8(SI)", "488106ff000000"},
				{101, "ADDQ $0xff, 1*8(SI)", "48814608ff000000"},
				{102, "ADDQ $-129, 100*8(SI)", "488186200300007fffffff"},
				{103, "ADDQ $1, AX", "4883c001"},
				{104, "ADDQ $-1, DI", "4883c7ff"},
			},
			run: func(asm *Assembler) {
				asm.AddlMemReg(RAX, RDX, 0)
				asm.AddlMemReg(RSI, RAX, 8)
				asm.AddlConst8Mem(7, RAX, 0)
				asm.AddlConst8Mem(-9, RDX, -8)
				asm.AddlConst8Reg(127, RCX)
				asm.AddlConst8Reg(-128, RBX)
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
			},
		},

		{
			name: "testMov",
			want: []expected{
				{108, "MOVL $0, 0*8(SI)", "c70600000000"},
				{109, "MOVL $1, 0*8(DI)", "c70701000000"},
				{110, "MOVL $1, 1*8(AX)", "c7400801000000"},
				{111, "MOVL $-50000, 40*8(SI)", "c78640010000b03cffff"},
				{112, "MOVL (AX), AX", "8b00"},
				{113, "MOVL -16(CX), DX", "8b51f0"},
				{114, "MOVL AX, (AX)", "8900"},
				{115, "MOVL DX, -16(CX)", "8951f0"},
				{116, "MOVL $1355, AX", "b84b050000"},
				{117, "MOVL $-6643, DX", "ba0de6ffff"},
				{118, "MOVQ 0*8(AX), BX", "488b18"},
				{119, "MOVQ 16*8(BX), AX", "488b8380000000"},
				{120, "MOVQ AX, 0*8(DI)", "488907"},
				{121, "MOVQ DX, 3*8(DI)", "48895718"},
				{122, "MOVQ AX, 0*8(AX)", "488900"},
				{123, "MOVQ $140038723203072, AX", "48b800f0594e5d7f0000"},
				{124, "MOVQ $9223372036854775807, DX", "48baffffffffffffff7f"},
				{125, "MOVQ $-9223372036854775800, SI", "48be0800000000000080"},
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
			},
		},

		{
			name: "testCmp",
			want: []expected{
				{129, "CMPL AX, 0*8(DI)", "3b07"},
				{130, "CMPL BX, 1*8(AX)", "3b5808"},
				{131, "CMPL 16(SI), $0", "837e1000"},
				{132, "CMPL (AX), $15", "83380f"},
				{133, "CMPL (DI), $242", "813ff2000000"},
				{134, "CMPL -8(BX), $-5343", "817bf821ebffff"},
				{135, "CMPQ 6*8(SI), $0", "48837e3000"},
			},
			run: func(asm *Assembler) {
				asm.CmplRegMem(RAX, RDI, 0*8)
				asm.CmplRegMem(RBX, RAX, 1*8)
				asm.CmplConst8Mem(0, RSI, 16)
				asm.CmplConst8Mem(15, RAX, 0)
				asm.CmplConst32Mem(242, RDI, 0)
				asm.CmplConst32Mem(-5343, RBX, -8)
				asm.CmpqConst8Mem(0, RSI, 6*8)
			},
		},

		{
			name: "testNeg",
			want: []expected{
				{139, "NEGQ 0*8(SI)", "48f71e"},
				{140, "NEGQ 5*8(AX)", "48f75828"},
				{141, "NEGL AX", "f7d8"},
				{142, "NEGL DX", "f7da"},
				{143, "NEGL (AX)", "f718"},
				{144, "NEGL 100(BX)", "f75b64"},
				{145, "NEGQ CX", "48f7d9"},
				{146, "NEGQ BX", "48f7db"},
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
				{150, "MOVL -16(CX), DX", "8b51f0"},
				{151, "JMP AX", "ffe0"},
				{152, "CMPQ 6*8(SI), $0", "48837e3000"},
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
				{156, "CALL AX", "ffd0"},
				{157, "CALL BX", "ffd3"},
			},
			run: func(asm *Assembler) {
				asm.CallReg(RAX)
				asm.CallReg(RBX)
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
