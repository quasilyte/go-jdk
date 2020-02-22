package ir

import (
	"github.com/quasilyte/go-jdk/jclass"
	"github.com/quasilyte/go-jdk/vmdat"
)

type Package struct {
	Classes []Class

	Out *vmdat.Package
}

type Class struct {
	Name string

	Methods []Method

	File *jclass.File
	Out  *vmdat.Class
}

type Method struct {
	Code        []Inst
	AccessFlags jclass.MethodAccessFlags

	Out *vmdat.Method
}
