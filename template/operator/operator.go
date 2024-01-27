package operator

import (
	"github.com/R-jim/Momentum/template/aggregate"
	"github.com/R-jim/Momentum/template/event"
	"github.com/google/uuid"
)

type AggregationSet struct {
	store      *event.Store
	aggregator aggregate.Aggregator
}

func NewAggregationSet(store *event.Store, aggregator aggregate.Aggregator) AggregationSet {
	return AggregationSet{
		store, aggregator,
	}
}

func (s AggregationSet) Store() event.Store {
	return *s.store
}

type EffectOperationSet struct {
	entityID       uuid.UUID
	effect         event.Effect
	data           event.Data
	aggregationSet AggregationSet
}

func NewEffectOperationSet(entityID uuid.UUID, effect event.Effect, data event.Data, aggregationSet AggregationSet) EffectOperationSet {
	return EffectOperationSet{
		entityID,
		effect,
		data,
		aggregationSet,
	}
}

func (s EffectOperationSet) NewEvent() event.Event {
	return s.aggregationSet.store.NewEvent(s.entityID, s.effect, s.data)
}

func (s EffectOperationSet) NewEventAndCommit() error {
	e := s.aggregationSet.store.NewEvent(s.entityID, s.effect, s.data)
	if err := s.Aggregate(e); err != nil {
		return err
	}
	return s.Commit(e)
}

func (s EffectOperationSet) Aggregate(e event.Event) error {
	return s.aggregationSet.aggregator.Aggregate(s.aggregationSet.store, e)
}

func (s EffectOperationSet) Commit(e event.Event) error {
	return s.aggregationSet.store.AppendEvent(e)
}

type EffectOperationSets []EffectOperationSet

func (s EffectOperationSets) PerformTxn() error {
	type eventAndStore struct {
		event event.Event
		store *event.Store
	}

	commitList := []eventAndStore{}
	for _, operationSet := range s {
		e := operationSet.NewEvent()
		if err := operationSet.Aggregate(e); err != nil {
			return err
		}
		commitList = append(commitList, eventAndStore{
			event: e,
			store: operationSet.aggregationSet.store,
		})
	}

	for _, set := range commitList {
		set.store.AppendEvent(set.event)
	}

	return nil
}
