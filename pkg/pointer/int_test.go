package pointer_test

import (
	"testing"

	"github.com/morphysm/famed-github-backend/pkg/pointer"
	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	t.Parallel()

	value := pointer.Int(1)
	assert.Equal(t, 1, *value)
}

func TestInt64(t *testing.T) {
	t.Parallel()

	value := pointer.Int64(1)
	assert.Equal(t, int64(1), *value)
}
