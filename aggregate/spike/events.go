package spike

type Data interface{}

type Effect string

var (
	InitEffect   Effect = "SPIKE_INITIALED"
	DamageEffect Effect = "SPIKE_DAMAGED"
	StrikeEffect Effect = "SPIKE_STROKE"

	MoveEffect Effect = "SPIKE_MOVED"
)

func (e Effect) IsValid() bool {
	return e == InitEffect ||
		e == DamageEffect ||
		e == StrikeEffect ||
		e == MoveEffect
}

type Event struct {
	ID     string
	Effect Effect
	Data   Data
}

func NewInitEvent(id string, artifactID string, health int) Event {
	state := State{
		Health: Health{
			Max:   health,
			Value: health,
		},
		ArtifactID: artifactID,
	}
	return Event{
		ID:     id,
		Effect: InitEffect,
		Data:   state,
	}
}

func NewDamageEvent(id string, damage int) Event {
	return Event{
		ID:     id,
		Effect: DamageEffect,
		Data:   damage,
	}
}

func NewStrikeEvent(id string, targetID string) Event {
	return Event{
		ID:     id,
		Effect: StrikeEffect,
		Data:   targetID,
	}
}

func NewMoveEffect(id string, x, y float64) Event {
	return Event{
		ID:     id,
		Effect: MoveEffect,
		Data: PositionState{
			X: x,
			Y: y,
		},
	}
}
