package jdeps

import (
	"github.com/quasilyte/GopherJRE/jclass"
)

func ClassDependencies(c *jclass.File) []string {
	visitor := makeVisitor(c)
	visitor.addClass(c.SuperClass)
	visitor.addClasses(c.Interfaces)
	visitor.addFields()
	visitor.addMethods()
	// TODO: deal with inner classes
	return visitor.keys()
}

type visitorCtx struct {
	deps      map[string]struct{}
	classFile *jclass.File
}

func makeVisitor(c *jclass.File) visitorCtx {
	return visitorCtx{map[string]struct{}{}, c}
}

func (ctx *visitorCtx) keys() []string {
	keys := make([]string, 0, len(ctx.deps))
	for k, _ := range ctx.deps {
		keys = append(keys, k)
	}
	return keys
}

func (ctx *visitorCtx) addDependency(d string) {
	ctx.deps[d] = struct{}{}
}

func (ctx *visitorCtx) addClass(i uint16) {
	if i != 0 {
		if class, ok := ctx.classFile.Consts[i].(*jclass.ClassConst); ok {
			ctx.addDependency(class.Name)
		}
	}
}

func (ctx *visitorCtx) addClasses(indices []uint16) {
	for _, i := range indices {
		ctx.addClass(i)
	}
}

func (ctx *visitorCtx) addType(typ jclass.DescriptorType) {
	if typ.IsReference() {
		ctx.addDependency(typ.Name)
	}
}

func (ctx *visitorCtx) addFields() {
	for _, f := range ctx.classFile.Fields {
		// TODO: walk attributes
		desc := jclass.FieldDescriptor(f.Descriptor)
		typ := desc.GetType()
		ctx.addType(typ)
	}
}

func (ctx *visitorCtx) addMethods() {
	for _, m := range ctx.classFile.Methods {
		desc := jclass.MethodDescriptor(m.Descriptor)
		// TODO: walk attributes

		desc.WalkParams(func(typ jclass.DescriptorType) {
			ctx.addType(typ)
		})
		ctx.addType(desc.ReturnType())
	}
}
