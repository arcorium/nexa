package util

func CastSlice[From, To any](slice []From, f func(*From) To) []To {
	if slice == nil || len(slice) == 0 {
		return nil
	}

	result := make([]To, 0, len(slice))
	for _, val := range slice {
		result = append(result, f(&val))
	}
	return result
}
