package x64

import (
	"fmt"
	"unsafe"

	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/jclass"
	"github.com/quasilyte/go-jdk/jit"
	"github.com/quasilyte/go-jdk/jit/x64"
	"github.com/quasilyte/go-jdk/symbol"
	"github.com/quasilyte/go-jdk/vmdat"
)

// Compiler implements a class compiler for x86-64 (amd64) architecture.
type Compiler struct {
	ctx jit.Context

	packageID uint64
	classID   uint64
	methodID  uint64

	asm          *x64.Assembler
	relocs       []relocation
	methodRelocs int
	method       *ir.Method
}

type relocation struct {
	sourceID     symbol.ID
	targetID     symbol.ID
	targetOffset int
}

// NewCompiler returns a compiler for amd64 platform.
func NewCompiler() *Compiler {
	return &Compiler{
		asm:    x64.NewAssembler(),
		relocs: make([]relocation, 0, 64),
	}
}

func (cl *Compiler) Compile(ctx jit.Context, packages []*ir.Package) error {
	cl.ctx = ctx
	cl.relocs = cl.relocs[:0]
	for _, p := range packages {
		cl.packageID = uint64(p.Out.ID)
		if err := cl.compilePackage(p); err != nil {
			return fmt.Errorf("package %s: %v", p.Out.Name, err)
		}
	}
	cl.link()
	return nil
}

func (cl *Compiler) link() {
	for _, rel := range cl.relocs {
		dstMethod := cl.getMethodByID(rel.targetID)
		srcMethod := cl.getMethodByID(rel.sourceID)
		dst := (*int64)(unsafe.Pointer(&dstMethod.Code[rel.targetOffset]))
		*dst = int64(uintptr(unsafe.Pointer(&srcMethod.Code[0])))
	}
}

func (cl *Compiler) compilePackage(p *ir.Package) error {
	for i := range p.Classes {
		cl.classID = uint64(i)
		c := &p.Classes[i]
		for j := range c.Methods {
			cl.methodID = uint64(j)
			m := &c.Methods[j]
			if err := cl.compileMethod(m); err != nil {
				// TODO: print descriptor in a more pretty way.
				return fmt.Errorf("%s.%s%s: %v",
					c.Name, m.Out.Name, m.Out.Descriptor, err)
			}
		}
	}
	return nil
}

func (cl *Compiler) compileMethod(m *ir.Method) error {
	if len(m.Code) == 0 {
		return nil
	}

	cl.asm.Reset()
	cl.methodRelocs = 0
	cl.method = m

	for i, inst := range m.Code {
		if inst.Flags.IsJumpTarget() {
			cl.asm.Label(int64(i))
		}
		if !cl.assembleInst(inst) {
			return fmt.Errorf("can't assemble: %s", inst)
		}
	}

	length := cl.asm.Link()
	if length == 0 {
		return fmt.Errorf("no machine code is generated")
	}
	code, err := cl.ctx.Mmap.AllocateExecutable(length)
	if err != nil {
		return fmt.Errorf("mmap(%d): %v", length, err)
	}
	cl.asm.Put(code)
	m.Out.Code = code

	relocs := cl.relocs[len(cl.relocs)-cl.methodRelocs:]
	for i := range relocs {
		rel := &relocs[i]
		const mov64width = 2
		rel.targetOffset = int(cl.asm.OffsetOf(rel.targetOffset) + mov64width)
	}

	return nil
}

func (cl *Compiler) assembleInst(inst ir.Inst) bool {
	asm := cl.asm

	var a1, a2 ir.Arg
	dst := inst.Dst
	if len(inst.Args) > 0 {
		a1 = inst.Args[0]
	}
	if len(inst.Args) > 1 {
		a2 = inst.Args[1]
	}

	switch inst.Kind {
	case ir.InstJumpGtEq:
		asm.Jge(a1.Value)
	case ir.InstJump:
		asm.Jmp(a1.Value)

	case ir.InstLload:
		switch a1.Kind {
		case ir.ArgReg:
			asm.MovqMemReg(x64.RSI, x64.RAX, regDisp(a1))
			asm.MovqRegMem(x64.RAX, x64.RSI, regDisp(dst))
		case ir.ArgIntConst:
			if !fits32bit(a1.Value) {
				return false
			}
			asm.MovqConst32Mem(int32(a1.Value), x64.RSI, regDisp(dst))
		default:
			return false
		}

	case ir.InstIload:
		switch a1.Kind {
		case ir.ArgReg:
			asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
			asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
		case ir.ArgIntConst:
			asm.MovlConstMem(a1.Value, x64.RSI, regDisp(dst))
		default:
			return false
		}

	case ir.InstIneg:
		if a1 == dst {
			asm.NeglMem(x64.RSI, regDisp(a1))
		} else {
			asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
			asm.NeglReg(x64.RAX)
			asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
		}
	case ir.InstLneg:
		if a1 == dst {
			asm.NegqMem(x64.RSI, regDisp(a1))
		} else {
			asm.MovqMemReg(x64.RSI, x64.RAX, regDisp(a1))
			asm.NegqReg(x64.RAX)
			asm.MovqRegMem(x64.RAX, x64.RSI, regDisp(dst))
		}

	case ir.InstIcmp:
		switch {
		case a1.Kind == ir.ArgReg && a2.Kind == ir.ArgIntConst:
			asm.CmplConstMem(a2.Value, x64.RSI, regDisp(a1))
		case a1.Kind == ir.ArgReg && a2.Kind == ir.ArgReg:
			asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
			asm.CmplRegMem(x64.RAX, x64.RSI, regDisp(a2))
		default:
			return false
		}
	case ir.InstLcmp:
		switch {
		case a1.Kind == ir.ArgReg && a2.Kind == ir.ArgIntConst:
			asm.CmpqConstMem(a1.Value, x64.RSI, regDisp(a1))
		default:
			return false
		}

	case ir.InstIret:
		switch a1.Kind {
		case ir.ArgReg:
			asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
		case ir.ArgIntConst:
			asm.MovlConstReg(a1.Value, x64.RAX)
		default:
			return false
		}
		asm.JmpMem(x64.RSI, -8)
	case ir.InstLret:
		asm.MovqMemReg(x64.RSI, x64.RAX, regDisp(a1))
		asm.JmpMem(x64.RSI, -8)
	case ir.InstRet:
		asm.JmpMem(x64.RSI, -8)

	case ir.InstCallGo:
		// TODO: refactor and optimize.
		sym := inst.Args[0].SymbolID()
		pkg := cl.ctx.State.Packages[sym.PackageIndex()]
		class := pkg.Classes[sym.ClassIndex()]
		method := class.Methods[sym.MemberIndex()]
		key := fmt.Sprintf("%s/%s.%s", pkg.Name, class.Name, method.Name)
		fnAddr := cl.ctx.State.GoFuncs[key]
		if fnAddr == 0 {
			return false
		}
		const tmp0offset = 24
		const envOffset = 16
		const arg0offset = -96
		offset := 0
		i := 1
		failed := false
		jclass.MethodDescriptor(method.Descriptor).WalkParams(func(typ jclass.DescriptorType) {
			arg := inst.Args[i]
			switch typ.Kind {
			case 'I':
				switch arg.Kind {
				case ir.ArgIntConst:
					asm.MovlConstMem(arg.Value, x64.RBP, int32(arg0offset+offset))
				case ir.ArgReg:
					asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(arg))
					asm.MovlRegMem(x64.RAX, x64.RBP, int32(arg0offset+offset))
				default:
					failed = true
				}
				offset += 4
			case 'J':
				switch arg.Kind {
				case ir.ArgIntConst:
					if !fits32bit(arg.Value) {
						failed = true
					}
					asm.MovqConst32Mem(int32(arg.Value), x64.RBP, int32(arg0offset+offset))
				case ir.ArgReg:
					asm.MovqMemReg(x64.RSI, x64.RAX, regDisp(arg))
					asm.MovqRegMem(x64.RAX, x64.RBP, int32(arg0offset+offset))
				default:
					failed = true
				}
				offset += 8
			default:
				failed = true // TODO: handle all other argument types
			}
			i++
		})
		if failed {
			return false
		}
		asm.MovqRegMem(x64.RSI, x64.RDX, tmp0offset) // Spill SI
		asm.MovlConstReg(int64(fnAddr), x64.RAX)     // TODO: handle addr higher than int32
		asm.CallReg(x64.RAX)
		asm.MovqMemReg(x64.RBP, x64.RDX, envOffset)  // Load DX
		asm.MovqMemReg(x64.RDX, x64.RSI, tmp0offset) // Load SI

	case ir.InstCallStatic:
		frameSize := cl.method.Out.FrameSlots*8 + 8
		sym := inst.Args[0].SymbolID()
		pkg := cl.ctx.State.Packages[sym.PackageIndex()]
		class := pkg.Classes[sym.ClassIndex()]
		method := class.Methods[sym.MemberIndex()]
		i := 1
		failed := false
		jclass.MethodDescriptor(method.Descriptor).WalkParams(func(typ jclass.DescriptorType) {
			arg := inst.Args[i]
			disp := int32(frameSize + (i-1)*8)
			switch typ.Kind {
			case 'I':
				switch arg.Kind {
				case ir.ArgIntConst:
					asm.MovlConstMem(arg.Value, x64.RSI, disp)
				case ir.ArgReg:
					asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(arg))
					asm.MovlRegMem(x64.RAX, x64.RSI, disp)
				default:
					failed = true
				}
				i++
			case 'J':
				switch arg.Kind {
				case ir.ArgIntConst:
					if !fits32bit(arg.Value) {
						failed = true
					}
					asm.MovqConst32Mem(int32(arg.Value), x64.RSI, disp)
				case ir.ArgReg:
					asm.MovqMemReg(x64.RSI, x64.RAX, regDisp(arg))
					asm.MovqRegMem(x64.RAX, x64.RSI, disp)
				default:
					failed = true
				}
				i++
			default:
				failed = true
			}
		})
		if failed {
			return false
		}
		asm.AddqConstReg(int64(frameSize), x64.RSI)
		{
			// The magic disp=16 is a width of instructions that
			// follow lea inside this block.
			asm.Raw(0x48, 0x8d, 0x05, 0x10, 0, 0, 0) // lea rax, [rip+16]
			asm.MovqRegMem(x64.RAX, x64.RSI, -8)
			index := asm.MovqFixup64Reg(x64.RAX)
			asm.JmpReg(x64.RAX)
			cl.pushReloc(a1.SymbolID(), index)
		}
		asm.AddqConstReg(int64(-frameSize), x64.RSI)
		if dst.Kind != 0 {
			// If function returns int, should use Movl.
			asm.MovqRegMem(x64.RAX, x64.RSI, regDisp(dst))
		}

	case ir.InstIsub:
		// We use negated argument for AddlConstMem for sub with constants.
		if a1 == dst {
			if a2.Kind == ir.ArgIntConst {
				asm.AddlConstMem(-a2.Value, x64.RSI, regDisp(dst))
			} else {
				return false
			}
		} else {
			switch {
			case a2.Kind == ir.ArgIntConst:
				asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
				asm.AddlConstReg(-a2.Value, x64.RAX)
				asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
			case a2.Kind == ir.ArgReg:
				asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
				asm.SublMemReg(x64.RSI, x64.RAX, regDisp(a2))
				asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
			default:
				return false
			}
		}
	case ir.InstIadd:
		if a1 == dst {
			if a2.Kind == ir.ArgIntConst {
				asm.AddlConstMem(a2.Value, x64.RSI, regDisp(dst))
			} else {
				return false
			}
		} else {
			switch {
			case a2.Kind == ir.ArgIntConst:
				asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
				asm.AddlConstReg(a2.Value, x64.RAX)
				asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
			case a2.Kind == ir.ArgReg:
				asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
				asm.AddlMemReg(x64.RSI, x64.RAX, regDisp(a2))
				asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
			default:
				return false
			}
		}

	default:
		return false
	}

	return true
}

func (cl *Compiler) pushReloc(src symbol.ID, offset int) {
	cl.methodRelocs++
	cl.relocs = append(cl.relocs, relocation{
		sourceID:     src,
		targetID:     symbol.NewID(cl.packageID, cl.classID, cl.methodID),
		targetOffset: offset,
	})
}

func (cl *Compiler) getMethodByID(id symbol.ID) *vmdat.Method {
	i1 := id.PackageIndex()
	i2 := id.ClassIndex()
	i3 := id.MemberIndex()
	return &cl.ctx.State.Packages[i1].Classes[i2].Methods[i3]
}
