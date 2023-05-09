package errs

import "github.com/pkg/errors"

var (
	InternalServer = errors.New("internal server error")
)
