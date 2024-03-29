package event

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
)

type Data interface{}

type Effect string

type Event struct {
	ID       uuid.UUID
	EntityID uuid.UUID
	Version  int
	Effect   Effect
	Data     Data

	CreatedAt time.Time
}

func (s *Store) NewEvent(entityID uuid.UUID, effect Effect, data Data) Event {
	return Event{
		ID:       uuid.New(),
		EntityID: entityID,
		Version:  len(s.GetEvents()[entityID]) + 1,
		Effect:   effect,
		Data:     data,

		CreatedAt: time.Now(),
	}
}

func ParseData[T Data](e Event) (T, error) {
	data, ok := e.Data.(T)
	if !ok {
		return data, pkgerrors.WithStack(fmt.Errorf("failed to parse data for effect: %s", reflect.TypeOf(data)))
	}
	return data, nil
}
