package automaton

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/store"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/operator"
	"github.com/google/uuid"
)

type WorkerAutomaton struct {
	EntityID uuid.UUID

	MapPaths []math.Path

	WorkerStore   *store.Store
	StreetStore   *store.Store
	BuildingStore *store.Store

	WorkerOperator operator.WorkerOperator
	StreetOperator operator.StreetOperator
}

func (w WorkerAutomaton) Automate() {
	w.PathFindingUpdate()
	w.Move()
}

func (w WorkerAutomaton) PathFindingUpdate() {
	events := (*w.WorkerStore).GetEvents()[w.EntityID]
	workerState, err := aggregator.GetWorkerState(events)
	if err != nil {
		fmt.Print(err)
	}

	if workerState.BuildingID == uuid.Nil {
		return
	}

	buildingEvents := (*w.BuildingStore).GetEvents()[workerState.BuildingID]
	buildingState, err := aggregator.GetBuildingState(buildingEvents)
	if err != nil {
		fmt.Print(err)
	}

	if workerState.Position == buildingState.Pos {
		err = w.WorkerOperator.ChangePlannedPoses(w.EntityID, nil)
		if err != nil {
			fmt.Print(err)
		}
		return
	}

	var plannedPath []math.Pos
	var lastPathCost *float64

	var testPath [][]math.Pos
	var testCost []float64

	matchedStreetIDs := getStreetIDsFromCurrentPosition(*w.StreetStore, workerState.Position)
	for _, streetID := range matchedStreetIDs {
		var streetState aggregator.StreetState
		{
			streetEvents := (*w.StreetStore).GetEvents()[streetID]
			streetState, err = aggregator.GetStreetState(streetEvents)
			if err != nil {
				fmt.Print(err)
			}
		}

		path := getStreetShortestPath(w.MapPaths, workerState.Position, buildingState.Pos, streetState, streetID)
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

	if workerState.Position == plannedPath[0] {
		plannedPath = plannedPath[1:]
	}

	// if duplicate with worker's old planned path, skip change planned poses
	isDiffPaths := len(workerState.PlannedPoses) != len(plannedPath)

	for i, pos := range workerState.PlannedPoses {
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

	err = w.WorkerOperator.ChangePlannedPoses(w.EntityID, plannedPath)
	if err != nil {
		fmt.Print(err)
	}
}

func (w WorkerAutomaton) Move() {
	events := (*w.WorkerStore).GetEvents()[w.EntityID]
	workerState, err := aggregator.GetWorkerState(events)
	if err != nil {
		fmt.Print(err)
	}

	if len(workerState.PlannedPoses) == 0 {
		return
	}

	pos := workerState.PlannedPoses[0]

	_, _, distanceSqrt := math.GetDistances(workerState.Position, pos)

	if distanceSqrt <= aggregator.MAX_WORKER_MOVE_DISTANT {
		nextPos := math.NewPos(math.GetNextStepXY(workerState.Position, 0, pos, 0, distanceSqrt, 180))
		err = w.WorkerOperator.Move(w.EntityID, nextPos)
	} else {
		nextPos := math.NewPos(math.GetNextStepXY(workerState.Position, 0, pos, 0, aggregator.MAX_WORKER_MOVE_DISTANT, 180))
		err = w.WorkerOperator.Move(w.EntityID, nextPos)
	}
	if err != nil {
		fmt.Print(err)
	}
}
