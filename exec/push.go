package exec

import "fmt"

type Push struct {
	key   string
	value Operation
}

func NewPush(key string, value Operation) Push {
	return Push{key: key, value: value}
}

func (r Push) Resolve(ctx Context) string {
	result := r.value.Resolve(ctx)
	ctx.Push(r.key, NewLiteral(result))
	return ""
}
func (r Push) String() string {
	return fmt.Sprintf("Push<%v:%v>", r.key, r.value)
}
