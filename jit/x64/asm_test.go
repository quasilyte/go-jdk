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
			name: "testAdd",
			want: []expected{
				{81, "ADDQ $0, 0*8(SI)", "48830600"},
				{82, "ADDQ $1, 0*8(SI)", "48830601"},
				{83, "ADDQ $1, 1*8(SI)", "4883460801"},
				{84, "ADDQ $-1, 3*8(SI)", "48834618ff"},
				{85, "ADDQ $14, 10*8(SI)", "488346500e"},
				{86, "ADDQ $14, 100*8(SI)", "488386200300000e"},
				{87, "ADDQ $0xff, 0*8(SI)", "488106ff000000"},
				{88, "ADDQ $0xff, 1*8(SI)", "48814608ff000000"},
				{89, "ADDQ $-129, 100*8(SI)", "488186200300007fffffff"},
				{90, "ADDQ $1, AX", "4883c001"},
				{91, "ADDQ $-1, DI", "4883c7ff"},
			},
			run: func(asm *Assembler) {
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
				{95, "MOVL $0, 0*8(SI)", "c70600000000"},
				{96, "MOVL $1, 0*8(DI)", "c70701000000"},
				{97, "MOVL $1, 1*8(AX)", "c7400801000000"},
				{98, "MOVL $-50000, 40*8(SI)", "c78640010000b03cffff"},
				{99, "MOVQ 0*8(AX), BX", "488b18"},
				{100, "MOVQ 16*8(BX), AX", "488b8380000000"},
				{101, "MOVQ AX, 0*8(DI)", "488907"},
				{102, "MOVQ DX, 3*8(DI)", "48895718"},
				{103, "MOVQ AX, 0*8(AX)", "488900"},
				{104, "MOVQ $140038723203072, AX", "48b800f0594e5d7f0000"},
				{105, "MOVQ $9223372036854775807, DX", "48baffffffffffffff7f"},
				{106, "MOVQ $-9223372036854775800, SI", "48be0800000000000080"},
			},
			run: func(asm *Assembler) {
				asm.MovlConst32Mem(0, RSI, 0*8)
				asm.MovlConst32Mem(1, RDI, 0*8)
				asm.MovlConst32Mem(1, RAX, 1*8)
				asm.MovlConst32Mem(-50000, RSI, 40*8)
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
				{110, "CMPL AX, 0*8(DI)", "3b07"},
				{111, "CMPL BX, 1*8(AX)", "3b5808"},
			},
			run: func(asm *Assembler) {
				asm.CmplRegMem(RAX, RDI, 0*8)
				asm.CmplRegMem(RBX, RAX, 1*8)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			asm := NewAssembler()
			test.run(asm)
			have := fmt.Sprintf("%x", asm.Link())
			var want strings.Builder
			for _, v := range test.want {
				want.WriteString(v.enc)
			}
			head := have
			for _, v := range test.want {
				if !strings.HasPrefix(head, v.enc) {
					t.Fatalf("asmtest.s:%d: %s:\nexpected: %s\nhead: %s\nhave: %s\nwant: %s",
						v.line, v.asm,
						v.enc, head,
						have, want.String())
				}
				head = head[len(v.enc):]
			}
			if len(head) != 0 {
				t.Fatalf("extra trailing bytes: %s", head)
			}
		})
	}
}
