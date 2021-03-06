package pointer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

func TestString(t *testing.T) {
	t.Parallel()

	value := pointer.String("FUN")
	assert.Equal(t, "FUN", *value)
}

func TestToString(t *testing.T) {
	t.Parallel()

	value := pointer.ToString(pointer.String("FUN"))
	assert.Equal(t, "FUN", value)
}

func TestToString_Nil(t *testing.T) {
	t.Parallel()

	value := pointer.ToString(nil)
	assert.Equal(t, "", value)
}
