package djson

import "errors"

var invalidPathError = errors.New("invalid path")
var unavailableError = errors.New("path func unavailable")
var failedToSortError = errors.New("failedToSortError")
