package error

type Status struct {
	Code uint
	err  error
}

func (s Status) Error() string {
	return s.err.Error()
}
