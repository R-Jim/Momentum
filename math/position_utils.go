package math

import (
	"math"

	"github.com/R-jim/Momentum/util"
)

type pos struct {
	x float64
	y float64
}

func GetNextStepXY(x, y, pivotDegree, desX, desY, desRad, step, maxTurnDegree float64) (resultX, resultY float64) {
	radius := math.RoundToEven((step / 2) / math.Sin(maxTurnDegree/2*math.Pi/180))
	maxRadius := radius - step
	if maxRadius < desRad {
		maxRadius = desRad
	}
	_, _, distance := util.GetDistances(x, y, desX, desY)

	positions := GetPositions(x, y, pivotDegree, maxTurnDegree, float64(60), step)
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

func GetPositions(x, y, pivotDegree, maxTurnDegree, degreeStep, step float64) []pos {
	result := make([]pos, 0)
	for i := maxTurnDegree; i >= float64(0); i -= degreeStep {
		if i == 0 {
			result = append(result, pos{
				x: math.RoundToEven((x+step*util.CosDegree(pivotDegree))*100) / 100,
				y: math.RoundToEven((y+step*util.SinDegree(pivotDegree))*100) / 100,
			})
		} else {
			result = append(result, pos{
				x: math.RoundToEven((x+step*util.CosDegree(pivotDegree+i))*100) / 100,
				y: math.RoundToEven((y+step*util.SinDegree(pivotDegree+i))*100) / 100,
			})
			result = append(result, pos{
				x: math.RoundToEven((x+step*util.CosDegree(pivotDegree+i*-1))*100) / 100,
				y: math.RoundToEven((y+step*util.SinDegree(pivotDegree+i*-1))*100) / 100,
			})
		}
	}
	return result
}
