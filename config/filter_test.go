package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsRelevantForFilterStrings(t *testing.T) {
	p := &ProcessConfig{
		Metadata: map[string]string{
			"group": "km",
		},
	}
	assert.True(t, isRelevantForFilter(t, p, "metadata.group=km"))
	assert.True(t, isRelevantForFilter(t, p, "metadata.group==km"))
	assert.True(t, isRelevantForFilter(t, p, "metadata.group!=kmc"))
}

func isRelevantForFilter(t *testing.T, p *ProcessConfig, filterPattern string) bool {
	isRelevant, err := p.IsRelevantForFilter(filterPattern)
	assert.NoError(t, err)
	return isRelevant
}
