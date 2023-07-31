package automaton

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
)

type MioAutomaton struct {
	EntityID uuid.UUID

	MapPaths []math.Path

	MioStore      *event.MioStore
	StreetStore   *event.StreetStore
	BuildingStore *event.BuildingStore

	MioOperator    operator.MioOperator
	StreetOperator operator.StreetOperator

	// For performance optimization
	prevSelectedBuilding uuid.UUID
}

type pathWithCost struct {
	poses []math.Pos
	cost  float64
}

func (m MioAutomaton) Automate() {
	// m.EnterStreetFromCurrentPosition() // Deprecated: should do it from street side
	m.MioMoodBehavior()
	m.PathFindingUpdate()
	m.Move()
	m.HourlyExhaustion()
}

func (m *MioAutomaton) MioMoodBehavior() {
	events := event.Store(*m.MioStore).GetEvents()[m.EntityID]
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
		err = m.MioOperator.UnselectBuilding(mioState.ID, mioState.BuildingID)
		if err != nil {
			fmt.Print(err)
		}
		return
	}

	var selectedBuildingID uuid.UUID
	var selectedBuildingDistance float64

	for buildingID, buildingEvents := range event.Store(*m.BuildingStore).GetEvents() {
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

	if m.prevSelectedBuilding == selectedBuildingID {
		return
	}

	m.prevSelectedBuilding = selectedBuildingID

	err = m.MioOperator.SelectBuilding(mioState.ID, selectedBuildingID)
	if err != nil {
		fmt.Print(err)
	}
}

func (m MioAutomaton) PathFindingUpdate() {
	events := event.Store(*m.MioStore).GetEvents()[m.EntityID]
	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		fmt.Print(err)
	}

	if mioState.SelectedBuildingID == uuid.Nil {
		return
	}

	buildingEvents := event.Store(*m.BuildingStore).GetEvents()[mioState.SelectedBuildingID]
	buildingState, err := aggregator.GetBuildingState(buildingEvents)
	if err != nil {
		fmt.Print(err)
	}

	if mioState.Position == buildingState.Pos {
		err = m.MioOperator.UnselectBuilding(m.EntityID, mioState.SelectedBuildingID)
		if err != nil {
			fmt.Print(err)
		}
		err = m.MioOperator.ChangePlannedPoses(m.EntityID, nil)
		if err != nil {
			fmt.Print(err)
		}
		return
	}

	var plannedPath []math.Pos
	var lastPathCost *float64

	var testPath [][]math.Pos
	var testCost []float64

	matchedStreetIDs := getStreetIDsFromCurrentPosition(event.Store(*m.StreetStore), mioState.Position)
	for _, streetID := range matchedStreetIDs {
		var streetState aggregator.StreetState
		{
			streetEvents := event.Store(*m.StreetStore).GetEvents()[streetID]
			streetState, err = aggregator.GetStreetState(streetEvents)
			if err != nil {
				fmt.Print(err)
			}
		}

		path := getStreetShortestPath(m.MapPaths, mioState.Position, buildingState.Pos, streetState, streetID)
		if lastPathCost == nil || *lastPathCost > path.cost {
			lastPathCost = &path.cost
			plannedPath = path.poses
		}

		testPath = append(testPath, path.poses)
		testCost = append(testCost, path.cost)
	}
	if plannedPath == nil {
		return
	}

	if mioState.Position == plannedPath[0] {
		plannedPath = plannedPath[1:]
	}

	// if duplicate with mio's old planned path, skip change planned poses
	isDiffPaths := len(mioState.PlannedPoses) != len(plannedPath)

	for i, pos := range mioState.PlannedPoses {
		if pos != plannedPath[i] {
			isDiffPaths = true
			break
		}
	}

	if !isDiffPaths {
		return
	}

	fmt.Printf("new path selected, cost: %f\n", *lastPathCost)
	for i, v := range testPath {
		fmt.Printf("path %v, cost: %f\n", v, testCost[i])
	}

	err = m.MioOperator.ChangePlannedPoses(m.EntityID, plannedPath)
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

func (m MioAutomaton) Move() {
	events := event.Store(*m.MioStore).GetEvents()[m.EntityID]
	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		fmt.Print(err)
	}

	if len(mioState.PlannedPoses) == 0 {
		return
	}

	pos := mioState.PlannedPoses[0]

	_, _, distanceSqrt := math.GetDistances(mioState.Position, pos)

	if distanceSqrt <= aggregator.MAX_WALK_DISTANT {
		nextPos := math.NewPos(math.GetNextStepXY(mioState.Position, 0, pos, 0, distanceSqrt, 180))
		err = m.MioOperator.Walk(m.EntityID, nextPos)
	} else if distanceSqrt <= aggregator.MAX_RUN_DISTANT {
		nextPos := math.NewPos(math.GetNextStepXY(mioState.Position, 0, pos, 0, distanceSqrt, 180))
		err = m.MioOperator.Run(m.EntityID, nextPos)
	} else {
		nextPos := math.NewPos(math.GetNextStepXY(mioState.Position, 0, pos, 0, aggregator.MAX_RUN_DISTANT, 180))
		err = m.MioOperator.Run(m.EntityID, nextPos)
	}
	if err != nil {
		fmt.Print(err)
	}
}

func (m MioAutomaton) HourlyExhaustion() {
	events := event.Store(*m.MioStore).GetEvents()[m.EntityID]
	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		fmt.Print(err)
	}

	err = m.MioOperator.Starve(mioState.ID, 5)
	if err != nil {
		fmt.Print(err)
	}
	err = m.MioOperator.Sweat(mioState.ID, 5)
	if err != nil {
		fmt.Print(err)
	}
}
