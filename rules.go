package tracery

type Rule interface {
	Resolve(ctx Context) string
}

type LiteralValue struct {
	value string
}

func (r LiteralValue) Resolve(ctx Context) string {
	return r.value
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

type PushOp struct {
	key   string
	value Rule
}

func (r PushOp) Resolve(ctx Context) string {
	result := r.value.Resolve(ctx)
	ctx.Set(r.key, LiteralValue{value: result})
	return ""
}
