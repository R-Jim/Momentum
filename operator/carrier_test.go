package operator

import (
	"errors"
	"testing"

	"github.com/R-jim/Momentum/aggregate/carrier"
	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/stretchr/testify/require"
)

func TestImpl_HouseJet(t *testing.T) {
	carrierID := "carrier_1"
	jet1ID := "jet_1"
	jet2ID := "jet_2"

	jet1InitEvent := jet.NewInitEvent(jet1ID)
	jet1LandingEvent := jet.NewLandingEvent(jet1ID)
	jet2InitEvent := jet.NewInitEvent(jet2ID)
	jet2TakeOffEvent := jet.NewTakeOffEvent(jet2ID)

	carrierInitEvent := carrier.NewInitEvent(carrierID)
	carrierHouseJet2Event := carrier.NewHouseJetEvent(carrierID, jet2ID)

	type arg struct {
		givenJetID            string
		givenJetEvents        []jet.Event
		givenCarrierEvents    []carrier.Event
		expErr                error
		expCarrierCombatState carrier.CombatState
		expJetCombatStatus    jet.Status
	}

	tcs := map[string]arg{
		"success, jet 1": {
			givenJetID: jet1ID,
			givenJetEvents: []jet.Event{
				jet1InitEvent,
			},
			givenCarrierEvents: []carrier.Event{
				carrierInitEvent,
			},
			expCarrierCombatState: carrier.CombatState{
				Status: carrier.IdleStatus,
				Jets:   carrier.JetIDs{jet1ID},
			},
			expJetCombatStatus: jet.LandedStatus,
		},
		"success, jet 2": {
			givenJetID: jet2ID,
			givenJetEvents: []jet.Event{
				jet2InitEvent,
				jet2TakeOffEvent,
			},
			givenCarrierEvents: []carrier.Event{
				carrierInitEvent,
			},
			expCarrierCombatState: carrier.CombatState{
				Status: carrier.IdleStatus,
				Jets:   carrier.JetIDs{jet2ID},
			},
			expJetCombatStatus: jet.LandedStatus,
		},
		"failure, jet 1 with landed status": {
			givenJetID: jet1ID,
			givenJetEvents: []jet.Event{
				jet1InitEvent,
				jet1LandingEvent,
			},
			givenCarrierEvents: []carrier.Event{
				carrierInitEvent,
			},
			expErr: errors.New("[JET_LANDED] aggregate fail"),
			expCarrierCombatState: carrier.CombatState{
				Status: carrier.IdleStatus,
				Jets:   carrier.JetIDs{},
			},
			expJetCombatStatus: jet.LandedStatus,
		},
		"failure with rollback, jet 2 already landed in carrier": {
			givenJetID: jet2ID,
			givenJetEvents: []jet.Event{
				jet2InitEvent,
				jet2TakeOffEvent,
			},
			givenCarrierEvents: []carrier.Event{
				carrierInitEvent,
				carrierHouseJet2Event,
			},
			expErr: errors.New("[CARRIER_JET_HOUSED] aggregate fail"),
			expCarrierCombatState: carrier.CombatState{
				Status: carrier.IdleStatus,
				Jets:   carrier.JetIDs{jet2ID},
			},
			expJetCombatStatus: jet.FlyingStatus,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			jetStore := jet.NewStore()
			carrierStore := carrier.NewStore()

			jetAggregator := jet.NewAggregator(jetStore)
			carrierAggregator := carrier.NewAggregator(carrierStore)

			for _, jetEvent := range tc.givenJetEvents {
				err := jetAggregator.Aggregate(jetEvent)
				require.NoError(t, err)
			}

			for _, carrierEvent := range tc.givenCarrierEvents {
				err := carrierAggregator.Aggregate(carrierEvent)
				require.NoError(t, err)
			}

			jetOperator := jetOperator{
				jetAggregator: jetAggregator,
			}

			carrierOperator := carrierOperator{
				jetOperator:       jetOperator,
				carrierAggregator: carrierAggregator,
			}
			// WHEN
			err := carrierOperator.HouseJet(carrierID, tc.givenJetID)
			// THEN
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}

			carrierCombatState, err := carrier.GetCombatState(carrierStore, carrierID)
			require.NoError(t, err)
			assertCarrierCombatState(t, tc.expCarrierCombatState, carrierCombatState)

			jetCombatState, err := jet.GetCombatState(jetStore, tc.givenJetID)
			require.NoError(t, err)
			require.Equal(t, tc.expJetCombatStatus, jetCombatState.Status)
		})
	}
}

func TestImpl_LaunchJet(t *testing.T) {
	carrierID := "carrier_1"
	jet1ID := "jet_1"

	jet1InitEvent := jet.NewInitEvent(jet1ID)
	jet1TakeOffEvent := jet.NewTakeOffEvent(jet1ID)

	carrierInitEvent := carrier.NewInitEvent(carrierID)
	carrierHouseJet1Event := carrier.NewHouseJetEvent(carrierID, jet1ID)

	type arg struct {
		givenJetID            string
		givenJetEvents        []jet.Event
		givenCarrierEvents    []carrier.Event
		expErr                error
		expCarrierCombatState carrier.CombatState
		expJetCombatStatus    jet.Status
	}

	tcs := map[string]arg{
		"success, jet 1": {
			givenJetID: jet1ID,
			givenJetEvents: []jet.Event{
				jet1InitEvent,
			},
			givenCarrierEvents: []carrier.Event{
				carrierInitEvent,
				carrierHouseJet1Event,
			},
			expCarrierCombatState: carrier.CombatState{
				Status: carrier.IdleStatus,
				Jets:   carrier.JetIDs{},
			},
			expJetCombatStatus: jet.FlyingStatus,
		},
		"failure, jet 1 with fly status": {
			givenJetID: jet1ID,
			givenJetEvents: []jet.Event{
				jet1InitEvent,
				jet1TakeOffEvent,
			},
			givenCarrierEvents: []carrier.Event{
				carrierInitEvent,
			},
			expErr: errors.New("[JET_TOOK_OFF] aggregate fail"),
			expCarrierCombatState: carrier.CombatState{
				Status: carrier.IdleStatus,
				Jets:   carrier.JetIDs{},
			},
			expJetCombatStatus: jet.FlyingStatus,
		},
		"failure with rollback, jet 1 not house in carrier": {
			givenJetID: jet1ID,
			givenJetEvents: []jet.Event{
				jet1InitEvent,
			},
			givenCarrierEvents: []carrier.Event{
				carrierInitEvent,
			},
			expErr: errors.New("[CARRIER_JET_LAUNCHED] aggregate fail"),
			expCarrierCombatState: carrier.CombatState{
				Status: carrier.IdleStatus,
				Jets:   carrier.JetIDs{},
			},
			expJetCombatStatus: jet.LandedStatus,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// GIVEN
			jetStore := jet.NewStore()
			carrierStore := carrier.NewStore()

			jetAggregator := jet.NewAggregator(jetStore)
			carrierAggregator := carrier.NewAggregator(carrierStore)

			for _, jetEvent := range tc.givenJetEvents {
				err := jetAggregator.Aggregate(jetEvent)
				require.NoError(t, err)
			}

			for _, carrierEvent := range tc.givenCarrierEvents {
				err := carrierAggregator.Aggregate(carrierEvent)
				require.NoError(t, err)
			}

			jetOperator := jetOperator{
				jetAggregator: jetAggregator,
			}

			carrierOperator := carrierOperator{
				jetOperator:       jetOperator,
				carrierAggregator: carrierAggregator,
			}
			// WHEN
			err := carrierOperator.LaunchJet(carrierID, tc.givenJetID)
			// THEN
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
			}

			carrierCombatState, err := carrier.GetCombatState(carrierStore, carrierID)
			require.NoError(t, err)
			assertCarrierCombatState(t, tc.expCarrierCombatState, carrierCombatState)

			jetCombatState, err := jet.GetCombatState(jetStore, tc.givenJetID)
			require.NoError(t, err)
			require.Equal(t, tc.expJetCombatStatus, jetCombatState.Status)
		})
	}
}

func assertCarrierCombatState(t *testing.T, expCombatState, result carrier.CombatState) {
	require.Equal(t, expCombatState.Status, result.Status)
	require.Equal(t, expCombatState.Jets, result.Jets)
}
