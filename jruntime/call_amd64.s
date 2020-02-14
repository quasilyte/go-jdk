#include "textflag.h"

// func jcall(e *Env, code *byte)
TEXT ·jcall(SB), 0, $8-16
        MOVQ e+0(FP), DX // env is always at DX
        MOVQ 1*8(DX), SI // stack is always at SI
        MOVQ code+8(FP), CX
        CALL callwrapper(SB)
        ADDQ $8, SP     // Because callwrapper has no RET
        MOVQ AX, -8(SI) // Save called function result, if any
        RET

// callwrapper captures return address written by a CALL
// and jumps into CX.
TEXT callwrapper(SB), 0, $0-0
        MOVQ (SP), AX
        MOVQ AX, (SI)
        ADDQ $8, SI
        JMP CX
