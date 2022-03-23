package pointer_test

import (
	"testing"

	"github.com/morphysm/famed-github-backend/pkg/pointer"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Parallel()

	value := pointer.String("FUN")
	assert.Equal(t, "FUN", *value)
}
