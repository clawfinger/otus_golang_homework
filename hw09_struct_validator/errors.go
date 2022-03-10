package hw09structvalidator

import "errors"

var (
	errNotAStruct     = errors.New("input is not a struct")
	errBadValidator   = errors.New("validation sequence is unusable")
	errIntNotInSet    = errors.New("int field value is not in a range")
	errIntLess        = errors.New("int field value is less")
	errIntGreater     = errors.New("int field value is greater")
	errStringLen      = errors.New("string field len is wrong")
	errStringNotInSet = errors.New("string field is not in set")
	errStringRegexp   = errors.New("string field is not satisfy the regexp")
)
