package app

import (
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkReadProcessEnvironment(b *testing.B) {
	defer func() {
		assert.NoError(b, Run([]string{"pctl", "kill", ":"}))
	}()
	assert.NoError(b, test.SetConfigEnvPath("benchmark.yaml"))
	assert.NoError(b, CurrentContext.InitializeRunningProcessInfo())

	assert.NoError(b, Run([]string{"pctl", "start", ":"}))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = CurrentContext.Cache.Refresh()
	}
}
