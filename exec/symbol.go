package exec

import "fmt"

type ModCall struct {
	key    string
	params []Operation
}

func NewModCallZero(key string) ModCall {
	return ModCall{key: key}
}

func NewModCall(key string, params []Operation) ModCall {
	return ModCall{key: key, params: params}
}

func (r ModCall) String() string {
	return fmt.Sprintf("ModCall	<%v:%d:%v>", r.key, len(r.params), r.params)
}

type Symbol struct {
	key  string
	mods []ModCall
}

func NewSymbol(key string) Symbol {
	return Symbol{key: key}
}
func NewSymbolWithMods(key string, mods []ModCall) Symbol {
	return Symbol{key: key, mods: mods}
}

func (r Symbol) Resolve(ctx Context) string {
	value := ctx.Lookup(r.key)
	if value == nil {
		return "((" + r.key + "))"
	}

	out := value.Resolve(ctx)

	for _, mod := range r.mods {
		m, ok := ctx.LookupModifier(mod.key)
		if !ok {
			out = out + "((." + mod.key + "))"
			continue
		}

		var params []string
		for _, rule := range mod.params {
			params = append(params, rule.Resolve(ctx))
		}

		out = m.Modify(out, params...)
	}

	return out
}
func (r Symbol) String() string {
	return fmt.Sprintf("Symbol<%v:%d:%v>", r.key, len(r.mods), r.mods)
}
