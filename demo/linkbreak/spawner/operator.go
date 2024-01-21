package spawner

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
)

type Operator struct {
	SpawnerStore *event.Store
}

func (s Operator) Init(id uuid.UUID, spawnType, faction, counter int, position math.Pos) error {
	initEvent := NewInitEvent(s.SpawnerStore, id, Spawner{
		spawnTypeCode: spawnType,
		faction:       faction,
		counter:       counter,
		position:      position,
	})

	if err := NewAggregator().Aggregate(s.SpawnerStore, initEvent); err != nil {
		return err
	}

	return s.SpawnerStore.AppendEvent(initEvent)
}

func (s Operator) CountDown(id uuid.UUID) error {
	countDownEvent := NewCountDownEvent(s.SpawnerStore, id)

	if err := NewAggregator().Aggregate(s.SpawnerStore, countDownEvent); err != nil {
		return err
	}

	return s.SpawnerStore.AppendEvent(countDownEvent)
}

func (s Operator) Spawn(id uuid.UUID) error {
	spawnEvent := NewSpawnEvent(s.SpawnerStore, id)

	if err := NewAggregator().Aggregate(s.SpawnerStore, spawnEvent); err != nil {
		return err
	}

	return s.SpawnerStore.AppendEvent(spawnEvent)
}
