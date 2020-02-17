package jclass

// AccessFlags is a mask of flags used to denote access permissions to and properties
// of this class or interface.
//
// The interpretation of each flag, when set, is specified below:
//	ACC_PUBLIC      0x0001  Declared public; may be accessed from outside its package.
//	ACC_FINAL       0x0010  Declared final; no subclasses allowed.
//	ACC_SUPER       0x0020  Treat superclass methods specially when invoked by the invokespecial instruction.
//	ACC_INTERFACE   0x0200  Is an interface, not a class.
//	ACC_ABSTRACT    0x0400  Declared abstract; must not be instantiated.
//	ACC_SYNTHETIC   0x1000  Declared synthetic; not present in the source code.
//	ACC_ANNOTATION  0x2000  Declared as an annotation type.
//	ACC_ENUM        0x4000  Declared as an enum type.
//	ACC_MODULE      0x8000  Is a module, not a class or interface.
type AccessFlags uint16

// IsPublic reports whether a symbol can be accessed from outside its package.
func (af AccessFlags) IsPublic() bool { return af&0x0001 != 0 }

// IsFinal reports whether a symbol is declared final (no subclasses allowed).
func (af AccessFlags) IsFinal() bool { return af&0x0010 != 0 }

// IsSuper reports whether to treat superclass methods specially when invoked
// by the invokespecial instruction.
func (af AccessFlags) IsSuper() bool { return af&0x0020 != 0 }

// IsInterface reports whether a symbol is interface.
func (af AccessFlags) IsInterface() bool { return af&0x0200 != 0 }

// IsAbstract reports whether a symbol is declared abstract (must not be instantiated).
func (af AccessFlags) IsAbstract() bool { return af&0x0400 != 0 }

// IsSynthetic reports whether a symbol is auto-generated (not present in the source code).
func (af AccessFlags) IsSynthetic() bool { return af&0x1000 != 0 }

// IsAnnotation reports whether a symbol is annotation.
func (af AccessFlags) IsAnnotation() bool { return af&0x2000 != 0 }

// IsEnum reports whether a symbol is enum.
func (af AccessFlags) IsEnum() bool { return af&0x4000 != 0 }

// IsModule reports whether a symbol is module.
func (af AccessFlags) IsModule() bool { return af&0x8000 != 0 }
