package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AddVector(t *testing.T) {
	v1 := Vector{
		AxisX: 2,
		AxisY: 1,
	}
	v2 := Vector{
		AxisX: 2,
		AxisY: 4,
	}

	result := v1.Add(v2)
	require.Equal(t, float64(4), result.AxisX)
	require.Equal(t, float64(5), result.AxisY)
	require.Equal(t, float64(0), result.AxisZ)
}

func Test_PartOfVector(t *testing.T) {
	v1 := Vector{
		AxisX: 3,
	}.PartOf(1)

	require.Equal(t, float64(1), v1.AxisX)
	require.Equal(t, float64(0), v1.AxisY)
	require.Equal(t, float64(0), v1.AxisZ)

	v2 := Vector{
		AxisX: 3,
		AxisY: 5,
	}.PartOf(1)

	require.Equal(t, float64(0.375), v2.AxisX)
	require.Equal(t, float64(0.625), v2.AxisY)
	require.Equal(t, float64(0), v2.AxisZ)

	v3 := Vector{
		AxisX: 1,
		AxisY: 2,
	}.PartOf(5)

	require.Equal(t, float64(1), v3.AxisX)
	require.Equal(t, float64(2), v3.AxisY)
	require.Equal(t, float64(0), v3.AxisZ)
}
