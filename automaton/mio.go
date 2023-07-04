package automaton

import (
	"fmt"
	"sync"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
)

type MioAutomaton struct {
	EntityID uuid.UUID

	MapPaths []math.Path
	MapGraph math.Graph

	MioStore      *store.Store
	StreetStore   *store.Store
	BuildingStore *store.Store

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

func (m MioAutomaton) getStreetIDsFromCurrentPosition() []uuid.UUID {
	events := (*m.MioStore).GetEvents()[m.EntityID]
	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		fmt.Print(err)
	}

	matchedStreetIDs := []uuid.UUID{}

	for streetID, streetEvents := range (*m.StreetStore).GetEvents() {
		streetState, err := aggregator.GetStreetState(streetEvents)
		if err != nil {
			fmt.Print(err)
		}
		if math.IsBetweenAAndB(mioState.Position, streetState.HeadA, streetState.HeadB, 0.2) {
			matchedStreetIDs = append(matchedStreetIDs, streetID)
		}
	}

	return matchedStreetIDs
}

func (m *MioAutomaton) MioMoodBehavior() {
	events := (*m.MioStore).GetEvents()[m.EntityID]
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

	for buildingID, buildingEvents := range (*m.BuildingStore).GetEvents() {
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
	events := (*m.MioStore).GetEvents()[m.EntityID]
	mioState, err := aggregator.GetMioState(events)
	if err != nil {
		fmt.Print(err)
	}

	if mioState.SelectedBuildingID == uuid.Nil {
		return
	}

	buildingEvents := (*m.BuildingStore).GetEvents()[mioState.SelectedBuildingID]
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

	getStreetShortestPath := func(streetID uuid.UUID) *pathWithCost {
		var mioPos, streetHeadAPos, streetHeadBPos, buildingPos math.Pos

		{
			streetEvents := (*m.StreetStore).GetEvents()[streetID]
			streetState, err := aggregator.GetStreetState(streetEvents)
			if err != nil {
				fmt.Print(err)
			}

			mioPos = mioState.Position
			streetHeadAPos = streetState.HeadA
			streetHeadBPos = streetState.HeadB
			buildingPos = buildingState.Pos
		}

		streetPathCost := 0.0
		for _, path := range m.MapPaths {
			if (path.Start == streetHeadAPos && path.End == streetHeadBPos) || (path.Start == streetHeadBPos && path.End == streetHeadAPos) {
				streetPathCost = path.Cost
			}
		}

		_, _, distFromA := math.GetDistances(mioPos, streetHeadAPos)
		_, _, distFromB := math.GetDistances(mioPos, streetHeadBPos)
		totalDist := distFromA + distFromB

		var shortestPathFromStreetA, shortestPathFromStreetB *pathWithCost
		var shortestPathFromStreetACost, shortestPathFromStreetBCost float64
		{
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				shortestPathFromStreetA = findShortestPath(m.MapGraph, m.MapPaths, streetHeadAPos, buildingPos)
			}()
			go func() {
				defer wg.Done()
				shortestPathFromStreetB = findShortestPath(m.MapGraph, m.MapPaths, streetHeadBPos, buildingPos)
			}()

			wg.Wait()
			if shortestPathFromStreetA != nil {
				shortestPathFromStreetACost = float64(streetPathCost)*distFromA/totalDist + float64(shortestPathFromStreetA.cost)
			}
			if shortestPathFromStreetB != nil {
				shortestPathFromStreetBCost = float64(streetPathCost)*distFromB/totalDist + float64(shortestPathFromStreetB.cost)
			}
			if shortestPathFromStreetA == nil && shortestPathFromStreetB == nil {
				return nil
			} else if shortestPathFromStreetA == nil {
				return &pathWithCost{
					poses: append([]math.Pos{streetHeadBPos}, shortestPathFromStreetB.poses...),
					cost:  shortestPathFromStreetBCost,
				}
			} else if shortestPathFromStreetB == nil {
				return &pathWithCost{
					poses: append([]math.Pos{streetHeadAPos}, shortestPathFromStreetA.poses...),
					cost:  shortestPathFromStreetACost,
				}
			}
		}
		{
			if shortestPathFromStreetACost <= shortestPathFromStreetBCost && len(shortestPathFromStreetA.poses) <= len(shortestPathFromStreetB.poses) {
				return &pathWithCost{
					poses: append([]math.Pos{streetHeadAPos}, shortestPathFromStreetA.poses...),
					cost:  shortestPathFromStreetACost,
				}
			} else {
				return &pathWithCost{
					poses: append([]math.Pos{streetHeadBPos}, shortestPathFromStreetB.poses...),
					cost:  shortestPathFromStreetBCost,
				}
			}
		}
	}

	var plannedPath []math.Pos
	var lastPathCost *float64

	var testPath [][]math.Pos
	var testCost []float64

	matchedStreetIDs := m.getStreetIDsFromCurrentPosition()
	for _, streetID := range matchedStreetIDs {
		path := getStreetShortestPath(streetID)
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

func findShortestPath(mapGraph math.Graph, mapPaths []math.Path, start, end math.Pos) *pathWithCost {
	if start == end {
		return &pathWithCost{}
	}

	pathWithCosts := []pathWithCost{}

	paths := mapGraph.FindPath(start, end)

	// TODO: might use goroutine to optimize, make sure pathWithCosts doesn't overlap when append
	for _, path := range paths {
		lastPos := start
		cost := 0.0

		for _, pos := range path {
			if len(pathWithCosts) != 0 && cost > pathWithCosts[0].cost {
				continue
			}

			for _, mp := range mapPaths {
				if (mp.Start == lastPos && mp.End == pos) || (mp.Start == pos && mp.End == lastPos) {
					cost += mp.Cost
					break
				}
			}

			lastPos = pos
		}

		if len(pathWithCosts) != 0 && cost > pathWithCosts[0].cost {
			continue
		}
		pathWithCosts = append(pathWithCosts, pathWithCost{
			poses: path,
			cost:  cost,
		})
	}

	var shortestPath *pathWithCost
	for _, path := range pathWithCosts {
		if shortestPath == nil || *&shortestPath.cost > path.cost {
			shortestPath = &path
		}
	}

	return shortestPath
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
	events := (*m.MioStore).GetEvents()[m.EntityID]
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
	events := (*m.MioStore).GetEvents()[m.EntityID]
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
