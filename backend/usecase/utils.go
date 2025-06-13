package usecase

import (
	"backend/domain"
	"time"
)

func isBeforeOrEqual[T1 domain.DateOnly | domain.DateTime, T2 domain.DateOnly | domain.DateTime](a T1, b T2) bool {
	return time.Time(a).Before(time.Time(b)) || time.Time(a).Equal(time.Time(b))
}
