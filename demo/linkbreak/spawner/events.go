package spawner

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/template/event"
)

const (
	InitEffect      event.Effect = "INIT_EFFECT"
	CountDownEffect event.Effect = "COUNT_DOWN_EFFECT"
	SpawnEffect     event.Effect = "SPAWN_EFFECT"
)

func NewInitEvent(s *event.Store, id uuid.UUID, spawner Spawner) event.Event {
	return s.NewEvent(id, InitEffect, spawner)
}

func NewCountDownEvent(s *event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, CountDownEffect, nil)
}

func NewSpawnEvent(s *event.Store, id uuid.UUID) event.Event {
	return s.NewEvent(id, SpawnEffect, nil)
}
