package runner

import (
	"github.com/google/uuid"
)

type Runner struct {
	id         uuid.UUID
	faction    int
	healthID   uuid.UUID
	positionID uuid.UUID
}

func (r Runner) ID() uuid.UUID {
	return r.id
}

func (r Runner) Faction() int {
	return r.faction
}

func (r Runner) PositionID() uuid.UUID {
	return r.positionID
}

func (r Runner) HealthID() uuid.UUID {
	return r.healthID
}
