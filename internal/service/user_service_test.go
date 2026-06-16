package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		dob  time.Time
		want int
	}{
		{
			name: "birthday already passed this year",
			dob:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: now.Year() - 1990,
		},
		{
			name: "birthday not yet this year",
			dob:  time.Date(1990, time.December, 31, 0, 0, 0, 0, time.UTC),
			want: now.Year() - 1990 - 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := CalculateAge(tc.dob); got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}
