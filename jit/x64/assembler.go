package x64

import (
	"math"
)

func NewAssembler() *Assembler {
	return &Assembler{
		pending: make([]instruction, 0, 256),
		jumps:   make([]int, 0, 32),
		labels:  make([]int, 0, 32),
	}
}

type Assembler struct {
	length  int
	pending []instruction
	jumps   []int
	labels  []int
}

func (a *Assembler) Reset() {
	a.length = 0
	a.pending = a.pending[:0]
	a.jumps = a.jumps[:0]
	a.labels = a.labels[:0]
}

func (a *Assembler) OffsetOf(index int) int32 {
	return a.pending[index].offset
}

func (a *Assembler) Link() int {
	a.link()
	return a.length
}

func (a *Assembler) Put(dst []byte) {
	offset := 0
	for _, inst := range a.pending {
		copy(dst[offset:], inst.Bytes())
		offset += int(inst.size)
	}
}

func (a *Assembler) Label(id int64) {
	a.labels = append(a.labels, len(a.pending))
	a.push(instruction{imm: id, flags: flagPseudo})
}

func (a *Assembler) Jmp(labelID int64) {
	a.pushJmp(jmp8op, labelID)
}

func (a *Assembler) Jge(labelID int64) {
	a.pushJcc(jge8op, labelID)
}

func (a *Assembler) Jlt(labelID int64) {
	a.pushJcc(jlt8op, labelID)
}

func (a *Assembler) Jgt(labelID int64) {
	a.pushJcc(jgt8op, labelID)
}

func (a *Assembler) Nop(length int) {
	// TODO(quasilyte): use wide NOPs.
	for i := 0; i < length; i++ {
		a.push(instruction{opcode: 0x90})
	}
}

func (a *Assembler) Cdq() {
	a.push(instruction{opcode: 0x99})
}

func (a *Assembler) CallReg(reg uint8) {
	a.push(instruction{
		opcode: 0xff,
		reg1:   op2,
		reg2:   reg,
		flags:  flagModRM,
	})
}

func (a *Assembler) JmpMem(reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0xFF,
		reg1:   op4,
		reg2:   reg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) JmpReg(reg uint8) {
	a.push(instruction{
		opcode: 0xFF,
		reg1:   op4,
		reg2:   reg,
		flags:  flagModRM,
	})
}

func (a *Assembler) Raw(enc ...byte) {
	var inst instruction
	copy(inst.buf[:], enc)
	inst.size = uint8(len(enc))
	inst.flags = flagPseudo
	a.push(inst)
}

func (a *Assembler) NegqReg(reg uint8) {
	a.push(instruction{
		opcode: 0xF7,
		reg1:   op3,
		reg2:   reg,
		flags:  flagModRM | flagRexW,
	})
}

func (a *Assembler) NeglReg(reg uint8) {
	a.push(instruction{
		opcode: 0xF7,
		reg1:   op3,
		reg2:   reg,
		flags:  flagModRM,
	})
}

func (a *Assembler) NeglMem(reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0xF7,
		reg1:   op3,
		reg2:   reg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) NegqMem(reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0xF7,
		reg1:   op3,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagRexW,
		disp:   disp,
	})
}

func (a *Assembler) CmplConstMem(v int64, reg uint8, disp int32) {
	if v >= -128 && v >= 127 {
		a.CmplConst8Mem(int8(v), reg, disp)
	} else {
		a.CmplConst32Mem(int32(v), reg, disp)
	}
}

func (a *Assembler) CmplConst8Mem(v int8, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x83,
		reg1:   op7,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm8,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) CmplConst32Mem(v int32, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x81,
		reg1:   op7,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm32,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) CmpqConstMem(v int64, reg uint8, disp int32) {
	if v >= -128 && v >= 127 {
		a.CmpqConst8Mem(int8(v), reg, disp)
	} else {
		a.CmpqConst32Mem(int32(v), reg, disp)
	}
}

func (a *Assembler) CmpqConst8Mem(v int8, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x83,
		reg1:   op7,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm8 | flagRexW,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) CmpqConst32Mem(v int32, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x81,
		reg1:   op7,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm32 | flagRexW,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) CmplRegMem(xreg uint8, yreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x3B,
		reg1:   xreg,
		reg2:   yreg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) MovbRegMem(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x88,
		reg1:   srcreg,
		reg2:   dstreg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) MovlRegMem(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x89,
		reg1:   srcreg,
		reg2:   dstreg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) MovlRegMemindex(srcreg, dstreg, index uint8) {
	a.push(instruction{
		opcode: 0x89,
		reg1:   srcreg,
		reg2:   dstreg,
		flags:  flagModSIB | flagMemory | flagScale4,
		index:  index,
	})
}

func (a *Assembler) MovlMemindexReg(srcreg, dstreg, index uint8) {
	a.push(instruction{
		opcode: 0x8B,
		reg1:   dstreg,
		reg2:   srcreg,
		index:  index,
		flags:  flagModSIB | flagMemory | flagScale4,
	})
}

func (a *Assembler) MovlMemReg(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x8B,
		reg1:   dstreg,
		reg2:   srcreg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) MovlConstMem(v int64, reg uint8, disp int32) {
	a.MovlConst32Mem(int32(v), reg, disp)
}

func (a *Assembler) MovlConstMemindex(v int64, reg uint8, index uint8) {
	a.MovlConst32Memindex(int32(v), reg, index)
}

func (a *Assembler) MovlConst32Memindex(v int32, reg uint8, index uint8) {
	a.push(instruction{
		opcode: 0xC7,
		reg1:   op0,
		reg2:   reg,
		index:  index,
		flags:  flagModSIB | flagMemory | flagImm32 | flagScale4,
		imm:    int64(v),
	})
}

func (a *Assembler) MovlConst32Mem(v int32, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0xC7,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm32,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) MovlConstReg(v int64, reg uint8) {
	a.MovlConst32Reg(int32(v), reg)
}

func (a *Assembler) MovlConst32Reg(v int32, reg uint8) {
	a.push(instruction{
		opcode: 0xB8,
		reg2:   reg,
		flags:  flagImm32 | flagReg2Op,
		imm:    int64(v),
	})
}

func (a *Assembler) MovqFixup64Reg(reg uint8) int {
	i := len(a.pending)
	a.MovqConst64Reg(0, reg)
	return i
}

func (a *Assembler) MovqConstReg(v int64, reg uint8) {
	if v >= math.MinInt32 && v <= math.MaxInt32 {
		a.MovqConst32Reg(int32(v), reg)
	} else {
		a.MovqConst64Reg(v, reg)
	}
}

func (a *Assembler) MovqConst32Reg(v int32, reg uint8) {
	a.push(instruction{
		opcode: 0xC7,
		reg1:   op0,
		reg2:   reg,
		flags:  flagImm32 | flagModRM | flagRexW,
		imm:    int64(v),
	})
}

func (a *Assembler) MovqConst64Reg(v int64, reg uint8) {
	a.push(instruction{
		opcode: 0xB8,
		reg2:   reg,
		flags:  flagImm64 | flagRexW | flagReg2Op,
		imm:    v,
	})
}

func (a *Assembler) MovqConst32Mem(v int32, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0xc7,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm32 | flagRexW,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) MovqRegMem(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x89,
		reg1:   srcreg,
		reg2:   dstreg,
		flags:  flagModRM | flagMemory | flagRexW,
		disp:   disp,
	})
}

func (a *Assembler) MovqMemReg(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x8B,
		reg1:   dstreg,
		reg2:   srcreg,
		flags:  flagModRM | flagMemory | flagRexW,
		disp:   disp,
	})
}

func (a *Assembler) SublMemReg(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x2b,
		reg1:   dstreg,
		reg2:   srcreg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) IdivlMem(reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0xF7,
		reg1:   op7,
		reg2:   reg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) ImullMemReg(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0xAF,
		reg1:   dstreg,
		reg2:   srcreg,
		flags:  flag0F | flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) AddlMemReg(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x03,
		reg1:   dstreg,
		reg2:   srcreg,
		flags:  flagModRM | flagMemory,
		disp:   disp,
	})
}

func (a *Assembler) AddlConstMem(v int64, reg uint8, disp int32) {
	if v >= -128 && v >= 127 {
		a.AddlConst8Mem(int8(v), reg, disp)
	} else {
		a.AddlConst32Mem(int32(v), reg, disp)
	}
}

func (a *Assembler) AddlConst8Mem(v int8, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x83,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm8,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) AddlConst32Mem(v int32, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x81,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm32,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) AddlConstReg(v int64, reg uint8) {
	if v >= -128 && v >= 127 {
		a.AddlConst8Reg(int8(v), reg)
	} else {
		a.AddlConst32Reg(int32(v), reg)
	}
}

func (a *Assembler) AddlConst8Reg(v int8, reg uint8) {
	a.push(instruction{
		opcode: 0x83,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagImm8,
		imm:    int64(v),
	})
}

func (a *Assembler) AddlConst32Reg(v int32, reg uint8) {
	a.push(instruction{
		opcode: 0x81,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagImm32,
		imm:    int64(v),
	})
}

func (a *Assembler) AddqConstReg(v int64, reg uint8) {
	if v >= -128 && v >= 127 {
		a.AddqConst8Reg(int8(v), reg)
	} else {
		a.AddqConst32Reg(int32(v), reg)
	}
}

func (a *Assembler) AddqConst8Reg(v int8, reg uint8) {
	a.push(instruction{
		opcode: 0x83,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagImm8 | flagRexW,
		imm:    int64(v),
	})
}

func (a *Assembler) AddqConst32Reg(v int32, reg uint8) {
	a.push(instruction{
		opcode: 0x81,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagImm32 | flagRexW,
		imm:    int64(v),
	})
}

func (a *Assembler) AddqConst8Mem(v int8, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x83,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm8 | flagRexW,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) AddqConst32Mem(v int32, reg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x81,
		reg1:   op0,
		reg2:   reg,
		flags:  flagModRM | flagMemory | flagImm32 | flagRexW,
		disp:   disp,
		imm:    int64(v),
	})
}

func (a *Assembler) MovlqsxMemReg(srcreg, dstreg uint8, disp int32) {
	a.push(instruction{
		opcode: 0x63,
		reg1:   dstreg,
		reg2:   srcreg,
		flags:  flagModRM | flagMemory | flagRexW,
		disp:   disp,
	})
}

func (a *Assembler) link() {
	// TODO: don't use a map here.
	id2index := make(map[int]int, len(a.labels))
	for _, labelIndex := range a.labels {
		id := a.pending[labelIndex].imm
		id2index[int(id)] = labelIndex
	}

	// Pass 1 (fast, only jumps): find short jumps, assign opcodes.
	pass2needed := false
	for _, jumpIndex := range a.jumps {
		jmp := &a.pending[jumpIndex]
		labelIndex := id2index[int(jmp.imm)]
		label := &a.pending[labelIndex]
		dist := label.offset - jmp.offset

		if dist < -120 || dist > 120 {
			pass2needed = true
			jmp.opcode = jumpRel8ToRel32[jmp.opcode]
			jmp.size = uint8(jmp.disp)
		}
	}

	// Pass 2: re-calculate offsets.
	//
	// For small functions where all jumps are 8-bit relative,
	// we don't need a second pass at all.
	if pass2needed {
		offset := int32(0)
		for i, inst := range a.pending {
			a.pending[i].offset = offset
			offset += int32(inst.size)
		}
		// Update length by using last instruction offset and size.
		last := a.pending[len(a.pending)-1]
		a.length = int(last.offset + int32(last.size))
	}

	// Pass 3 (fast, only jumps): assemble jumps.
	for _, jumpIndex := range a.jumps {
		jmp := &a.pending[jumpIndex]
		labelIndex := id2index[int(jmp.imm)]
		label := &a.pending[labelIndex]
		target := (label.offset - jmp.offset) - int32(jmp.size)

		// We're using a fact that all jumps have very similar encoding.
		// It's either 2, 5 or 6 bytes.
		// 2 bytes are used for all rel8 forms.
		// 5 bytes are unconditional jump (JMP) rel32 form.
		// 6 bytes is a conditional jump (Jcc) rel32 that requires 0x0F prefix.
		buf := jmp.buf[:0]
		if jmp.size == 6 {
			buf = append(buf, 0x0F)
		}
		buf = append(buf, jmp.opcode)
		if jmp.size == 2 {
			buf = append(buf, byte(target))
		} else {
			buf = appendInt32(buf, target)
		}
		if len(buf) != int(jmp.size) {
			panic("len(buf) != jmp.size")
		}
	}
}

func (a *Assembler) pushJcc(op uint8, labelID int64) {
	inst := instruction{
		opcode: op,
		imm:    labelID,
		size:   2, // rel8 form size
		disp:   6, // rel32 form size
		flags:  flagPseudo,
	}
	a.jumps = append(a.jumps, len(a.pending))
	a.push(inst)
}

func (a *Assembler) pushJmp(op uint8, labelID int64) {
	inst := instruction{
		opcode: op,
		imm:    labelID,
		size:   2, // rel8 form size
		disp:   5, // rel32 form size
		flags:  flagPseudo,
	}
	a.jumps = append(a.jumps, len(a.pending))
	a.push(inst)
}

func (a *Assembler) push(inst instruction) {
	if inst.flags&flagPseudo == 0 {
		a.encode(&inst)
	}
	inst.offset = int32(a.length)
	a.length += int(inst.size)
	a.pending = append(a.pending, inst)
}

func (a *Assembler) encode(inst *instruction) {
	buf := inst.buf[:0]

	// Encode prefix.
	var rexPrefix byte
	if inst.flags&flagRexW != 0 {
		rexPrefix |= rexW
	}
	if inst.reg1 >= R8 {
		rexPrefix |= rexR
	}
	if inst.reg2 >= R8 {
		rexPrefix |= rexB
	}
	if rexPrefix != 0 {
		buf = append(buf, rexPrefix)
	}

	reg1 := inst.reg1 & regBitMask
	reg2 := inst.reg2 & regBitMask

	// Encode opcode.
	if inst.flags&flag0F != 0 {
		buf = append(buf, 0x0F)
	}
	opcode := inst.opcode
	if inst.flags&flagReg2Op != 0 {
		opcode += reg2
	}
	buf = append(buf, opcode)

	// Encode ModRM.
	// If ModSIB is set, we encode ModRM along with SIB.
	if inst.flags&flagModRM != 0 {
		if inst.flags&flagMemory != 0 {
			var mod byte
			switch {
			case inst.disp == 0:
				mod = disp0
			case fitsInt8(int64(inst.disp)):
				mod = disp8
			default:
				mod = disp32
			}
			buf = append(buf, modrm(mod, reg1, reg2))
		} else {
			buf = append(buf, modrm(regreg, reg1, reg2))
		}
	} else if inst.flags&flagModSIB != 0 {
		buf = append(buf, modrm(0, reg1, 4))
		switch {
		case inst.flags&flagScale4 != 0:
			buf = append(buf, sib(scale4, inst.index, reg2))
		case inst.flags&flagScale8 != 0:
			buf = append(buf, sib(scale8, inst.index, reg2))
		}
	}

	// Encode displacement.
	if inst.disp != 0 {
		if fitsInt8(int64(inst.disp)) {
			buf = append(buf, byte(inst.disp))
		} else {
			buf = appendInt32(buf, int32(inst.disp))
		}
	}

	// Encode immediate.
	switch {
	case inst.flags&flagImm8 != 0:
		buf = append(buf, byte(inst.imm))
	case inst.flags&flagImm32 != 0:
		buf = appendInt32(buf, int32(inst.imm))
	case inst.flags&flagImm64 != 0:
		buf = appendInt64(buf, inst.imm)
	}

	inst.size = uint8(len(buf))
}
