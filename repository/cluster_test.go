package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClusterRepo(t *testing.T) {
	h1 := NewClusterRepo()
	h2 := NewClusterRepo()
	assert.Equal(t, h1, h2)
}
