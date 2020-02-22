package vmdat

import (
	"sort"

	"github.com/quasilyte/go-jdk/jclass"
	"github.com/quasilyte/go-jdk/symbol"
)

type GoFunc struct {
	Name string
	Addr uintptr
}

type State struct {
	Packages      []*Package
	pkgname2index map[string]uint32
	GoFuncs       []GoFunc
	goname2index  map[string]uint32
}

func (st *State) Init() {
	st.Packages = make([]*Package, 0, 32)
	st.pkgname2index = map[string]uint32{}
	st.goname2index = map[string]uint32{}
}

func (st *State) FindPackage(name string) *Package {
	index, ok := st.pkgname2index[name]
	if !ok {
		return nil
	}
	return st.Packages[index]
}

func (st *State) AddGoFunc(fn GoFunc) {
	index := uint32(len(st.GoFuncs))
	st.goname2index[fn.Name] = index
	st.GoFuncs = append(st.GoFuncs, fn)
}

func (st *State) GoFuncIndex(name string) int {
	return int(st.goname2index[name])
}

func (st *State) NewPackage(name string) *Package {
	index := uint32(len(st.Packages))
	pkg := &Package{
		ID:   index,
		Name: name,
	}
	st.pkgname2index[name] = index
	st.Packages = append(st.Packages, pkg)
	return pkg
}

type Package struct {
	ID      uint32
	Name    string
	Classes []Class
}

type Class struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name        string
	Descriptor  string
	AccessFlags jclass.MethodAccessFlags
	FrameSlots  int
	ID          symbol.ID
	Code        []byte
}

func (p *Package) FindClass(name string) *Class {
	i := sort.Search(len(p.Classes), func(i int) bool {
		return p.Classes[i].Name >= name
	})
	if i >= len(p.Classes) || p.Classes[i].Name != name {
		return nil // Not found
	}
	return &p.Classes[i]
}

func (c *Class) FindMethod(name, descriptor string) *Method {
	// Can use binary search because methods are sorted by name.
	i := sort.Search(len(c.Methods), func(i int) bool {
		return c.Methods[i].Name >= name
	})
	if i >= len(c.Methods) || c.Methods[i].Name != name {
		return nil // Not found
	}
	// Prefer a method that matches a specified descriptor.
	for j := i; j < len(c.Methods) && c.Methods[j].Name == name; j++ {
		if c.Methods[j].Descriptor == descriptor {
			return &c.Methods[j]
		}
	}
	// As a special case, allow "" to match any signature.
	if descriptor == "" {
		return &c.Methods[i]
	}
	return nil
}
