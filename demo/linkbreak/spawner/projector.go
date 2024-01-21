package spawner

import (
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
)

type SpawnerProjection struct {
	Type      int
	Faction   int
	Counter   int
	Position  math.Pos
	IsSpawned bool
}

type Projector struct {
	SpawnerStore *event.Store
}

func (p Projector) SpawnerProjections() ([]SpawnerProjection, error) {
	projections := []SpawnerProjection{}

	for _, spawnerEvents := range p.SpawnerStore.GetEvents() {
		var projection SpawnerProjection

		for _, e := range spawnerEvents {
			switch e.Effect {
			case InitEffect:
				spawner, err := event.ParseData[Spawner](e)
				if err != nil {
					return nil, err
				}
				projection.Faction = spawner.faction
				projection.Type = spawner.spawnTypeCode
				projection.Counter = spawner.counter
				projection.Position = spawner.position
			case CountDownEffect:
				projection.Counter--
			case SpawnEffect:
				projection.IsSpawned = true
			}
		}

		projections = append(projections, projection)
	}

	return projections, nil
}
