package automaton

import (
	"math"

	"github.com/R-jim/Momentum/aggregate/jet"
)

func getNextStep(start, destination, step float64) float64 {
	if start > destination {
		tmp := start
		start = destination - 1
		destination = tmp

	}

	if start+step > destination {
		return destination
	}
	return start + step
}

func getDistances(startX, startY, desX, desY float64) (distanceX, distanceY, distanceSqrt float64) {
	distanceX = math.Abs(desX - startX)
	distanceY = math.Abs(desY - startY)

	if distanceX == 0 && distanceY == 0 {
		return 0, 0, 0
	} else if distanceX == 0 {
		distanceSqrt = distanceY
	} else if distanceY == 0 {
		distanceSqrt = distanceX
	} else {
		distanceSqrt = math.RoundToEven(math.Sqrt(math.Pow(distanceX, 2)+math.Pow(distanceX, 2))*100) / 100
	}

	return distanceX, distanceY, distanceSqrt
}

func getSteps(startX, startY, desX, desY float64) (stepX, stepY float64) {
	distanceX, distanceY, distanceSqrt := getDistances(startX, startY, desX, desY)
	if distanceX == 0 && distanceY == 0 && distanceSqrt == 0 {
		return 0.5, 0.5
	}

	stepX = math.RoundToEven(distanceX/distanceSqrt*100) / 100
	stepY = math.RoundToEven(distanceY/distanceSqrt*100) / 100
	return stepX, stepY
}

func getNextStepXY(positionState jet.PositionState, desX, desY, desRad, step float64) (x, y float64) {
	stepX, stepY := getSteps(positionState.X, positionState.Y, desX, desY)
	stepX *= step
	stepY *= step

	x = getNextStep(positionState.X, desX, stepX)
	y = getNextStep(positionState.Y, desY, stepY)

	_, _, distanceSqrt := getDistances(x, y, desX, desY)

	if distanceSqrt < desRad {
		var farthestPos pos
		var farthestDistance float64
		for _, position := range getPositions(positionState, step) {
			_, _, dis := getDistances(position.x, position.y, desX, desY)
			if dis-step >= desRad {
				return position.x, position.y
			}
			if dis > farthestDistance {
				farthestPos = position
				farthestDistance = dis
			}
		}
		return farthestPos.x, farthestPos.y
	}

	return x, y
}

type pos struct {
	x float64
	y float64
}

// Assume [0,0] POI varies in 0, 45 degree
var (
	varies = []float64{0, 60, -60}
)

func getPositions(positionState jet.PositionState, step float64) []pos {
	result := make([]pos, 0, len(varies))
	for _, v := range varies {
		d := math.Atan(positionState.HeadPivotY/positionState.HeadPivotX) * 180 / math.Pi

		result = append(result, pos{
			x: positionState.X + step*math.Cos(d+v),
			y: positionState.Y + step*math.Sin(d+v),
		})
	}
	return result
}
