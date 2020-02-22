package jclass

import (
	"fmt"
	"strings"
)

type MethodDescriptor string

func (d MethodDescriptor) ReturnType() DescriptorType {
	s := string(d)
	start := strings.LastIndexByte(s, ')') + 1
	var result DescriptorType
	walkDescriptor(s[start:], func(typ DescriptorType) {
		result = typ
	})
	return result
}

func (d MethodDescriptor) WalkParams(visit func(typ DescriptorType)) {
	s := string(d)
	walkDescriptor(s[len("("):strings.IndexByte(s, ')')], visit)
}

func (d MethodDescriptor) SignatureString(name string) string {
	var params []string
	d.WalkParams(func(typ DescriptorType) {
		params = append(params, typ.String())
	})
	return fmt.Sprintf("%s %s(%s)",
		d.ReturnType().String(), name, strings.Join(params, ", "))
}

func walkDescriptor(s string, visit func(typ DescriptorType)) {
	var typ DescriptorType
	i := 0
	for i < len(s) {
		switch s[i] {
		case ')':
			return
		case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', 'V':
			typ.Kind = s[i]
			visit(typ)
			i++
			typ.Dims = 0
		case 'L':
			end := strings.IndexByte(s[i:], ';') + i
			typ.Kind = 'L'
			typ.Name = s[i+1 : end]
			visit(typ)
			i = end + len(";")
			typ.Dims = 0
		case '[':
			typ.Dims++
			i++
		}
	}
}

type FieldDescriptor string

func (d FieldDescriptor) GetType() DescriptorType {
	var result DescriptorType
	walkDescriptor(string(d), func(typ DescriptorType) {
		result = typ
	})
	return result
}

type DescriptorType struct {
	Name string
	Kind byte
	Dims int
}

func (typ DescriptorType) IsReference() bool {
	return typ.Kind == 'L'
}

func (typ DescriptorType) String() string {
	var p string
	switch typ.Kind {
	case 'B':
		p = "byte"
	case 'C':
		p = "char"
	case 'D':
		p = "double"
	case 'F':
		p = "float"
	case 'I':
		p = "int"
	case 'J':
		p = "long"
	case 'S':
		p = "short"
	case 'Z':
		p = "boolean"
	case 'V':
		p = "void"
	case 'L':
		p = typ.Name
	}
	if typ.Dims != 0 {
		p += strings.Repeat("[]", typ.Dims)
	}
	return p
}
