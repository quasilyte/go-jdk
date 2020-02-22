package jdeps

import (
	"github.com/quasilyte/go-jdk/jclass"
)

// ClassDependencies returns a list of package names that given class depends on.
// The list is not sorted.
func ClassDependencies(c *jclass.File) []string {
	return nil
}

func classDependencies(c *jclass.File) []string {
	// FIXME: should return package names instead of type names.

	visitor := makeVisitor(c)
	visitor.addClass(c.SuperClass)
	visitor.addClasses(c.Interfaces)
	visitor.addFields()
	visitor.addMethods()
	// TODO: deal with inner classes
	return visitor.keys()
}

type depsFinder struct {
	deps      map[string]struct{}
	classFile *jclass.File
}

func makeVisitor(c *jclass.File) depsFinder {
	return depsFinder{map[string]struct{}{}, c}
}

func (f *depsFinder) keys() []string {
	keys := make([]string, 0, len(f.deps))
	for k := range f.deps {
		keys = append(keys, k)
	}
	return keys
}

func (f *depsFinder) addDependency(d string) {
	f.deps[d] = struct{}{}
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

func (f *depsFinder) addFields() {
	for _, field := range f.classFile.Fields {
		// TODO: walk attributes
		desc := jclass.FieldDescriptor(field.Descriptor)
		typ := desc.GetType()
		f.addType(typ)
	}
}

func (f *depsFinder) addMethods() {
	for _, m := range f.classFile.Methods {
		desc := jclass.MethodDescriptor(m.Descriptor)
		// TODO: walk attributes

		desc.WalkParams(func(typ jclass.DescriptorType) {
			f.addType(typ)
		})
		f.addType(desc.ReturnType())
	}
}
