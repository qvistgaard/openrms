package race

import "github.com/qvistgaard/openrms/internal/state/observable"

func filterRaceStatusChange() observable.Filter[Status] {
	return func(oldStatus Status, newStatus Status) bool {
		return oldStatus != newStatus
	}
}

func filterTotalLapsCountChange() observable.Filter[uint16] {
	return func(oldLapCount uint16, newLapCount uint16) bool {
		return newLapCount > oldLapCount
	}
}
