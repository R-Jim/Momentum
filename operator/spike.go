package operator

import (
	"github.com/R-jim/Momentum/aggregate/artifact"
	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/animator"
)

type spikeOperator struct {
	spikeAggregator    spike.Aggregator
	artifactAggregator artifact.Aggregator
	knightAggregator   knight.Aggregator

	animator animator.Animator
}

func (j spikeOperator) Init(spikeID string, artifactID string) error {
	spikeInitEvent := spike.NewInitEvent(spikeID, artifactID, 100)
	err := j.spikeAggregator.Aggregate(spikeInitEvent)
	if err != nil {
		return err
	}
	return nil
}

func (j spikeOperator) Move(spikeID string, toPosition spike.PositionState) error {
	spikeMoveEvent := spike.NewMoveEffect(spikeID, toPosition.X, toPosition.Y)
	err := j.spikeAggregator.Aggregate(spikeMoveEvent)
	if err != nil {
		return err
	}

	j.animator.AppendEvent(spikeMoveEvent)
	return nil
}

func (j spikeOperator) StrikeKnight(spikeID string, targetKnightID string) error {
	spikeStrikeEvent := spike.NewStrikeEvent(spikeID, targetKnightID)
	err := j.spikeAggregator.Aggregate(spikeStrikeEvent)
	if err != nil {
		return err
	}

	knightDamageEvent := knight.NewDamageEvent(targetKnightID, 5)
	err = j.knightAggregator.Aggregate(knightDamageEvent)
	if err != nil {
		return err
	}

	j.animator.AppendEvent(spikeStrikeEvent)
	j.animator.AppendEvent(knightDamageEvent)
	return nil
}
