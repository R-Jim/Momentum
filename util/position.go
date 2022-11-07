package util

import "math"

func GetDistances(startX, startY, desX, desY float64) (distanceX, distanceY, distanceSqrt float64) {
	distanceX = math.Abs(desX - startX)
	distanceY = math.Abs(desY - startY)

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

func GetDegree(startX, startY, desX, desY float64) float64 {
	x := desX - startX
	y := desY - startY

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
