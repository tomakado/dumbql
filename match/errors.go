package match

import "errors"

var (
	errFieldNotFound = errors.New("field not found")
	errNotAStruct    = errors.New("not a struct")
)
