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

type mioAutomaton struct {
	entityID uuid.UUID

	mapPaths []math.Path
	mapGraph math.Graph

	mioStore      *store.Store
	streetStore   *store.Store
	buildingStore *store.Store

	mioOperator    operator.MioOperator
	streetOperator operator.StreetOperator

	// For performance optimization
	prevSelectedBuilding uuid.UUID
}

type shortestPath struct {
	poses []math.Pos
	cost  int
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

func (m *mioAutomaton) MioMoodBehavior() {
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

	if m.prevSelectedBuilding == selectedBuildingID {
		return
	}

	m.prevSelectedBuilding = selectedBuildingID

	err = m.mioOperator.SelectBuilding(mioState.ID, selectedBuildingID)
	if err != nil {
		fmt.Print(err)
	}

	m.mioBuildingPathFinding()
}

func (m mioAutomaton) mioBuildingPathFinding() {
	if m.prevSelectedBuilding.String() == uuid.Nil.String() {
		return
	}

	var plannedPath []math.Pos

	var mioPos, streetHeadAPos, streetHeadBPos, buildingPos math.Pos

	{
		buildingEvents := (*m.buildingStore).GetEvents()[m.prevSelectedBuilding]
		buildingState, err := aggregator.GetBuildingState(buildingEvents)
		if err != nil {
			fmt.Print(err)
		}

		events := (*m.mioStore).GetEvents()[m.entityID]
		mioState, err := aggregator.GetMioState(events)
		if err != nil {
			fmt.Print(err)
		}

		streetEvents := (*m.streetStore).GetEvents()[mioState.StreetID]
		streetState, err := aggregator.GetStreetState(streetEvents)
		if err != nil {
			fmt.Print(err)
		}

		mioPos = mioState.Position
		streetHeadAPos = streetState.HeadA
		streetHeadBPos = streetState.HeadB
		buildingPos = buildingState.Pos
	}

	streetPathCost := 0
	for _, path := range m.mapPaths {
		if (path.Start == streetHeadAPos && path.End == streetHeadBPos) || (path.Start == streetHeadBPos && path.End == streetHeadAPos) {
			streetPathCost = path.Cost
		}
	}

	_, _, distFromA := math.GetDistances(mioPos, streetHeadAPos)
	_, _, distFromB := math.GetDistances(mioPos, streetHeadBPos)
	totalDist := distFromA + distFromB

	var shortestPathFromStreetA, shortestPathFromStreetB shortestPath
	{
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			shortestPathFromStreetA = findShortestPath(m.mapGraph, m.mapPaths, streetHeadAPos, buildingPos)
		}()
		go func() {
			defer wg.Done()
			shortestPathFromStreetB = findShortestPath(m.mapGraph, m.mapPaths, streetHeadBPos, buildingPos)
		}()

		wg.Wait()
	}

	{
		shortestPathFromStreetACost := float64(streetPathCost)*distFromA/totalDist + float64(shortestPathFromStreetA.cost)
		shortestPathFromStreetBCost := float64(streetPathCost)*distFromB/totalDist + float64(shortestPathFromStreetB.cost)
		if shortestPathFromStreetACost <= shortestPathFromStreetBCost {
			plannedPath = append([]math.Pos{streetHeadAPos}, shortestPathFromStreetA.poses...)
		} else {
			plannedPath = append([]math.Pos{streetHeadBPos}, shortestPathFromStreetB.poses...)
		}
	}

	err := m.mioOperator.ChangePlannedPoses(m.entityID, plannedPath)
	if err != nil {
		fmt.Print(err)
	}
}

func findShortestPath(mapGraph math.Graph, mapPaths []math.Path, start, end math.Pos) shortestPath {
	var sp *shortestPath

	paths := mapGraph.FindPath(start, end)

	var wg sync.WaitGroup
	wg.Add(len(paths))
	for _, path := range paths {
		p := path
		go func() {
			defer wg.Done()

			lastPos := start
			cost := 0

			for _, pos := range p {
				if sp != nil && cost > sp.cost {
					return
				}

				for _, mp := range mapPaths {
					if (mp.Start == lastPos && mp.End == pos) || (mp.Start == pos && mp.End == lastPos) {
						cost += mp.Cost
						break
					}
				}

				lastPos = pos
			}

			if sp != nil && cost > sp.cost {
				return
			}

			if sp == nil {
				sp = &shortestPath{}
			}

			sp.poses = p
			sp.cost = cost
		}()
	}
	wg.Wait()

	if sp == nil {
		return shortestPath{}
	}

	return *sp
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
