package graph

import (
	"errors"
	"testing"
)

func TestInitializeContextByGivenGraph(t *testing.T) {

	graph := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": nil,
	}
	n := len(graph)

	ctx := initialize(&graph)

	if c := len(ctx.marked); c != n {
		t.Fatalf("list length of marked nodes must be equal number of nodes: %d != %d", c, n)
	}
	for _, v := range ctx.marked {
		if v {
			t.Fatalf("each node must not be marked")
		}
	}

	if c := len(ctx.finished); c != n {
		t.Fatalf("list length of finished nodes must be equal number of nodes: %d != %d", c, n)
	}
	for _, v := range ctx.finished {
		if v {
			t.Fatalf("each node must no be finished")
		}
	}

	if len(ctx.path) != 0 {
		t.Fatalf("(current) path must be empty")
	}

	if len(ctx.result) != 0 {
		t.Fatalf("result must be empty")
	}
}

func TestMarkNonMarkedNodeSucceeds(t *testing.T) {

	graph := map[string][]string{
		"A": nil,
	}
	ctx := initialize(&graph)

	err := ctx.mark("A")
	if err != nil {
		t.Fatalf("mark node must succeed")
	}
}

func TestMarkAlreadyMarkedNodeFails(t *testing.T) {

	graph := map[string][]string{
		"A": nil,
	}
	ctx := initialize(&graph)

	_ = ctx.mark("A")
	err := ctx.mark("A")
	if err == nil || err.Error() != "cycle detected" {
		t.Fatalf("mark already marked node must fail")
	}
}

func TestFinishNodeShouldAppendToResult(t *testing.T) {

	graph := map[string][]string{
		"A": nil,
	}
	ctx := initialize(&graph)

	ctx.finish("A")
	if len(ctx.result) != 1 || ctx.result[0] != "A" {
		t.Fatalf("finished node must appear on result")
	}
}

func TestVisitShouldVisitAllAdjacentNodesOfGivenNode(t *testing.T) {

	graph := map[string][]string{
		"A": {"B", "D"},
		"B": nil,
		"C": nil,
		"D": nil,
	}
	ctx := initialize(&graph)

	if err := ctx.visit("A"); err != nil {
		t.Fatalf("visit fails unexpected with: %s", err)
	}

	expected := map[string]bool{
		"A": true,
		"B": true,
		"C": false,
		"D": true,
	}

	equals := true
	for k, v := range expected {
		equals = equals && ctx.marked[k] == v
	}
	if !equals {
		t.Fatalf("marked nodes must be equal: %v != %v", ctx.marked, expected)
	}
}

func TestAfterVisitReturnNodeShouldBeFinished(t *testing.T) {

	graph := map[string][]string{
		"A": nil,
	}
	ctx := initialize(&graph)

	_ = ctx.visit("A")
	if !ctx.finished["A"] {
		t.Fatalf("node must be finished after visit call returned")
	}
}

func TestAfterVisitReturnPathShouldBeUnmodified(t *testing.T) {

	graph := map[string][]string{
		"A": {"B"},
		"B": nil,
	}
	ctx := initialize(&graph)

	ctx.marked["A"] = true
	ctx.path = append(ctx.path, "A")

	_ = ctx.visit("B")
	if !ctx.finished["B"] {
		t.Fatalf("node must be finished after visit call returned")
	}
	if len(ctx.path) != 1 || ctx.path[0] != "A" {
		t.Fatalf("(current) path must be unmodified after finished node visit")
	}
}

func TestTopsortSingle(t *testing.T) {

	graph := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {"D", "E"},
		"D": nil,
		"E": nil,
	}

	sorted, err := Topsort(&graph)
	if err != nil {
		t.Fatalf("topological order must exist")
	}

	assertIsContainedInOrder(t, &sorted, &[]string{"D", "C", "B", "A"})
	assertIsContainedInOrder(t, &sorted, &[]string{"E", "C", "B", "A"})
}

func TestTopsortMultiple(t *testing.T) {

	graph := map[string][]string{
		"A": {"B"},
		"B": nil,
		"C": {"D"},
		"D": nil,
	}

	sorted, err := Topsort(&graph)
	if err != nil {
		t.Fatalf("topological order must exist")
	}

	assertIsContainedInOrder(t, &sorted, &[]string{"B", "A"})
	assertIsContainedInOrder(t, &sorted, &[]string{"D", "C"})
}

func TestTopsortShouldDetectCycle(t *testing.T) {

	graph := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {"D", "E"},
		"D": {"A"},
		"E": nil,
	}

	if _, err := Topsort(&graph); err == nil {
		t.Fatalf("must detect cycle")
	}
}

func TestTopsortShouldReturnErrorWithCycle(t *testing.T) {

	graph := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {"D", "E"},
		"D": {"A"},
		"E": nil,
	}

	var err *ErrCycle
	if _, e := Topsort(&graph); !errors.As(e, &err) {
		t.Fatalf("must detect cycle")
	}
	if err.Error() != "cycle detected" {
		t.Fatalf("error message must be 'cycle detected'")
	}
	if len(err.Cycle) != 5 {
		t.Fatalf("error must contain cycle, but is: %v", err.Cycle)
	}
}

func indexOf(slice *[]string, item string) int {
	for i, v := range *slice {
		if v == item {
			return i
		}
	}
	return -1
}

func assertIsContainedInOrder(t *testing.T, slice *[]string, items *[]string) {

	var positions = make([]int, len(*items))
	for k, v := range *items {
		positions[k] = indexOf(slice, v)
	}

	prev := -1
	for _, v := range positions {
		if v < prev {
			t.Fatalf("must be contained in order")
		}
		prev = v
	}
}
