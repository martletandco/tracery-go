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
	LookupModifier(key string) (Modifier, bool)
}

type Modifier interface {
	Modify(value string, params ...string) string
}

/** ModifierFn is provided as a convience for using plain functions as Modifiers
 */
type ModifierFn func(value string, params ...string) string

/** Modify implements the single member of Modifier for ModifierFn
 */
func (f ModifierFn) Modify(value string, params ...string) string {
	if f == nil {
		return value
	}
	return f(value, params...)
}

type Grammar struct {
	Rand      func(n int) int
	value     map[string][]Rule
	modifiers map[string]Modifier
}

func NewGrammar() Grammar {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	return Grammar{
		Rand:      r.Intn,
		value:     make(map[string][]Rule),
		modifiers: make(map[string]Modifier),
	}
}

// Flatten resolves a grammer tree
func (g *Grammar) Flatten(input string) string {
	tree := parse(input)
	return tree.Resolve(g)
}

// @cleanup: Give this more informative name incl. target and strings
func (g *Grammar) PushRules(key string, inputs ...string) {
	for _, input := range inputs {
		rule := parse(input)
		g.Push(key, rule)
	}
}

func (g *Grammar) AddModifier(name string, mod Modifier) {
	g.modifiers[name] = mod
}

// Context implementation below

func (c *Grammar) Lookup(key string) Rule {
	rules, ok := c.value[key]
	if !ok {
		return nil
	}
	return rules[len(rules)-1]
}
func (c *Grammar) Push(key string, value Rule) {
	rules, ok := c.value[key]
	if !ok {
		c.value[key] = []Rule{value}
		return
	}
	c.value[key] = append(rules, value)
}
func (c *Grammar) Pop(key string) {
	rules, ok := c.value[key]
	if !ok || len(rules) == 1 {
		// Nothing left to pop (there is a different action to clear)
		// @enhance: warning about empty stack?
		return
	}

	c.value[key] = rules[:len(rules)-1]
}
func (c *Grammar) Intn(n int) int {
	return c.Rand(n)
}

func (c *Grammar) LookupModifier(key string) (Modifier, bool) {
	mod, ok := c.modifiers[key]
	return mod, ok
}
