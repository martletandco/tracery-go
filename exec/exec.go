package exec

type Operation interface {
	Resolve(ctx Context) string
}

type Modifier interface {
	Modify(value string, params ...string) string
}

type Context interface {
	Lookup(key string) Operation
	Push(key string, value Operation)
	Pop(key string)
	// https://golang.org/pkg/math/rand/#Intn
	Intn(n int) int
	LookupModifier(key string) (Modifier, bool)
}
