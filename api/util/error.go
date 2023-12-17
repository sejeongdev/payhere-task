package util

import (
	"context"
	"net/http"

	"github.com/juju/errors"
)

// CauseError ...
func CauseError(err error) (int, string) {
	switch {
	case errors.Is(err, context.Canceled):
		return 444, errors.Cause(err).Error() // Nginx Non-Standard error for `No Response`
	case errors.Is(err, errors.NotFound), errors.Is(err, errors.UserNotFound):
		return http.StatusNotFound, errors.Cause(err).Error()
	case errors.Is(err, errors.AlreadyExists):
		return http.StatusConflict, errors.Cause(err).Error()
	case errors.Is(err, errors.Unauthorized):
		return http.StatusUnauthorized, errors.Cause(err).Error()
	case errors.Is(err, errors.BadRequest):
		return http.StatusBadRequest, errors.Cause(err).Error()
	case errors.Is(err, errors.Forbidden):
		return http.StatusForbidden, errors.Cause(err).Error()
	case errors.Is(err, errors.NotSupported):
		return http.StatusGone, errors.Cause(err).Error()
	case errors.Is(err, errors.NotValid):
		return http.StatusPreconditionFailed, errors.Cause(err).Error()
	case errors.Is(err, errors.QuotaLimitExceeded):
		return http.StatusLocked, errors.Cause(err).Error()
	}
	return http.StatusInternalServerError, err.Error()
}
