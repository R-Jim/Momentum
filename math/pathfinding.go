package math

import (
	"fmt"
	"sync"
)

type Path struct {
	Start Pos
	End   Pos
	Cost  int
}

type Graph map[Pos][]Pos

func NewGraph(paths []Path) Graph {
	graph := map[Pos][]Pos{}
	for _, path := range paths {
		startValue, ok := graph[path.Start]
		if !ok {
			graph[path.Start] = []Pos{path.End}
		} else {
			graph[path.Start] = append(startValue, path.End)
		}

		endValue, ok := graph[path.End]
		if !ok {
			graph[path.End] = []Pos{path.Start}
		} else {
			graph[path.End] = append(endValue, path.Start)
		}
	}
	return graph
}

func (g Graph) FindPath(start, end Pos) [][]Pos {
	return g.findPath([]Pos{}, start, end)
}

func (g Graph) findPath(prevPoses []Pos, start, end Pos) [][]Pos {
	startList, ok := g[start]
	if !ok {
		fmt.Printf("no Pos found in graph, %v\n", start)
		return nil
	}

	prevPoses = append(prevPoses, start)

	result := [][]Pos{}
	var wg sync.WaitGroup

	wg.Add(len(startList))

	for _, pos := range startList {
		p := pos
		go func() {
			defer wg.Done()

			if end == p {
				result = append(result, []Pos{p})
				return
			}

			{
				isLoop := false
				for _, prevPos := range prevPoses {
					if p == prevPos {
						isLoop = true
						break
					}
				}
				if isLoop == true {
					return
				}
			}

			r := g.findPath(prevPoses, p, end)
			if len(r) == 0 {
				return
			}

			for _, v := range r {
				result = append(result, append([]Pos{p}, v...))
			}
		}()
	}

	wg.Wait()
	return result
}
