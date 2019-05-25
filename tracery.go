package tracery

type Context interface {
	Lookup(key string) Rule
	Set(key string, value Rule)
}

type MapContext struct {
	value map[string][]Rule
}

func newMapContext() MapContext {
	return MapContext{
		value: make(map[string][]Rule),
	}
}
func (c *MapContext) Lookup(key string) Rule {
	rules, ok := c.value[key]
	if !ok {
		return nil
	}
	return rules[len(rules)-1]
}
func (c *MapContext) Set(key string, value Rule) {
	rules, ok := c.value[key]
	if !ok {
		c.value[key] = []Rule{value}
		return
	}
	c.value[key] = append(rules, value)
}

type Grammar struct {
	ctx Context
}

func NewGrammar() Grammar {
	ctx := newMapContext()
	return Grammar{
		ctx: &ctx,
	}
}

// Flatten resolves a grammer tree
func (g *Grammar) Flatten(input string) string {
	return parse(input).Resolve(g.ctx)
}

func (g *Grammar) PushRules(key string, inputs ...string) {
	for _, input := range inputs {
		rule := parse(input)
		g.ctx.Set(key, rule)
	}
}
