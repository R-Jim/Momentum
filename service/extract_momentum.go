package service

import (
	"github.com/R-jim/Momentum/entity"
	"github.com/R-jim/Momentum/valueobject"
)

const minimumExtractValue valueobject.AxisValue = 1.0

func ExtractMomentum(core *entity.Core) valueobject.Momentum {
	extractedX := extractAxisValue(core.Momentum.X)
	extractedY := extractAxisValue(core.Momentum.Y)

	core.Momentum.X -= extractedX
	core.Momentum.Y -= extractedY

	return valueobject.Momentum{
		X: extractedX,
		Y: extractedY,
	}
}

func extractAxisValue(axisValue valueobject.AxisValue) valueobject.AxisValue {
	if axisValue > 0 && axisValue-minimumExtractValue >= 0 {
		return minimumExtractValue
	} else if axisValue < 0 && axisValue+minimumExtractValue <= 0 {
		return minimumExtractValue * -1
	}
	return 0
}
