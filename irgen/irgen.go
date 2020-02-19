package irgen

import (
	"fmt"

	"github.com/quasilyte/GopherJRE/ir"
	"github.com/quasilyte/GopherJRE/vmdat"
)

func Generate(st *vmdat.State, packages []*ir.Package) error {
	var g generator
	g.state = st
	for _, pkg := range packages {
		for i := range pkg.Classes {
			c := &pkg.Classes[i]
			g.f = c.File
			for j := range c.Methods {
				m := &c.Methods[j]
				if m.Out.Name == "<init>" {
					continue
				}
				if err := g.Generate(j, m); err != nil {
					return fmt.Errorf("%s: %s.%s: %v",
						pkg.Out.Name, c.Name, m.Out.Name, err)
				}
			}
		}
	}
	return nil
}
