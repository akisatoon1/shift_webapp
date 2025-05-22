package model

import (
	"testing"
)

func TestIsBeforeOrEqual(t *testing.T) {
	tests := []struct {
		a    DateOnly
		b    DateOnly
		want bool
	}{
		{mustNewDateOnly("2024-06-01"), mustNewDateOnly("2024-06-01"), true},
		{mustNewDateOnly("2024-06-01"), mustNewDateOnly("2024-06-02"), true},
		{mustNewDateOnly("2024-06-01"), mustNewDateOnly("2024-05-31"), false},
	}
	for _, test := range tests {
		got := isBeforeOrEqual(test.a, test.b)
		if got != test.want {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
}
