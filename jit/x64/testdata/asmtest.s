#define NOP1 BYTE $0x90

TEXT testJge1(SB), 0, $0-0
        JGE forward2 // asm.Jge(2)
forward2:
        JGE forward1 // asm.Label(2); asm.Jge(1)
        NOP1         // asm.Nop(1)
forward1:
        NOP1 // asm.Label(1); asm.Nop(1)
        RET

TEXT testJge2(SB), 0, $0-0
l2:
        NOP1   // asm.Label(2); asm.Nop(1)
        JGE l1 // asm.Jge(1)
l3:
        NOP1   // asm.Label(3); asm.Nop(1)
        JGE l2 // asm.Jge(2)
l1:
        NOP1   // asm.Label(1); asm.Nop(1)
        JGE l3 // asm.Jge(3)
        RET

TEXT testJmp1(SB), 0, $0-0
        JMP forward // asm.Jmp(7); asm.Label(7)
forward:
        RET

TEXT testJmp2(SB), 0, $0-0
looping:
        JMP looping // asm.Label(3); asm.Jmp(3)
        RET

TEXT testJmp3(SB), 0, $0-0
backward:
        NOP1 // asm.Label(0); asm.Nop(1)
        JMP backward // asm.Jmp(0)
        RET

TEXT testJmp4(SB), 0, $0-0
        JMP sharedlabel // asm.Jmp(0)
        NOP1 // asm.Nop(1)
        JMP sharedlabel // asm.Jmp(0)
        NOP1 // asm.Nop(1)
sharedlabel:
        NOP1 // asm.Label(0); asm.Nop(1)
        RET

TEXT testJmp5(SB), 0, $0-0
        NOP1   // asm.Nop(1)
        JMP l1 // asm.Jmp(1)
l3:
        NOP1   // asm.Label(3); asm.Nop(1)
        JMP l2 // asm.Jmp(2)
l1:
        NOP1   // asm.Label(1); asm.Nop(1)
        JMP l3 // asm.Jmp(3)
l2:
        NOP1 // asm.Label(2); asm.Nop(1)
        RET

TEXT testJmp6(SB), 0, $0-0
l2:
        NOP1   // asm.Label(2); asm.Nop(1)
        JMP l1 // asm.Jmp(1)
l3:
        NOP1   // asm.Label(3); asm.Nop(1)
        JMP l2 // asm.Jmp(2)
l1:
        NOP1   // asm.Label(1); asm.Nop(1)
        JMP l3 // asm.Jmp(3)
        RET

TEXT testAddqConst(SB), 0, $0-0
        ADDQ $0, 0*8(SI) // asm.AddqConst(0, 0)
        ADDQ $1, 0*8(SI) // asm.AddqConst(1, 0)
        ADDQ $1, 1*8(SI) // asm.AddqConst(1, 1)
        ADDQ $-1, 3*8(SI) // asm.AddqConst(-1, 3)
        ADDQ $14, 10*8(SI) // asm.AddqConst(14, 10)
        ADDQ $14, 100*8(SI) // asm.AddqConst(14, 100)
        ADDQ $0xff, 0*8(SI) // asm.AddqConst(0xff, 0)
        ADDQ $0xff, 1*8(SI) // asm.AddqConst(0xff, 1)
        ADDQ $-129, 100*8(SI) // asm.AddqConst(-129, 100)
        RET
