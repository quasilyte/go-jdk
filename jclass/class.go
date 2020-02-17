package jclass

type Method struct {
	AccessFlags     uint16
	Name            string
	DescriptorIndex uint16
	Attrs           []Attribute
}

type Field struct {
	AccessFlags     uint16
	Name            string
	DescriptorIndex uint16
	Attrs           []Attribute
}

type File struct {
	Ver           Version
	Consts        []Const
	AccessFlags   uint16
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
