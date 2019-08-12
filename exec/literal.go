package exec

import "fmt"

type Literal struct {
	value string
}

func NewLiteral(value string) Literal {
	return Literal{value: value}
}

func (r Literal) Resolve(ctx Context) string {
	return r.value
}
func (r Literal) String() string {
	return fmt.Sprintf("Literal<%v>", r.value)
}
