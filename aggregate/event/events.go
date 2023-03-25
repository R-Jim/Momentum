package event

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
)

type data interface{}

type Effect string

type Event struct {
	ID       uuid.UUID
	EntityID uuid.UUID
	Version  int
	Effect   Effect
	Data     data
}

func ParseData[T data](e Event) (T, error) {
	data, ok := e.Data.(T)
	if !ok {
		return data, pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect: %s", reflect.TypeOf(data)))
	}
	return data, nil
}
