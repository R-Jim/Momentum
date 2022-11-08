package knight

type Data interface{}

type Effect string

var (
	InitEffect   Effect = "KNIGHT_INITIALED"
	DamageEffect Effect = "KNIGHT_DAMAGED"
	StrikeEffect Effect = "KNIGHT_STROKE"

	MoveEffect Effect = "KNIGHT_MOVED"

	ChangeWeaponEffect Effect = "KNIGHT_WEAPON_CHANGED"
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

func NewStrikeEvent(id string, weaponID string) Event {
	return Event{
		ID:     id,
		Effect: StrikeEffect,
		Data:   weaponID,
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

func NewChangeWeaponEvent(id string, weaponID string) Event {
	return Event{
		ID:     id,
		Effect: ChangeWeaponEffect,
		Data:   weaponID,
	}
}
