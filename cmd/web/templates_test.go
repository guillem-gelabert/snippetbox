package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	testCases := []struct {
		desc     string
		target   time.Time
		expected string
	}{
		{
			desc:     "UTC",
			target:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			expected: "17 Dec 2020 at 10:00",
		},
		{
			desc:     "Empty",
			target:   time.Time{},
			expected: "",
		},
		{
			desc:     "CET",
			target:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "17 Dec 2020 at 09:00",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := humanDate(tC.target)
			if actual != tC.expected {
				t.Errorf("expected %q; got %q", tC.expected, actual)
			}
		})
	}
}
