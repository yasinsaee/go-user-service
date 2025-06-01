package util

import (
	"fmt"
	"strconv"
	"strings"
)

func TimeToMinutes(t string) (int, error) {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	// If the time is 00:00, return 1440 minutes (end of day)
	if hours == 0 && minutes == 0 {
		return 1440, nil
	}

	return hours*60 + minutes, nil
}
