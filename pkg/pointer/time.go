package pointer

import "time"

func Time(t time.Time) *time.Time {
	return &t
}
