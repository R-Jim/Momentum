package weapon

type Data interface{}

type Effect string

var (
	InitEffect Effect = "ARTIFACT_INITIALED"
)

func (e Effect) IsValid() bool {
	return e == InitEffect
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
