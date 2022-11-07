package carrier

type Data interface{}

type Effect string

var (
	// Combat
	InitEffect      Effect = "CARRIER_INITIALED"
	LaunchJetEffect Effect = "CARRIER_JET_LAUNCHED"
	HouseJetEffect  Effect = "CARRIER_JET_HOUSED"

	// Position
	MoveEffect Effect = "CARRIER_MOVED"
	IdleEffect Effect = "CARRIER_IDLED"
)

func (e Effect) IsValid() bool {
	return e == InitEffect ||
		e == LaunchJetEffect ||
		e == HouseJetEffect ||
		e == MoveEffect ||
		e == IdleEffect
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

func NewLaunchJetEvent(id string, jetID string) Event {
	return Event{
		ID:     id,
		Effect: LaunchJetEffect,
		Data:   jetID,
	}
}

func NewHouseJetEvent(id string, jetID string) Event {
	return Event{
		ID:     id,
		Effect: HouseJetEffect,
		Data:   jetID,
	}
}

func NewMoveEvent(id string, x, y float64) Event {
	return Event{
		ID:     id,
		Effect: MoveEffect,
		Data: PositionState{
			X: x,
			Y: y,
		},
	}
}

func NewIdleEvent(id string) Event {
	return Event{
		ID:     id,
		Effect: IdleEffect,
	}
}
