package optimizer

import (
	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/symbol"
)

func contains(a []symbol.ID, s symbol.ID) bool {
	for _, el := range a {
		if el == s {
			return true
		}
	}
	return false
}

type CallGraph struct {
	graph           map[symbol.ID][]symbol.ID
	accessableNodes []symbol.ID
}

func (g *CallGraph) Nodes() []symbol.ID {
	return g.accessableNodes
}

func (g *CallGraph) Neighbours(node symbol.ID) []symbol.ID {
	return g.graph[node]
}

func BuildCallGraph(pkg *ir.Package) *CallGraph {
	g := CallGraph{map[symbol.ID][]symbol.ID{}, []symbol.ID{}}
	for _, c := range pkg.Classes {
		g.walkClass(&c)
	}
	return &g
}

func (g *CallGraph) walkClass(c *ir.Class) {
	for _, m := range c.Methods {
		g.accessableNodes = append(g.accessableNodes, m.Out.ID)
	}
	for _, m := range c.Methods {
		g.walkMethod(&m)
	}
}

func (g *CallGraph) walkMethod(m *ir.Method) {
	neighbours := []symbol.ID{}
	for _, inst := range m.Code {
		if inst.Kind == ir.InstCallStatic {
			callee := inst.Args[0].SymbolID()
			if !contains(neighbours, callee) && contains(g.accessableNodes, callee) {
				neighbours = append(neighbours, callee)
			}
		}
	}
	g.graph[m.Out.ID] = neighbours
}
