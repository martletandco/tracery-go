package exec

import (
	"fmt"
	"strings"
)

type Concat struct {
	rules []Operation
}

func NewConcat(ops []Operation) Concat {
	return Concat{rules: ops}
}

func (r Concat) Resolve(ctx Context) string {
	out := []string{}
	for _, rule := range r.rules {
		out = append(out, rule.Resolve(ctx))
	}
	return strings.Join(out, "")
}
func (r Concat) String() string {
	return fmt.Sprintf("Concat<%d:%v>", len(r.rules), r.rules)
}
