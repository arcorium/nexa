package errors

type IndexedError struct {
  Index int
  Err   error
}
