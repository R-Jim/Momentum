package operator

type Command string

var (
	FlyCommand Command = "fly"
)

func (c Command) IsValid() bool {
	return c == FlyCommand
}
