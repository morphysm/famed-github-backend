package arrays_test

import (
	"testing"

	"github.com/morphysm/famed-github-backend/pkg/arrays"

	"github.com/stretchr/testify/assert"
)

func TestArraysRemoveElement(t *testing.T) {
	t.Parallel()

	slice := []int{1, 2, 3, 4, 5}

	slice = arrays.Remove(slice, 2)

	assert.Equal(t, []int{1, 2, 4, 5}, slice)
}
