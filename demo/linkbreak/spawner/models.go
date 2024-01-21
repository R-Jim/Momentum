package spawner

import "github.com/R-jim/Momentum/math"

type Spawner struct {
	spawnTypeCode int
	faction       int
	counter       int
	position      math.Pos
}

func (s Spawner) SpawnTypeData() interface{} {
	return getSpawnTypeData(SpawnTypeCode(s.spawnTypeCode))
}

func (s Spawner) Faction() int {
	return s.faction
}

func (s Spawner) Counter() int {
	return s.counter
}

func (s Spawner) Position() math.Pos {
	return s.position
}
