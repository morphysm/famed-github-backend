package sync_test

import (
	"sync"
	"testing"

	libSync "github.com/morphysm/famed-github-backend/pkg/sync"
	"github.com/stretchr/testify/assert"
)

func TestWaitGroups(t *testing.T) {
	t.Parallel()

	wg := sync.WaitGroup{}
	value := false
	waitGroups := libSync.NewWaitGroups()

	waitGroups.Wait(1)
	wg.Add(1)
	go func() {
		waitGroups.Wait(1)
		value = true
		wg.Done()
		waitGroups.Done(1)
	}()

	waitGroups.Done(1)
	wg.Wait()
	waitGroups.Wait(1)

	assert.True(t, value)
}
