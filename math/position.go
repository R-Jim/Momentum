package math

import (
	"math"
)

type Pos struct {
	X float64
	Y float64
}

func NewPos(x, y float64) Pos {
	return Pos{x, y}
}

func GetNextStepXY(posStart Pos, pivotDegree float64, posEnd Pos, desRad, step, maxTurnDegree float64) (resultX, resultY float64) {
	radius := math.RoundToEven((step / 2) / math.Sin(maxTurnDegree/2*math.Pi/180))
	maxRadius := radius - step
	if maxRadius < desRad {
		maxRadius = desRad
	}
	_, _, distance := GetDistances(posStart, posEnd)

	positions := GetPositions(posStart, pivotDegree, maxTurnDegree, float64(60), step)
	if distance >= maxRadius {
		var nearestPos Pos
		var nearestDistance float64
		for _, position := range positions {
			_, _, dis := GetDistances(position, posEnd)
			if nearestDistance == 0 {
				nearestPos = position
				nearestDistance = dis
			} else if dis < nearestDistance && dis >= maxRadius {
				nearestPos = position
				nearestDistance = dis
			}
		}

		if nearestDistance >= maxRadius {
			return nearestPos.X, nearestPos.Y
		}
	}

	var farthestPos Pos
	var farthestDistance float64
	for _, position := range positions {
		_, _, dis := GetDistances(position, posEnd)
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
	return farthestPos.X, farthestPos.Y
}

func GetPositions(pos Pos, pivotDegree, maxTurnDegree, degreeStep, step float64) []Pos {
	result := make([]Pos, 0)
	for i := maxTurnDegree; i >= float64(0); i -= degreeStep {
		if i == 0 {
			result = append(result, Pos{
				X: math.RoundToEven((pos.X+step*CosDegree(pivotDegree))*100) / 100,
				Y: math.RoundToEven((pos.Y+step*SinDegree(pivotDegree))*100) / 100,
			})
		} else {
			result = append(result, Pos{
				X: math.RoundToEven((pos.X+step*CosDegree(pivotDegree+i))*100) / 100,
				Y: math.RoundToEven((pos.Y+step*SinDegree(pivotDegree+i))*100) / 100,
			})
			result = append(result, Pos{
				X: math.RoundToEven((pos.X+step*CosDegree(pivotDegree+i*-1))*100) / 100,
				Y: math.RoundToEven((pos.Y+step*SinDegree(pivotDegree+i*-1))*100) / 100,
			})
		}
	}
	return result
}

func GetDistances(posStart Pos, posEnd Pos) (distanceX, distanceY, distanceSqrt float64) {
	distanceX = math.Abs(posEnd.X - posStart.X)
	distanceY = math.Abs(posEnd.Y - posStart.Y)

	if distanceX == 0 && distanceY == 0 {
		return 0, 0, 0
	} else if distanceX == 0 {
		distanceSqrt = distanceY
	} else if distanceY == 0 {
		distanceSqrt = distanceX
	} else {
		powX := distanceX * distanceX
		powY := distanceY * distanceY
		a := math.Sqrt(powX + powY)
		distanceSqrt = math.RoundToEven(a*100) / 100
	}

	return distanceX, distanceY, distanceSqrt
}

func (p Pos) GetDegree(end Pos) float64 {
	x := end.X - p.X
	y := end.Y - p.Y

	deg := math.RoundToEven(math.Atan(y/x) * 180 / math.Pi)
	if x < 0 && y < 0 {
		return deg - 180
	} else if x < 0 {
		return deg + 180
	}
	return deg
}

func SinDegree(d float64) float64 {
	return math.Sin(d * math.Pi / 180)
}

func CosDegree(d float64) float64 {
	return math.Cos(d * math.Pi / 180)
}

// buffer represent zone created by connecting A,B's circle with r = buffer
func IsBetweenAAndB(currentPos, a, b Pos, buffer float64) bool {
	_, _, distAB := GetDistances(a, b)
	_, _, distFromA := GetDistances(currentPos, a)
	_, _, distFromB := GetDistances(currentPos, b)

	fixedBuffer := 0.02

	return (distFromA+distFromB-distAB)-((buffer+fixedBuffer)*2) <= 0
}

func (p Pos) IsEqualRound(input Pos) bool {
	xDiff := p.X - input.X
	yDiff := p.Y - input.Y

	return xDiff+yDiff > -1 && xDiff+yDiff < 1
}
