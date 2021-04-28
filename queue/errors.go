package queue

import "errors"

/**
 * Error definitions
 */

var FullError = errors.New("FULL_ERROR: attempt to Put while Queue is Full")
var NilPointerError = errors.New("NilPointerError: attempt to store a nil into Queue")
var NoSuchElementError = errors.New("NoSuchElementError")
var IllegalStateError = errors.New("IllegalStateError, could cause by container full")
var IllegalArgumentError = errors.New("IllegalArgumentError ")