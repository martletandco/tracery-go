package exec

import "fmt"

type Select struct {
	ops []Operation
}

func NewSelect(ops []Operation) Select {
	return Select{ops: ops}
}

func (r Select) Resolve(ctx Context) string {
	i := ctx.Intn(len(r.ops))
	return r.ops[i].Resolve(ctx)
}
func (r Select) String() string {
	return fmt.Sprintf("Select<%d:%v>", len(r.ops), r.ops)
}
