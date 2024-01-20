package runner

import "errors"

var (
	ErrPositionIDRequired error = errors.New("position ID required")
	ErrPositionIDNotFound error = errors.New("position ID not found")
)
