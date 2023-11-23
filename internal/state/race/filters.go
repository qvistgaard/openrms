package race

import "github.com/qvistgaard/openrms/internal/state/observable"

func filterRaceStatusChange() observable.Filter[RaceStatus] {
	return func(oldStatus RaceStatus, newStatus RaceStatus) bool {
		return oldStatus != newStatus
	}
}

func filterTotalLapsCountChange() observable.Filter[uint16] {
	return func(oldLapCount uint16, newLapCount uint16) bool {
		return newLapCount > oldLapCount
	}
}
