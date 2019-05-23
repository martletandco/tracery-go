package tracery

import "regexp"

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
	return g.parse(input).Resolve(g.ctx)
}

func (g *Grammar) PushRules(key string, inputs ...string) {
	for _, input := range inputs {
		rule := g.parse(input)
		g.ctx.Set(key, rule)
	}
}

func (g Grammar) parse(input string) Rule {
	actionRe := regexp.MustCompile(`^\[(.*?):(.*?)\]`)
	plainRe := regexp.MustCompile(`^([^\[#]+)`)
	tagRe := regexp.MustCompile(`^#([^#]+)#`)

	rules := []Rule{}
	var index []int
	for {
		if len(input) == 0 {
			break
		}
		index = actionRe.FindStringIndex(input)
		if index != nil {
			match := actionRe.FindStringSubmatch(input[index[0]:index[1]])
			rules = append(rules, PushOp{key: match[1], value: LiteralValue{value: match[2]}})
			input = input[index[1]:]
			continue
		}
		index = tagRe.FindStringIndex(input)
		if index != nil {
			match := tagRe.FindStringSubmatch(input[index[0]:index[1]])
			rules = append(rules, SymbolValue{key: match[1]})
			input = input[index[1]:]
			continue
		}
		index = plainRe.FindStringIndex(input)
		if index != nil {
			match := plainRe.FindStringSubmatch(input[index[0]:index[1]])
			if match != nil {
				rules = append(rules, LiteralValue{value: match[1]})
			}
			input = input[index[1]:]
			continue
		}
	}
	return VariadicRule{rules: rules}
}
