package jet

type Data interface{}

type Effect string

var (
	// Combat
	InitEffect         Effect = "JET_INITIALED"
	AttackEffect       Effect = "JET_TARGET_ACQUIRED"
	CancelAttackEffect Effect = "JET_TARGET_RELEASED"
	EngageEffect       Effect = "JET_TARGET_ENGAGED"
	DisengageEffect    Effect = "JET_TARGET_DISENGAGED"
	TakeOffEffect      Effect = "JET_TOOK_OFF"
	LandingEffect      Effect = "JET_LANDED"

	// Position
	FlyEffect Effect = "JET_FLEW"

	// Inventory
	FuelTankChangedEffect Effect = "JET_FUEL_TANK_CHANGED"
)

func (e Effect) IsValid() bool {
	return e == InitEffect ||
		e == AttackEffect ||
		e == CancelAttackEffect ||
		e == EngageEffect ||
		e == DisengageEffect ||
		e == TakeOffEffect ||
		e == LandingEffect ||
		e == FlyEffect ||
		e == FuelTankChangedEffect
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

func NewAttackEvent(id string, targetID string) Event {
	return Event{
		ID:     id,
		Effect: AttackEffect,
		Data:   targetID,
	}
}

func NewCancelAttackEvent(id string) Event {
	return Event{
		ID:     id,
		Effect: CancelAttackEffect,
	}
}

func NewEngageEvent(id string) Event {
	return Event{
		ID:     id,
		Effect: EngageEffect,
	}
}

func NewDisengageEvent(id string) Event {
	return Event{
		ID:     id,
		Effect: DisengageEffect,
	}
}

func NewTakeOffEvent(id string) Event {
	return Event{
		ID:     id,
		Effect: TakeOffEffect,
	}
}

func NewLandingEvent(id string) Event {
	return Event{
		ID:     id,
		Effect: LandingEffect,
	}
}

func NewFlyEvent(id string, to PositionState) Event {
	return Event{
		ID:     id,
		Effect: FlyEffect,
		Data:   to,
	}
}

func NewChangeFuelTankEvent(id string, fuelTankID string) Event {
	return Event{
		ID:     id,
		Effect: FuelTankChangedEffect,
		Data:   fuelTankID,
	}
}
