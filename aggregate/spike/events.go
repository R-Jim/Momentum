package spike

type Data interface{}

type Effect string

var (
	InitEffect   Effect = "SPIKE_INITIALED"
	DamageEffect Effect = "SPIKE_DAMAGED"

	MoveEffect Effect = "SPIKE_MOVED"
)

func (e Effect) IsValid() bool {
	return e == InitEffect ||
		e == MoveEffect
}

type Event struct {
	ID     string
	Effect Effect
	Data   Data
}

func NewInitEvent(id string) Event {
	return Event{
		ID:     id,
		Effect: InitEffect,
	}
}

func NewDamageEvent(id string, damage int) Event {
	return Event{
		ID:     id,
		Effect: DamageEffect,
		Data:   damage,
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
