package model

import (
	"time"
)

func isBeforeOrEqual[T1 DateOnly | DateTime, T2 DateOnly | DateTime](a T1, b T2) bool {
	return time.Time(a).Before(time.Time(b)) || time.Time(a).Equal(time.Time(b))
}
