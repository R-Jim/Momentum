package link

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/template/event"
	"github.com/R-jim/Momentum/template/operator"
)

type Operator struct {
	linkAggregationSet operator.AggregationSet
}

func NewOperator(linkStore *event.Store) Operator {
	return Operator{
		linkAggregationSet: operator.NewAggregationSet(linkStore, NewAggregator()),
	}
}

func (o Operator) New(id, source, target uuid.UUID) error {
	return operator.NewEffectOperationSet(id, InitEffect, Link{
		source: source,
		target: target,
	}, o.linkAggregationSet).NewEventAndCommit()
}

func (o Operator) Destroy(id uuid.UUID) error {
	return operator.NewEffectOperationSet(id, DestroyEffect, nil, o.linkAggregationSet).NewEventAndCommit()
}

func (o Operator) Strengthen(id uuid.UUID) error {
	return operator.NewEffectOperationSet(id, StrengthenEffect, nil, o.linkAggregationSet).NewEventAndCommit()
}
