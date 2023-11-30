package observable

type Filter[T any] func(T, T) bool

// DistinctPercentageChange returns a filter function that checks if a change in a uint8 percentage value is valid.
// It returns true if the new percentage value is within the valid range of 0 to 100 (inclusive) and if the current and new values are different.
// If the new value is less than 0 or greater than 100, it returns false to indicate an invalid change.
//
// Parameters:
//   - current: The current uint8 percentage value.
//   - new: The new uint8 percentage value.
//
// Returns:
//   - bool: true if the change is valid, false otherwise.
func DistinctPercentageChange() Filter[uint8] {
	return func(current uint8, new uint8) bool {
		if new < 0 {
			return false
		}
		if new > 100 {
			return false
		}
		return current != new
	}
}

// DistinctBooleanChange returns a filter function that checks if a change in a boolean value is distinct.
// It returns true if the current and new boolean values are different, indicating a distinct change.
// If the current and new values are the same, it returns false.
//
// Parameters:
//   - current: The current boolean value.
//   - new: The new boolean value.
//
// Returns:
//   - bool: true if the change is distinct, false if the values are the same.
func DistinctBooleanChange() Filter[bool] {
	return DistictComparableChange[bool]()
}

func DistictComparableChange[T comparable]() Filter[T] {
	return func(current T, new T) bool {
		return current != new
	}
}
