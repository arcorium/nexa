package status

import (
	"google.golang.org/grpc/codes"
)

func MapGRPCCode(code Code) codes.Code {
	switch code {
	case INTERNAL_SERVER_ERROR, EXTERNAL_SERVICE_ERROR, REPOSITORY_ERROR:
		return codes.Internal
	case SUCCESS, CREATED, UPDATED, DELETED:
		return codes.OK
	case OBJECT_NOT_FOUND:
		return codes.NotFound
	case BAD_REQUEST_ERROR:
		return codes.InvalidArgument
	case FIELD_VALIDATION_ERROR:
		return codes.FailedPrecondition
	case NOT_AUTHENTICATED_ERROR:
		return codes.Unauthenticated
	case NOT_AUTHORIZED_ERROR:
		return codes.PermissionDenied
	case SERVICE_UNAVAILABLE_ERROR:
		return codes.Unavailable
	}
	return codes.Unknown
}
