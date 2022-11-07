package artifact

type Data interface{}

type Effect string

var (
	InitEffect Effect = "ARTIFACT_INITIALED"

	MoveEffect Effect = "ARTIFACT_MOVED"
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
