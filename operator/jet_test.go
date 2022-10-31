package operator

import (
	"context"
	"errors"
	"testing"

	"github.com/R-jim/Momentum/domain/jet"
	"github.com/R-jim/Momentum/domain/storage"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestImpl_Fly(t *testing.T) {
	jetID := "jet_1"
	fuelTankID := "fuel_tank_1"

	fuelTankInitEvent := storage.NewInitEvent(fuelTankID)
	jetInitEvent := jet.NewInitEvent(jetID)
	jetChangeFuelTankEvent := jet.NewChangeFuelTankEvent(jetID, fuelTankID)
	fuelTankRefillEvent := storage.NewRefillEvent(fuelTankID, 10)
	jetTakeOffEvent := jet.NewTakeOffEvent(jetID)

	type arg struct {
		givenFuelTankEvents      []storage.Event
		givenJetEvents           []jet.Event
		givenFuelConsumeQuantity int
		givenToPosition          jet.PositionState
		expFuelState             storage.State
		expJetPositionState      jet.PositionState
		expErr                   error
	}

	tcs := map[string]arg{
		"success": {
			givenFuelTankEvents: []storage.Event{
				fuelTankInitEvent,
				fuelTankRefillEvent,
			},
			givenJetEvents: []jet.Event{
				jetInitEvent,
				jetChangeFuelTankEvent,
				jetTakeOffEvent,
			},
			givenFuelConsumeQuantity: 1,
			givenToPosition: jet.PositionState{
				X: 1,
				Y: 1,
			},
			expFuelState: storage.State{
				Quantity: 9,
			},
			expJetPositionState: jet.PositionState{
				X: 1,
				Y: 1,
			},
		},
		"failure, fuel consume exceed": {
			givenFuelTankEvents: []storage.Event{
				fuelTankInitEvent,
			},
			givenJetEvents: []jet.Event{
				jetInitEvent,
				jetChangeFuelTankEvent,
				jetTakeOffEvent,
			},
			givenFuelConsumeQuantity: 1,
			givenToPosition: jet.PositionState{
				X: 1,
				Y: 1,
			},
			expFuelState: storage.State{
				Quantity: 0,
			},
			expJetPositionState: jet.PositionState{
				X: 0,
				Y: 0,
			},
			expErr: errors.New("[FUEL_TANK_CONSUMED] aggregate fail"),
		},
		"failure, jet landed": {
			givenFuelTankEvents: []storage.Event{
				fuelTankInitEvent,
				fuelTankRefillEvent,
			},
			givenJetEvents: []jet.Event{
				jetInitEvent,
				jetChangeFuelTankEvent,
			},
			givenFuelConsumeQuantity: 1,
			givenToPosition: jet.PositionState{
				X: 1,
				Y: 1,
			},
			expFuelState: storage.State{
				Quantity: 10,
			},
			expJetPositionState: jet.PositionState{
				X: 0,
				Y: 0,
			},
			expErr: errors.New("[JET_FLEW] aggregate fail"),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			jetStore := jet.NewStore()
			fuelTankStore := storage.NewStore()

			jetAggregator := jet.NewAggregator(jetStore)
			fuelTankAggregator := storage.NewAggregator(fuelTankStore)

			for _, fuelTankEvent := range tc.givenFuelTankEvents {
				err := fuelTankAggregator.Aggregate(fuelTankEvent)
				require.NoError(t, err)
			}

			for _, jetEvent := range tc.givenJetEvents {
				err := jetAggregator.Aggregate(jetEvent)
				require.NoError(t, err)
			}

			jetOperator := jetOperator{
				jetAggregator:      jetAggregator,
				fuelTankAggregator: fuelTankAggregator,
			}
			// WHEN
			err := jetOperator.Fly(jetID, fuelTankID, tc.givenFuelConsumeQuantity, tc.givenToPosition)
			// THEN
			if tc.expErr != nil {
				require.NotNil(t, err)
				require.EqualError(t, tc.expErr, err.Error())
			} else {
				require.NoError(t, err)
			}

			fuelTankState, err := storage.GetState(fuelTankStore, fuelTankID)
			require.NoError(t, err)
			assertFuelState(t, tc.expFuelState, fuelTankState)

			jetPositionState, err := jet.GetPositionState(jetStore, jetID)
			require.NoError(t, err)
			assertJetPositionState(t, tc.expJetPositionState, jetPositionState)
		})
	}
}

func Test_Concurrency(t *testing.T) {
	jetID := "jet_1"
	fuelTankID := "fuel_tank_1"

	fuelTankInitEvent := storage.NewInitEvent(fuelTankID)
	jetInitEvent := jet.NewInitEvent(jetID)
	jetChangeFuelTankEvent := jet.NewChangeFuelTankEvent(jetID, fuelTankID)
	fuelTankRefillEvent := storage.NewRefillEvent(fuelTankID, 10)
	jetTakeOffEvent := jet.NewTakeOffEvent(jetID)

	type arg struct {
		givenFuelTankEvents      []storage.Event
		givenJetEvents           []jet.Event
		givenFuelConsumeQuantity int
		givenToPosition          jet.PositionState
		expFuelState             storage.State
		expJetPositionState      jet.PositionState
	}

	tcs := map[string]arg{
		"success": {
			givenFuelTankEvents: []storage.Event{
				fuelTankInitEvent,
				fuelTankRefillEvent,
			},
			givenJetEvents: []jet.Event{
				jetInitEvent,
				jetChangeFuelTankEvent,
				jetTakeOffEvent,
			},
			givenFuelConsumeQuantity: 1,
			givenToPosition: jet.PositionState{
				X: 1,
				Y: 1,
			},
			expFuelState: storage.State{
				Quantity: 8,
			},
			expJetPositionState: jet.PositionState{
				X: 1,
				Y: 1,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			jetStore := jet.NewStore()
			fuelTankStore := storage.NewStore()

			jetAggregator := jet.NewAggregator(jetStore)
			fuelTankAggregator := storage.NewAggregator(fuelTankStore)

			for _, fuelTankEvent := range tc.givenFuelTankEvents {
				err := fuelTankAggregator.Aggregate(fuelTankEvent)
				require.NoError(t, err)
			}

			for _, jetEvent := range tc.givenJetEvents {
				err := jetAggregator.Aggregate(jetEvent)
				require.NoError(t, err)
			}

			jetOperator := jetOperator{
				jetAggregator:      jetAggregator,
				fuelTankAggregator: fuelTankAggregator,
			}
			// WHEN
			g, _ := errgroup.WithContext(context.Background())
			g.SetLimit(1)

			operations := []func() error{
				func() error {
					return jetOperator.Fly(jetID, fuelTankID, tc.givenFuelConsumeQuantity, tc.givenToPosition)
				},
				func() error {
					return jetOperator.Fly(jetID, fuelTankID, tc.givenFuelConsumeQuantity, tc.givenToPosition)
				},
				func() error {
					return jetOperator.Landing(jetID)
				},
				func() error {
					return jetOperator.Fly(jetID, fuelTankID, tc.givenFuelConsumeQuantity, tc.givenToPosition)
				},
			}

			for _, operation := range operations {
				op := operation
				g.Go(func() error {
					return op()
				})
			}
			// THEN
			g.Wait()

			fuelTankState, err := storage.GetState(fuelTankStore, fuelTankID)
			require.NoError(t, err)
			assertFuelState(t, tc.expFuelState, fuelTankState)

			jetPositionState, err := jet.GetPositionState(jetStore, jetID)
			require.NoError(t, err)
			assertJetPositionState(t, tc.expJetPositionState, jetPositionState)
		})
	}
}

func assertFuelState(t *testing.T, expFuelState, result storage.State) {
	// require.Equal(t, expFuelState.ID, result.ID) // WARNING: should not be used
	require.Equal(t, expFuelState.Quantity, result.Quantity)
}

func assertJetPositionState(t *testing.T, expPositionState, result jet.PositionState) {
	require.Equal(t, expPositionState.X, result.X)
	require.Equal(t, expPositionState.Y, result.Y)
}
