package utils

import (
	"fmt"
	"strings"
	"time"
)

// IsWithinOpenHours checks if the current time is within any of the open-close time ranges
func IsWithinOpenHours(openHours []string, now time.Time) bool {
	for _, entry := range openHours {
		// Example entry: "0000,2359,00:00AM-23:59PM"
		parts := strings.Split(entry, ",")
		if len(parts) < 2 {
			continue
		}

		openStr := parts[0]
		closeStr := parts[1]

		openTime, err1 := time.Parse("1504", openStr)
		closeTime, err2 := time.Parse("1504", closeStr)

		if err1 != nil || err2 != nil {
			fmt.Println("Error parsing time:", err1, err2)
			continue
		}

		// Get the time in HHMM format from current time
		nowFormatted := now.Format("1504")
		currentTime, err3 := time.Parse("1504", nowFormatted)
		if err3 != nil {
			fmt.Println("Error parsing current time:", err3)
			continue
		}

		// Check if currentTime is within the open-close range
		if currentTime.After(openTime) && currentTime.Before(closeTime) ||
			currentTime.Equal(openTime) || currentTime.Equal(closeTime) {
			return true
		}
	}
	return false
}
