package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchSimpleName(t *testing.T) {
	nsp, err := NewGroupNameSpecifier("name")
	assert.NoError(t, err)

	assert.True(t, nsp.IsMatchingGroupAndName("", "name"))
	assert.True(t, nsp.IsMatchingGroupAndName("some-group", "name"))

	assert.False(t, nsp.IsMatchingGroupAndName("", "noname"))
	assert.False(t, nsp.IsMatchingGroupAndName("some-group", "noname"))
}

func TestMatchCompleteGroup(t *testing.T) {
	nsp, err := NewGroupNameSpecifier("group:")
	assert.NoError(t, err)

	assert.True(t, nsp.IsMatchingGroupAndName("group", "1"))
	assert.True(t, nsp.IsMatchingGroupAndName("group", "2"))

	assert.False(t, nsp.IsMatchingGroupAndName("other-group", "1"))
	assert.False(t, nsp.IsMatchingGroupAndName("other-group", "2"))
}

func TestMatchGroupAndName(t *testing.T) {
	nsp, err := NewGroupNameSpecifier("group:name")
	assert.NoError(t, err)

	assert.True(t, nsp.IsMatchingGroupAndName("group", "name"))
	assert.False(t, nsp.IsMatchingGroupAndName("group", "other-name"))
	assert.False(t, nsp.IsMatchingGroupAndName("other-group", "name"))
}

func TestMatchStar(t *testing.T) {
	nsp, err := NewGroupNameSpecifier("*")
	assert.NoError(t, err)

	assert.True(t, nsp.IsMatchingGroupAndName("group", "name"))
	assert.True(t, nsp.IsMatchingGroupAndName("", "name"))
}
