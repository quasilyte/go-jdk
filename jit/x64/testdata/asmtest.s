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

TEXT testJmpReg(SB), 0, $0-0
        JMP AX // asm.JmpReg(RAX)
        JMP DX // asm.JmpReg(RDX)
        JMP CX // asm.JmpReg(RCX)
        RET

TEXT testJmpMem(SB), 0, $0-0
        JMP (AX) // asm.JmpMem(RAX, 0)
        JMP -8(DI) // asm.JmpMem(RDI, -8)
        JMP 13935(CX) // asm.JmpMem(RCX, 13935)
        RET

TEXT testAdd(SB), 0, $0-0
        ADDQ $0, 0*8(SI) // asm.AddqConst8Mem(0, RSI, 0*8)
        ADDQ $1, 0*8(SI) // asm.AddqConst8Mem(1, RSI, 0*8)
        ADDQ $1, 1*8(SI) // asm.AddqConst8Mem(1, RSI, 1*8)
        ADDQ $-1, 3*8(SI) // asm.AddqConst8Mem(-1, RSI, 3*8)
        ADDQ $14, 10*8(SI) // asm.AddqConst8Mem(14, RSI, 10*8)
        ADDQ $14, 100*8(SI) // asm.AddqConst8Mem(14, RSI, 100*8)
        ADDQ $0xff, 0*8(SI) // asm.AddqConst32Mem(0xff, RSI, 0*8)
        ADDQ $0xff, 1*8(SI) // asm.AddqConst32Mem(0xff, RSI, 1*8)
        ADDQ $-129, 100*8(SI) // asm.AddqConst32Mem(-129, RSI, 100*8)
        ADDQ $1, AX // asm.AddqConst8Reg(1, RAX)
        ADDQ $-1, DI // asm.AddqConst8Reg(-1, RDI)
        RET

TEXT testMov(SB), 0, $0-0
        MOVL $0, 0*8(SI) // asm.MovlConst32Mem(0, RSI, 0*8)
        MOVL $1, 0*8(DI) // asm.MovlConst32Mem(1, RDI, 0*8)
        MOVL $1, 1*8(AX) // asm.MovlConst32Mem(1, RAX, 1*8)
        MOVL $-50000, 40*8(SI) // asm.MovlConst32Mem(-50000, RSI, 40*8)
        MOVQ 0*8(AX), BX // asm.MovqMemReg(RAX, RBX, 0*8)
        MOVQ 16*8(BX), AX // asm.MovqMemReg(RBX, RAX, 16*8)
        MOVQ AX, 0*8(DI) // asm.MovqRegMem(RAX, RDI, 0*8)
        MOVQ DX, 3*8(DI) // asm.MovqRegMem(RDX, RDI, 3*8)
        MOVQ AX, 0*8(AX) // asm.MovqRegMem(RAX, RAX, 0*8)
        MOVQ $140038723203072, AX // asm.MovqConst64Reg(140038723203072, RAX)
        MOVQ $9223372036854775807, DX // asm.MovqConst64Reg(9223372036854775807, RDX)
        MOVQ $-9223372036854775800, SI // asm.MovqConst64Reg(-9223372036854775800, RSI)
        RET

TEXT testCmp(SB), 0, $0-0
        CMPL AX, 0*8(DI) // asm.CmplRegMem(RAX, RDI, 0*8)
        CMPL BX, 1*8(AX) // asm.CmplRegMem(RBX, RAX, 1*8)
        RET
