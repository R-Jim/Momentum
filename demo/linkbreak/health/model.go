package health

import "github.com/google/uuid"

type Health struct {
	ownerID   uuid.UUID
	baseValue int
}

func NewHealth(ownerID uuid.UUID, baseValue int) Health {
	return Health{
		ownerID, baseValue,
	}
}

func (h Health) OwnerID() uuid.UUID {
	return h.ownerID
}

func (h Health) BaseValue() int {
	return h.baseValue
}
