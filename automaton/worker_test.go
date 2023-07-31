package automaton

import (
	"testing"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Worker_PathFindingUpdate_simple(t *testing.T) {
	posA := math.NewPos(0, 0)
	posB := math.NewPos(5, 0)
	posC := math.NewPos(8, 0)
	buildingPos := math.NewPos(8, 2)

	workerPos := math.NewPos(2, 0)

	mapPaths := []math.Path{
		{Start: posA, End: posB, Cost: 10}, // street1
		{Start: posB, End: posC},           // street2
		{Start: posC, End: buildingPos},    // building street
	}

	workerID := uuid.New()
	street1ID := uuid.New()
	street2ID := uuid.New()
	buildingStreetID := uuid.New()

	workerStore := event.NewWorkerStore()
	streetStore := event.NewStreetStore()
	buildingStore := event.NewBuildingStore()

	workerOperator := operator.NewWorker(&workerStore, &buildingStore)
	BuildingOperator := operator.NewBuilding(&buildingStore, nil)
	streetOperator := operator.NewStreet(&streetStore)

	w := WorkerAutomaton{
		EntityID: workerID,
		MapPaths: mapPaths,

		WorkerStore:   &workerStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		WorkerOperator: workerOperator,
		StreetOperator: streetOperator,
	}

	err := workerOperator.Init(workerID, workerPos)
	require.NoError(t, err)

	events, err := event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	mioState, err := aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, mioState.ID)

	err = streetOperator.Init(street1ID, posA, posB)
	require.NoError(t, err)
	err = streetOperator.Init(street2ID, posB, posC)
	require.NoError(t, err)
	err = streetOperator.Init(buildingStreetID, posC, buildingPos)
	require.NoError(t, err)

	drinkStoreID := uuid.New()
	require.NoError(t, BuildingOperator.Init(drinkStoreID, event.BuildingTypeDrinkStore, buildingPos))

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	mioState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.NoError(t, workerOperator.AssignBuilding(workerID, drinkStoreID))

	w.PathFindingUpdate()

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	mioState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.Equal(t, []math.Pos{posB, posC, buildingPos}, mioState.PlannedPoses)
}

func Test_Worker_PathFindingUpdate(t *testing.T) {
	posA := math.NewPos(0, 0)
	posB := math.NewPos(5, 0)
	posC := math.NewPos(5, -2)
	posD := math.NewPos(0, -2)
	posE := math.NewPos(2, 2)
	building1Pos := math.NewPos(4, 4)
	building2Pos := math.NewPos(2, -1)

	workerPos := math.NewPos(2, 0)

	mapPaths := []math.Path{
		{Start: posA, End: posB, Cost: 2},
		{Start: posA, End: posD, Cost: 1},
		{Start: posA, End: posE, Cost: 3},

		{Start: posB, End: posC, Cost: 2},
		{Start: posB, End: building1Pos, Cost: 8},
		{Start: posB, End: posE, Cost: 3},

		{Start: posC, End: posD, Cost: 4},
		{Start: posC, End: building2Pos, Cost: 2},

		{Start: posD, End: building2Pos, Cost: 1},

		{Start: posE, End: building1Pos, Cost: 5},
	}

	workerID := uuid.New()
	streetAB_ID := uuid.New()
	streetAD_ID := uuid.New()
	streetAE_ID := uuid.New()

	streetBC_ID := uuid.New()
	streetBBuilding1ID := uuid.New()
	streetBE_ID := uuid.New()

	streetCD_ID := uuid.New()
	streetCBuilding2ID := uuid.New()

	streetDBuilding2ID := uuid.New()

	streetEBuilding1ID := uuid.New()

	building1ID := uuid.New()
	building2ID := uuid.New()

	workerStore := event.NewWorkerStore()
	streetStore := event.NewStreetStore()
	buildingStore := event.NewBuildingStore()

	workerOperator := operator.NewWorker(&workerStore, &buildingStore)
	buildingOperator := operator.NewBuilding(&buildingStore, nil)
	streetOperator := operator.NewStreet(&streetStore)

	w := WorkerAutomaton{
		EntityID: workerID,
		MapPaths: mapPaths,

		WorkerStore:   &workerStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		WorkerOperator: workerOperator,
		StreetOperator: streetOperator,
	}

	err := workerOperator.Init(workerID, workerPos)
	require.NoError(t, err)

	events, err := event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err := aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, workerState.ID)

	{
		err = streetOperator.Init(streetAB_ID, posA, posB)
		require.NoError(t, err)
		err = streetOperator.Init(streetAD_ID, posA, posD)
		require.NoError(t, err)
		err = streetOperator.Init(streetAE_ID, posA, posE)
		require.NoError(t, err)

		err = streetOperator.Init(streetBC_ID, posB, posC)
		require.NoError(t, err)
		err = streetOperator.Init(streetBE_ID, posB, posE)
		require.NoError(t, err)
		err = streetOperator.Init(streetBBuilding1ID, posB, building1Pos)
		require.NoError(t, err)

		err = streetOperator.Init(streetCD_ID, posC, posD)
		require.NoError(t, err)
		err = streetOperator.Init(streetCBuilding2ID, posC, building2Pos)
		require.NoError(t, err)

		err = streetOperator.Init(streetDBuilding2ID, posC, building2Pos)
		require.NoError(t, err)

		err = streetOperator.Init(streetEBuilding1ID, posC, building1Pos)
		require.NoError(t, err)
	}

	require.NoError(t, buildingOperator.Init(building1ID, event.BuildingTypeDrinkStore, building1Pos))
	require.NoError(t, buildingOperator.Init(building2ID, event.BuildingTypeDrinkStore, building2Pos))

	require.NoError(t, workerOperator.AssignBuilding(workerID, building1ID))

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	w.PathFindingUpdate()

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.Equal(t, []math.Pos{posA, posE, building1Pos}, workerState.PlannedPoses)
}

func Test_Worker_Move(t *testing.T) {
	posA := math.NewPos(0, 0)
	posB := math.NewPos(5, 0)

	workerPos := math.NewPos(2, 0)

	workerID := uuid.New()

	workerStore := event.NewWorkerStore()
	buildingStore := event.NewBuildingStore()

	workerOperator := operator.NewWorker(&workerStore, &buildingStore)

	err := workerOperator.Init(workerID, workerPos)
	require.NoError(t, err)

	events, err := event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err := aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, workerState.ID)

	w := WorkerAutomaton{
		EntityID: workerID,

		WorkerStore:   &workerStore,
		BuildingStore: &buildingStore,

		WorkerOperator: workerOperator,
	}

	require.NoError(t, w.WorkerOperator.ChangePlannedPoses(workerID, []math.Pos{posA, posB}))

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.Equal(t, []math.Pos{posA, posB}, workerState.PlannedPoses)

	w.Move()

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(0, 0), workerState.Position)

	require.NoError(t, w.WorkerOperator.ChangePlannedPoses(workerID, []math.Pos{posB}))
	w.Move()

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(2, 0), workerState.Position)

	w.Move()

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(4, 0), workerState.Position)

	w.Move()

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(5, 0), workerState.Position)

	// No planned position left, stay in place
	w.Move()

	events, err = event.Store(workerStore).GetEventsByEntityID(workerID)
	require.NoError(t, err)

	workerState, err = aggregator.GetWorkerState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(5, 0), workerState.Position)
}
