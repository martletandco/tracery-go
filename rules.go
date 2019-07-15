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

type SymbolModifier struct {
	key    string
	params []Rule
}

func (r SymbolModifier) String() string {
	return fmt.Sprintf("SymbolModifier<%v:%d:%v>", r.key, len(r.params), r.params)
}

type SymbolValue struct {
	key       string
	modifiers []SymbolModifier
}

func (r SymbolValue) Resolve(ctx Context) string {
	value := ctx.Lookup(r.key)
	if value == nil {
		return "((" + r.key + "))"
	}

	out := value.Resolve(ctx)

	for _, mod := range r.modifiers {
		fn := ctx.LookupModifier(mod.key)
		if fn == nil {
			continue
		}

		var params []string
		for _, rule := range mod.params {
			params = append(params, rule.Resolve(ctx))
		}

		out = fn(out, params...)
	}

	return out
}
func (r SymbolValue) String() string {
	return fmt.Sprintf("SymbolValue<%v:%d:%v>", r.key, len(r.modifiers), r.modifiers)
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

type RandomRule struct {
	rules []Rule
}

func (r RandomRule) Resolve(ctx Context) string {
	i := ctx.Intn(len(r.rules))
	return r.rules[i].Resolve(ctx)
}
func (r RandomRule) String() string {
	return fmt.Sprintf("RandomRule<%d:%v>", len(r.rules), r.rules)
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
