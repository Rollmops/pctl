package app

import (
	"fmt"
	"github.com/yourbasic/graph"
)

func ValidateAcyclicDependencies() error {
	mapping := make(map[string]int)
	for i, _config := range CurrentContext.Config.ProcessConfigs {
		mapping[_config.Name] = i
	}

	gm := graph.New(len(CurrentContext.Config.ProcessConfigs))
	for _, _config := range CurrentContext.Config.ProcessConfigs {
		for _, n := range _config.DependsOn {
			gm.Add(mapping[_config.Name], mapping[n])
		}
	}
	if !graph.Acyclic(gm) {
		return fmt.Errorf("process dependency configuration is not acyclic")
	}
	return nil
}
