package aggregate

import "errors"

var (
	ErrAggregateFail      error = errors.New("aggregate fail")
	ErrEffectNotSupported error = errors.New("effect not supported")
)
