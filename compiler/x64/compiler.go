package x64

import (
	"fmt"

	"github.com/quasilyte/GopherJRE/codegen"
	"github.com/quasilyte/GopherJRE/ir"
	"github.com/quasilyte/GopherJRE/jit/x64"
)

// Compiler implements a class compiler for x86-64 (amd64) architecture.
type Compiler struct {
	asm *x64.Assembler

	result [][]byte
}

// NewCompiler returns a compiler for amd64 platform.
func NewCompiler() *Compiler {
	return &Compiler{
		asm:    x64.NewAssembler(),
		result: make([][]byte, 0, 40),
	}
}

func (cl *Compiler) CompileMethods(methods []ir.Method) (codegen.Result, error) {
	var res codegen.Result
	res.Methods = make([][]byte, len(methods))

	for i := range methods {
		m := &methods[i]
		code, err := cl.compileMethod(m)
		if err != nil {
			return res, fmt.Errorf("compile %s: %v", m.Name, err)
		}
		res.Methods[i] = code
	}

	return res, nil
}

func (cl *Compiler) compileMethod(m *ir.Method) ([]byte, error) {
	asm := cl.asm
	asm.Reset()

	for i, inst := range m.Code {
		var a1, a2 ir.Arg
		dst := inst.Dst
		if len(inst.Args) > 0 {
			a1 = inst.Args[0]
		}
		if len(inst.Args) > 1 {
			a2 = inst.Args[1]
		}

		if inst.Flags.IsJumpTarget() {
			asm.Label(int64(i))
		}

		switch inst.Kind {
		case ir.InstLoad:
			switch a1.Kind {
			case ir.ArgIntConst:
				asm.MovlConst32Mem(int32(a1.Value), x64.RSI, int32(dst.Value*8))
			}

		case ir.InstIneg:
			if a1 == dst {
				asm.NeglMem(x64.RSI, int32(a1.Value*8))
			} else {
				asm.MovlMemReg(x64.RSI, x64.RAX, int32(a1.Value*8))
				asm.NeglReg(x64.RAX)
				asm.MovlRegMem(x64.RAX, x64.RSI, int32(dst.Value*8))
			}
		case ir.InstLneg:
			if a1 == dst {
				asm.NegqMem(x64.RSI, int32(a1.Value*8))
			} else {
				asm.MovqMemReg(x64.RSI, x64.RAX, int32(a1.Value*8))
				asm.NegqReg(x64.RAX)
				asm.MovqRegMem(x64.RAX, x64.RSI, int32(dst.Value*8))
			}

		case ir.InstJumpGtEq:
			asm.Jge(a1.Value)
		case ir.InstJump:
			asm.Jmp(a1.Value)

		case ir.InstIcmp:
			switch {
			case a1.Kind == ir.ArgReg && a2.Kind == ir.ArgIntConst:
				asm.CmplConst8Mem(int8(a1.Value), x64.RSI, int32(a1.Value*8))
			case a1.Kind == ir.ArgReg && a2.Kind == ir.ArgReg:
				asm.MovlMemReg(x64.RSI, x64.RAX, int32(a1.Value*8))
				asm.CmplRegMem(x64.RAX, x64.RSI, int32(a2.Value*8))
			}
		case ir.InstLcmp:
			switch {
			case a1.Kind == ir.ArgReg && a2.Kind == ir.ArgIntConst:
				asm.CmpqConst8Mem(int8(a1.Value), x64.RSI, int32(a1.Value*8))
			}

		case ir.InstIret:
			asm.MovlMemReg(x64.RSI, x64.RAX, int32(a1.Value*8))
			asm.JmpMem(x64.RSI, -8)
		case ir.InstLret:
			asm.MovqMemReg(x64.RSI, x64.RAX, int32(a1.Value*8))
			asm.JmpMem(x64.RSI, -8)

		case ir.InstRet:
			asm.JmpMem(x64.RSI, -8)

		case ir.InstIadd:
			if a1 == dst {
				if a2.Kind == ir.ArgIntConst {
					asm.AddqConst8Mem(int8(a1.Value), x64.RSI, int32(dst.Value*8))
				}
			} else {
				if a2.Kind == ir.ArgReg {
					asm.MovlMemReg(x64.RSI, x64.RAX, int32(a1.Value*8))
					asm.AddlMemReg(x64.RSI, x64.RAX, int32(a2.Value*8))
					asm.MovlRegMem(x64.RAX, x64.RSI, int32(dst.Value*8))
				}
			}

		default:
			return nil, fmt.Errorf("can't assemble %s", inst)
		}
	}

	return asm.Link(), nil
}
