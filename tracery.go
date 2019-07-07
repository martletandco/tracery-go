package tracery

import (
	"math/rand"
	"time"
)

type Context interface {
	Lookup(key string) Rule
	Push(key string, value Rule)
	Pop(key string)
	// https://golang.org/pkg/math/rand/#Intn
	Intn(n int) int
}

type MapContext struct {
	Rand  *rand.Rand
	value map[string][]Rule
}

func newMapContext() MapContext {
	s := rand.NewSource(time.Now().Unix())
	return MapContext{
		Rand:  rand.New(s),
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
func (c *MapContext) Push(key string, value Rule) {
	rules, ok := c.value[key]
	if !ok {
		c.value[key] = []Rule{value}
		return
	}
	c.value[key] = append(rules, value)
}
func (c *MapContext) Pop(key string) {
	rules, ok := c.value[key]
	if !ok || len(rules) == 1 {
		// Nothing left to pop (there is a different action to clear)
		// @enhance: warning about empty stack?
		return
	}

	c.value[key] = rules[:len(rules)-1]
}
func (c *MapContext) Intn(n int) int {
	// incomplete: use c.Rand.Intn
	return 0
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
		g.ctx.Push(key, rule)
	}
}
