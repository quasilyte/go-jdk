package irgen

import (
	"sort"

	"github.com/quasilyte/GopherJRE/ir"
	"github.com/quasilyte/GopherJRE/jclass"
)

func ConvertClass(f *jclass.File) (*ir.Class, error) {
	fullname := f.ThisClassName
	name, pkgName := splitName(fullname)
	c := &ir.Class{Name: name, PkgName: pkgName}
	c.Methods = convertMethods(f)
	sort.Slice(c.Methods, func(i, j int) bool {
		return c.Methods[i].Name < c.Methods[j].Name
	})
	return c, nil
}

func convertMethods(f *jclass.File) []ir.Method {
	var g generator
	methods := make([]ir.Method, len(f.Methods))
	for i := range f.Methods {
		methods[i] = convertMethod(f, &f.Methods[i], &g)
	}
	return methods
}

func convertMethod(f *jclass.File, m *jclass.Method, g *generator) ir.Method {
	name := m.Name
	var code []ir.Inst
	if name != "<init>" {
		code = g.ConvertMethod(f, m)
	}
	return ir.Method{
		Name: name,
		Code: code,
	}
}
