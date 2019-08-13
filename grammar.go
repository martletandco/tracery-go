package tracery

import (
	"math/rand"
	"time"

	"github.com/martletandco/tracery-go/exec"
	"github.com/martletandco/tracery-go/parse"
)

type Grammar struct {
	Rand      func(n int) int
	value     map[string][]exec.Operation
	modifiers map[string]exec.Modifier
}

func NewGrammar() Grammar {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	return Grammar{
		Rand:      r.Intn,
		value:     make(map[string][]exec.Operation),
		modifiers: make(map[string]exec.Modifier),
	}
}

// Flatten resolves a grammar tree
func (g *Grammar) Flatten(input string) string {
	tree := parse.String(input)
	return tree.Resolve(g)
}

// @cleanup: Give this more informative name incl. target and strings
func (g *Grammar) PushRules(key string, inputs ...string) {
	for _, input := range inputs {
		rule := parse.String(input)
		g.Push(key, rule)
	}
}

func (g *Grammar) AddModifier(name string, mod exec.Modifier) {
	g.modifiers[name] = mod
}

func (g *Grammar) AddModifyFunc(name string, mod func(value string, params ...string) string) {
	g.AddModifier(name, ModifierFunc(mod))
}

// Context implementation below

func (c *Grammar) Lookup(key string) exec.Operation {
	rules, ok := c.value[key]
	if !ok {
		return nil
	}
	return rules[len(rules)-1]
}
func (c *Grammar) Push(key string, value exec.Operation) {
	rules, ok := c.value[key]
	if !ok {
		c.value[key] = []exec.Operation{value}
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

func (c *Grammar) LookupModifier(key string) (exec.Modifier, bool) {
	mod, ok := c.modifiers[key]
	return mod, ok
}
