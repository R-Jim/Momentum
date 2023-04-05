package automaton

import (
	"testing"

	"github.com/R-jim/Momentum/aggregate/aggregator"
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

	automaton.Automate()

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

	automaton.Automate()

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
