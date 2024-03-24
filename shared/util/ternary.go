package util

func Ternary[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}

func TernaryF[T any](cond bool, trueFunc func() T, falseFunc func() T) T {
	if cond {
		return trueFunc()
	}
	return falseFunc()
}
