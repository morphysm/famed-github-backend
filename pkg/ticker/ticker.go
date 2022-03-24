package ticker

import (
	"time"
)

func NewTicker(interval time.Duration, function func()) {
	ticker := time.NewTicker(interval)
	go func() {
		for ; true; <-ticker.C {
			go function()
		}
	}()
}
