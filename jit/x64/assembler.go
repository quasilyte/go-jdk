package x64

import (
	"bytes"
	"encoding/binary"
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
	a.pushJmp(instruction{
		opcode: jmp32op,
		size:   5, // Size of rel32 jump
		imm:    labelID,
		flags:  flagPseudo,
	})
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

	// Pass 1 (short, only jumps): find short jumps, assign opcodes.
	pass2needed := false
	for _, jumpIndex := range a.jumps {
		jmp := &a.pending[jumpIndex]
		labelIndex := id2index[int(jmp.imm)]
		label := &a.pending[labelIndex]
		dist := label.offset - jmp.offset

		if dist > -120 && dist < 120 {
			pass2needed = true
			jmp.size = 2
			jmp.buf[0] = jmp8op
		} else {
			jmp.buf[0] = jmp32op
		}
	}

	// Pass 2: re-calculate offsets.
	if pass2needed {
		offset := int32(0)
		for i, inst := range a.pending {
			a.pending[i].offset = offset
			offset += int32(inst.size)
		}
	}

	// Pass 3 (short, only jumps): write jump relative targets.
	for _, jumpIndex := range a.jumps {
		jmp := &a.pending[jumpIndex]
		labelIndex := id2index[int(jmp.imm)]
		label := &a.pending[labelIndex]
		dist := label.offset - jmp.offset

		if jmp.buf[0] == jmp8op {
			jmp.buf[1] = byte(dist - 2)
		} else {
			binary.LittleEndian.PutUint32(jmp.buf[1:], uint32(dist-5))
		}
	}
}

func (a *Assembler) pushJmp(inst instruction) {
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
