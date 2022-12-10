package automaton

import (
	"math/rand"

	"github.com/R-jim/Momentum/aggregate/artifact"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/operator"
)

type ArtifactAutomaton interface {
	Auto(id string) error
}

type artifactImpl struct {
	artifactStore artifact.Store
	spikeStore    spike.Store

	operator operator.Operator
}

func NewArtifactAutomaton(artifactStore artifact.Store, operator operator.Operator) ArtifactAutomaton {
	return artifactImpl{
		artifactStore: artifactStore,

		operator: operator,
	}
}

func (i artifactImpl) Auto(id string) error {
	return i.autoSpawnSpike(id)
}

func (i artifactImpl) autoSpawnSpike(id string) error {
	state, err := artifact.GetState(i.artifactStore, id)
	if err != nil {
		return err
	}

	if state.Status != artifact.PlantedStatus || state.Energy < 10 {
		return nil
	}

	// TODO: check alive spike
	if len(state.SpikeIDs) > 5 {
		return nil
	}

	position, err := artifact.GetPositionState(i.artifactStore, id)
	if err != nil {
		return err
	}

	maxRadius := float64(50)
	minRadius := float64(40)

	radius := rand.Float64()*(maxRadius-minRadius) + minRadius

	positions := getPositions(position.X, position.Y, 0, 360, 5, radius)
	pos := positions[rand.Intn(len(positions))]

	i.operator.Artifact.SpawnSpike(id, pos.x, pos.y)
	return nil
}
