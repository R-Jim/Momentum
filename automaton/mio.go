package automaton

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
)

type mioAutomaton struct {
	entityID uuid.UUID

	mioStore    *store.Store
	streetStore *store.Store

	mioOperator operator.MioOperator
}

func (m mioAutomaton) Automate() {
	m.EnterStreetFromCurrentPosition()
}

func (m mioAutomaton) EnterStreetFromCurrentPosition() {
	events := (*m.mioStore).GetEvents()[m.entityID]
	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		fmt.Print(err)
	}

	for streetID, streetEvents := range (*m.streetStore).GetEvents() {
		streetState, err := aggregator.GetStreetState(streetEvents)
		if err != nil {
			fmt.Print(err)
		}
		if math.IsBetweenAAndB(mioState.Position, streetState.HeadA, streetState.HeadB, 1) {
			m.mioOperator.EnterStreet(mioState.ID, streetID)
		}
	}
}
