package automaton

import (
	"fmt"
	"sync"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/R-jim/Momentum/math"
	"github.com/google/uuid"
)

func getStreetShortestPath(paths []math.Path, entityPos, buildingPos math.Pos, streetState aggregator.StreetState, streetID uuid.UUID) *pathWithCost {
	mapGraph := math.NewGraph(paths)

	streetHeadAPos := streetState.HeadA
	streetHeadBPos := streetState.HeadB

	streetPathCost := 0.0
	for _, path := range paths {
		if (path.Start == streetHeadAPos && path.End == streetHeadBPos) || (path.Start == streetHeadBPos && path.End == streetHeadAPos) {
			streetPathCost = path.Cost
		}
	}

	_, _, distFromA := math.GetDistances(entityPos, streetHeadAPos)
	_, _, distFromB := math.GetDistances(entityPos, streetHeadBPos)
	totalDist := distFromA + distFromB

	var shortestPathFromStreetA, shortestPathFromStreetB *pathWithCost
	var shortestPathFromStreetACost, shortestPathFromStreetBCost float64
	{
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			shortestPathFromStreetA = findShortestPath(mapGraph, paths, streetHeadAPos, buildingPos)
		}()
		go func() {
			defer wg.Done()
			shortestPathFromStreetB = findShortestPath(mapGraph, paths, streetHeadBPos, buildingPos)
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
		if shortestPath == nil || *&shortestPath.cost > path.cost { // DO NOT REMOVE *&
			shortestPath = &path
		}
	}

	return shortestPath
}

func getStreetIDsFromCurrentPosition(streetStore event.Store, entityPos math.Pos) []uuid.UUID {
	matchedStreetIDs := []uuid.UUID{}

	for streetID, streetEvents := range streetStore.GetEvents() {
		streetState, err := aggregator.GetStreetState(streetEvents)
		if err != nil {
			fmt.Print(err)
		}
		if math.IsBetweenAAndB(entityPos, streetState.HeadA, streetState.HeadB, 0.2) {
			matchedStreetIDs = append(matchedStreetIDs, streetID)
		}
	}

	return matchedStreetIDs
}
