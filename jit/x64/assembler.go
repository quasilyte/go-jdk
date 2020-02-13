package x64

import (
	"bytes"
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

func (a *Assembler) Link() []byte {
	a.link()
	var buf bytes.Buffer
	for _, inst := range a.pending {
		buf.Write(inst.Bytes())
	}
	return buf.Bytes()
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

func (a *Assembler) Nop(length int) {
	// TODO(quasilyte): use wide NOPs.
	for i := 0; i < length; i++ {
		a.push(instruction{opcode: 0x90})
	}
}

func (a *Assembler) AddqConst(v, slot int64) {
	i := instruction{
		prefix: rexW,
		reg1:   op0,
		reg2:   rsi,
		flags:  flagModRM | flagMemory | flagImm,
	}
	i.disp = slot * 8
	i.imm = v
	if fitsInt8(v) {
		i.opcode = 0x83
	} else {
		i.opcode = 0x81
	}
	a.push(i)
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
	if inst.prefix != 0 {
		buf = append(buf, inst.prefix)
	}

	// Encode opcode.
	if inst.flags&flag0F != 0 {
		buf = append(buf, 0x0F)
	}
	buf = append(buf, inst.opcode)

	// Encode ModRM.
	if inst.flags&flagModRM != 0 {
		if inst.flags&flagMemory != 0 {
			var mod byte
			switch {
			case inst.disp == 0:
				mod = disp0
			case fitsInt8(inst.disp):
				mod = disp8
			default:
				mod = disp32
			}
			buf = append(buf, modrm(mod, inst.reg1, inst.reg2))
		} else {
			buf = append(buf, modrm(regreg, inst.reg1, inst.reg2))
		}
	}

	// Encode displacement.
	if inst.disp != 0 {
		if fitsInt8(inst.disp) {
			buf = append(buf, byte(inst.disp))
		} else {
			buf = appendInt32(buf, int32(inst.disp))
		}
	}

	// Encode immediate.
	if inst.flags&flagImm != 0 {
		if fitsInt8(inst.imm) {
			buf = append(buf, byte(inst.imm))
		} else {
			buf = appendInt32(buf, int32(inst.imm))
		}
	}

	inst.size = uint8(len(buf))
}
