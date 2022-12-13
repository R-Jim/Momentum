package operator

import (
	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/animator"
)

type knightOperator struct {
	knightAggregator knight.Aggregator

	animator animator.Animator
}

func (k knightOperator) Init(knightID string, health knight.Health) error {
	knightInitEvent := knight.NewInitEvent(knightID, health)
	err := k.knightAggregator.Aggregate(knightInitEvent)
	if err != nil {
		return err
	}
	return nil
}

func (k knightOperator) Move(knightID string, toPosition knight.PositionState) error {
	knightMoveEvent := knight.NewMoveEvent(knightID, toPosition.X, toPosition.Y)
	err := k.knightAggregator.Aggregate(knightMoveEvent)
	if err != nil {
		return err
	}

	if k.animator != nil {
		k.animator.AppendEvent(knightMoveEvent)
	}
	return nil
}

func (k knightOperator) ChangeTarget(knightID string, target knight.Target) error {
	knightChangeTargetEvent := knight.NewChangeTargetEvent(knightID, target)
	err := k.knightAggregator.Aggregate(knightChangeTargetEvent)
	if err != nil {
		return err
	}

	if k.animator != nil {
		k.animator.AppendEvent(knightChangeTargetEvent)
	}
	return nil
}
