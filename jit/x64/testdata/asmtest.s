#define NOP1 BYTE $0x90
#define NOSPLIT 4

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
        ADDL (AX), DX // asm.AddlMemReg(RAX, RDX, 0)
        ADDL 8(SI), AX // asm.AddlMemReg(RSI, RAX, 8)
        ADDL $7, (AX) // asm.AddlConst8Mem(7, RAX, 0)
        ADDL $-9, -8(DX) // asm.AddlConst8Mem(-9, RDX, -8)
        ADDL $9300, 16(SI) // asm.AddlConst32Mem(9300, RSI, 16)
        ADDL $127, CX // asm.AddlConst8Reg(127, RCX)
        ADDL $-128, BX // asm.AddlConst8Reg(-128, RBX)
        ADDL $200, BP // asm.AddlConst32Reg(200, RBP)
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
        ADDQ $5000, SP // asm.AddqConst32Reg(5000, RSP)
        RET

TEXT testMov(SB), 0, $0-0
        MOVL $0, 0*8(SI) // asm.MovlConst32Mem(0, RSI, 0*8)
        MOVL $1, 0*8(DI) // asm.MovlConst32Mem(1, RDI, 0*8)
        MOVL $1, 1*8(AX) // asm.MovlConst32Mem(1, RAX, 1*8)
        MOVL $-50000, 40*8(SI) // asm.MovlConst32Mem(-50000, RSI, 40*8)
        MOVL (AX), AX // asm.MovlMemReg(RAX, RAX, 0)
        MOVL -16(CX), DX // asm.MovlMemReg(RCX, RDX, -16)
        MOVL AX, (AX) // asm.MovlRegMem(RAX, RAX, 0)
        MOVL DX, -16(CX) // asm.MovlRegMem(RDX, RCX, -16)
        MOVL $1355, AX // asm.MovlConst32Reg(1355, RAX)
        MOVL $-6643, DX // asm.MovlConst32Reg(-6643, RDX)
        MOVQ 0*8(AX), BX // asm.MovqMemReg(RAX, RBX, 0*8)
        MOVQ 16*8(BX), AX // asm.MovqMemReg(RBX, RAX, 16*8)
        MOVQ AX, 0*8(DI) // asm.MovqRegMem(RAX, RDI, 0*8)
        MOVQ DX, 3*8(DI) // asm.MovqRegMem(RDX, RDI, 3*8)
        MOVQ AX, 0*8(AX) // asm.MovqRegMem(RAX, RAX, 0*8)
        MOVQ $140038723203072, AX // asm.MovqConst64Reg(140038723203072, RAX)
        MOVQ $9223372036854775807, DX // asm.MovqConst64Reg(9223372036854775807, RDX)
        MOVQ $-9223372036854775800, SI // asm.MovqConst64Reg(-9223372036854775800, RSI)
        MOVQ $1423, AX // asm.MovqConst32Reg(1423, RAX)
        MOVQ $-23, CX // asm.MovqConst32Reg(-23, RCX)
        MOVQ $1, DX // asm.MovqConst32Reg(1, RDX)
        MOVQ 100(BP), DX // asm.MovqMemReg(RBP, RDX, 100)
        MOVQ $1, 1(AX) // asm.MovqConst32Mem(1, RAX, 1)
        MOVQ $-1, 2(AX) // asm.MovqConst32Mem(-1, RAX, 2)
        MOVQ $0, -96(BP) // asm.MovqConst32Mem(0, RBP, -96)
        MOVQ $100, -96(BP) // asm.MovqConst32Mem(100, RBP, -96)
        RET

TEXT testCmp(SB), 0, $0-0
        CMPL AX, 0*8(DI)    // asm.CmplRegMem(RAX, RDI, 0*8)
        CMPL BX, 1*8(AX)    // asm.CmplRegMem(RBX, RAX, 1*8)
        CMPL 16(SI), $0     // asm.CmplConst8Mem(0, RSI, 16)
        CMPL (AX), $15      // asm.CmplConst8Mem(15, RAX, 0)
        CMPL (DI), $242     // asm.CmplConst32Mem(242, RDI, 0)
        CMPL -8(BX), $-5343 // asm.CmplConst32Mem(-5343, RBX, -8)
        CMPQ 6*8(SI), $0    // asm.CmpqConst8Mem(0, RSI, 6*8)
        CMPQ (SI), $999     // asm.CmpqConst32Mem(999, RSI, 0)
        CMPQ 8(DI), $-999   // asm.CmpqConst32Mem(-999, RDI, 8)
        RET

TEXT testNeg(SB), 0, $0-0
        NEGQ 0*8(SI) // asm.NegqMem(RSI, 0*8)
        NEGQ 5*8(AX) // asm.NegqMem(RAX, 5*8)
        NEGL AX // asm.NeglReg(RAX)
        NEGL DX // asm.NeglReg(RDX)
        NEGL (AX) // asm.NeglMem(RAX, 0)
        NEGL 100(BX) // asm.NeglMem(RBX, 100)
        NEGQ CX // asm.NegqReg(RCX)
        NEGQ BX // asm.NegqReg(RBX)
        RET

TEXT testRaw(SB), 0, $0-0
        MOVL -16(CX), DX // asm.Raw(0x8b, 0x51, 0xf0)
        JMP AX // asm.Raw(0xff, 0xe0)
        CMPQ 6*8(SI), $0 // asm.Raw(0x48, 0x83, 0x7e, 0x30, 0x00)
        RET

TEXT testCall(SB), NOSPLIT, $0-0
        CALL AX // asm.CallReg(RAX)
        CALL BX // asm.CallReg(RBX)
        RET

TEXT testSub(SB), 0, $0-0
        SUBL (AX), DI // asm.SublMemReg(RAX, RDI, 0)
        SUBL 16(SI), AX // asm.SublMemReg(RSI, RAX, 16)
        SUBL 640(BX), DX // asm.SublMemReg(RBX, RDX, 640)
        RET

TEXT testJgt1(SB), 0, $0-0
        JGT forward2 // asm.Jgt(2)
forward2:
        JGT forward1 // asm.Label(2); asm.Jgt(1)
        NOP1         // asm.Nop(1)
forward1:
        NOP1 // asm.Label(1); asm.Nop(1)
        RET

TEXT testJgt2(SB), 0, $0-0
l2:
        NOP1   // asm.Label(2); asm.Nop(1)
        JGT l1 // asm.Jgt(1)
l3:
        NOP1   // asm.Label(3); asm.Nop(1)
        JGT l2 // asm.Jgt(2)
l1:
        NOP1   // asm.Label(1); asm.Nop(1)
        JGT l3 // asm.Jgt(3)
        RET

TEXT testImul(SB), 0, $0-0
        IMULL (AX), CX // asm.ImullMemReg(RAX, RCX, 0)
        IMULL 4(SI), CX // asm.ImullMemReg(RSI, RCX, 4)
        IMULL -8(DX), AX // asm.ImullMemReg(RDX, RAX, -8)
        RET

TEXT testMovlqsx(SB), 0, $0-0
        MOVLQSX 4(AX), BX // asm.MovlqsxMemReg(RAX, RBX, 4)
        MOVLQSX 8(AX), AX // asm.MovlqsxMemReg(RAX, RAX, 8)
        RET
