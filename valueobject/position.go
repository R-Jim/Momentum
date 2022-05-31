package valueobject

import "fmt"

type Position struct {
	X AxisValue
	Y AxisValue
}

func (p Position) String() string {
	return fmt.Sprintf("[%v, %v]", p.X, p.Y)
}
