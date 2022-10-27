package fueltank

type Data interface{}

type Effect string

var (
	InitEffect    Effect = "FUEL_TANK_INITIALED"
	RefillEffect  Effect = "FUEL_TANK_REFILLED"
	ConsumeEffect Effect = "FUEL_TANK_CONSUMED"
)

func (e Effect) IsValid() bool {
	return e == InitEffect ||
		e == RefillEffect ||
		e == ConsumeEffect
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

func NewRefillEvent(id string, quantity int) Event {
	return Event{
		ID:     id,
		Effect: RefillEffect,
		Data:   quantity,
	}
}

func NewConsumeEvent(id string, quantity int) Event {
	return Event{
		ID:     id,
		Effect: ConsumeEffect,
		Data:   quantity,
	}
}
