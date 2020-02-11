package bytecode

// Op represents JVM bytecode operation code (opcode).
type Op byte

//go:generate stringer -type=Op -linecomment=true
const (
	// Nop - perform no operation.
	//
	// Encoding: 00
	// Operands stack effect: no changes
	Nop Op = 0 // nop

	// Aconstnull - push a `null` reference onto the stack.
	//
	// Encoding: 01
	// Operands stack effect: (→ null)
	Aconstnull Op = 1 // aconst_null

	// Iconstm1 - load the int value −1 onto the stack.
	//
	// Encoding: 02
	// Operands stack effect: (→ -1)
	Iconstm1 Op = 2 // iconst_m1

	// Iconst0 - load the int value 0 onto the stack.
	//
	// Encoding: 03
	// Operands stack effect: (→ 0)
	Iconst0 Op = 3 // iconst_0

	// Iconst1 - load the int value 1 onto the stack.
	//
	// Encoding: 04
	// Operands stack effect: (→ 1)
	Iconst1 Op = 4 // iconst_1

	// Iconst2 - load the int value 2 onto the stack.
	//
	// Encoding: 05
	// Operands stack effect: (→ 2)
	Iconst2 Op = 5 // iconst_2

	// Iconst3 - load the int value 3 onto the stack.
	//
	// Encoding: 06
	// Operands stack effect: (→ 3)
	Iconst3 Op = 6 // iconst_3

	// Iconst4 - load the int value 4 onto the stack.
	//
	// Encoding: 07
	// Operands stack effect: (→ 4)
	Iconst4 Op = 7 // iconst_4

	// Iconst5 - load the int value 5 onto the stack.
	//
	// Encoding: 08
	// Operands stack effect: (→ 5)
	Iconst5 Op = 8 // iconst_5

	// Lconst0 - push `0L` (the number [[zero]] with type `long`) onto the stack.
	//
	// Encoding: 09
	// Operands stack effect: (→ 0L)
	Lconst0 Op = 9 // lconst_0

	// Lconst1 - push `1L` (the number [[one]] with type `long`) onto the stack.
	//
	// Encoding: 0a
	// Operands stack effect: (→ 1L)
	Lconst1 Op = 10 // lconst_1

	// Fconst0 - push `0.0f` on the stack.
	//
	// Encoding: 0b
	// Operands stack effect: (→ 0.0f)
	Fconst0 Op = 11 // fconst_0

	// Fconst1 - push `1.0f` on the stack.
	//
	// Encoding: 0c
	// Operands stack effect: (→ 1.0f)
	Fconst1 Op = 12 // fconst_1

	// Fconst2 - push `2.0f` on the stack.
	//
	// Encoding: 0d
	// Operands stack effect: (→ 2.0f)
	Fconst2 Op = 13 // fconst_2

	// Dconst0 - push the constant `0.0` (a `double`) onto the stack.
	//
	// Encoding: 0e
	// Operands stack effect: (→ 0.0)
	Dconst0 Op = 14 // dconst_0

	// Dconst1 - push the constant `1.0` (a `double`) onto the stack.
	//
	// Encoding: 0f
	// Operands stack effect: (→ 1.0)
	Dconst1 Op = 15 // dconst_1

	// Bipush - push a `byte` onto the stack as an integer `value`.
	//
	// Encoding: 10 + 1: byte
	// Operands stack effect: (→ value)
	Bipush Op = 16 // bipush

	// Sipush - push a short onto the stack as an integer `value`.
	//
	// Encoding: 11 + 2: byte1, byte2
	// Operands stack effect: (→ value)
	Sipush Op = 17 // sipush

	// Ldc - push a constant `#index` from a constant pool (String, int, float, Class, java.lang.invoke.MethodType, java.lang.invoke.MethodHandle, or a dynamically-computed constant) onto the stack.
	//
	// Encoding: 12 + 1: index
	// Operands stack effect: (→ value)
	Ldc Op = 18 // ldc

	// Ldcw - push a constant `#index` from a constant pool (String, int,
	// float, Class, java.lang.invoke.MethodType, java.lang.invoke.MethodHandle,
	// or a dynamically-computed constant) onto the stack.
	// Wide index is constructed as indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: 13 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (→ value)
	Ldcw Op = 19 // ldc_w

	// Ldc2w - push a constant `#index` from a constant pool (double, long,
	// or a dynamically-computed constant) onto the stack.
	// Wide index is constructed as indexbyte1 << 8 + indexbyte2.
	//
	// Encoding: 14 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (→ value)
	Ldc2w Op = 20 // ldc2_w

	// Iload - load an int `value` from a local variable `#index`.
	//
	// Encoding: 15 + 1: index
	// Operands stack effect: (→ value)
	Iload Op = 21 // iload

	// Lload - load a long value from a local variable `#index`.
	//
	// Encoding: 16 + 1: index
	// Operands stack effect: (→ value)
	Lload Op = 22 // lload

	// Fload - load a float `value` from a local variable `#index`.
	//
	// Encoding: 17 + 1: index
	// Operands stack effect: (→ value)
	Fload Op = 23 // fload

	// Dload - load a double `value` from a local variable `#index`.
	//
	// Encoding: 18 + 1: index
	// Operands stack effect: (→ value)
	Dload Op = 24 // dload

	// Aload - load a reference onto the stack from a local variable `#index`.
	//
	// Encoding: 19 + 1: index
	// Operands stack effect: (→ objectref)
	Aload Op = 25 // aload

	// Iload0 - load an int `value` from local variable 0.
	//
	// Encoding: 1a
	// Operands stack effect: (→ value)
	Iload0 Op = 26 // iload_0

	// Iload1 - load an int `value` from local variable 1.
	//
	// Encoding: 1b
	// Operands stack effect: (→ value)
	Iload1 Op = 27 // iload_1

	// Iload2 - load an int `value` from local variable 2.
	//
	// Encoding: 1c
	// Operands stack effect: (→ value)
	Iload2 Op = 28 // iload_2

	// Iload3 - load an int `value` from local variable 3.
	//
	// Encoding: 1d
	// Operands stack effect: (→ value)
	Iload3 Op = 29 // iload_3

	// Lload0 - load a long value from a local variable 0.
	//
	// Encoding: 1e
	// Operands stack effect: (→ value)
	Lload0 Op = 30 // lload_0

	// Lload1 - load a long value from a local variable 1.
	//
	// Encoding: 1f
	// Operands stack effect: (→ value)
	Lload1 Op = 31 // lload_1

	// Lload2 - load a long value from a local variable 2.
	//
	// Encoding: 20
	// Operands stack effect: (→ value)
	Lload2 Op = 32 // lload_2

	// Lload3 - load a long value from a local variable 3.
	//
	// Encoding: 21
	// Operands stack effect: (→ value)
	Lload3 Op = 33 // lload_3

	// Fload0 - load a float `value` from local variable 0.
	//
	// Encoding: 22
	// Operands stack effect: (→ value)
	Fload0 Op = 34 // fload_0

	// Fload1 - load a float `value` from local variable 1.
	//
	// Encoding: 23
	// Operands stack effect: (→ value)
	Fload1 Op = 35 // fload_1

	// Fload2 - load a float `value` from local variable 2.
	//
	// Encoding: 24
	// Operands stack effect: (→ value)
	Fload2 Op = 36 // fload_2

	// Fload3 - load a float `value` from local variable 3.
	//
	// Encoding: 25
	// Operands stack effect: (→ value)
	Fload3 Op = 37 // fload_3

	// Dload0 - load a double from local variable 0.
	//
	// Encoding: 26
	// Operands stack effect: (→ value)
	Dload0 Op = 38 // dload_0

	// Dload1 - load a double from local variable 1.
	//
	// Encoding: 27
	// Operands stack effect: (→ value)
	Dload1 Op = 39 // dload_1

	// Dload2 - load a double from local variable 2.
	//
	// Encoding: 28
	// Operands stack effect: (→ value)
	Dload2 Op = 40 // dload_2

	// Dload3 - load a double from local variable 3.
	//
	// Encoding: 29
	// Operands stack effect: (→ value)
	Dload3 Op = 41 // dload_3

	// Aload0 - load a reference onto the stack from local variable 0.
	//
	// Encoding: 2a
	// Operands stack effect: (→ objectref)
	Aload0 Op = 42 // aload_0

	// Aload1 - load a reference onto the stack from local variable 1.
	//
	// Encoding: 2b
	// Operands stack effect: (→ objectref)
	Aload1 Op = 43 // aload_1

	// Aload2 - load a reference onto the stack from local variable 2.
	//
	// Encoding: 2c
	// Operands stack effect: (→ objectref)
	Aload2 Op = 44 // aload_2

	// Aload3 - load a reference onto the stack from local variable 3.
	//
	// Encoding: 2d
	// Operands stack effect: (→ objectref)
	Aload3 Op = 45 // aload_3

	// Iaload - load an int from an array.
	//
	// Encoding: 2e
	// Operands stack effect: (arrayref, index → value)
	Iaload Op = 46 // iaload

	// Laload - load a long from an array.
	//
	// Encoding: 2f
	// Operands stack effect: (arrayref, index → value)
	Laload Op = 47 // laload

	// Faload - load a float from an array.
	//
	// Encoding: 30
	// Operands stack effect: (arrayref, index → value)
	Faload Op = 48 // faload

	// Daload - load a double from an array.
	//
	// Encoding: 31
	// Operands stack effect: (arrayref, index → value)
	Daload Op = 49 // daload

	// Aaload - load onto the stack a reference from an array.
	//
	// Encoding: 32
	// Operands stack effect: (arrayref, index → value)
	Aaload Op = 50 // aaload

	// Baload - load a byte or Boolean value from an array.
	//
	// Encoding: 33
	// Operands stack effect: (arrayref, index → value)
	Baload Op = 51 // baload

	// Caload - load a char from an array.
	//
	// Encoding: 34
	// Operands stack effect: (arrayref, index → value)
	Caload Op = 52 // caload

	// Saload - load short from array.
	//
	// Encoding: 35
	// Operands stack effect: (arrayref, index → value)
	Saload Op = 53 // saload

	// Istore - store int `value` into variable `#index`.
	//
	// Encoding: 36 + 1: index
	// Operands stack effect: (value →)
	Istore Op = 54 // istore

	// Lstore - store a long `value` in a local variable `#index`.
	//
	// Encoding: 37 + 1: index
	// Operands stack effect: (value →)
	Lstore Op = 55 // lstore

	// Fstore - store a float `value` into a local variable `#index`.
	//
	// Encoding: 38 + 1: index
	// Operands stack effect: (value →)
	Fstore Op = 56 // fstore

	// Dstore - store a double `value` into a local variable `#index`.
	//
	// Encoding: 39 + 1: index
	// Operands stack effect: (value →)
	Dstore Op = 57 // dstore

	// Astore - store a reference into a local variable `#index`.
	//
	// Encoding: 3a + 1: index
	// Operands stack effect: (objectref →)
	Astore Op = 58 // astore

	// Istore0 - store int `value` into variable 0.
	//
	// Encoding: 3b
	// Operands stack effect: (value →)
	Istore0 Op = 59 // istore_0

	// Istore1 - store int `value` into variable 1.
	//
	// Encoding: 3c
	// Operands stack effect: (value →)
	Istore1 Op = 60 // istore_1

	// Istore2 - store int `value` into variable 2.
	//
	// Encoding: 3d
	// Operands stack effect: (value →)
	Istore2 Op = 61 // istore_2

	// Istore3 - store int `value` into variable 3.
	//
	// Encoding: 3e
	// Operands stack effect: (value →)
	Istore3 Op = 62 // istore_3

	// Lstore0 - store a long `value` in a local variable 0.
	//
	// Encoding: 3f
	// Operands stack effect: (value →)
	Lstore0 Op = 63 // lstore_0

	// Lstore1 - store a long `value` in a local variable 1.
	//
	// Encoding: 40
	// Operands stack effect: (value →)
	Lstore1 Op = 64 // lstore_1

	// Lstore2 - store a long `value` in a local variable 2.
	//
	// Encoding: 41
	// Operands stack effect: (value →)
	Lstore2 Op = 65 // lstore_2

	// Lstore3 - store a long `value` in a local variable 3.
	//
	// Encoding: 42
	// Operands stack effect: (value →)
	Lstore3 Op = 66 // lstore_3

	// Fstore0 - store a float `value` into local variable 0.
	//
	// Encoding: 43
	// Operands stack effect: (value →)
	Fstore0 Op = 67 // fstore_0

	// Fstore1 - store a float `value` into local variable 1.
	//
	// Encoding: 44
	// Operands stack effect: (value →)
	Fstore1 Op = 68 // fstore_1

	// Fstore2 - store a float `value` into local variable 2.
	//
	// Encoding: 45
	// Operands stack effect: (value →)
	Fstore2 Op = 69 // fstore_2

	// Fstore3 - store a float `value` into local variable 3.
	//
	// Encoding: 46
	// Operands stack effect: (value →)
	Fstore3 Op = 70 // fstore_3

	// Dstore0 - store a double into local variable 0.
	//
	// Encoding: 47
	// Operands stack effect: (value →)
	Dstore0 Op = 71 // dstore_0

	// Dstore1 - store a double into local variable 1.
	//
	// Encoding: 48
	// Operands stack effect: (value →)
	Dstore1 Op = 72 // dstore_1

	// Dstore2 - store a double into local variable 2.
	//
	// Encoding: 49
	// Operands stack effect: (value →)
	Dstore2 Op = 73 // dstore_2

	// Dstore3 - store a double into local variable 3.
	//
	// Encoding: 4a
	// Operands stack effect: (value →)
	Dstore3 Op = 74 // dstore_3

	// Astore0 - store a reference into local variable 0.
	//
	// Encoding: 4b
	// Operands stack effect: (objectref →)
	Astore0 Op = 75 // astore_0

	// Astore1 - store a reference into local variable 1.
	//
	// Encoding: 4c
	// Operands stack effect: (objectref →)
	Astore1 Op = 76 // astore_1

	// Astore2 - store a reference into local variable 2.
	//
	// Encoding: 4d
	// Operands stack effect: (objectref →)
	Astore2 Op = 77 // astore_2

	// Astore3 - store a reference into local variable 3.
	//
	// Encoding: 4e
	// Operands stack effect: (objectref →)
	Astore3 Op = 78 // astore_3

	// Iastore - store an int into an array.
	//
	// Encoding: 4f
	// Operands stack effect: (arrayref, index, value →)
	Iastore Op = 79 // iastore

	// Lastore - store a long to an array.
	//
	// Encoding: 50
	// Operands stack effect: (arrayref, index, value →)
	Lastore Op = 80 // lastore

	// Fastore - store a float in an array.
	//
	// Encoding: 51
	// Operands stack effect: (arrayref, index, value →)
	Fastore Op = 81 // fastore

	// Dastore - store a double into an array.
	//
	// Encoding: 52
	// Operands stack effect: (arrayref, index, value →)
	Dastore Op = 82 // dastore

	// Aastore - store a reference in an array.
	//
	// Encoding: 53
	// Operands stack effect: (arrayref, index, value →)
	Aastore Op = 83 // aastore

	// Bastore - store a byte or Boolean value into an array.
	//
	// Encoding: 54
	// Operands stack effect: (arrayref, index, value →)
	Bastore Op = 84 // bastore

	// Castore - store a char into an array.
	//
	// Encoding: 55
	// Operands stack effect: (arrayref, index, value →)
	Castore Op = 85 // castore

	// Sastore - store short to array.
	//
	// Encoding: 56
	// Operands stack effect: (arrayref, index, value →)
	Sastore Op = 86 // sastore

	// Pop - discard the top value on the stack.
	//
	// Encoding: 57
	// Operands stack effect: (value →)
	Pop Op = 87 // pop

	// Pop2 - discard the top two values on the stack (or one value, if it is a double or long).
	//
	// Encoding: 58
	// Operands stack effect: ({value2, value1} →)
	Pop2 Op = 88 // pop2

	// Dup - duplicate the value on top of the stack.
	//
	// Encoding: 59
	// Operands stack effect: (value → value, value)
	Dup Op = 89 // dup

	// Dupx1 - insert a copy of the top value into the stack two values from the top. value1 and value2 must not be of the type double or long..
	//
	// Encoding: 5a
	// Operands stack effect: (value2, value1 → value1, value2, value1)
	Dupx1 Op = 90 // dup_x1

	// Dupx2 - insert a copy of the top value into the stack two (if value2 is double or long it takes up the entry of value3, too) or three values (if value2 is neither double nor long) from the top.
	//
	// Encoding: 5b
	// Operands stack effect: (value3, value2, value1 → value1, value3, value2, value1)
	Dupx2 Op = 91 // dup_x2

	// Dup2 - duplicate top two stack words (two values, if value1 is not double nor long; a single value, if value1 is double or long).
	//
	// Encoding: 5c
	// Operands stack effect: ({value2, value1} → {value2, value1}, {value2, value1})
	Dup2 Op = 92 // dup2

	// Dup2x1 - duplicate two words and insert beneath third word (see explanation above).
	//
	// Encoding: 5d
	// Operands stack effect: (value3, {value2, value1} → {value2, value1}, value3, {value2, value1})
	Dup2x1 Op = 93 // dup2_x1

	// Dup2x2 - duplicate two words and insert beneath fourth word.
	//
	// Encoding: 5e
	// Operands stack effect: ({value4, value3}, {value2, value1} → {value2, value1}, {value4, value3}, {value2, value1})
	Dup2x2 Op = 94 // dup2_x2

	// Swap - swaps two top words on the stack (note that value1 and value2 must not be double or long).
	//
	// Encoding: 5f
	// Operands stack effect: (value2, value1 → value1, value2)
	Swap Op = 95 // swap

	// Iadd - add two ints.
	//
	// Encoding: 60
	// Operands stack effect: (value1, value2 → result)
	Iadd Op = 96 // iadd

	// Ladd - add two longs.
	//
	// Encoding: 61
	// Operands stack effect: (value1, value2 → result)
	Ladd Op = 97 // ladd

	// Fadd - add two floats.
	//
	// Encoding: 62
	// Operands stack effect: (value1, value2 → result)
	Fadd Op = 98 // fadd

	// Dadd - add two doubles.
	//
	// Encoding: 63
	// Operands stack effect: (value1, value2 → result)
	Dadd Op = 99 // dadd

	// Isub - int subtract.
	//
	// Encoding: 64
	// Operands stack effect: (value1, value2 → result)
	Isub Op = 100 // isub

	// Lsub - subtract two longs.
	//
	// Encoding: 65
	// Operands stack effect: (value1, value2 → result)
	Lsub Op = 101 // lsub

	// Fsub - subtract two floats.
	//
	// Encoding: 66
	// Operands stack effect: (value1, value2 → result)
	Fsub Op = 102 // fsub

	// Dsub - subtract a double from another.
	//
	// Encoding: 67
	// Operands stack effect: (value1, value2 → result)
	Dsub Op = 103 // dsub

	// Imul - multiply two integers.
	//
	// Encoding: 68
	// Operands stack effect: (value1, value2 → result)
	Imul Op = 104 // imul

	// Lmul - multiply two longs.
	//
	// Encoding: 69
	// Operands stack effect: (value1, value2 → result)
	Lmul Op = 105 // lmul

	// Fmul - multiply two floats.
	//
	// Encoding: 6a
	// Operands stack effect: (value1, value2 → result)
	Fmul Op = 106 // fmul

	// Dmul - multiply two doubles.
	//
	// Encoding: 6b
	// Operands stack effect: (value1, value2 → result)
	Dmul Op = 107 // dmul

	// Idiv - divide two integers.
	//
	// Encoding: 6c
	// Operands stack effect: (value1, value2 → result)
	Idiv Op = 108 // idiv

	// Ldiv - divide two longs.
	//
	// Encoding: 6d
	// Operands stack effect: (value1, value2 → result)
	Ldiv Op = 109 // ldiv

	// Fdiv - divide two floats.
	//
	// Encoding: 6e
	// Operands stack effect: (value1, value2 → result)
	Fdiv Op = 110 // fdiv

	// Ddiv - divide two doubles.
	//
	// Encoding: 6f
	// Operands stack effect: (value1, value2 → result)
	Ddiv Op = 111 // ddiv

	// Irem - logical int remainder.
	//
	// Encoding: 70
	// Operands stack effect: (value1, value2 → result)
	Irem Op = 112 // irem

	// Lrem - remainder of division of two longs.
	//
	// Encoding: 71
	// Operands stack effect: (value1, value2 → result)
	Lrem Op = 113 // lrem

	// Frem - get the remainder from a division between two floats.
	//
	// Encoding: 72
	// Operands stack effect: (value1, value2 → result)
	Frem Op = 114 // frem

	// Drem - get the remainder from a division between two doubles.
	//
	// Encoding: 73
	// Operands stack effect: (value1, value2 → result)
	Drem Op = 115 // drem

	// Ineg - negate int.
	//
	// Encoding: 74
	// Operands stack effect: (value → result)
	Ineg Op = 116 // ineg

	// Lneg - negate a long.
	//
	// Encoding: 75
	// Operands stack effect: (value → result)
	Lneg Op = 117 // lneg

	// Fneg - negate a float.
	//
	// Encoding: 76
	// Operands stack effect: (value → result)
	Fneg Op = 118 // fneg

	// Dneg - negate a double.
	//
	// Encoding: 77
	// Operands stack effect: (value → result)
	Dneg Op = 119 // dneg

	// Ishl - int shift left.
	//
	// Encoding: 78
	// Operands stack effect: (value1, value2 → result)
	Ishl Op = 120 // ishl

	// Lshl - bitwise shift left of a long `value1` by int `value2` positions.
	//
	// Encoding: 79
	// Operands stack effect: (value1, value2 → result)
	Lshl Op = 121 // lshl

	// Ishr - int arithmetic shift right.
	//
	// Encoding: 7a
	// Operands stack effect: (value1, value2 → result)
	Ishr Op = 122 // ishr

	// Lshr - bitwise shift right of a long `value1` by int `value2` positions.
	//
	// Encoding: 7b
	// Operands stack effect: (value1, value2 → result)
	Lshr Op = 123 // lshr

	// Iushr - int logical shift right.
	//
	// Encoding: 7c
	// Operands stack effect: (value1, value2 → result)
	Iushr Op = 124 // iushr

	// Lushr - bitwise shift right of a long `value1` by int `value2` positions, unsigned.
	//
	// Encoding: 7d
	// Operands stack effect: (value1, value2 → result)
	Lushr Op = 125 // lushr

	// Iand - perform a bitwise AND on two integers.
	//
	// Encoding: 7e
	// Operands stack effect: (value1, value2 → result)
	Iand Op = 126 // iand

	// Land - [[bitwise operation.
	//
	// Encoding: 7f
	// Operands stack effect: (value1, value2 → result)
	Land Op = 127 // land

	// Ior - bitwise int OR.
	//
	// Encoding: 80
	// Operands stack effect: (value1, value2 → result)
	Ior Op = 128 // ior

	// Lor - bitwise OR of two longs.
	//
	// Encoding: 81
	// Operands stack effect: (value1, value2 → result)
	Lor Op = 129 // lor

	// Ixor - int xor.
	//
	// Encoding: 82
	// Operands stack effect: (value1, value2 → result)
	Ixor Op = 130 // ixor

	// Lxor - bitwise XOR of two longs.
	//
	// Encoding: 83
	// Operands stack effect: (value1, value2 → result)
	Lxor Op = 131 // lxor

	// Iinc - increment local variable `#index` by signed byte `const`.
	//
	// Encoding: 84 + 2: index, const
	// Operands stack effect: no changes
	Iinc Op = 132 // iinc

	// I2l - convert an int into a long.
	//
	// Encoding: 85
	// Operands stack effect: (value → result)
	I2l Op = 133 // i2l

	// I2f - convert an int into a float.
	//
	// Encoding: 86
	// Operands stack effect: (value → result)
	I2f Op = 134 // i2f

	// I2d - convert an int into a double.
	//
	// Encoding: 87
	// Operands stack effect: (value → result)
	I2d Op = 135 // i2d

	// L2i - convert a long to a int.
	//
	// Encoding: 88
	// Operands stack effect: (value → result)
	L2i Op = 136 // l2i

	// L2f - convert a long to a float.
	//
	// Encoding: 89
	// Operands stack effect: (value → result)
	L2f Op = 137 // l2f

	// L2d - convert a long to a double.
	//
	// Encoding: 8a
	// Operands stack effect: (value → result)
	L2d Op = 138 // l2d

	// F2i - convert a float to an int.
	//
	// Encoding: 8b
	// Operands stack effect: (value → result)
	F2i Op = 139 // f2i

	// F2l - convert a float to a long.
	//
	// Encoding: 8c
	// Operands stack effect: (value → result)
	F2l Op = 140 // f2l

	// F2d - convert a float to a double.
	//
	// Encoding: 8d
	// Operands stack effect: (value → result)
	F2d Op = 141 // f2d

	// D2i - convert a double to an int.
	//
	// Encoding: 8e
	// Operands stack effect: (value → result)
	D2i Op = 142 // d2i

	// D2l - convert a double to a long.
	//
	// Encoding: 8f
	// Operands stack effect: (value → result)
	D2l Op = 143 // d2l

	// D2f - convert a double to a float.
	//
	// Encoding: 90
	// Operands stack effect: (value → result)
	D2f Op = 144 // d2f

	// I2b - convert an int into a byte.
	//
	// Encoding: 91
	// Operands stack effect: (value → result)
	I2b Op = 145 // i2b

	// I2c - convert an int into a character.
	//
	// Encoding: 92
	// Operands stack effect: (value → result)
	I2c Op = 146 // i2c

	// I2s - convert an int into a short.
	//
	// Encoding: 93
	// Operands stack effect: (value → result)
	I2s Op = 147 // i2s

	// Lcmp - push 0 if the two longs are the same, 1 if value1 is greater than value2, -1 otherwise.
	//
	// Encoding: 94
	// Operands stack effect: (value1, value2 → result)
	Lcmp Op = 148 // lcmp

	// Fcmpl - compare two floats.
	//
	// Encoding: 95
	// Operands stack effect: (value1, value2 → result)
	Fcmpl Op = 149 // fcmpl

	// Fcmpg - compare two floats.
	//
	// Encoding: 96
	// Operands stack effect: (value1, value2 → result)
	Fcmpg Op = 150 // fcmpg

	// Dcmpl - compare two doubles.
	//
	// Encoding: 97
	// Operands stack effect: (value1, value2 → result)
	Dcmpl Op = 151 // dcmpl

	// Dcmpg - compare two doubles.
	//
	// Encoding: 98
	// Operands stack effect: (value1, value2 → result)
	Dcmpg Op = 152 // dcmpg

	// Ifeq - if `value` is 0, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: 99 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value →)
	Ifeq Op = 153 // ifeq

	// Ifne - if `value` is not 0, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: 9a + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value →)
	Ifne Op = 154 // ifne

	// Iflt - if `value` is less than 0, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: 9b + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value →)
	Iflt Op = 155 // iflt

	// Ifge - if `value` is greater than or equal to 0, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: 9c + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value →)
	Ifge Op = 156 // ifge

	// Ifgt - if `value` is greater than 0, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: 9d + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value →)
	Ifgt Op = 157 // ifgt

	// Ifle - if `value` is less than or equal to 0, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: 9e + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value →)
	Ifle Op = 158 // ifle

	// Ificmpeq - if ints are equal, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: 9f + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value1, value2 →)
	Ificmpeq Op = 159 // if_icmpeq

	// Ificmpne - if ints are not equal, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a0 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value1, value2 →)
	Ificmpne Op = 160 // if_icmpne

	// Ificmplt - if `value1` is less than `value2`, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a1 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value1, value2 →)
	Ificmplt Op = 161 // if_icmplt

	// Ificmpge - if `value1` is greater than or equal to `value2`, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a2 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value1, value2 →)
	Ificmpge Op = 162 // if_icmpge

	// Ificmpgt - if `value1` is greater than `value2`, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a3 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value1, value2 →)
	Ificmpgt Op = 163 // if_icmpgt

	// Ificmple - if `value1` is less than or equal to `value2`, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a4 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value1, value2 →)
	Ificmple Op = 164 // if_icmple

	// Ifacmpeq - if references are equal, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a5 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value1, value2 →)
	Ifacmpeq Op = 165 // if_acmpeq

	// Ifacmpne - if references are not equal, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a6 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value1, value2 →)
	Ifacmpne Op = 166 // if_acmpne

	// Goto - goes to another instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a7 + 2: branchbyte1, branchbyte2
	// Operands stack effect: no changes
	Goto Op = 167 // goto

	// Jsr - jump to subroutine at `branchoffset` and place the return address on the stack.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: a8 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (→ address)
	Jsr Op = 168 // jsr

	// Ret - continue execution from address taken from a local variable `#index`
	// The asymmetry with jsr is intentional.
	//
	// Encoding: a9 + 1: index
	// Operands stack effect: no changes
	Ret Op = 169 // ret

	// Ireturn - return an integer from a method.
	//
	// Encoding: ac
	// Operands stack effect: (value → [empty])
	Ireturn Op = 172 // ireturn

	// Lreturn - return a long value.
	//
	// Encoding: ad
	// Operands stack effect: (value → [empty])
	Lreturn Op = 173 // lreturn

	// Freturn - return a float.
	//
	// Encoding: ae
	// Operands stack effect: (value → [empty])
	Freturn Op = 174 // freturn

	// Dreturn - return a double from a method.
	//
	// Encoding: af
	// Operands stack effect: (value → [empty])
	Dreturn Op = 175 // dreturn

	// Areturn - return a reference from a method.
	//
	// Encoding: b0
	// Operands stack effect: (objectref → [empty])
	Areturn Op = 176 // areturn

	// Return - return void from method.
	//
	// Encoding: b1
	// Operands stack effect: (→ [empty])
	Return Op = 177 // return

	// Getstatic - get a static field `value` of a class, where the field is
	// identified by field reference in the constant pool `index` (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: b2 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (→ value)
	Getstatic Op = 178 // getstatic

	// Putstatic - set static field to `value` in a class, where the field is
	// identified by a field reference `index` in constant pool (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: b3 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (value →)
	Putstatic Op = 179 // putstatic

	// Getfield - get a field `value` of an object `objectref`,
	// where the field is identified by field reference in the constant pool `index` (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: b4 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (objectref → value)
	Getfield Op = 180 // getfield

	// Putfield - set field to `value` in an object `objectref`,
	// where the field is identified by a field reference `index` in constant pool (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: b5 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (objectref, value →)
	Putfield Op = 181 // putfield

	// Invokevirtual - invoke virtual method on object `objectref` and puts the result on the stack (might be void).
	// The method is identified by method reference `index` in constant pool (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: b6 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (objectref, [arg1, arg2, ...] → result)
	Invokevirtual Op = 182 // invokevirtual

	// Invokespecial - invoke instance method on object `objectref` and puts the result on the stack (might be void).
	// The method is identified by method reference `index` in constant pool (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: b7 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (objectref, [arg1, arg2, ...] → result)
	Invokespecial Op = 183 // invokespecial

	// Invokestatic - invoke a static method and puts the result on the stack (might be void).
	// The method is identified by method reference `index` in constant pool (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: b8 + 2: indexbyte1, indexbyte2
	// Operands stack effect: ([arg1, arg2, ...] → result)
	Invokestatic Op = 184 // invokestatic

	// Invokeinterface - invokes an interface method on object `objectref` and puts the result on the stack (might be void).
	// The interface method is identified by method reference `index` in constant pool (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: b9 + 4: indexbyte1, indexbyte2, count, 0
	// Operands stack effect: (objectref, [arg1, arg2, ...] → result)
	Invokeinterface Op = 185 // invokeinterface

	// Invokedynamic - invokes a dynamic method and puts the result on the stack (might be void).
	// The method is identified by method reference `index` in constant pool (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: ba + 4: indexbyte1, indexbyte2, 0, 0
	// Operands stack effect: ([arg1, [arg2 ...]] → result)
	Invokedynamic Op = 186 // invokedynamic

	// New - create new object of type identified by class reference in constant pool `index` (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: bb + 2: indexbyte1, indexbyte2
	// Operands stack effect: (→ objectref)
	New Op = 187 // new

	// Newarray - create new array with `count` elements of primitive type identified by `atype`.
	//
	// Encoding: bc + 1: atype
	// Operands stack effect: (count → arrayref)
	Newarray Op = 188 // newarray

	// Anewarray - create a new array of references of length `count` and component type
	// identified by the class reference `index` (indexbyte1 << 8 + indexbyte2) in the constant pool.
	//
	// Encoding: bd + 2: indexbyte1, indexbyte2
	// Operands stack effect: (count → arrayref)
	Anewarray Op = 189 // anewarray

	// Arraylength - get the length of an array.
	//
	// Encoding: be
	// Operands stack effect: (arrayref → length)
	Arraylength Op = 190 // arraylength

	// Athrow - throws an error or exception (notice that the rest of the stack is cleared, leaving only a reference to the Throwable).
	//
	// Encoding: bf
	// Operands stack effect: (objectref → [empty], objectref)
	Athrow Op = 191 // athrow

	// Checkcast - checks whether an `objectref` is of a certain type, the class reference
	// of which is in the constant pool at `index` (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: c0 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (objectref → objectref)
	Checkcast Op = 192 // checkcast

	// Instanceof - determines if an object `objectref` is of a given type,
	// identified by class reference `index` in constant pool (indexbyte1 << 8 + indexbyte2).
	//
	// Encoding: c1 + 2: indexbyte1, indexbyte2
	// Operands stack effect: (objectref → result)
	Instanceof Op = 193 // instanceof

	// Monitorenter - enter monitor for object ("grab the lock" – start of synchronized() section).
	//
	// Encoding: c2
	// Operands stack effect: (objectref →)
	Monitorenter Op = 194 // monitorenter

	// Monitorexit - exit monitor for object ("release the lock" – end of synchronized() section).
	//
	// Encoding: c3
	// Operands stack effect: (objectref →)
	Monitorexit Op = 195 // monitorexit

	// Multianewarray - create a new array of `dimensions` dimensions of type
	// identified by class reference in constant pool `index` (indexbyte1 << 8 + indexbyte2);
	// the sizes of each dimension is identified by count1, [count2, etc.]
	//
	// Encoding: c5 + 3: indexbyte1, indexbyte2, dimensions
	// Operands stack effect: (count1, [count2,...] → arrayref)
	Multianewarray Op = 197 // multianewarray

	// Ifnull - if value is null, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: c6 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value →)
	Ifnull Op = 198 // ifnull

	// Ifnonnull - if value is not null, branch to instruction at `branchoffset`.
	// Signed short constructed from unsigned bytes branchbyte1 << 8 + branchbyte2.
	//
	// Encoding: c7 + 2: branchbyte1, branchbyte2
	// Operands stack effect: (value →)
	Ifnonnull Op = 199 // ifnonnull

	// Gotow - goes to another instruction at `branchoffset`.
	// Signed int constructed from unsigned bytes branchbyte1 << 24 + branchbyte2 << 16 + branchbyte3 << 8 + branchbyte4).
	//
	// Encoding: c8 + 4: branchbyte1, branchbyte2, branchbyte3, branchbyte4
	// Operands stack effect: no changes
	Gotow Op = 200 // goto_w

	// Jsrw - jump to subroutine at `branchoffset` and place the return address on the stack.
	// Signed int constructed from unsigned bytes branchbyte1 << 24 + branchbyte2 << 16 + branchbyte3 << 8 + branchbyte4.
	//
	// Encoding: c9 + 4: branchbyte1, branchbyte2, branchbyte3, branchbyte4
	// Operands stack effect: (→ address)
	Jsrw Op = 201 // jsr_w

	// Breakpoint - reserved for breakpoints in Java debuggers; should not appear in any class file.
	//
	// Encoding: ca
	// Operands stack effect: no changes
	Breakpoint Op = 202 // breakpoint

	// Impdep1 - reserved for implementation-dependent operations within debuggers; should not appear in any class file.
	//
	// Encoding: fe
	// Operands stack effect: no changes
	Impdep1 Op = 254 // impdep1

	// Impdep2 - reserved for implementation-dependent operations within debuggers; should not appear in any class file.
	//
	// Encoding: ff
	// Operands stack effect: no changes
	Impdep2 Op = 255 // impdep2
)
