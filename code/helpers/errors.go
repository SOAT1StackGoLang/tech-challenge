package helpers

import "errors"

var ErrUnauthorized = errors.New("user is not authorize to access this resource")
var ErrInvalidCurrencyFormat = errors.New("invalid currency format. Expected 'R$ X,XX'")
var ErrBadRequest = errors.New("bad request")
var ErrInvalidInput = errors.New("invalid input")
