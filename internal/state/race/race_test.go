package race

import (
	"testing"
	"time"
)

func TestCalculateRaceDuration(t *testing.T) {
	// Create a reference time
	startTime := time.Date(2023, time.November, 18, 12, 0, 0, 0, time.UTC)

	// Test case 1: Race is ongoing (positive duration)
	callTime1 := startTime.Add(10 * time.Second)
	duration1 := calculateRaceDuration(time.Second*0, startTime, callTime1)
	expected1 := 10 * time.Second
	if duration1 != expected1 {
		t.Errorf("Test case 1: Expected race duration to be %v, but got %v", expected1, duration1)
	}

}
