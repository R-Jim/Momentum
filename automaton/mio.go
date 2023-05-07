package automaton

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
)

type mioAutomaton struct {
	entityID uuid.UUID

	mioStore      *store.Store
	streetStore   *store.Store
	buildingStore *store.Store

	mioOperator    operator.MioOperator
	streetOperator operator.StreetOperator
}

func (m mioAutomaton) Automate() {
	m.EnterStreetFromCurrentPosition()
	m.MioMoodBehavior()
}

func (m mioAutomaton) EnterStreetFromCurrentPosition() {
	events := (*m.mioStore).GetEvents()[m.entityID]
	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		fmt.Print(err)
	}

	oldStreetID := mioState.StreetID

	for streetID, streetEvents := range (*m.streetStore).GetEvents() {
		streetState, err := aggregator.GetStreetState(streetEvents)
		if err != nil {
			fmt.Print(err)
		}
		if math.IsBetweenAAndB(mioState.Position, streetState.HeadA, streetState.HeadB, 1) {
			if oldStreetID != streetID {
				err := m.mioOperator.EnterStreet(mioState.ID, streetID)
				if err != nil {
					fmt.Print(err)
				}
				err = m.streetOperator.EntityEnter(streetID, mioState.ID)
				if err != nil {
					fmt.Print(err)
				}
				if oldStreetID.String() != uuid.Nil.String() {
					err = m.streetOperator.EntityLeave(oldStreetID, mioState.ID)
					if err != nil {
						fmt.Print(err)
					}
				}
				break
			}
		}
	}
}

func (m mioAutomaton) MioMoodBehavior() {
	events := (*m.mioStore).GetEvents()[m.entityID]
	mioActivityState, err := aggregator.GetMioActivityState(events)
	if err != nil {
		fmt.Print(err)
	}

	isBored := mioActivityState.Mood < 30
	isHungry := mioActivityState.Energy < 30
	isThirsty := mioActivityState.Dehydration < 30

	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		fmt.Print(err)
	}

	if !isBored && !isHungry && !isThirsty {
		err = m.mioOperator.UnselectBuilding(mioState.ID, mioState.BuildingID)
		if err != nil {
			fmt.Print(err)
		}
		return
	}

	var selectedBuildingID uuid.UUID
	var selectedBuildingDistance float64

	for buildingID, buildingEvents := range (*m.buildingStore).GetEvents() {
		buildingState, err := aggregator.GetBuildingState(buildingEvents)
		if err != nil {
			fmt.Print(err)
		}

		if isBuildingFitMood(buildingState, isBored, isHungry, isThirsty) {
			_, _, distance := math.GetDistances(mioState.Position, buildingState.Pos) // TODO: should check by path
			if selectedBuildingID.String() == uuid.Nil.String() || distance < selectedBuildingDistance {
				selectedBuildingID = buildingID
				selectedBuildingDistance = distance
			}
		}
	}

	if selectedBuildingID.String() == uuid.Nil.String() {
		return
	}

	err = m.mioOperator.SelectBuilding(mioState.ID, selectedBuildingID)
	if err != nil {
		fmt.Print(err)
	}
}

func isBuildingFitMood(buildingState aggregator.BuildingState, isBored, isHungry, isThirsty bool) bool {
	if isBored && buildingState.Type == event.BuildingTypeMioHouse.String() {
		return true
	}
	if isHungry && buildingState.Type == event.BuildingTypeFoodStore.String() {
		return true
	}
	if isThirsty && buildingState.Type == event.BuildingTypeDrinkStore.String() {
		return true
	}

	return false
}
