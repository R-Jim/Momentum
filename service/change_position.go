package service

import "github.com/R-jim/Momentum/entity"

func ChangePosition(core *entity.Core) {
	extractMomentum := ExtractMomentum(core)
	core.Position.X += extractMomentum.X
	core.Position.Y += extractMomentum.Y
}
