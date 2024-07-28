package util

import (
  "errors"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/sony/gobreaker"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "slices"
)

// CastBreakerError will cast error due to breaker error into internal service error. otherwise it
// will return the error unchanged
func CastBreakerError(err error) error {
  if errors.Is(err, gobreaker.ErrOpenState) {
    return sharedErr.ErrServiceUnavailable
  }
  if errors.Is(err, gobreaker.ErrTooManyRequests) {
    return sharedErr.ErrServiceRecovering
  }
  return err
}

// IsGrpcConnectivityError will return true if the error is not related to grpc or the connectivity issue,
// otherwise it will return false
func IsGrpcConnectivityError(err error) bool {
  s, ok := status.FromError(err)
  if !ok {
    return true
  }
  return slices.Contains([]codes.Code{
    codes.Unavailable,
    codes.ResourceExhausted,
    codes.Internal,
    codes.DeadlineExceeded,
    codes.Unknown,
  }, s.Code())
}
