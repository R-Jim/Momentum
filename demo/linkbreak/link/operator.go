package link

import (
	"github.com/R-jim/Momentum/template/event"
	"github.com/google/uuid"
)

type Operator struct {
	LinkStore *event.Store
}

func (o Operator) New(id, source, target uuid.UUID) error {
	initEvent := NewInitEvent(o.LinkStore, id, Link{
		source: source,
		target: target,
	})

	if err := NewAggregator().Aggregate(o.LinkStore, initEvent); err != nil {
		return err
	}

	o.LinkStore.AppendEvent(initEvent)

	return nil
}

func (o Operator) Destroy(id uuid.UUID) error {
	destroyEvent := NewDestroyEvent(o.LinkStore, id)

	if err := NewAggregator().Aggregate(o.LinkStore, destroyEvent); err != nil {
		return err
	}

	o.LinkStore.AppendEvent(destroyEvent)

	return nil
}

func (o Operator) Strengthen(id uuid.UUID) error {
	strengthenEvent := NewStrengthenEvent(o.LinkStore, id)

	if err := NewAggregator().Aggregate(o.LinkStore, strengthenEvent); err != nil {
		return err
	}

	o.LinkStore.AppendEvent(strengthenEvent)

	return nil
}
