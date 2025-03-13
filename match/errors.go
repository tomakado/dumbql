package match

import "errors"

var (
	ErrFieldNotFound = errors.New("field not found")
	ErrNotAStruct    = errors.New("not a struct")
)
