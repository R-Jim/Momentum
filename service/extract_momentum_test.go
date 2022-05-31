package service

import (
	"testing"

	"github.com/R-jim/Momentum/entity"
	"github.com/R-jim/Momentum/valueobject"
)

func Test_extractAxisValue(t *testing.T) {
	type arg struct {
		givenValue  valueobject.AxisValue
		expectValue valueobject.AxisValue
	}

	tcs := map[string]arg{
		"axisValue = 2, minimumExtractValue = 1": {
			givenValue:  2,
			expectValue: 1,
		},
		"axisValue = -2, minimumExtractValue = 1": {
			givenValue:  -2,
			expectValue: -1,
		},
		"axisValue = 1, minimumExtractValue = 1": {
			givenValue:  1,
			expectValue: 1,
		},
		"axisValue = -1, minimumExtractValue = 1": {
			givenValue:  -1,
			expectValue: -1,
		},
		"axisValue = 0.5, minimumExtractValue = 1": {
			givenValue:  0.5,
			expectValue: 0,
		},
		"axisValue = -0.5, minimumExtractValue = 1": {
			givenValue:  -0.5,
			expectValue: 0,
		},
		"axisValue = 0, minimumExtractValue = 1": {
			givenValue:  0,
			expectValue: 0,
		},
	}

	for key, tc := range tcs {
		resultValue := extractAxisValue(tc.givenValue)
		if resultValue != tc.expectValue {
			t.Fatalf(`extractAxisValue("%v") want %v, result %v`, key, tc.expectValue, resultValue)
		}
	}
}

func Test_ExtractMomentum(t *testing.T) {
	type arg struct {
		givenCore               *entity.Core
		expectExtractedMomentum valueobject.Momentum
		expectFinalCoreMomentum valueobject.Momentum
	}

	tcs := map[string]arg{
		"Core.momentum[3, 2]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: 3,
					Y: 2,
				},
			},
			expectExtractedMomentum: valueobject.Momentum{
				X: 1,
				Y: 1,
			},
			expectFinalCoreMomentum: valueobject.Momentum{
				X: 2,
				Y: 1,
			},
		},
		"Core.momentum[2,-3]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: 2,
					Y: -3,
				},
			},
			expectExtractedMomentum: valueobject.Momentum{
				X: 1,
				Y: -1,
			},
			expectFinalCoreMomentum: valueobject.Momentum{
				X: 1,
				Y: -2,
			},
		},
		"Core.momentum[0,-3]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: 0,
					Y: -3,
				},
			},
			expectExtractedMomentum: valueobject.Momentum{
				X: 0,
				Y: -1,
			},
			expectFinalCoreMomentum: valueobject.Momentum{
				X: 0,
				Y: -2,
			},
		},
		"Core.momentum[3, 0]": {
			givenCore: &entity.Core{
				Momentum: valueobject.Momentum{
					X: 3,
					Y: 0,
				},
			},
			expectExtractedMomentum: valueobject.Momentum{
				X: 1,
				Y: 0,
			},
			expectFinalCoreMomentum: valueobject.Momentum{
				X: 2,
				Y: 0,
			},
		},
	}

	for key, tc := range tcs {
		result := ExtractMomentum(tc.givenCore)
		if result.X != tc.expectExtractedMomentum.X || result.Y != tc.expectExtractedMomentum.Y {
			t.Fatalf(`ExtractMomentum("%v") extract momentum want %v, result %v`, key, tc.expectExtractedMomentum.String(), result.String())
		}
		if tc.givenCore.Momentum.X != tc.expectFinalCoreMomentum.X || tc.givenCore.Momentum.Y != tc.expectFinalCoreMomentum.Y {
			t.Fatalf(`ExtractMomentum("%v") final core momentum want %v, result %v`, key, tc.expectFinalCoreMomentum.String(), tc.givenCore.Momentum.String())
		}
	}
}
