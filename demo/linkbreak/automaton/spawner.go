package automaton

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"

	"github.com/R-jim/Momentum/demo/linkbreak/runner"
	"github.com/R-jim/Momentum/demo/linkbreak/spawner"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
)

type SpawnerAutomaton struct {
	spawnerStore *event.Store

	spawnerOperator spawner.Operator
	runnerOperator  runner.OperatorV2
}

func NewSpawnerAutomaton(spawnerStore, runnerStore, positionStore, healthStore *event.Store) SpawnerAutomaton {
	return SpawnerAutomaton{
		spawnerStore: spawnerStore,

		spawnerOperator: spawner.Operator{
			SpawnerStore: spawnerStore,
		},
		runnerOperator: runner.NewOperatorV2(runnerStore, healthStore, positionStore),
	}
}

func (s SpawnerAutomaton) NewSpawner(minPosition, maxPosition math.Pos) error {
	x := rand.Float64() * (maxPosition.X - minPosition.X)
	y := rand.Float64() * (maxPosition.Y - minPosition.Y)

	return s.spawnerOperator.Init(uuid.New(), 1, 2, 3, math.NewPos(x, y))
}

func (s SpawnerAutomaton) SpawnOrCountDown() error {
	for _, spawnerEvents := range s.spawnerStore.GetEvents() {
		if err := spawnOrCountDown(spawnerEvents, s.spawnerOperator, s.runnerOperator); err != nil {
			return err
		}
	}

	return nil
}

func spawnOrCountDown(events []event.Event, spawnerOperator spawner.Operator, runnerOperator runner.OperatorV2) error {
	var counter int
	var spawnerID uuid.UUID
	var spawnerData spawner.Spawner

	for _, e := range events {
		switch e.Effect {
		case spawner.InitEffect:
			var err error
			spawnerData, err = event.ParseData[spawner.Spawner](e)
			if err != nil {
				return err
			}
			counter = spawnerData.Counter()
			spawnerID = e.EntityID
		case spawner.CountDownEffect:
			counter--
		case spawner.SpawnEffect:
			return nil
		}
	}

	if counter > 0 {
		return spawnerOperator.CountDown(spawnerID)
	} else {
		return spawnRunner(spawnerOperator, runnerOperator, spawnerID, spawnerData)
	}
}

func spawnRunner(spawnerOperator spawner.Operator, runnerOperator runner.OperatorV2, spawnerID uuid.UUID, spawnerData spawner.Spawner) error {
	data, ok := spawnerData.SpawnTypeData().(spawner.SpawnRunnerData)
	if !ok {
		return errors.New("failed to parse spawn runner data")
	}

	if err := runnerOperator.NewRunner(uuid.New(), data.HealthValue, spawnerData.Faction(), spawnerData.Position()); err != nil {
		return err
	}

	return spawnerOperator.Spawn(spawnerID)
}
