package status

import "google.golang.org/grpc/status"

func Error(code Code, err error) Object {
	return New(code, err)
}

func ErrorC(code Code) Object {
	return New(code, nil)
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
	return s.Error != nil || s.Codes > DELETED
}

func (s *Object) ToGRPCError() error {
	return status.Error(MapGRPCCode(s.Codes), s.Error.Error())
}
