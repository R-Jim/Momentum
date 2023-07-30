package automaton

import (
	"testing"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// func Test_EnterStreetFromCurrentPosition(t *testing.T) {
// 	mioID := uuid.New()
// 	mioStore := store.NewStore()

// 	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}

// 	err := mioOperator.Init(mioID, math.NewPos(2, 2))
// 	require.NoError(t, err)

// 	events, err := mioStore.GetEventsByEntityID(mioID)
// 	require.NoError(t, err)

// 	mioState, err := aggregator.GetMioState(events)
// 	require.NoError(t, err)

// 	require.Equal(t, math.NewPos(2, 2), mioState.Position)

// 	streetStore := store.NewStore()
// 	streetOperator := operator.NewStreet(aggregator.NewStreetAggregator(&streetStore), nil)

// 	streetID := uuid.New()
// 	err = streetOperator.Init(streetID, math.NewPos(2, 0), math.NewPos(2, 4))
// 	require.NoError(t, err)

// 	//
// 	automaton := MioAutomaton{
// 		EntityID: mioID,

// 		MioStore:    &mioStore,
// 		StreetStore: &streetStore,

// 		MioOperator:    mioOperator,
// 		StreetOperator: streetOperator,
// 	}

// 	automaton.EnterStreetFromCurrentPosition()

// 	events, err = mioStore.GetEventsByEntityID(mioID)
// 	mioState, err = aggregator.GetMioState(events)
// 	require.NoError(t, err)

// 	require.Equal(t, streetID, mioState.StreetID)

// 	events, err = streetStore.GetEventsByEntityID(streetID)
// 	streetState, err := aggregator.GetStreetState(events)

// 	require.True(t, streetState.EntityMap[mioID])

// 	// Check moving to new pos
// 	newStreetID := uuid.New()
// 	err = streetOperator.Init(newStreetID, math.NewPos(4, 0), math.NewPos(4, 4))
// 	require.NoError(t, err)

// 	err = mioOperator.Walk(mioID, math.NewPos(3, 2))
// 	require.NoError(t, err)
// 	err = mioOperator.Walk(mioID, math.NewPos(4, 2))
// 	require.NoError(t, err)

// 	// automaton.EnterStreetFromCurrentPosition()

// 	events, err = mioStore.GetEventsByEntityID(mioID)
// 	mioState, err = aggregator.GetMioState(events)
// 	require.NoError(t, err)
// 	require.Equal(t, newStreetID, mioState.StreetID)

// 	events, err = streetStore.GetEventsByEntityID(streetID)
// 	streetState, err = aggregator.GetStreetState(events)

// 	require.False(t, streetState.EntityMap[mioID])

// 	events, err = streetStore.GetEventsByEntityID(newStreetID)
// 	newStreetState, err := aggregator.GetStreetState(events)

// 	require.True(t, newStreetState.EntityMap[mioID])
// }

func Test_MioMoodBehavior_Mood(t *testing.T) {
	mioID := uuid.New()
	mioStore := store.NewStore()
	streetStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}
	BuildingOperator := operator.BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	err := mioOperator.Init(mioID, math.NewPos(2, 2))
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, mioState.ID)

	mioHouseID := uuid.New()
	require.NoError(t, BuildingOperator.Init(mioHouseID, event.BuildingTypeMioHouse, math.NewPos(2, 0)))

	//
	automaton := MioAutomaton{
		EntityID: mioID,

		MioStore:      &mioStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		MioOperator: mioOperator,
	}

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 100, mioActivityState.MaxMood)
	require.Equal(t, 70, mioActivityState.Mood)

	require.NoError(t, mioOperator.Stream(mioID, 50))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 20, mioActivityState.Mood)

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, mioHouseID, mioState.SelectedBuildingID)

	require.NoError(t, mioOperator.Eat(mioID, 200))

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)
}

func Test_MioMoodBehavior_Energy(t *testing.T) {
	mioID := uuid.New()
	mioStore := store.NewStore()
	streetStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}
	BuildingOperator := operator.BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	err := mioOperator.Init(mioID, math.NewPos(2, 2))
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, mioState.ID)

	foodStoreID := uuid.New()
	require.NoError(t, BuildingOperator.Init(foodStoreID, event.BuildingTypeFoodStore, math.NewPos(2, 6)))

	//
	automaton := MioAutomaton{
		EntityID: mioID,

		MioStore:      &mioStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		MioOperator: mioOperator,
	}

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 100, mioActivityState.MaxEnergy)
	require.Equal(t, 70, mioActivityState.Energy)

	require.NoError(t, mioOperator.Starve(mioID, 50))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 20, mioActivityState.Energy)

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, foodStoreID, mioState.SelectedBuildingID)

	require.NoError(t, mioOperator.Eat(mioID, 50))

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)
}

func Test_MioMoodBehavior_Drink(t *testing.T) {
	mioID := uuid.New()
	mioStore := store.NewStore()
	streetStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}
	BuildingOperator := operator.BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	err := mioOperator.Init(mioID, math.NewPos(2, 2))
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, mioState.ID)

	drinkStoreID := uuid.New()
	require.NoError(t, BuildingOperator.Init(drinkStoreID, event.BuildingTypeDrinkStore, math.NewPos(4, 4)))

	//
	automaton := MioAutomaton{
		EntityID: mioID,

		MioStore:      &mioStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		MioOperator: mioOperator,
	}

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 100, mioActivityState.MaxDehydration)
	require.Equal(t, 70, mioActivityState.Dehydration)

	require.NoError(t, mioOperator.Sweat(mioID, 50))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 20, mioActivityState.Dehydration)

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, drinkStoreID, mioState.SelectedBuildingID)

	require.NoError(t, mioOperator.Drink(mioID, 50))

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)
}

func Test_MioMoodBehavior(t *testing.T) {
	mioID := uuid.New()
	mioStore := store.NewStore()
	streetStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}
	BuildingOperator := operator.BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}

	err := mioOperator.Init(mioID, math.NewPos(2, 2))
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, mioState.ID)

	mioHouseID := uuid.New()
	require.NoError(t, BuildingOperator.Init(mioHouseID, event.BuildingTypeMioHouse, math.NewPos(2, 0)))
	foodStoreID := uuid.New()
	require.NoError(t, BuildingOperator.Init(foodStoreID, event.BuildingTypeFoodStore, math.NewPos(2, 6)))
	drinkStoreID := uuid.New()
	require.NoError(t, BuildingOperator.Init(drinkStoreID, event.BuildingTypeDrinkStore, math.NewPos(4, 4)))

	//
	automaton := MioAutomaton{
		EntityID: mioID,

		MioStore:      &mioStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		MioOperator: mioOperator,
	}

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, uuid.Nil, mioState.SelectedBuildingID)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.NoError(t, mioOperator.Starve(mioID, 50))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 20, mioActivityState.Energy)

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, foodStoreID, mioState.SelectedBuildingID)

	require.NoError(t, mioOperator.Sweat(mioID, 50))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 20, mioActivityState.Dehydration)

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, drinkStoreID, mioState.SelectedBuildingID)

	require.NoError(t, mioOperator.Stream(mioID, 50))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 20, mioActivityState.Mood)

	automaton.MioMoodBehavior()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, mioHouseID, mioState.SelectedBuildingID)
}

func Test_Mio_PathFindingUpdate_simple(t *testing.T) {
	posA := math.NewPos(0, 0)
	posB := math.NewPos(5, 0)
	posC := math.NewPos(8, 0)
	buildingPos := math.NewPos(8, 2)

	mioPos := math.NewPos(2, 0)

	mapPaths := []math.Path{
		{Start: posA, End: posB, Cost: 10}, // street1
		{Start: posB, End: posC},           // street2
		{Start: posC, End: buildingPos},    // building street
	}

	mioID := uuid.New()
	street1ID := uuid.New()
	street2ID := uuid.New()
	buildingStreetID := uuid.New()

	mioStore := store.NewStore()
	buildingStore := store.NewStore()
	streetStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}
	BuildingOperator := operator.BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}
	streetOperator := operator.NewStreet(aggregator.NewStreetAggregator(&streetStore))

	m := MioAutomaton{
		EntityID: mioID,
		MapPaths: mapPaths,

		MioStore:      &mioStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		MioOperator:    mioOperator,
		StreetOperator: streetOperator,
	}

	err := mioOperator.Init(mioID, mioPos)
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
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

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.NoError(t, mioOperator.SelectBuilding(mioID, drinkStoreID))

	m.prevSelectedBuilding = drinkStoreID

	m.PathFindingUpdate()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, []math.Pos{posB, posC, buildingPos}, mioState.PlannedPoses)
}

func Test_Mio_PathFindingUpdate(t *testing.T) {
	posA := math.NewPos(0, 0)
	posB := math.NewPos(5, 0)
	posC := math.NewPos(5, -2)
	posD := math.NewPos(0, -2)
	posE := math.NewPos(2, 2)
	building1Pos := math.NewPos(4, 4)
	building2Pos := math.NewPos(2, -1)

	mioPos := math.NewPos(2, 0)

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

	mioID := uuid.New()
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

	mioStore := store.NewStore()
	buildingStore := store.NewStore()
	streetStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}
	buildingOperator := operator.BuildingOperator{
		BuildingAggregator: aggregator.NewBuildingAggregator(&buildingStore),
	}
	streetOperator := operator.NewStreet(aggregator.NewStreetAggregator(&streetStore))

	m := MioAutomaton{
		EntityID: mioID,
		MapPaths: mapPaths,

		MioStore:      &mioStore,
		StreetStore:   &streetStore,
		BuildingStore: &buildingStore,

		MioOperator:    mioOperator,
		StreetOperator: streetOperator,
	}

	err := mioOperator.Init(mioID, mioPos)
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, mioState.ID)

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

	require.NoError(t, mioOperator.SelectBuilding(mioID, building1ID))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	m.prevSelectedBuilding = building1ID

	m.PathFindingUpdate()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, []math.Pos{posA, posE, building1Pos}, mioState.PlannedPoses)
}

func Test_Mio_Move(t *testing.T) {
	posA := math.NewPos(0, 0)
	posB := math.NewPos(5, 0)

	mioPos := math.NewPos(2, 0)

	mioID := uuid.New()

	mioStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}

	err := mioOperator.Init(mioID, mioPos)
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, mioState.ID)

	m := MioAutomaton{
		EntityID: mioID,

		MioStore:      &mioStore,
		BuildingStore: &buildingStore,

		MioOperator: mioOperator,
	}

	require.NoError(t, m.MioOperator.ChangePlannedPoses(mioID, []math.Pos{posA, posB}))

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, []math.Pos{posA, posB}, mioState.PlannedPoses)

	m.Move()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(0, 0), mioState.Position)

	require.NoError(t, m.MioOperator.ChangePlannedPoses(mioID, []math.Pos{posB}))
	m.Move()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(2, 0), mioState.Position)

	m.Move()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(4, 0), mioState.Position)

	m.Move()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(5, 0), mioState.Position)

	// No planned position left, stay in place
	m.Move()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(5, 0), mioState.Position)
}

func Test_HourlyExhaustion(t *testing.T) {
	mioPos := math.NewPos(2, 0)

	mioID := uuid.New()

	mioStore := store.NewStore()
	buildingStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}

	err := mioOperator.Init(mioID, mioPos)
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, mioState.ID)

	m := MioAutomaton{
		EntityID: mioID,

		MioStore:      &mioStore,
		BuildingStore: &buildingStore,

		MioOperator: mioOperator,
	}

	m.HourlyExhaustion()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioActivityState, err := aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 65, mioActivityState.Energy)
	require.Equal(t, 65, mioActivityState.Dehydration)

	m.HourlyExhaustion()

	events, err = mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioActivityState, err = aggregator.GetMioActivityState(events)
	require.NoError(t, err)

	require.Equal(t, 60, mioActivityState.Energy)
	require.Equal(t, 60, mioActivityState.Dehydration)
}
