package automaton

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getNextStep(t *testing.T) {
	require.Equal(t, 5, getNextStep(1, 8, 4))
	require.Equal(t, 2, getNextStep(-2, 8, 4))
	require.Equal(t, 2, getNextStep(5, 0, 3))
	require.Equal(t, 5, getNextStep(5, 5, 3))
	require.Equal(t, -2, getNextStep(-5, 0, 3))
	require.Equal(t, -3, getNextStep(-5, -3, 3))
	require.Equal(t, -4, getNextStep(-1, -6, 3))
}
