package pointers

import "time"

func Time(t time.Time) *time.Time {
	return &t
}
