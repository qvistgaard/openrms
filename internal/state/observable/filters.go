package observable

// Filter represents a function type used for filtering changes in observable values.
//
// This function type is intended for use in determining whether changes to an observable value are significant
// and should be propagated to observers. It takes two parameters: the original value and the new value of type T.
// The function should return a boolean value, which indicates whether the change is substantial enough to be accepted and published.
//
// Parameters:
//   - original: The original value of the observable, of type T.
//   - new: The new value proposed for the observable, of type T.
//
// Returns:
//   - bool: A boolean value indicating whether the change between the original and new values should be accepted.
//     Returns true if the change is significant and should be propagated; false if the change should be ignored.
//
// Usage:
// This type is commonly used in the context of an Observable interface, where it can be employed to filter changes
// before they are communicated to observers. It enables the implementation of logic to ignore minor, insignificant,
// or invalid changes, ensuring that only meaningful updates are broadcast.
//
// Example:
// A Filter[uint8] might be used to determine if a numeric value has changed by a certain threshold before notifying observers.
//
// See also DistinctPercentageChange and DistinctComparableChange for examples
type Filter[T any] func(original T, new T) bool

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
	return DistinctComparableChange[bool]()
}

// DistinctComparableChange returns a filter function for comparable types.
// It checks if there's a distinct change between the current and new values.
//
// Parameters:
//   - current: The current value of type T.
//   - new: The new value of type T.
//
// Returns:
//   - bool: true if the change is distinct, false otherwise.
func DistinctComparableChange[T comparable]() Filter[T] {
	return func(current T, new T) bool {
		return current != new
	}
}
