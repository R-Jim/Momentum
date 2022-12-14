package operator

import (
	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/animator"
)

type knightOperator struct {
	knightAggregator knight.Aggregator
	spikeAggregator  spike.Aggregator

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

func (k knightOperator) StrikeSpike(spikeID string, targetSpikeID string) error {
	knightStrikeEvent := knight.NewStrikeEvent(spikeID, targetSpikeID)
	err := k.knightAggregator.Aggregate(knightStrikeEvent)
	if err != nil {
		return err
	}
	spikeDamageEvent := spike.NewDamageEvent(targetSpikeID, 20)
	err = k.spikeAggregator.Aggregate(spikeDamageEvent)
	if err != nil {
		return err
	}

	k.animator.AppendEvent(knightStrikeEvent)
	k.animator.AppendEvent(spikeDamageEvent)
	return nil
}
