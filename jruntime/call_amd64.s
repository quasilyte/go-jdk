#include "textflag.h"
#include "funcdata.h"

// Go asm has no real rip-based addressing, this is why we
// encode `lea rax, [rip+N]` instruction manually, byte-by-byte.
// N is the length of all instructions in this macro block.
#define JCALL(fnreg) \
        BYTE $0x48; BYTE $0x8d; BYTE $0x05; BYTE $0x09; BYTE $0; BYTE $0; BYTE $0; \
        MOVQ AX, (SI) \
        ADDQ $8, SI \
        JMP fnreg

// $96 bytes for Go call arguments space.
// $16 bytes for 2 pointer arguments.
//
// Register roles:
// AX - tmp register; also used for return value
// BX - <unused> (pin to "this" pointer?)
// CX - <unused>
// DX - env pointer (no need to spill, always in an argument slot)
// SI - stack pointer
// DI - <unused>
//
// func jcall(e *Env, code *byte)
TEXT ·jcall(SB), 0, $96-16
        NO_LOCAL_POINTERS // TODO: Can we do without it?
        MOVQ e+0(FP), DX // env is always at DX
        MOVQ 1*8(DX), SI // stack is always at SI
        MOVQ code+8(FP), CX
        JCALL(CX)
        MOVQ AX, -8(SI) // Save called function result, if any
        RET
