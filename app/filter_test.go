package app

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
	filter, err := NewFilter(filterPattern)
	assert.NoError(t, err)
	isRelevant, err := filter.IsMatchingProcess(&Process{Config: p})
	assert.NoError(t, err)
	return isRelevant
}
