package operator

import (
	"testing"

	"github.com/R-jim/Momentum/aggregate/aggregator"
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

	err = mioOperator.Run(mioID, math.NewPos(4, 3))
	require.NoError(t, err)

	events, err := store.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(4, 3), mioState.Position)

	// invalid run distant, too short
	err = mioOperator.Run(mioID, math.NewPos(5, 3))
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
	require.NoError(t, BuildingOperator.Init(buildingID, math.NewPos(2, 2)))

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
	require.NoError(t, BuildingOperator.Init(buildingID, math.NewPos(2, 2)))

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
