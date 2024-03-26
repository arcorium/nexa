package util

import "nexa/shared/errors"

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

func CastSlice2[From, To any](slice []From, f func(*From) (To, error)) ([]To, error) {
	if slice == nil || len(slice) == 0 {
		return nil, errors.ErrEmptySlice
	}

	result := make([]To, 0, len(slice))
	for _, val := range slice {
		res, err := f(&val)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}
