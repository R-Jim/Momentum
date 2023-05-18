package operator

import (
	"testing"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Mio_Init(t *testing.T) {
	mioID := uuid.New()
	store := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&store),
	}

	err := mioOperator.Init(mioID, math.NewPos(10, 5))
	require.NoError(t, err)

	events, err := store.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(10, 5), mioState.Position)
}

func Test_Mio_Walk(t *testing.T) {
	mioID := uuid.New()
	store := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&store),
	}

	err := mioOperator.Init(mioID, math.NewPos(2, 2))
	require.NoError(t, err)

	err = mioOperator.Walk(mioID, math.NewPos(3, 2))
	require.NoError(t, err)

	events, err := store.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(3, 2), mioState.Position)

	// invalid walk distant
	err = mioOperator.Walk(mioID, math.NewPos(5, 2))
	require.ErrorContains(t, err, aggregator.ErrAggregateFail.Error())
}

func Test_Mio_Run(t *testing.T) {
	mioID := uuid.New()
	store := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&store),
	}

	err := mioOperator.Init(mioID, math.NewPos(2, 2))
	require.NoError(t, err)

	err = mioOperator.Run(mioID, math.NewPos(2, 4))
	require.NoError(t, err)

	events, err := store.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(2, 4), mioState.Position)

	// invalid run distant, too short
	err = mioOperator.Run(mioID, math.NewPos(2, 4.5))
	require.ErrorContains(t, err, aggregator.ErrAggregateFail.Error())

	// invalid run distant, too long
	err = mioOperator.Run(mioID, math.NewPos(8, 8))
	require.ErrorContains(t, err, aggregator.ErrAggregateFail.Error())
}

func Test_Mio_Idle(t *testing.T) {
	mioID := uuid.New()
	store := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&store),
	}

	err := mioOperator.Init(mioID, math.NewPos(2, 2))
	require.NoError(t, err)

	err = mioOperator.Idle(mioID)
	require.NoError(t, err)
}

func Test_Mio_EnterBuilding(t *testing.T) {
	mioStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator:      aggregator.NewMioAggregator(&mioStore),
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	BuildingOperator := BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	mioID := uuid.New()
	buildingID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))
	require.NoError(t, BuildingOperator.Init(buildingID, event.BuildingTypeMioHouse, math.NewPos(2, 2)))

	require.NoError(t, mioOperator.EnterBuilding(mioID, buildingID))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, buildingID, mioState.BuildingID)

	events, err = buildingStore.GetEventsByEntityID(buildingID)
	require.NoError(t, err)
	buildingState, err := aggregator.GetBuildingState(events)
	require.NoError(t, err)
	require.True(t, buildingState.EntityMap[mioID])

	require.Error(t, mioOperator.EnterBuilding(mioID, buildingID))
}

func Test_Mio_SelectBuilding(t *testing.T) {
	mioStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&mioStore),
	}

	mioID := uuid.New()
	buildingID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))

	require.NoError(t, mioOperator.SelectBuilding(mioID, buildingID))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, buildingID, mioState.SelectedBuildingID)

	require.NoError(t, mioOperator.SelectBuilding(mioID, buildingID))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, buildingID, mioState.SelectedBuildingID)
}

func Test_Mio_UnselectBuilding(t *testing.T) {
	mioStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&mioStore),
	}

	mioID := uuid.New()
	buildingID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))

	require.NoError(t, mioOperator.SelectBuilding(mioID, buildingID))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, buildingID, mioState.SelectedBuildingID)

	require.NoError(t, mioOperator.UnselectBuilding(mioID, buildingID))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)

	require.NoError(t, mioOperator.UnselectBuilding(mioID, buildingID))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)
}

func Test_Mio_LeaveBuilding(t *testing.T) {
	mioStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator:      aggregator.NewMioAggregator(&mioStore),
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	BuildingOperator := BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	mioID := uuid.New()
	buildingID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))
	require.NoError(t, BuildingOperator.Init(buildingID, event.BuildingTypeMioHouse, math.NewPos(2, 2)))

	require.NoError(t, mioOperator.EnterBuilding(mioID, buildingID))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, buildingID, mioState.BuildingID)

	events, err = buildingStore.GetEventsByEntityID(buildingID)
	require.NoError(t, err)
	buildingState, err := aggregator.GetBuildingState(events)
	require.NoError(t, err)
	require.True(t, buildingState.EntityMap[mioID])

	require.NoError(t, mioOperator.LeaveBuilding(mioID, buildingID))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, uuid.Nil, mioState.BuildingID)

	events, err = buildingStore.GetEventsByEntityID(buildingID)
	require.NoError(t, err)
	buildingState, err = aggregator.GetBuildingState(events)
	require.NoError(t, err)
	require.False(t, buildingState.EntityMap[mioID])

	require.Error(t, mioOperator.LeaveBuilding(mioID, buildingID))
}

func Test_Mio_Act(t *testing.T) {
	mioStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator:      aggregator.NewMioAggregator(&mioStore),
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	BuildingOperator := BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	mioID := uuid.New()
	buildingID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))
	require.NoError(t, BuildingOperator.Init(buildingID, event.BuildingTypeMioHouse, math.NewPos(2, 2)))

	require.NoError(t, mioOperator.EnterBuilding(mioID, buildingID))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, buildingID, mioState.BuildingID)

	events, err = buildingStore.GetEventsByEntityID(buildingID)
	require.NoError(t, err)
	buildingState, err := aggregator.GetBuildingState(events)
	require.NoError(t, err)
	require.True(t, buildingState.EntityMap[mioID])

	require.NoError(t, mioOperator.Act(mioID, buildingID))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	require.Equal(t, event.MioActEffect, events[len(events)-1].Effect)

	events, err = buildingStore.GetEventsByEntityID(buildingID)
	require.NoError(t, err)

	require.Equal(t, event.BuildingEntityActEffect, events[len(events)-1].Effect)
}

func Test_Mio_Stream(t *testing.T) {
	mioStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&mioStore),
	}

	mioID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, mioID, mioState.ID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 100, mioActivityState.MaxMood)
	require.Equal(t, 70, mioActivityState.Mood)

	require.NoError(t, mioOperator.Stream(mioID, 12))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 58, mioActivityState.Mood)

	require.NoError(t, mioOperator.Stream(mioID, 60))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 0, mioActivityState.Mood)
}

func Test_Mio_Eat(t *testing.T) {
	mioStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&mioStore),
	}

	mioID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, mioID, mioState.ID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 100, mioActivityState.MaxEnergy)
	require.Equal(t, 70, mioActivityState.Energy)

	require.NoError(t, mioOperator.Eat(mioID, 12))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 73, mioActivityState.Mood)
	require.Equal(t, 82, mioActivityState.Energy)

	require.NoError(t, mioOperator.Eat(mioID, 60))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 88, mioActivityState.Mood)
	require.Equal(t, 100, mioActivityState.Energy)
}

func Test_Mio_Starve(t *testing.T) {
	mioStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&mioStore),
	}

	mioID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, mioID, mioState.ID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 100, mioActivityState.MaxEnergy)
	require.Equal(t, 70, mioActivityState.Energy)

	require.NoError(t, mioOperator.Starve(mioID, 12))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 58, mioActivityState.Energy)

	require.NoError(t, mioOperator.Starve(mioID, 60))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 0, mioActivityState.Energy)
}

func Test_Mio_Drink(t *testing.T) {
	mioStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&mioStore),
	}

	mioID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, mioID, mioState.ID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 100, mioActivityState.MaxDehydration)
	require.Equal(t, 70, mioActivityState.Dehydration)

	require.NoError(t, mioOperator.Drink(mioID, 12))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 73, mioActivityState.Mood)
	require.Equal(t, 82, mioActivityState.Dehydration)

	require.NoError(t, mioOperator.Drink(mioID, 60))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 88, mioActivityState.Mood)
	require.Equal(t, 100, mioActivityState.Dehydration)
}

func Test_Mio_Sweat(t *testing.T) {
	mioStore := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&mioStore),
	}

	mioID := uuid.New()

	require.NoError(t, mioOperator.Init(mioID, math.NewPos(2, 2)))

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, mioID, mioState.ID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 100, mioActivityState.MaxDehydration)
	require.Equal(t, 70, mioActivityState.Dehydration)

	require.NoError(t, mioOperator.Sweat(mioID, 12))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 58, mioActivityState.Dehydration)

	require.NoError(t, mioOperator.Sweat(mioID, 60))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 0, mioActivityState.Dehydration)
}

func Test_Mio_ChangePlannedPoses(t *testing.T) {
	mioID := uuid.New()
	store := store.NewStore()

	mioOperator := MioOperator{
		MioAggregator: aggregator.NewMioAggregator(&store),
	}

	err := mioOperator.Init(mioID, math.NewPos(10, 5))
	require.NoError(t, err)

	plannedPos := []math.Pos{
		math.NewPos(3, 3),
		math.NewPos(2, 5),
	}
	mioOperator.ChangePlannedPoses(mioID, plannedPos)

	events, err := store.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, plannedPos, mioState.PlannedPoses)

	plannedPos = []math.Pos{
		math.NewPos(2, 5),
		math.NewPos(1, 5),
	}
	mioOperator.ChangePlannedPoses(mioID, plannedPos)

	events, err = store.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, plannedPos, mioState.PlannedPoses)
}
