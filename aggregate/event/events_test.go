package event

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseData(t *testing.T) {
	e := Event{
		Data: "TEST_STRING_TYPE",
	}

	d, err := ParseData[string](e)
	require.NoError(t, err)
	require.Equal(t, "TEST_STRING_TYPE", d)

	_, err = ParseData[int](e)
	require.EqualError(t, err, "failed to parse data for effect: int")
}
