package pointer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/morphysm/famed-github-backend/pkg/pointer"
)

func TestTime(t *testing.T) {
	t.Parallel()

	value := pointer.Time(time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC))
	assert.Equal(t, time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), *value)
}
