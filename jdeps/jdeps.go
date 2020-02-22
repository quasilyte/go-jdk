package jdeps

import (
	"strings"

	"github.com/quasilyte/go-jdk/jclass"
)

// ClassDependencies returns a list of package names that given class depends on.
// The list is not sorted.
func ClassDependencies(c *jclass.File) []string {
	visitor := depsFinder{map[string]struct{}{}, c}
	visitor.addClass(c.SuperClass)
	visitor.addClasses(c.Interfaces)
	visitor.walkConsts()
	visitor.walkFields()
	visitor.walkMethods()
	// TODO: deal with inner classes
	return visitor.keys()
}

type depsFinder struct {
	deps      map[string]struct{}
	classFile *jclass.File
}

func (f *depsFinder) keys() []string {
	keys := make([]string, 0, len(f.deps))
	for k := range f.deps {
		keys = append(keys, k)
	}
	return keys
}

func (f *depsFinder) walkConsts() {
	for _, c := range f.classFile.Consts {
		switch c := c.(type) {
		case *jclass.MethodrefConst:
			f.addDependency(c.ClassName)
			f.walkMethodDescriptor(c.Descriptor)
		}
	}
}

func (f *depsFinder) addDependency(d string) {
	i := strings.LastIndexByte(d, '/')
	if i != -1 {
		f.deps[d[:i]] = struct{}{}
	}
}

func (f *depsFinder) addClass(i uint16) {
	if i != 0 {
		if class, ok := f.classFile.Consts[i].(*jclass.ClassConst); ok {
			f.addDependency(class.Name)
		}
	}
}

func (f *depsFinder) addClasses(indices []uint16) {
	for _, i := range indices {
		f.addClass(i)
	}
}

func (f *depsFinder) addType(typ jclass.DescriptorType) {
	if typ.IsReference() {
		f.addDependency(typ.Name)
	}
}

func (f *depsFinder) walkFields() {
	for _, field := range f.classFile.Fields {
		// TODO: walk attributes
		desc := jclass.FieldDescriptor(field.Descriptor)
		typ := desc.GetType()
		f.addType(typ)
	}
}

func (f *depsFinder) walkMethods() {
	for _, m := range f.classFile.Methods {
		// TODO: walk attributes
		f.walkMethodDescriptor(m.Descriptor)
	}
}

func (f *depsFinder) walkMethodDescriptor(s string) {
	desc := jclass.MethodDescriptor(s)
	desc.WalkParams(func(typ jclass.DescriptorType) {
		f.addType(typ)
	})
	f.addType(desc.ReturnType())
}
