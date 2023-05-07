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

func Test_EnterStreetFromCurrentPosition(t *testing.T) {
	mioID := uuid.New()
	mioStore := store.NewStore()

	mioOperator := operator.MioOperator{MioAggregator: aggregator.NewMioAggregator(&mioStore)}

	err := mioOperator.Init(mioID, math.NewPos(2, 2))
	require.NoError(t, err)

	events, err := mioStore.GetEventsByEntityID(mioID)
	require.NoError(t, err)

	mioState, err := aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, math.NewPos(2, 2), mioState.Position)

	streetStore := store.NewStore()
	streetOperator := operator.NewStreet(aggregator.NewStreetAggregator(&streetStore), nil)

	streetID := uuid.New()
	err = streetOperator.Init(streetID, math.NewPos(2, 0), math.NewPos(2, 4))
	require.NoError(t, err)

	//
	automaton := mioAutomaton{
		entityID: mioID,

		mioStore:    &mioStore,
		streetStore: &streetStore,

		mioOperator:    mioOperator,
		streetOperator: streetOperator,
	}

	automaton.EnterStreetFromCurrentPosition()

	events, err = mioStore.GetEventsByEntityID(mioID)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, streetID, mioState.StreetID)

	events, err = streetStore.GetEventsByEntityID(streetID)
	streetState, err := aggregator.GetStreetState(events)

	require.True(t, streetState.EntityMap[mioID])

	// Check moving to new pos
	newStreetID := uuid.New()
	err = streetOperator.Init(newStreetID, math.NewPos(4, 0), math.NewPos(4, 4))
	require.NoError(t, err)

	err = mioOperator.Walk(mioID, math.NewPos(3, 2))
	require.NoError(t, err)
	err = mioOperator.Walk(mioID, math.NewPos(4, 2))
	require.NoError(t, err)

	automaton.EnterStreetFromCurrentPosition()

	events, err = mioStore.GetEventsByEntityID(mioID)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)
	require.Equal(t, newStreetID, mioState.StreetID)

	events, err = streetStore.GetEventsByEntityID(streetID)
	streetState, err = aggregator.GetStreetState(events)

	require.False(t, streetState.EntityMap[mioID])

	events, err = streetStore.GetEventsByEntityID(newStreetID)
	newStreetState, err := aggregator.GetStreetState(events)

	require.True(t, newStreetState.EntityMap[mioID])
}

func Test_MioMoodBehavior_Mood(t *testing.T) {
	mioID := uuid.New()
	mioStore := store.NewStore()
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
	automaton := mioAutomaton{
		entityID: mioID,

		mioStore:      &mioStore,
		buildingStore: &buildingStore,

		mioOperator: mioOperator,
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
	automaton := mioAutomaton{
		entityID: mioID,

		mioStore:      &mioStore,
		buildingStore: &buildingStore,

		mioOperator: mioOperator,
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
	automaton := mioAutomaton{
		entityID: mioID,

		mioStore:      &mioStore,
		buildingStore: &buildingStore,

		mioOperator: mioOperator,
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
	automaton := mioAutomaton{
		entityID: mioID,

		mioStore:      &mioStore,
		buildingStore: &buildingStore,

		mioOperator: mioOperator,
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
