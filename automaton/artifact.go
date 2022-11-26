package automaton

import (
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
	return i.autoSpawnEnemy(id)
}

func (i artifactImpl) autoSpawnEnemy(id string) error {
	state, err := artifact.GetState(i.artifactStore, id)
	if err != nil {
		return err
	}

	if state.Status != artifact.PlantedStatus {
		return nil
	}

	if len(state.SpikeIDs) > 5 {
		return nil
	}

	_, err = artifact.GetPositionState(i.artifactStore, id)
	if err != nil {
		return err
	}

	// TODO: spawn spike
	return nil
}
