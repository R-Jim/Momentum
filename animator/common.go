package animator

import (
	"github.com/R-jim/Momentum/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	HitEffectImage   *ebiten.Image
	DeathEffectImage *ebiten.Image
)

func initCommonImages() {
	HitEffectImage = getAssetImage(asset.Hit_png)
	DeathEffectImage = getAssetImage(asset.Death_png)
}
