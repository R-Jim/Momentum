package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsBetweenAAndB(t *testing.T) {
	a := Pos{
		x: 1,
		y: 1,
	}

	b := Pos{
		x: 0,
		y: 4,
	}

	require.True(t, IsBetweenAAndB(Pos{x: 0.47, y: 2.6}, a, b, 0))

	// on same line as A, B
	require.True(t, IsBetweenAAndB(Pos{x: -0.3, y: 4.9}, a, b, 1))
	require.True(t, IsBetweenAAndB(Pos{x: 1.31, y: 0.07}, a, b, 1))

	require.False(t, IsBetweenAAndB(Pos{x: -0.2, y: 4.59}, a, b, 0))
	require.False(t, IsBetweenAAndB(Pos{x: 1.21, y: 0.36}, a, b, 0))
	// on 2 side of AB
	require.True(t, IsBetweenAAndB(Pos{x: 1.25, y: 2.75}, a, b, 1))
	require.True(t, IsBetweenAAndB(Pos{x: 0, y: 2.25}, a, b, 1))

	require.False(t, IsBetweenAAndB(Pos{x: 1.35, y: 3.46}, a, b, 0))
	require.False(t, IsBetweenAAndB(Pos{x: -1.02, y: 2.67}, a, b, 0))
}
