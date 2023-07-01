package helpers

import "errors"

var ErrUnauthorized = errors.New("user is not authorize to access this resource")
