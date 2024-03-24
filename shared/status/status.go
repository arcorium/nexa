package status

func Success() Object {
	return New(SUCCESS, nil)
}

func Created() Object {
	return New(CREATED, nil)
}

func Updated() Object {
	return New(UPDATED, nil)
}

func Deleted() Object {
	return New(DELETED, nil)
}

func Error(code Code, err error) Object {
	return New(code, err)
}

func Internal(err error) Object {
	return New(INTERNAL_SERVER_ERROR, err)
}

func New(code Code, err error) Object {
	return Object{
		Codes: code,
		Error: err,
	}
}

type Object struct {
	Codes Code
	Error error
}

func (s *Object) IsError() bool {
	return s.Error != nil
}
