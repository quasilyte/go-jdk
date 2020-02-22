package loader

import (
	"fmt"
	"sort"

	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/jclass"
	"github.com/quasilyte/go-jdk/symbol"
	"github.com/quasilyte/go-jdk/vmdat"
)

type Config struct {
	ClassPath []string
}

func LoadClass(st *vmdat.State, filename string, cfg *Config) ([]*ir.Package, error) {
	classfile, err := decodeClassFile(filename)
	if err != nil {
		return nil, err
	}
	_, pkgName := splitName(classfile.ThisClassName)
	pkgFiles := []*jclass.File{classfile}
	return loadPackageSet(st, pkgName, pkgFiles, cfg)
}

// LoadPackage loads a package with name pkgName as well as all its dependencies.
// Loaded packages are returned to be compiled.
// All of these packages are also added to a provided VM state.
//
// If a package was already loaded, nil is returned.
func LoadPackage(st *vmdat.State, pkgName string, cfg *Config) ([]*ir.Package, error) {
	if st.FindPackage(pkgName) != nil {
		return nil, nil // Already loaded
	}

	pkgFiles, err := readClassFiles(pkgName, cfg)
	if err != nil {
		return nil, fmt.Errorf("read %q class files: %v", pkgName, err)
	}
	return loadPackageSet(st, pkgName, pkgFiles, cfg)
}

func loadPackageSet(st *vmdat.State, pkgName string, initial []*jclass.File, cfg *Config) ([]*ir.Package, error) {
	pkg := createPackage(st, pkgName, initial)
	deps := findDependencies(st, initial)

	toLoad := make([]*ir.Package, 0, len(deps)+1)
	toLoad = append(toLoad, pkg)
	for len(deps) > 0 {
		d := deps[len(deps)-1]
		if d == "java/lang" {
			// We don't support stdlib packages yet.
			// FIXME: remove this after we handle them properly.
			deps = deps[:len(deps)-1]
			continue
		}
		files, err := readClassFiles(d, cfg)
		if err != nil {
			return nil, fmt.Errorf("find %q package: %v", d, err)
		}
		pkg := createPackage(st, d, files)
		toLoad = append(toLoad, pkg)
		depDeps := findDependencies(st, files)
		deps = append(deps[:len(deps)-1], depDeps...)
	}

	return toLoad, nil
}

func createPackage(st *vmdat.State, name string, files []*jclass.File) *ir.Package {
	pkg := ir.Package{Out: st.NewPackage(name)}

	for _, f := range files {
		sort.Slice(f.Methods, func(i, j int) bool {
			return f.Methods[i].Name < f.Methods[j].Name
		})

		switch {
		case f.AccessFlags.IsEnum():
			panic("enums types are not implemented")
		case f.AccessFlags.IsInterface():
			panic("interface types are not implemented")
		case f.AccessFlags.IsAnnotation():
			panic("annotation types are not implemented")
		default: // Otherwise it's a normal class
			methods := make([]ir.Method, len(f.Methods))
			className, _ := splitName(f.ThisClassName)
			pkg.Classes = append(pkg.Classes, ir.Class{
				Name:    className,
				File:    f,
				Methods: methods,
			})
		}
	}

	sort.Slice(pkg.Classes, func(i, j int) bool {
		return pkg.Classes[i].Name < pkg.Classes[j].Name
	})

	pkg.Out.Classes = make([]vmdat.Class, len(pkg.Classes))
	for i := range pkg.Classes {
		irClass := &pkg.Classes[i]
		c := &pkg.Out.Classes[i]
		c.Name = irClass.Name
		c.Methods = make([]vmdat.Method, len(irClass.Methods))
		f := irClass.File
		for j := range irClass.Methods {
			m := &c.Methods[j]
			m.Name = f.Methods[j].Name
			m.Descriptor = f.Methods[j].Descriptor
			m.AccessFlags = f.Methods[j].AccessFlags
			m.ID = symbol.NewID(uint64(pkg.Out.ID), uint64(i), uint64(j))
			irClass.Methods[j].Out = m
		}
		irClass.Out = c
	}

	return &pkg
}
