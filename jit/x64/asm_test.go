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
			name: "testAddqConst",
			want: []expected{
				{75, "ADDQ $0, 0*8(SI)", "48830600"},
				{76, "ADDQ $1, 0*8(SI)", "48830601"},
				{77, "ADDQ $1, 1*8(SI)", "4883460801"},
				{78, "ADDQ $-1, 3*8(SI)", "48834618ff"},
				{79, "ADDQ $14, 10*8(SI)", "488346500e"},
				{80, "ADDQ $14, 100*8(SI)", "488386200300000e"},
				{81, "ADDQ $0xff, 0*8(SI)", "488106ff000000"},
				{82, "ADDQ $0xff, 1*8(SI)", "48814608ff000000"},
				{83, "ADDQ $-129, 100*8(SI)", "488186200300007fffffff"},
			},
			run: func(asm *Assembler) {
				asm.AddqConst(0, 0)
				asm.AddqConst(1, 0)
				asm.AddqConst(1, 1)
				asm.AddqConst(-1, 3)
				asm.AddqConst(14, 10)
				asm.AddqConst(14, 100)
				asm.AddqConst(0xff, 0)
				asm.AddqConst(0xff, 1)
				asm.AddqConst(-129, 100)
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
