package operator

import (
	"testing"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Worker_Init(t *testing.T) {
	workerID := uuid.New()
	store := store.NewStore()

	workerOperator := WorkerOperator{
		WorkerAggregator: aggregator.NewWorkerAggregator(&store),
	}

	err := workerOperator.Init(workerID)
	require.NoError(t, err)

	events, err := store.GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err := aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, workerState.ID)
}

func Test_Worker_AssignBuilding(t *testing.T) {
	workerStore := store.NewStore()
	buildingStore := store.NewStore()

	workerOperator := WorkerOperator{
		WorkerAggregator:   aggregator.NewWorkerAggregator(&workerStore),
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	BuildingOperator := BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	workerID := uuid.New()
	buildingID := uuid.New()

	require.NoError(t, workerOperator.Init(workerID))
	require.NoError(t, BuildingOperator.Init(buildingID, math.NewPos(2, 2)))

	require.NoError(t, workerOperator.AssignBuilding(workerID, buildingID))

	events, err := workerStore.GetEventsByEntityID(workerID)
	require.NoError(t, err)
	workerState, err := aggregator.GetWorkerState(events)
	require.NoError(t, err)
	require.Equal(t, buildingID, workerState.BuildingID)

	events, err = buildingStore.GetEventsByEntityID(buildingID)
	require.NoError(t, err)
	buildingState, err := aggregator.GetBuildingState(events)
	require.NoError(t, err)
	require.True(t, buildingState.WorkerMap[workerID])

	require.Error(t, workerOperator.AssignBuilding(workerID, buildingID))
}

func Test_Worker_UnAssignBuilding(t *testing.T) {
	workerStore := store.NewStore()
	buildingStore := store.NewStore()

	workerOperator := WorkerOperator{
		WorkerAggregator:   aggregator.NewWorkerAggregator(&workerStore),
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	BuildingOperator := BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	workerID := uuid.New()
	buildingID := uuid.New()

	require.NoError(t, workerOperator.Init(workerID))
	require.NoError(t, BuildingOperator.Init(buildingID, math.NewPos(2, 2)))

	require.NoError(t, workerOperator.AssignBuilding(workerID, buildingID))

	events, err := workerStore.GetEventsByEntityID(workerID)
	require.NoError(t, err)
	workerState, err := aggregator.GetWorkerState(events)
	require.NoError(t, err)
	require.Equal(t, buildingID, workerState.BuildingID)

	events, err = buildingStore.GetEventsByEntityID(buildingID)
	require.NoError(t, err)
	buildingState, err := aggregator.GetBuildingState(events)
	require.NoError(t, err)
	require.True(t, buildingState.WorkerMap[workerID])

	require.NoError(t, workerOperator.UnassignBuilding(workerID, buildingID))

	events, err = workerStore.GetEventsByEntityID(workerID)
	require.NoError(t, err)
	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)
	require.Equal(t, uuid.Nil, workerState.BuildingID)

	events, err = buildingStore.GetEventsByEntityID(buildingID)
	require.NoError(t, err)
	buildingState, err = aggregator.GetBuildingState(events)
	require.NoError(t, err)
	require.False(t, buildingState.WorkerMap[workerID])

	require.Error(t, workerOperator.UnassignBuilding(workerID, buildingID))
}
