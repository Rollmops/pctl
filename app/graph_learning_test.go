package app

import (
	"fmt"
	"github.com/yourbasic/graph"
	"testing"
)

func TestGraphLearning(t *testing.T) {
	gm := graph.New(5)
	gm.AddBoth(0, 1) // km -> mls

	gm.AddBoth(1, 2) // pes -> km
	gm.AddBoth(1, 3) // pec -> km

	gm.AddBoth(3, 4) // bla -> pec

	g := graph.Sort(gm)
	dist := make([]int, g.Order())
	graph.BFS(g, 0, func(v, w int, _ int64) {
		fmt.Println(v, "to", w)
		dist[w] = dist[v] + 1
	})
	fmt.Println("dist:", dist)

	fmt.Println(graph.MST(gm))
}
