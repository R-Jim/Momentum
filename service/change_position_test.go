package service

import (
	"testing"

	"github.com/R-jim/Momentum/entity"
	"github.com/R-jim/Momentum/valueobject"
)

func Test_ChangePosition(t *testing.T) {
	type arg struct {
		givenCore           *entity.Core
		expectFinalPosition valueobject.Position
	}

	tcs := map[string]arg{
		"Core.momentum[3, 2], position[0, 0]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: 3,
					Y: 2,
				},
			},
			expectFinalPosition: valueobject.Position{
				X: 1,
				Y: 1,
			},
		},
		"Core.momentum[3, 2], position[1, 0]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: 3,
					Y: 2,
				},
				Position: valueobject.Position{
					X: 1,
				},
			},
			expectFinalPosition: valueobject.Position{
				X: 2,
				Y: 1,
			},
		},
		"Core.momentum[3, 2], position[1, -1]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: 3,
					Y: 2,
				},
				Position: valueobject.Position{
					X: 1,
					Y: -1,
				},
			},
			expectFinalPosition: valueobject.Position{
				X: 2,
				Y: 0,
			},
		},
		"Core.momentum[-3, 0], position[1, -1]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: -3,
					Y: 0,
				},
				Position: valueobject.Position{
					X: 1,
					Y: -1,
				},
			},
			expectFinalPosition: valueobject.Position{
				X: 0,
				Y: -1,
			},
		},
		"Core.momentum[0, -3], position[1, -1]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: 0,
					Y: -3,
				},
				Position: valueobject.Position{
					X: 1,
					Y: -1,
				},
			},
			expectFinalPosition: valueobject.Position{
				X: 1,
				Y: -2,
			},
		},
	}

	for key, tc := range tcs {
		ChangePosition(tc.givenCore)
		if tc.givenCore.Position.X != tc.expectFinalPosition.X || tc.givenCore.Position.Y != tc.expectFinalPosition.Y {
			t.Fatalf(`ChangePosition("%v") final position want %v, result %v`, key, tc.expectFinalPosition.String(), tc.givenCore.Position.String())
		}
	}
}
