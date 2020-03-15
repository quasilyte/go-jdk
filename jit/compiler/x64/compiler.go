package x64

import (
	"encoding/binary"
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
		addr := uint64(uintptr(unsafe.Pointer(&srcMethod.Code[0])))
		binary.LittleEndian.PutUint64(dstMethod.Code[rel.targetOffset:], addr)
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
	case ir.InstJumpGt:
		asm.Jgt(a1.Value)
	case ir.InstJumpLtEq:
		asm.Jle(a1.Value)
	case ir.InstJumpLt:
		asm.Jlt(a1.Value)
	case ir.InstJumpNotEqual:
		asm.Jne(a1.Value)
	case ir.InstJump:
		asm.Jmp(a1.Value)

	case ir.InstArrayLen:
		asm.MovqMemReg(x64.RSI, x64.RAX, ptrDisp(a1))
		asm.MovlMemReg(x64.RAX, x64.RAX, 16)
		asm.MovlRegMem(x64.RAX, x64.RSI, scalarDisp(dst))
	case ir.InstIntArrayGet:
		aref := a1
		index := a2
		asm.MovqMemReg(x64.RSI, x64.RAX, ptrDisp(aref))
		asm.MovqMemReg(x64.RAX, x64.RAX, 8)
		switch index.Kind {
		case ir.ArgIntConst:
			asm.MovlMemReg(x64.RAX, x64.RAX, int32(index.Value)*4)
		case ir.ArgReg:
			asm.MovlMemReg(x64.RSI, x64.RCX, scalarDisp(index))
			asm.MovlMemindexReg(x64.RAX, x64.RAX, x64.RCX)
		}
		asm.MovlRegMem(x64.RAX, x64.RSI, scalarDisp(dst))
	case ir.InstIntArraySet:
		aref := a1
		index := a2
		v := inst.Args[2]
		asm.MovqMemReg(x64.RSI, x64.RAX, ptrDisp(aref))
		asm.MovqMemReg(x64.RAX, x64.RAX, 8)
		switch {
		case v.Kind == ir.ArgIntConst && index.Kind == ir.ArgIntConst:
			asm.MovlConstMem(v.Value, x64.RAX, int32(index.Value)*4)
		case v.Kind == ir.ArgReg && index.Kind == ir.ArgIntConst:
			asm.MovlMemReg(x64.RSI, x64.RCX, scalarDisp(v))
			asm.MovlRegMem(x64.RCX, x64.RAX, int32(index.Value)*4)
		case v.Kind == ir.ArgIntConst && index.Kind == ir.ArgReg:
			asm.MovlMemReg(x64.RSI, x64.RCX, scalarDisp(index))
			asm.MovlConstMemindex(v.Value, x64.RAX, x64.RCX)
		case v.Kind == ir.ArgReg && index.Kind == ir.ArgReg:
			asm.MovlMemReg(x64.RSI, x64.R8, scalarDisp(v))
			asm.MovlMemReg(x64.RSI, x64.RCX, scalarDisp(index))
			asm.MovlRegMemindex(x64.R8, x64.RAX, x64.RCX)
		default:
			return false
		}
	case ir.InstNewIntArray:
		fnAddr := cl.ctx.Funcs.NewIntArray
		ok := cl.assembleCallGo(uintptr(fnAddr), "($I)[I", inst.Dst, inst.Args)
		if !ok {
			return false
		}

	case ir.InstAload:
		switch a1.Kind {
		case ir.ArgReg:
			asm.MovqMemReg(x64.RSI, x64.RAX, ptrDisp(a1))
			asm.MovqRegMem(x64.RAX, x64.RSI, ptrDisp(dst))
		default:
			return false
		}
	case ir.InstLload:
		switch a1.Kind {
		case ir.ArgReg:
			asm.MovqMemReg(x64.RSI, x64.RAX, scalarDisp(a1))
			asm.MovqRegMem(x64.RAX, x64.RSI, scalarDisp(dst))
		case ir.ArgIntConst:
			if !fits32bit(a1.Value) {
				return false
			}
			asm.MovqConst32Mem(int32(a1.Value), x64.RSI, scalarDisp(dst))
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
		asm.JmpMem(x64.RSI, -16)
	case ir.InstAret:
		if a1.Kind != ir.ArgReg {
			return false
		}
		asm.MovqMemReg(x64.RSI, x64.RAX, ptrDisp(a1))
		asm.JmpMem(x64.RSI, -16)
	case ir.InstLret:
		switch a1.Kind {
		case ir.ArgReg:
			asm.MovqMemReg(x64.RSI, x64.RAX, scalarDisp(a1))
		case ir.ArgIntConst:
			if !fits32bit(a1.Value) {
				return false
			}
			asm.MovqConstReg(a1.Value, x64.RAX)
		}
		asm.JmpMem(x64.RSI, -16)
	case ir.InstRet:
		asm.JmpMem(x64.RSI, -16)

	case ir.InstCallGo:
		sym := inst.Args[0].SymbolID()
		pkg := cl.ctx.State.Packages[sym.PackageIndex()]
		class := pkg.Classes[sym.ClassIndex()]
		method := class.Methods[sym.MemberIndex()]
		key := fmt.Sprintf("%s/%s.%s", pkg.Name, class.Name, method.Name)
		fnAddr := cl.ctx.State.GoFuncs[key]
		return cl.assembleCallGo(fnAddr, method.Descriptor, inst.Dst, inst.Args[1:])

	case ir.InstCallStatic:
		return cl.assembleCallStatic(inst)

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
	case ir.InstImul:
		switch {
		case a2.Kind == ir.ArgReg:
			asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
			asm.ImullMemReg(x64.RSI, x64.RAX, regDisp(a2))
			asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
		default:
			return false
		}
	case ir.InstIdiv:
		switch {
		case a1.Kind == ir.ArgReg && a2.Kind == ir.ArgReg:
			asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
			asm.Cdq()
			asm.IdivlMem(x64.RSI, regDisp(a2))
			asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
		case a1.Kind == ir.ArgReg && a2.Kind == ir.ArgIntConst:
			asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
			asm.Cdq()
			asm.MovlConstReg(a2.Value, x64.RCX)
			asm.IdivlReg(x64.RCX)
			asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
		default:
			return false
		}

	case ir.InstConvI2L:
		asm.MovlqsxMemReg(x64.RSI, x64.RAX, regDisp(a1))
		asm.MovqRegMem(x64.RAX, x64.RSI, regDisp(dst))
	case ir.InstConvI2B:
		asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(a1))
		asm.MovbRegMem(x64.RAX, x64.RSI, regDisp(dst))
	default:
		return false
	}

	return true
}

func (cl *Compiler) assembleCallStatic(inst ir.Inst) bool {
	asm := cl.asm

	// Frame size is 16 bytes per every stack slot plus
	// one extra slot to store the return address.
	frameSize := cl.method.Out.FrameSlots*16 + 16
	sym := inst.Args[0].SymbolID()
	pkg := cl.ctx.State.Packages[sym.PackageIndex()]
	class := pkg.Classes[sym.ClassIndex()]
	method := class.Methods[sym.MemberIndex()]
	i := 1
	failed := false
	signature := jclass.MethodDescriptor(method.Descriptor)
	signature.WalkParams(func(typ jclass.DescriptorType) {
		arg := inst.Args[i]
		disp := int32(frameSize + (i-1)*16)
		switch {
		case typ.Dims != 0:
			if arg.Kind != ir.ArgReg {
				failed = true
				return
			}
			asm.MovqMemReg(x64.RSI, x64.RAX, ptrDisp(arg))
			asm.MovqRegMem(x64.RAX, x64.RSI, disp+8)
		case typ.Kind == 'I':
			switch arg.Kind {
			case ir.ArgIntConst:
				asm.MovlConstMem(arg.Value, x64.RSI, disp)
			case ir.ArgReg:
				asm.MovlMemReg(x64.RSI, x64.RAX, regDisp(arg))
				asm.MovlRegMem(x64.RAX, x64.RSI, disp)
			default:
				failed = true
			}
		case typ.Kind == 'J':
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
		default:
			failed = true
		}
		i++
	})
	if failed {
		return false
	}

	asm.AddqConstReg(int64(frameSize), x64.RSI)
	{
		// The magic disp=16 is a width of instructions that
		// follow lea inside this block.
		asm.Raw(0x48, 0x8d, 0x05, 0x10, 0, 0, 0) // lea rax, [rip+16]
		asm.MovqRegMem(x64.RAX, x64.RSI, -16)
		index := asm.MovqFixup64Reg(x64.RAX)
		asm.JmpReg(x64.RAX)
		cl.pushReloc(inst.Args[0].SymbolID(), index)
	}
	asm.AddqConstReg(int64(-frameSize), x64.RSI)
	if inst.Dst.Kind != 0 {
		typ := signature.ReturnType()
		switch {
		case typ.Dims != 0:
			asm.MovqRegMem(x64.RAX, x64.RSI, ptrDisp(inst.Dst))
		case typ.Kind == 'I':
			asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(inst.Dst))
		case typ.Kind == 'J':
			asm.MovqRegMem(x64.RAX, x64.RSI, regDisp(inst.Dst))
		default:
			return false
		}
	}

	return true
}

func (cl *Compiler) assembleCallGo(fnAddr uintptr, desc string, dst ir.Arg, args []ir.Arg) bool {
	// TODO: refactor and optimize.

	const (
		arg0offset   = -96
		gocallOffset = 73
		frameSize    = 96
		tmp0offset   = 24
		envOffset    = 16
	)

	asm := cl.asm // Just for convenience

	offset := 0
	i := 0
	failed := false
	signature := jclass.MethodDescriptor(desc)

	signature.WalkParams(func(typ jclass.DescriptorType) {
		arg := args[i]
		switch {
		case typ.Dims != 0:
			if rem := offset % 8; rem != 0 {
				offset += rem
			}
			if arg.Kind != ir.ArgReg {
				failed = true
				return
			}
			asm.MovqMemReg(x64.RSI, x64.RAX, ptrDisp(arg))
			asm.MovqRegMem(x64.RAX, x64.RBP, int32(arg0offset+offset))
			offset += 8
		case typ.Kind == '$':
			// Dollar ($) is our special marker for env argument.
			if rem := offset % 8; rem != 0 {
				offset += rem
			}
			asm.MovqMemReg(x64.RBP, x64.RAX, envOffset)
			asm.MovqRegMem(x64.RAX, x64.RBP, int32(arg0offset+offset))
			offset += 8
		case typ.Kind == 'I':
			if rem := offset % 4; rem != 0 {
				offset += rem
			}
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
		case typ.Kind == 'J':
			if rem := offset % 8; rem != 0 {
				offset += rem
			}
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

	asm.MovqRegMem(x64.RSI, x64.RDI, tmp0offset) // Spill SI
	asm.MovlConstReg(int64(fnAddr), x64.RCX)
	asm.MovlConstReg(int64(cl.ctx.Funcs.JcallScalar+gocallOffset), x64.RDI)
	asm.Raw(0x48, 0x8d, 0x05, 4+2, 0, 0, 0)
	asm.MovqRegMem(x64.RAX, x64.RBP, -8)
	asm.JmpReg(x64.RDI)
	asm.MovqMemReg(x64.RBP, x64.RDI, envOffset)  // Load DI
	asm.MovqMemReg(x64.RDI, x64.RSI, tmp0offset) // Load SI
	if dst.Kind != 0 {
		// Return values start from a location aligned to a pointer size.
		if rem := offset % 8; rem != 0 {
			offset += rem
		}
		typ := signature.ReturnType()
		switch {
		case typ.Dims != 0 || typ.Kind == 'L':
			asm.MovqMemReg(x64.RBP, x64.RAX, int32(arg0offset+offset))
			asm.MovqRegMem(x64.RAX, x64.RSI, ptrDisp(dst))
		case typ.Kind == 'I':
			asm.MovlMemReg(x64.RBP, x64.RAX, int32(arg0offset+offset))
			asm.MovlRegMem(x64.RAX, x64.RSI, regDisp(dst))
		case typ.Kind == 'J':
			asm.MovqMemReg(x64.RBP, x64.RAX, int32(arg0offset+offset))
			asm.MovqRegMem(x64.RAX, x64.RSI, regDisp(dst))
		default:
			return false
		}
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
