package valueobject

import (
	"fmt"
)

type Momentum struct {
	X AxisValue
	Y AxisValue
}

func (m Momentum) String() string {
	return fmt.Sprintf("[%v, %v]", m.X, m.Y)
}
