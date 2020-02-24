package optimizer

import (
	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/symbol"
)

type inliner struct {
	pkg       *ir.Package
	callGraph *CallGraph
	candidats []symbol.ID
}

func RunInliner(pkg *ir.Package, g *CallGraph) {
	inliner := inliner{pkg, g, []symbol.ID{}}
	inliner.findInlineCandidats()
	for _, c := range pkg.Classes {
		for _, m := range c.Methods {
			inliner.tryInline(&m)
		}
	}
}

func (inl *inliner) findMethod(id symbol.ID) *ir.Method {
	return &inl.pkg.Classes[id.ClassIndex()].Methods[id.MemberIndex()]
}

func (inl *inliner) findRecursion(visited map[symbol.ID]bool, cur symbol.ID, original symbol.ID) bool {
	visited[cur] = true
	if cur == original {
		return true
	}
	for _, n := range inl.callGraph.Neighbours(cur) {
		if visited[n] {
			continue
		}
		if inl.findRecursion(visited, n, original) {
			return true
		}
	}
	return false
}

func (inl *inliner) hasCyclicDependency(s symbol.ID) bool {
	visited := map[symbol.ID]bool{}
	for _, n := range inl.callGraph.Neighbours(s) {
		if visited[n] {
			return true
		}
		if inl.findRecursion(visited, n, s) {
			return true
		}
	}
	return false
}

func (inl *inliner) tryInline(m *ir.Method) {
	for ind, inst := range m.Code {
		if inst.Kind == ir.InstCallStatic {
			callee := inst.Args[0]
			if contains(inl.candidats, callee.SymbolID()) {
				inl.inline(m, inl.findMethod(callee.SymbolID()), ind)
			}
		}
	}
}

func (inl *inliner) inline(m, m2 *ir.Method, ind int) {
	insts2 := make([]ir.Inst, len(m2.Code))
	copy(insts2, m2.Code)
	for i := range insts2 {
		if insts2[i].Kind == ir.InstRet {
			insts2[i] = ir.Inst{
				Kind: ir.InstJump,
				Args: []ir.Arg{
					ir.Arg{
						Kind:  ir.ArgBranch,
						Value: int64(ind),
					},
				},
			}
		}
	}
	m.Code = append(m.Code[:ind], append(insts2, m.Code[ind+1:]...)...)

	for _, inst := range m.Code {
		if inst.Kind == ir.InstJump {
			if inst.Args[0].Value >= int64(ind) {
				inst.Args[0].Value += int64(len(insts2))
			}
		}
	}
	inl.tryInline(m)
}

func (inl *inliner) findInlineCandidats() {
	for _, c := range inl.pkg.Classes {
		for _, m := range c.Methods {
			if !inl.hasCyclicDependency(m.Out.ID) {
				inl.candidats = append(inl.candidats, m.Out.ID)
			}
		}
	}
}
