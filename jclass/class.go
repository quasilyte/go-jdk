package jclass

type Method struct {
	AccessFlags MethodAccessFlags
	Name        string
	Descriptor  string
	Attrs       []Attribute
}

type Field struct {
	AccessFlags FieldAccessFlags
	Name        string
	Descriptor  string
	Attrs       []Attribute
}

type File struct {
	Ver           Version
	Consts        []Const
	AccessFlags   AccessFlags
	ThisClassName string
	SuperClass    uint16
	Interfaces    []uint16
	Fields        []Field
	Methods       []Method
	Attrs         []Attribute
}

type Version struct {
	Minor uint16
	Major uint16
}

type ExceptionHandler struct {
	StartPC   uint16
	EndPC     uint16
	HandlerPC uint16
	CatchType uint16
}

type StackMapFrame struct {
	Offset     uint32
	StackDepth uint16
}
