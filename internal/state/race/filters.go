package race

import "github.com/qvistgaard/openrms/internal/state/observable"

func filterRaceStatusChange() observable.Filter[Status] {
	return func(oldStatus Status, newStatus Status) bool {
		return oldStatus != newStatus
	}
}

func filterTotalLapsCountChange() observable.Filter[uint32] {
	return func(oldLapCount uint32, newLapCount uint32) bool {
		return newLapCount > oldLapCount
	}
}
