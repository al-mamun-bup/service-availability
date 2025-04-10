package tests

import (
	"testing"
	"time"

	"service-availability/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestIsWithinOpenHours(t *testing.T) {
	type testCase struct {
		name      string
		openHours []string
		current   string // in "1504" format
		expected  bool
	}

	cases := []testCase{
		{
			name:      "Valid time within range",
			openHours: []string{"0900,1800"},
			current:   "1200",
			expected:  true,
		},
		{
			name:      "Time before range",
			openHours: []string{"0900,1800"},
			current:   "0800",
			expected:  false,
		},
		{
			name:      "Time after range",
			openHours: []string{"0900,1800"},
			current:   "1900",
			expected:  false,
		},
		{
			name:      "Time exactly at opening time",
			openHours: []string{"0900,1800"},
			current:   "0900",
			expected:  true,
		},
		{
			name:      "Time exactly at closing time",
			openHours: []string{"0900,1800"},
			current:   "1800",
			expected:  true,
		},
		{
			name:      "Invalid entry format (no comma)",
			openHours: []string{"09001800"},
			current:   "1200",
			expected:  false,
		},
		{
			name:      "Invalid entry format (non-numeric)",
			openHours: []string{"abcd,efgh"},
			current:   "1200",
			expected:  false,
		},
		{
			name:      "Multiple slots, one valid",
			openHours: []string{"0000,0100", "1100,1300", "2200,2300"},
			current:   "1200",
			expected:  true,
		},
		{
			name:      "Multiple slots, all invalid",
			openHours: []string{"0000-0100", "hello,world"},
			current:   "1200",
			expected:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			now, err := time.Parse("1504", tc.current)
			assert.NoError(t, err)

			result := utils.IsWithinOpenHours(tc.openHours, now)
			assert.Equal(t, tc.expected, result)
		})
	}
}
