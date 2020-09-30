package graph

type ErrCycle struct {
	Cycle []string
}

func (e *ErrCycle) Error() string {
	return "cycle detected"
}

type context struct {
	graph    *map[string][]string
	marked   map[string]bool
	finished map[string]bool
	path     []string
	result   []string
}

func initialize(graph *map[string][]string) context {

	ctx := context{
		graph:    graph,
		marked:   make(map[string]bool),
		finished: make(map[string]bool),
		result:   make([]string, 0, len(*graph)),
	}

	for node := range *graph {
		ctx.marked[node] = false
		ctx.finished[node] = false
	}

	return ctx
}

func (ctx *context) mark(node string) error {

	if ctx.marked[node] {
		return &ErrCycle{
			ctx.path,
		}
	}
	ctx.marked[node] = true

	return nil
}

func (ctx *context) finish(node string) {
	ctx.finished[node] = true
	ctx.result = append(ctx.result, node)
}

func (ctx *context) isUnfinished(node string) bool {
	return !ctx.finished[node]
}

func (ctx *context) findFirstUnfinishedNode() (string, bool) {
	for node, isFinished := range ctx.finished {
		if !isFinished {
			return node, true
		}
	}
	return "", false
}

func (ctx *context) visit(node string) error {

	if ctx.isUnfinished(node) {

		ctx.path = append(ctx.path, node)
		if err := ctx.mark(node); err != nil {
			return err
		}

		for _, neighbor := range (*ctx.graph)[node] {
			if err := ctx.visit(neighbor); err != nil {
				return err
			}
		}
		ctx.finish(node)
		ctx.path = ctx.path[:len(ctx.path)-1]
	}
	return nil
}

func Topsort(graph *map[string][]string) ([]string, error) {

	ctx := initialize(graph)

	node, isUnfinished := ctx.findFirstUnfinishedNode()
	for isUnfinished {
		if err := ctx.visit(node); err != nil {
			return nil, err
		}
		node, isUnfinished = ctx.findFirstUnfinishedNode()
	}

	return ctx.result, nil
}
