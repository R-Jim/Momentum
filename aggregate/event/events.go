package event

import (
	"fmt"
	"reflect"
	"time"

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

	CreatedAt time.Time
}

func newEvent(entityID uuid.UUID, version int, effect Effect, data data) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  version,
		Effect:   effect,
		Data:     data,

		CreatedAt: time.Now(),
	}
}

func ParseData[T data](e Event) (T, error) {
	data, ok := e.Data.(T)
	if !ok {
		return data, pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect: %s", reflect.TypeOf(data)))
	}
	return data, nil
}
