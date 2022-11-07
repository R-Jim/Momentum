package automaton

import (
	"math"

	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/util"
)

func getNextStepXY(positionState jet.PositionState, desX, desY, desRad, step float64) (x, y float64) {
	maxTurnDegree := float64(60)

	radius := math.RoundToEven((step / 2) / math.Sin(maxTurnDegree/2*math.Pi/180))
	maxRadius := radius - step
	if maxRadius < desRad {
		maxRadius = desRad
	}
	_, _, distance := util.GetDistances(positionState.X, positionState.Y, desX, desY)

	positions := getPositions(positionState, maxTurnDegree, step)
	if distance >= maxRadius {
		var nearestPos pos
		var nearestDistance float64
		for _, position := range positions {
			_, _, dis := util.GetDistances(position.x, position.y, desX, desY)
			if nearestDistance == 0 {
				nearestPos = position
				nearestDistance = dis
			} else if dis < nearestDistance && dis >= maxRadius {
				nearestPos = position
				nearestDistance = dis
			}
		}

		if nearestDistance >= maxRadius {
			return nearestPos.x, nearestPos.y
		}
	}

	var farthestPos pos
	var farthestDistance float64
	for _, position := range positions {
		_, _, dis := util.GetDistances(position.x, position.y, desX, desY)
		if dis >= maxRadius && dis < farthestDistance {
			farthestPos = position
			farthestDistance = dis
			continue
		}
		if dis > farthestDistance {
			farthestPos = position
			farthestDistance = dis
		}
	}
	return farthestPos.x, farthestPos.y

}

type pos struct {
	x float64
	y float64
}

func getPositions(positionState jet.PositionState, maxTurnDegree, step float64) []pos {
	degreeStep := float64(60)

	result := make([]pos, 0)
	for i := maxTurnDegree; i >= float64(0); i -= degreeStep {
		if i == 0 {
			result = append(result, pos{
				x: math.RoundToEven((positionState.X+step*util.CosDegree(positionState.HeadDegree))*100) / 100,
				y: math.RoundToEven((positionState.Y+step*util.SinDegree(positionState.HeadDegree))*100) / 100,
			})
		} else {
			result = append(result, pos{
				x: math.RoundToEven((positionState.X+step*util.CosDegree(positionState.HeadDegree+i))*100) / 100,
				y: math.RoundToEven((positionState.Y+step*util.SinDegree(positionState.HeadDegree+i))*100) / 100,
			})
			result = append(result, pos{
				x: math.RoundToEven((positionState.X+step*util.CosDegree(positionState.HeadDegree+i*-1))*100) / 100,
				y: math.RoundToEven((positionState.Y+step*util.SinDegree(positionState.HeadDegree+i*-1))*100) / 100,
			})
		}
	}
	return result
}
