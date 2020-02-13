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
			name: "testJmp1",
			want: []expected{
				{4, "JMP forward", "eb00"},
			},
			run: func(asm *Assembler) {
				asm.Jmp(7)
				asm.Label(7)
			},
		},

		{
			name: "testJmp2",
			want: []expected{
				{10, "JMP looping", "ebfe"},
			},
			run: func(asm *Assembler) {
				asm.Label(3)
				asm.Jmp(3)
			},
		},

		{
			name: "testJmp3",
			want: []expected{
				{15, "NOP1", "90"},
				{16, "JMP backward", "ebfd"},
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
				{20, "JMP sharedlabel", "eb04"},
				{21, "NOP1", "90"},
				{22, "JMP sharedlabel", "eb01"},
				{23, "NOP1", "90"},
				{25, "NOP1", "90"},
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
				{29, "NOP1", "90"},
				{30, "JMP l1", "eb03"},
				{32, "NOP1", "90"},
				{33, "JMP l2", "eb03"},
				{35, "NOP1", "90"},
				{36, "JMP l3", "ebfa"},
				{38, "NOP1", "90"},
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
				{44, "NOP1", "90"},
				{45, "JMP l1", "eb03"},
				{47, "NOP1", "90"},
				{48, "JMP l2", "ebfa"},
				{50, "NOP1", "90"},
				{51, "JMP l3", "ebfa"},
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
				{55, "ADDQ $0, 0*8(SI)", "48830600"},
				{56, "ADDQ $1, 0*8(SI)", "48830601"},
				{57, "ADDQ $1, 1*8(SI)", "4883460801"},
				{58, "ADDQ $-1, 3*8(SI)", "48834618ff"},
				{59, "ADDQ $14, 10*8(SI)", "488346500e"},
				{60, "ADDQ $14, 100*8(SI)", "488386200300000e"},
				{61, "ADDQ $0xff, 0*8(SI)", "488106ff000000"},
				{62, "ADDQ $0xff, 1*8(SI)", "48814608ff000000"},
				{63, "ADDQ $-129, 100*8(SI)", "488186200300007fffffff"},
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
