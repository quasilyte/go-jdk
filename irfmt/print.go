package irfmt

import (
	"strings"

	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/vmdat"
)

func Sprint(st *vmdat.State, inst ir.Inst) string {
	var buf strings.Builder
	if inst.Dst.Kind != 0 {
		buf.WriteString(inst.Dst.String())
		buf.WriteString(" = ")
	}
	buf.WriteString(inst.Kind.String())
	for _, arg := range inst.Args {
		buf.WriteByte(' ')
		switch arg.Kind {
		case ir.ArgSymbolID:
			id := arg.SymbolID()
			pkg := st.Packages[id.PackageIndex()]
			class := pkg.Classes[id.ClassIndex()]
			name := class.Methods[id.MemberIndex()].Name
			buf.WriteString(name)
		default:
			buf.WriteString(arg.String())
		}
	}
	return buf.String()
}
