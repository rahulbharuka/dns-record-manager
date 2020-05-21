package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerRepo(t *testing.T) {
	h1 := NewServerRepo()
	h2 := NewServerRepo()
	assert.Equal(t, h1, h2)
}
