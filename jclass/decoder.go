package jclass

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Decoder struct {
	r *bufio.Reader
	f *File

	deferred struct {
		names        map[*string]uint16
		names2       map[*string]uint16
		descriptors  map[*string]uint16
		descriptors2 map[*string]uint16
		classNames   map[*string]uint16
	}
}

func (d *Decoder) Decode(r io.Reader) (*File, error) {
	d.r = bufio.NewReader(r)
	d.f = &File{}
	d.deferred.names = map[*string]uint16{}
	d.deferred.names2 = map[*string]uint16{}
	d.deferred.descriptors = map[*string]uint16{}
	d.deferred.descriptors2 = map[*string]uint16{}
	d.deferred.classNames = map[*string]uint16{}
	err := d.decode()
	return d.f, err
}

func (d *Decoder) resolveDeferred() {
	for s, index := range d.deferred.names {
		*s = d.f.Consts[index].(*Utf8Const).Value
	}
	for s, index := range d.deferred.names2 {
		*s = d.f.Consts[index].(*NameAndTypeConst).Name
	}
	for s, index := range d.deferred.descriptors {
		*s = d.f.Consts[index].(*Utf8Const).Value
	}
	for s, index := range d.deferred.descriptors2 {
		*s = d.f.Consts[index].(*NameAndTypeConst).Descriptor
	}
	for s, index := range d.deferred.classNames {
		*s = d.f.Consts[index].(*ClassConst).Name
	}
}

func (d *Decoder) deferNameResolving(index uint16, s *string) {
	d.deferred.names[s] = index
}

func (d *Decoder) deferClassNameResolving(index uint16, s *string) {
	d.deferred.classNames[s] = index
}

func (d *Decoder) deferDescriptorResolving(index uint16, desc *string) {
	d.deferred.descriptors[desc] = index
}

func (d *Decoder) deferNameAndTypeResolving(index uint16, name, desc *string) {
	d.deferred.names2[name] = index
	d.deferred.descriptors2[desc] = index
}

func (d *Decoder) decode() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"magic", d.decodeMagic},
		{"version", d.decodeVersion},
		{"constant pool", d.decodeConstantPool},
		{"access flags", d.decodeClassAccessFlags},
		{"this class", d.decodeClassName},
		{"super class", d.makeUint16Decode(&d.f.SuperClass)},
		{"interfaces", d.decodeInterfaces},
		{"fields", d.decodeFields},
		{"methods", d.decodeMethods},
		{"attributes", d.decodeAttributes},
	}
	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("decode %s: %w", step.name, err)
		}
	}
	return nil
}

func (d *Decoder) decodeClassAccessFlags() error {
	v, err := d.readUint16()
	if err != nil {
		return err
	}
	d.f.AccessFlags = AccessFlags(v)
	return nil
}

func (d *Decoder) decodeClassName() error {
	v, err := d.readUint16()
	if err != nil {
		return err
	}
	d.f.ThisClassName = d.f.Consts[v].(*ClassConst).Name
	return nil
}

func (d *Decoder) decodeMagic() error {
	magic, err := d.readUint32()
	if err != nil {
		return err
	}
	if magic != 0xCAFEBABE {
		return fmt.Errorf("invalid value (want 0xCAFEBABE, got 0x%X)", magic)
	}
	return nil
}

func (d *Decoder) decodeVersion() error {
	minor, err := d.readUint16()
	if err != nil {
		return fmt.Errorf("read minor: %w", err)
	}
	major, err := d.readUint16()
	if err != nil {
		return fmt.Errorf("read major: %w", err)
	}
	d.f.Ver.Minor = minor
	d.f.Ver.Major = major
	return nil
}

func (d *Decoder) decodeConstantPool() error {
	n, err := d.readUint16()
	if err != nil {
		return fmt.Errorf("read count: %w", err)
	}
	cp := make([]Const, n)
	i := 1 // Constant at 0 index is undefined
	for i < len(cp) {
		c, skip, err := d.readConst()
		if err != nil {
			return fmt.Errorf("const%d: %w", i, err)
		}
		cp[i] = c
		i += skip
	}
	d.f.Consts = cp
	d.resolveDeferred()
	return nil
}

func (d *Decoder) decodeInterfaces() error {
	n, err := d.readUint16()
	if err != nil {
		return fmt.Errorf("read count: %w", err)
	}
	ifaces := make([]uint16, n)
	for i := range ifaces {
		index, err := d.readUint16()
		if err != nil {
			return fmt.Errorf("interface%d: %w", i, err)
		}
		ifaces[i] = index
	}
	d.f.Interfaces = ifaces
	return nil
}

func (d *Decoder) decodeFields() error {
	n, err := d.readUint16()
	if err != nil {
		return fmt.Errorf("read count: %w", err)
	}
	fields := make([]Field, n)
	for i := range fields {
		field, err := d.readField()
		if err != nil {
			return fmt.Errorf("field%d: %w", i, err)
		}
		fields[i] = field
	}
	d.f.Fields = fields
	return nil
}

func (d *Decoder) decodeMethods() error {
	n, err := d.readUint16()
	if err != nil {
		return fmt.Errorf("read count: %w", err)
	}
	methods := make([]Method, n)
	for i := range methods {
		field, err := d.readField()
		if err != nil {
			return fmt.Errorf("method%d: %w", i, err)
		}
		methods[i] = Method(field)
	}
	d.f.Methods = methods
	return nil
}

func (d *Decoder) decodeAttributes() error {
	attrs, err := d.readAttributes()
	if err != nil {
		return err
	}
	d.f.Attrs = attrs
	return nil
}

func (d *Decoder) makeUint16Decode(dst *uint16) func() error {
	return func() error {
		v, err := d.readUint16()
		if err != nil {
			return err
		}
		*dst = v
		return nil
	}
}

func (d *Decoder) readExceptionHandler() (ExceptionHandler, error) {
	var h ExceptionHandler
	startPC, err := d.readUint16()
	if err != nil {
		return h, fmt.Errorf("read start_pc: %w", err)
	}
	endPC, err := d.readUint16()
	if err != nil {
		return h, fmt.Errorf("read end_pc: %w", err)
	}
	handlerPC, err := d.readUint16()
	if err != nil {
		return h, fmt.Errorf("read handler_pc: %w", err)
	}
	catchType, err := d.readUint16()
	if err != nil {
		return h, fmt.Errorf("read catch_type: %w", err)
	}
	h.StartPC = startPC
	h.EndPC = endPC
	h.HandlerPC = handlerPC
	h.CatchType = catchType
	return h, nil
}

func (d *Decoder) readAttributes() ([]Attribute, error) {
	n, err := d.readUint16()
	if err != nil {
		return nil, fmt.Errorf("read attrs count: %w", err)
	}
	attrs := make([]Attribute, n)
	for i := range attrs {
		attr, err := d.readAttr()
		if err != nil {
			return nil, fmt.Errorf("attr%d: %w", i, err)
		}
		attrs[i] = attr
	}
	return attrs, nil
}

func (d *Decoder) readAttr() (Attribute, error) {
	nameIndex, err := d.readUint16()
	if err != nil {
		return nil, fmt.Errorf("read name_index: %w", err)
	}

	length, err := d.readUint32()
	if err != nil {
		return nil, fmt.Errorf("read attribute_length: %w", err)
	}

	var attr Attribute
	switch d.f.Consts[nameIndex].(*Utf8Const).Value {
	case "Code":
		maxStack, err := d.readUint16()
		if err != nil {
			return nil, fmt.Errorf("read max_stack: %w", err)
		}
		maxLocals, err := d.readUint16()
		if err != nil {
			return nil, fmt.Errorf("read max_locals: %w", err)
		}
		code, err := d.readByteSlice32()
		if err != nil {
			return nil, fmt.Errorf("read code: %w", err)
		}
		handlersCount, err := d.readUint16()
		if err != nil {
			return nil, fmt.Errorf("read exception_table_length")
		}
		handlers := make([]ExceptionHandler, handlersCount)
		for i := range handlers {
			h, err := d.readExceptionHandler()
			if err != nil {
				return nil, fmt.Errorf("handler%d: %w", i, err)
			}
			handlers[i] = h
		}
		attrs, err := d.readAttributes()
		if err != nil {
			return nil, fmt.Errorf("read attributes: %w", err)
		}
		attr = CodeAttribute{
			MaxStack:       maxStack,
			MaxLocals:      maxLocals,
			Code:           code,
			ExceptionTable: handlers,
			Attrs:          attrs,
		}

	default:
		buf := make([]byte, length)
		_, err = io.ReadFull(d.r, buf)
		if err != nil {
			return nil, fmt.Errorf("read info: %w", err)
		}
		attr = RawAttribute{
			NameIndex: nameIndex,
			Data:      buf,
		}
	}

	return attr, nil
}

func (d *Decoder) readField() (Field, error) {
	var f Field
	accessFlags, err := d.readUint16()
	if err != nil {
		return f, fmt.Errorf("read access_flags: %w", err)
	}
	nameIndex, err := d.readUint16()
	if err != nil {
		return f, fmt.Errorf("read name_index: %w", err)
	}
	descriptorIndex, err := d.readUint16()
	if err != nil {
		return f, fmt.Errorf("read descriptor_index: %w", err)
	}
	attrs, err := d.readAttributes()
	if err != nil {
		return f, err
	}
	f.AccessFlags = accessFlags
	f.Name = d.f.Consts[nameIndex].(*Utf8Const).Value
	f.Descriptor = d.f.Consts[descriptorIndex].(*Utf8Const).Value
	f.Attrs = attrs
	return f, nil
}

func (d *Decoder) readConst() (Const, int, error) {
	tag, err := d.r.ReadByte()
	if err != nil {
		return nil, 0, fmt.Errorf("read tag: %w", err)
	}

	var c Const
	skip := 1 // Almost all consts occupy 1 index
	switch tag {
	case 1:
		buf, err := d.readByteSlice16()
		if err != nil {
			return nil, 0, fmt.Errorf("read bytes: %w", err)
		}
		c = &Utf8Const{Value: string(buf)}
	case 3:
		v, err := d.readUint32()
		if err != nil {
			return nil, 0, fmt.Errorf("read bytes: %w", err)
		}
		c = &IntConst{Value: int32(v)}
	case 5:
		v, err := d.readUint64()
		if err != nil {
			return nil, 0, fmt.Errorf("read bytes: %w", err)
		}
		c = &LongConst{Value: int64(v)}
		skip = 2
	case 6:
		v, err := d.readUint64()
		if err != nil {
			return nil, 0, fmt.Errorf("read bytes: %w", err)
		}
		c = &DoubleConst{Value: math.Float64frombits(v)}
		skip = 2
	case 7:
		nameIndex, err := d.readUint16()
		if err != nil {
			return nil, 0, fmt.Errorf("read name_index: %w", err)
		}
		cc := &ClassConst{}
		d.deferNameResolving(nameIndex, &cc.Name)
		c = cc
	case 9, 10:
		classIndex, err := d.readUint16()
		if err != nil {
			return nil, 0, fmt.Errorf("read class_index: %w", err)
		}
		nameAndTypeIndex, err := d.readUint16()
		if err != nil {
			return nil, 0, fmt.Errorf("read name_and_type_index: %w", err)
		}
		switch tag {
		case 9:
			fc := &FieldrefConst{}
			d.deferClassNameResolving(classIndex, &fc.ClassName)
			d.deferNameAndTypeResolving(nameAndTypeIndex, &fc.Name, &fc.Descriptor)
			c = fc
		case 10:
			mc := &MethodrefConst{}
			d.deferClassNameResolving(classIndex, &mc.ClassName)
			d.deferNameAndTypeResolving(nameAndTypeIndex, &mc.Name, &mc.Descriptor)
			c = mc
		}
	case 12:
		nameIndex, err := d.readUint16()
		if err != nil {
			return nil, 0, fmt.Errorf("read name_index: %w", err)
		}
		descriptorIndex, err := d.readUint16()
		if err != nil {
			return nil, 0, fmt.Errorf("read descriptor_index: %w", err)
		}
		ntc := &NameAndTypeConst{}
		d.deferNameResolving(nameIndex, &ntc.Name)
		d.deferDescriptorResolving(descriptorIndex, &ntc.Descriptor)
		c = ntc
	default:
		return nil, 0, fmt.Errorf("unexpected tag: %d", tag)
	}

	return c, skip, nil
}

func (d *Decoder) readUint64() (uint64, error) {
	var buf [8]byte
	n, err := d.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != 8 {
		return 0, fmt.Errorf("not enough input bytes (want 8, got %d)", n)
	}
	v := binary.BigEndian.Uint64(buf[:])
	return v, nil
}

func (d *Decoder) readUint32() (uint32, error) {
	var buf [4]byte
	n, err := d.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, fmt.Errorf("not enough input bytes (want 4, got %d)", n)
	}
	v := binary.BigEndian.Uint32(buf[:])
	return v, nil
}

func (d *Decoder) readUint16() (uint16, error) {
	var buf [2]byte
	n, err := d.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != 2 {
		return 0, fmt.Errorf("not enough input bytes (want 2, got %d)", n)
	}
	v := binary.BigEndian.Uint16(buf[:])
	return v, nil
}

func (d *Decoder) readByteSlice32() ([]byte, error) {
	length, err := d.readUint32()
	if err != nil {
		return nil, fmt.Errorf("read length: %w", err)
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(d.r, buf)
	return buf, err
}

func (d *Decoder) readByteSlice16() ([]byte, error) {
	length, err := d.readUint16()
	if err != nil {
		return nil, fmt.Errorf("read length: %w", err)
	}
	buf := make([]byte, length)
	_, err = io.ReadFull(d.r, buf)
	return buf, err
}
