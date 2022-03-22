package errors

import "errors"

var (
	ErrDateBusy    = errors.New("event for this date already created")
	ErrNoSuchEvent = errors.New("no event for this date")
)
