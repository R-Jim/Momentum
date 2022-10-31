package animator

import "errors"

var (
	ErrEffectNotSupported           error = errors.New("effect not supported")
	ErrStateNotSupported            error = errors.New("state not supported")
	ErrCanNotAppendUnsupportedEvent error = errors.New("can not append unsupported event")
	ErrNoPendingEvents              error = errors.New("no pending event")
)
