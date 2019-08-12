package exec

import "fmt"

type Pop struct {
	key string
}

func NewPop(key string) Pop {
	return Pop{key: key}
}

func (r Pop) Resolve(ctx Context) string {
	ctx.Pop(r.key)
	return ""
}
func (r Pop) String() string {
	return fmt.Sprintf("Pop<%s>", r.key)
}
