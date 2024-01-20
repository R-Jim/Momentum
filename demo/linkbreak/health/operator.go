package health

import (
	"github.com/R-jim/Momentum/template/event"
	"github.com/google/uuid"
)

type Operator struct {
	HealthStore *event.Store
}

func (o Operator) Modify(id uuid.UUID, value int) error {
	modifyEvent := NewModifyEvent(o.HealthStore, id, value)

	if err := NewAggregator().Aggregate(o.HealthStore, modifyEvent); err != nil {
		return err
	}

	o.HealthStore.AppendEvent(modifyEvent)

	return nil
}
