package tracery

import (
	"fmt"
	"strings"
)

// @rename: Might want to separate operations and value, or at least get a better term
type Rule interface {
	Resolve(ctx Context) string
}

type LiteralValue struct {
	value string
}

func (r LiteralValue) Resolve(ctx Context) string {
	return r.value
}
func (r LiteralValue) String() string {
	return fmt.Sprintf("LiteralValue<%v>", r.value)
}

type SymbolValue struct {
	key string
}

func (r SymbolValue) Resolve(ctx Context) string {
	value := ctx.Lookup(r.key)
	if value == nil {
		return "((" + r.key + "))"
	}

	return value.Resolve(ctx)
}
func (r SymbolValue) String() string {
	return fmt.Sprintf("SymbolValue<%v>", r.key)
}

type ListRule struct {
	rules []Rule
}

func (r ListRule) Resolve(ctx Context) string {
	out := []string{}
	for _, rule := range r.rules {
		out = append(out, rule.Resolve(ctx))
	}
	return strings.Join(out, "")
}
func (r ListRule) String() string {
	return fmt.Sprintf("ListRule<%d:%v>", len(r.rules), r.rules)
}

type PushOp struct {
	key   string
	value Rule
}

func (r PushOp) Resolve(ctx Context) string {
	result := r.value.Resolve(ctx)
	ctx.Push(r.key, LiteralValue{value: result})
	return ""
}
func (r PushOp) String() string {
	return fmt.Sprintf("PushOp<%v:%v>", r.key, r.value)
}

type PopOp struct {
	key string
}

func (r PopOp) Resolve(ctx Context) string {
	ctx.Pop(r.key)
	return ""
}
func (r PopOp) String() string {
	return fmt.Sprintf("PopOp<%s>", r.key)
}
