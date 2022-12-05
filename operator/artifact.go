package operator

import (
	"github.com/R-jim/Momentum/aggregate/artifact"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/animator"
	"github.com/google/uuid"
)

type artifactOperator struct {
	spikeAggregator    spike.Aggregator
	artifactAggregator artifact.Aggregator

	animator animator.Animator
}

func (j artifactOperator) Init(artifactID string) error {
	artifactInitEvent := artifact.NewInitEvent(artifactID)
	err := j.artifactAggregator.Aggregate(artifactInitEvent)
	if err != nil {
		return err
	}
	return nil
}

// TODO: might need to check if position is already occupied
func (j artifactOperator) SpawnSpike(artifactID string, x, y float64) error {
	spikeID := uuid.New().String()

	artifactSpawnSpikeEvent := artifact.NewSpawnSpikeEffect(artifactID, spikeID, x, y)
	err := j.artifactAggregator.Aggregate(artifactSpawnSpikeEvent)
	if err != nil {
		return err
	}

	spikeInitEvent := spike.NewInitEvent(spikeID, artifactID, 100)
	err = j.spikeAggregator.Aggregate(spikeInitEvent)
	if err != nil {
		return err
	}

	spikeMoveEvent := spike.NewMoveEffect(spikeID, x, y)
	err = j.spikeAggregator.Aggregate(spikeMoveEvent)
	if err != nil {
		return err
	}

	j.animator.AppendEvent(artifactSpawnSpikeEvent)
	return nil
}

func (j artifactOperator) Move(artifactID string, toPosition artifact.PositionState) error {
	artifactMoveEvent := artifact.NewMoveEffect(artifactID, toPosition.X, toPosition.Y)
	err := j.artifactAggregator.Aggregate(artifactMoveEvent)
	if err != nil {
		return err
	}

	j.animator.AppendEvent(artifactMoveEvent)
	return nil
}
