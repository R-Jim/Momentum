package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FindPath(t *testing.T) {
	posA := NewPos(0, 0)
	posB := NewPos(1, 1)
	posC := NewPos(2, 2)
	posD := NewPos(3, 3)

	paths := []Path{
		{
			Start: posA,
			End:   posB,
		},
		{
			Start: posB,
			End:   posC,
		},
		{
			Start: posA,
			End:   posD,
		},
		{
			Start: posC,
			End:   posD,
		},
	}

	graph := NewGraph(paths)

	expPaths := [][]Pos{
		{posB, posC},
		{posD, posC},
	}

	result := graph.FindPath(posA, posC)

	assertPaths(t, expPaths, result)

	expPaths = [][]Pos{
		{posA, posD},
		{posC, posD},
	}

	result = graph.FindPath(posB, posD)

	assertPaths(t, expPaths, result)

	posE := NewPos(4, 4)
	paths = append(paths,
		Path{
			Start: posE,
			End:   posD,
		},
		Path{
			Start: posE,
			End:   posA,
		},
		Path{
			Start: posE,
			End:   posC,
		},
	)

	expPaths = [][]Pos{
		{posA, posD, posC, posE},
		{posA, posD, posE},
		{posA, posE},
		{posC, posD, posA, posE},
		{posC, posD, posE},
		{posC, posE},
	}

	graph = NewGraph(paths)

	result = graph.FindPath(posB, posE)

	assertPaths(t, expPaths, result)
}

func assertPaths(t *testing.T, expPaths, result [][]Pos) {
	expPathKeys := map[int]bool{}
	require.Equal(t, len(expPaths), len(result))
	for _, path := range result {
		for i := 0; i < len(expPaths); i++ {
			if isEqualPoses(path, expPaths[i]) {
				expPathKeys[i] = true
				break
			}
		}
	}
	require.Equal(t, len(expPaths), len(expPathKeys))
}

func isEqualPoses(a, b []Pos) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
