package sync

import (
	"github.com/phuslu/log"
	"sync"
)

// WaitGroups represents a mutex map of a ID with an associated sync.WaitGroup slice.
type WaitGroups struct {
	wGs map[int64][]*sync.WaitGroup
	mu  sync.Mutex
}

// NewWaitGroups return a new waitGroups pointer.
func NewWaitGroups() *WaitGroups {
	return &WaitGroups{wGs: make(map[int64][]*sync.WaitGroup)}
}

// Wait waits for the last wait group in the wait group slice with the given ID.
func (wG *WaitGroups) Wait(iD int64) {
	wG.mu.Lock()

	wGs, ok := wG.wGs[iD]
	if ok && len(wGs) > 0 {
		log.Info().Msgf("[Wait] waiting for wg at position %d in queue %d", len(wGs), iD)

		defer log.Info().Msgf("[Wait] done waiting at position %d in queue %d", len(wGs), iD)
		defer wGs[len(wGs)-1].Wait()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wGs = append(wGs, &wg)
	wG.wGs[iD] = wGs
	wG.mu.Unlock()
}

// Done sets the first wait group in the wait group slice with the given ID to done and removes it from the slice
func (wG *WaitGroups) Done(iD int64) {
	wG.mu.Lock()
	defer wG.mu.Unlock()

	wGs, ok := wG.wGs[iD]
	if !ok {
		log.Info().Msg("[Done] no wait group found")
	}

	log.Info().Msgf("[Done] done in queue %d", iD)
	wGs[0].Done()
	wGs = wGs[1:]
	wG.wGs[iD] = wGs
}
