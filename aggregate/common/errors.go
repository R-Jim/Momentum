package common

import "errors"

var (
	ErrEventNotValid      error = errors.New("event not valid")
	ErrAggregateFail      error = errors.New("aggregate fail")
	ErrEffectNotSupported error = errors.New("effect not supported")
	ErrEntityNotFound     error = errors.New("entity not found")
)
