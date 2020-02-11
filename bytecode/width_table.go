package bytecode

// OpWidth table maps Op to its encoding width.
//
// For example, OpWidth[Bupush] returns 2 as you need two bytes
// to encode bipush instruction.
//
// For variadic-length instructions like tableswitch it returns 0.
var OpWidth = [256]byte{
	Nop:             1 + 0,
	Aconstnull:      1 + 0,
	Iconstm1:        1 + 0,
	Iconst0:         1 + 0,
	Iconst1:         1 + 0,
	Iconst2:         1 + 0,
	Iconst3:         1 + 0,
	Iconst4:         1 + 0,
	Iconst5:         1 + 0,
	Lconst0:         1 + 0,
	Lconst1:         1 + 0,
	Fconst0:         1 + 0,
	Fconst1:         1 + 0,
	Fconst2:         1 + 0,
	Dconst0:         1 + 0,
	Dconst1:         1 + 0,
	Bipush:          1 + 1, // 1: byte
	Sipush:          1 + 2, // 2: byte1, byte2
	Ldc:             1 + 1, // 1: index
	Ldcw:            1 + 2, // 2: indexbyte1, indexbyte2
	Ldc2w:           1 + 2, // 2: indexbyte1, indexbyte2
	Iload:           1 + 1, // 1: index
	Lload:           1 + 1, // 1: index
	Fload:           1 + 1, // 1: index
	Dload:           1 + 1, // 1: index
	Aload:           1 + 1, // 1: index
	Iload0:          1 + 0,
	Iload1:          1 + 0,
	Iload2:          1 + 0,
	Iload3:          1 + 0,
	Lload0:          1 + 0,
	Lload1:          1 + 0,
	Lload2:          1 + 0,
	Lload3:          1 + 0,
	Fload0:          1 + 0,
	Fload1:          1 + 0,
	Fload2:          1 + 0,
	Fload3:          1 + 0,
	Dload0:          1 + 0,
	Dload1:          1 + 0,
	Dload2:          1 + 0,
	Dload3:          1 + 0,
	Aload0:          1 + 0,
	Aload1:          1 + 0,
	Aload2:          1 + 0,
	Aload3:          1 + 0,
	Iaload:          1 + 0,
	Laload:          1 + 0,
	Faload:          1 + 0,
	Daload:          1 + 0,
	Aaload:          1 + 0,
	Baload:          1 + 0,
	Caload:          1 + 0,
	Saload:          1 + 0,
	Istore:          1 + 1, // 1: index
	Lstore:          1 + 1, // 1: index
	Fstore:          1 + 1, // 1: index
	Dstore:          1 + 1, // 1: index
	Astore:          1 + 1, // 1: index
	Istore0:         1 + 0,
	Istore1:         1 + 0,
	Istore2:         1 + 0,
	Istore3:         1 + 0,
	Lstore0:         1 + 0,
	Lstore1:         1 + 0,
	Lstore2:         1 + 0,
	Lstore3:         1 + 0,
	Fstore0:         1 + 0,
	Fstore1:         1 + 0,
	Fstore2:         1 + 0,
	Fstore3:         1 + 0,
	Dstore0:         1 + 0,
	Dstore1:         1 + 0,
	Dstore2:         1 + 0,
	Dstore3:         1 + 0,
	Astore0:         1 + 0,
	Astore1:         1 + 0,
	Astore2:         1 + 0,
	Astore3:         1 + 0,
	Iastore:         1 + 0,
	Lastore:         1 + 0,
	Fastore:         1 + 0,
	Dastore:         1 + 0,
	Aastore:         1 + 0,
	Bastore:         1 + 0,
	Castore:         1 + 0,
	Sastore:         1 + 0,
	Pop:             1 + 0,
	Pop2:            1 + 0,
	Dup:             1 + 0,
	Dupx1:           1 + 0,
	Dupx2:           1 + 0,
	Dup2:            1 + 0,
	Dup2x1:          1 + 0,
	Dup2x2:          1 + 0,
	Swap:            1 + 0,
	Iadd:            1 + 0,
	Ladd:            1 + 0,
	Fadd:            1 + 0,
	Dadd:            1 + 0,
	Isub:            1 + 0,
	Lsub:            1 + 0,
	Fsub:            1 + 0,
	Dsub:            1 + 0,
	Imul:            1 + 0,
	Lmul:            1 + 0,
	Fmul:            1 + 0,
	Dmul:            1 + 0,
	Idiv:            1 + 0,
	Ldiv:            1 + 0,
	Fdiv:            1 + 0,
	Ddiv:            1 + 0,
	Irem:            1 + 0,
	Lrem:            1 + 0,
	Frem:            1 + 0,
	Drem:            1 + 0,
	Ineg:            1 + 0,
	Lneg:            1 + 0,
	Fneg:            1 + 0,
	Dneg:            1 + 0,
	Ishl:            1 + 0,
	Lshl:            1 + 0,
	Ishr:            1 + 0,
	Lshr:            1 + 0,
	Iushr:           1 + 0,
	Lushr:           1 + 0,
	Iand:            1 + 0,
	Land:            1 + 0,
	Ior:             1 + 0,
	Lor:             1 + 0,
	Ixor:            1 + 0,
	Lxor:            1 + 0,
	Iinc:            1 + 2, // 2: index, const
	I2l:             1 + 0,
	I2f:             1 + 0,
	I2d:             1 + 0,
	L2i:             1 + 0,
	L2f:             1 + 0,
	L2d:             1 + 0,
	F2i:             1 + 0,
	F2l:             1 + 0,
	F2d:             1 + 0,
	D2i:             1 + 0,
	D2l:             1 + 0,
	D2f:             1 + 0,
	I2b:             1 + 0,
	I2c:             1 + 0,
	I2s:             1 + 0,
	Lcmp:            1 + 0,
	Fcmpl:           1 + 0,
	Fcmpg:           1 + 0,
	Dcmpl:           1 + 0,
	Dcmpg:           1 + 0,
	Ifeq:            1 + 2, // 2: branchbyte1, branchbyte2
	Ifne:            1 + 2, // 2: branchbyte1, branchbyte2
	Iflt:            1 + 2, // 2: branchbyte1, branchbyte2
	Ifge:            1 + 2, // 2: branchbyte1, branchbyte2
	Ifgt:            1 + 2, // 2: branchbyte1, branchbyte2
	Ifle:            1 + 2, // 2: branchbyte1, branchbyte2
	Ificmpeq:        1 + 2, // 2: branchbyte1, branchbyte2
	Ificmpne:        1 + 2, // 2: branchbyte1, branchbyte2
	Ificmplt:        1 + 2, // 2: branchbyte1, branchbyte2
	Ificmpge:        1 + 2, // 2: branchbyte1, branchbyte2
	Ificmpgt:        1 + 2, // 2: branchbyte1, branchbyte2
	Ificmple:        1 + 2, // 2: branchbyte1, branchbyte2
	Ifacmpeq:        1 + 2, // 2: branchbyte1, branchbyte2
	Ifacmpne:        1 + 2, // 2: branchbyte1, branchbyte2
	Goto:            1 + 2, // 2: branchbyte1, branchbyte2
	Jsr:             1 + 2, // 2: branchbyte1, branchbyte2
	Ret:             1 + 1, // 1: index
	Ireturn:         1 + 0,
	Lreturn:         1 + 0,
	Freturn:         1 + 0,
	Dreturn:         1 + 0,
	Areturn:         1 + 0,
	Return:          1 + 0,
	Getstatic:       1 + 2, // 2: indexbyte1, indexbyte2
	Putstatic:       1 + 2, // 2: indexbyte1, indexbyte2
	Getfield:        1 + 2, // 2: indexbyte1, indexbyte2
	Putfield:        1 + 2, // 2: indexbyte1, indexbyte2
	Invokevirtual:   1 + 2, // 2: indexbyte1, indexbyte2
	Invokespecial:   1 + 2, // 2: indexbyte1, indexbyte2
	Invokestatic:    1 + 2, // 2: indexbyte1, indexbyte2
	Invokeinterface: 1 + 4, // 4: indexbyte1, indexbyte2, count, 0
	Invokedynamic:   1 + 4, // 4: indexbyte1, indexbyte2, 0, 0
	New:             1 + 2, // 2: indexbyte1, indexbyte2
	Newarray:        1 + 1, // 1: atype
	Anewarray:       1 + 2, // 2: indexbyte1, indexbyte2
	Arraylength:     1 + 0,
	Athrow:          1 + 0,
	Checkcast:       1 + 2, // 2: indexbyte1, indexbyte2
	Instanceof:      1 + 2, // 2: indexbyte1, indexbyte2
	Monitorenter:    1 + 0,
	Monitorexit:     1 + 0,
	Multianewarray:  1 + 3, // 3: indexbyte1, indexbyte2, dimensions
	Ifnull:          1 + 2, // 2: branchbyte1, branchbyte2
	Ifnonnull:       1 + 2, // 2: branchbyte1, branchbyte2
	Gotow:           1 + 4, // 4: branchbyte1, branchbyte2, branchbyte3, branchbyte4
	Jsrw:            1 + 4, // 4: branchbyte1, branchbyte2, branchbyte3, branchbyte4
	Breakpoint:      1 + 0,
	Impdep1:         1 + 0,
	Impdep2:         1 + 0,
}
