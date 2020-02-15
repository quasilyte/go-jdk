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
				{4, "JGE forward2", "7d00"},
				{6, "JGE forward1", "7d01"},
				{7, "NOP1", "90"},
				{9, "NOP1", "90"},
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
				{14, "NOP1", "90"},
				{15, "JGE l1", "7d03"},
				{17, "NOP1", "90"},
				{18, "JGE l2", "7dfa"},
				{20, "NOP1", "90"},
				{21, "JGE l3", "7dfa"},
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
				{25, "JMP forward", "eb00"},
			},
			run: func(asm *Assembler) {
				asm.Jmp(7)
				asm.Label(7)
			},
		},

		{
			name: "testJmp2",
			want: []expected{
				{31, "JMP looping", "ebfe"},
			},
			run: func(asm *Assembler) {
				asm.Label(3)
				asm.Jmp(3)
			},
		},

		{
			name: "testJmp3",
			want: []expected{
				{36, "NOP1", "90"},
				{37, "JMP backward", "ebfd"},
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
				{41, "JMP sharedlabel", "eb04"},
				{42, "NOP1", "90"},
				{43, "JMP sharedlabel", "eb01"},
				{44, "NOP1", "90"},
				{46, "NOP1", "90"},
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
				{50, "NOP1", "90"},
				{51, "JMP l1", "eb03"},
				{53, "NOP1", "90"},
				{54, "JMP l2", "eb03"},
				{56, "NOP1", "90"},
				{57, "JMP l3", "ebfa"},
				{59, "NOP1", "90"},
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
				{64, "NOP1", "90"},
				{65, "JMP l1", "eb03"},
				{67, "NOP1", "90"},
				{68, "JMP l2", "ebfa"},
				{70, "NOP1", "90"},
				{71, "JMP l3", "ebfa"},
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
				{75, "JMP AX", "ffe0"},
				{76, "JMP DX", "ffe2"},
				{77, "JMP CX", "ffe1"},
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
				{81, "JMP (AX)", "ff20"},
				{82, "JMP -8(DI)", "ff67f8"},
				{83, "JMP 13935(CX)", "ffa16f360000"},
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
				{87, "ADDL (AX), DX", "0310"},
				{88, "ADDL 8(SI), AX", "034608"},
				{89, "ADDL $7, (AX)", "830007"},
				{90, "ADDL $-9, -8(DX)", "8342f8f7"},
				{91, "ADDL $127, CX", "83c17f"},
				{92, "ADDL $-128, BX", "83c380"},
				{93, "ADDQ $0, 0*8(SI)", "48830600"},
				{94, "ADDQ $1, 0*8(SI)", "48830601"},
				{95, "ADDQ $1, 1*8(SI)", "4883460801"},
				{96, "ADDQ $-1, 3*8(SI)", "48834618ff"},
				{97, "ADDQ $14, 10*8(SI)", "488346500e"},
				{98, "ADDQ $14, 100*8(SI)", "488386200300000e"},
				{99, "ADDQ $0xff, 0*8(SI)", "488106ff000000"},
				{100, "ADDQ $0xff, 1*8(SI)", "48814608ff000000"},
				{101, "ADDQ $-129, 100*8(SI)", "488186200300007fffffff"},
				{102, "ADDQ $1, AX", "4883c001"},
				{103, "ADDQ $-1, DI", "4883c7ff"},
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
				{107, "MOVL $0, 0*8(SI)", "c70600000000"},
				{108, "MOVL $1, 0*8(DI)", "c70701000000"},
				{109, "MOVL $1, 1*8(AX)", "c7400801000000"},
				{110, "MOVL $-50000, 40*8(SI)", "c78640010000b03cffff"},
				{111, "MOVL (AX), AX", "8b00"},
				{112, "MOVL -16(CX), DX", "8b51f0"},
				{113, "MOVL AX, (AX)", "8900"},
				{114, "MOVL DX, -16(CX)", "8951f0"},
				{115, "MOVQ 0*8(AX), BX", "488b18"},
				{116, "MOVQ 16*8(BX), AX", "488b8380000000"},
				{117, "MOVQ AX, 0*8(DI)", "488907"},
				{118, "MOVQ DX, 3*8(DI)", "48895718"},
				{119, "MOVQ AX, 0*8(AX)", "488900"},
				{120, "MOVQ $140038723203072, AX", "48b800f0594e5d7f0000"},
				{121, "MOVQ $9223372036854775807, DX", "48baffffffffffffff7f"},
				{122, "MOVQ $-9223372036854775800, SI", "48be0800000000000080"},
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
				{126, "CMPL AX, 0*8(DI)", "3b07"},
				{127, "CMPL BX, 1*8(AX)", "3b5808"},
				{128, "CMPL 16(SI), $0", "837e1000"},
				{129, "CMPL (AX), $15", "83380f"},
				{130, "CMPQ 6*8(SI), $0", "48837e3000"},
			},
			run: func(asm *Assembler) {
				asm.CmplRegMem(RAX, RDI, 0*8)
				asm.CmplRegMem(RBX, RAX, 1*8)
				asm.CmplConst8Mem(0, RSI, 16)
				asm.CmplConst8Mem(15, RAX, 0)
				asm.CmpqConst8Mem(0, RSI, 6*8)
			},
		},

		{
			name: "testNeg",
			want: []expected{
				{134, "NEGQ 0*8(SI)", "48f71e"},
				{135, "NEGQ 5*8(AX)", "48f75828"},
				{136, "NEGL AX", "f7d8"},
				{137, "NEGL DX", "f7da"},
				{138, "NEGL (AX)", "f718"},
				{139, "NEGL 100(BX)", "f75b64"},
				{140, "NEGQ CX", "48f7d9"},
				{141, "NEGQ BX", "48f7db"},
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			asm := NewAssembler()
			test.run(asm)
			have := fmt.Sprintf("%x", asm.Link())
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
			if len(head) != 0 {
				t.Fatalf("extra trailing bytes: %s", head)
			}
		})
	}
}
