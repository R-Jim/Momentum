package link

import "github.com/google/uuid"

type Link struct {
	source uuid.UUID
	target uuid.UUID
}

func (l Link) Source() uuid.UUID {
	return l.source
}

func (l Link) Target() uuid.UUID {
	return l.target
}
