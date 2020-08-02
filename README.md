# tracery-go
An (currently incomplete) implementation of the [Tracery](http://tracery.io/) text-expansion library by [GalaxyKate](http://www.galaxykate.com/)

```
import "github.com/martletandco/tracery-go"
```

```
g := tracery.NewGrammar()
g.PushRule("toggle", "flop")
g.PushRule("animal", "fox", "dog", "snail", "whale", "cow", "emu")
g.PushRule("pace", "quick", "slow", "sluggard", "rapid", "brisk")
g.PushRule("colour", "brown", "red", "purple", "orange")
g.PushRule("disposition", "lazy", "alert", "bored", "distracted", "eager")
out := g.Flatten("[toggle:flip]The #pace# #colour# #animal# jumped over the #disposition# #animal# to #toggle#[toggle:POP] #toggle#")
fmt.Println(out)
```

Outputs
```
The rapid red fox jumped over the eager emu
```

[See in the Go Playground](https://play.golang.org/p/wwn5d-L9iFC)

_Note that due to caching and other reasons random numbers to not really work in the playground_

## List of important features missing
- Load/save to/from JSON
- Parse errors
- Expand for grammar debugging*
- CBDQ compatibility†
- _(Probably many others that have escaped me just now)_

## List of less important but still missing features
- Convenience functions for serialising/deserialising rules
- A command line wrapper for ad-hock usage
- Adding rules in bulk
- Adding modifiers in bulk

## List of ideas to explore
- Symbol/modifier missing handler
- Improve indefinite article application<sup>[1](https://stackoverflow.com/a/4558514)</sup>
- Mode option to support new versions or different compatibility
- Add short hand for an in-place random selection based on `[x:1,2,3]#x#`

- *Parses and 'expands' eagerly, rather than lazily as the original does, so I'm not sure the same interface works
- †Some features are implemented, such as the _Random Push_ (i.e. `[x:1,2]`)
