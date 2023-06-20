package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsBetweenAAndB(t *testing.T) {
	a := Pos{
		X: 1,
		Y: 1,
	}

	b := Pos{
		X: 0,
		Y: 4,
	}

	require.True(t, IsBetweenAAndB(Pos{X: 0.47, Y: 2.6}, a, b, 0))

	// on same line as A, B
	require.True(t, IsBetweenAAndB(Pos{X: -0.3, Y: 4.9}, a, b, 1))
	require.True(t, IsBetweenAAndB(Pos{X: 1.31, Y: 0.07}, a, b, 1))

	require.False(t, IsBetweenAAndB(Pos{X: -0.2, Y: 4.59}, a, b, 0))
	require.False(t, IsBetweenAAndB(Pos{X: 1.21, Y: 0.36}, a, b, 0))
	// on 2 side of AB
	require.True(t, IsBetweenAAndB(Pos{X: 1.25, Y: 2.75}, a, b, 1))
	require.True(t, IsBetweenAAndB(Pos{X: 0, Y: 2.25}, a, b, 1))

	require.False(t, IsBetweenAAndB(Pos{X: 1.35, Y: 3.46}, a, b, 0))
	require.False(t, IsBetweenAAndB(Pos{X: -1.02, Y: 2.67}, a, b, 0))
}
