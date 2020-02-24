package optimizer

import (
	"fmt"

	"github.com/quasilyte/go-jdk/ir"
	"github.com/quasilyte/go-jdk/symbol"
)

type testPass struct {
	pkg       *ir.Package
	callGraph *CallGraph
}

func ExecuteTestPass(pkg *ir.Package, g *CallGraph) {
	pass := testPass{pkg, g}
	for _, node := range pass.callGraph.Nodes() {
		pass.printNeighbours(node)
	}
}

func (p *testPass) printNeighbours(id symbol.ID) {
	fmt.Printf("%v\n", p.callGraph.Neighbours(id))
}
