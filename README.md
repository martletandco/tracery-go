# tracery-go
An (currently incomplete) implementation of the [Tracery](http://tracery.io/) text-expansion library by [GalaxyKate](http://www.galaxykate.com/)

```
g := tracery.NewGrammar()
g.PushRules("animal", "[x:fox,dog,snail,whale,cow,emu]#x#")
g.PushRules("pace", "[x:quick,slow,sluggard,rapid,brisk]#x#")
g.PushRules("colour", "[x:brown,red,purple,orange]#x#")
g.PushRules("disposition", "[x:lazy,alert,bored,distracted,eager]#x#")
out := g.Flatten("The #pace# #colour# #animal# jumped over the #disposition# #animal#")
fmt.Println(out)
```

Outputs
```
The rapid red fox jumped over the eager emu
```

[See in the Go Playground](https://play.golang.org/p/oItZ8pOX4ZW)

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
