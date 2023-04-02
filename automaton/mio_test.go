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

	mioOperator := operator.NewMio(aggregator.NewMioAggregator(&mioStore), nil)

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

		mioOperator: mioOperator,
	}

	automaton.Automate()

	events, err = mioStore.GetEventsByEntityID(mioID)
	mioState, err = aggregator.GetMioState(events)
	require.NoError(t, err)

	require.Equal(t, streetID, mioState.StreetID)
}
