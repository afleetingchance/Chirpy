package auth

import "errors"

var ErrMissingHeader = errors.New("need Authorization header")
var ErrInvalidHeader = errors.New("invalid Authorization header")
