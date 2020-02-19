package ir

import (
	"github.com/quasilyte/GopherJRE/jclass"
	"github.com/quasilyte/GopherJRE/vmdat"
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
	Code []Inst

	Out *vmdat.Method
}
