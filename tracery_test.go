package tracery

import "testing"

/**
Literals
*/
func TestFlattenLiterals(t *testing.T) {
	t.Run("it returns empty when given an empty rule", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten("")
		want := ""
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a plain value when given a plain value", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten("a")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a plain value when given a plain unicode value", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten("🌻")
		want := "🌻"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns an escaped tag", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten(`\#notakey\#`)
		want := "#notakey#"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns an escaped action", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten(`\[not:an,action\]`)
		want := "[not:an,action]"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

/**
Push and read (inline)
*/
func TestFlattenPushAndReadInline(t *testing.T) {
	// @enhance: should return error or warning when configured to
	t.Run("it returns wrapped symbol when given a non-assigned symbol", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten("#x#")
		want := "((x))"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten("[x:a]#x#")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol twice", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten("[x:a]#x##x#")
		want := "aa"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol before and after assignment", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten("#x#[x:a]#x#")
		want := "((x))a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns literals assigned and read from two symbols", func(t *testing.T) {
		g := NewGrammar()
		got := g.Flatten("[x:a][y:b]#y# #x#")
		want := "b a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

func TestFlattenPushAndReadContext(t *testing.T) {
	t.Run("it returns a literal assigned and read from a symbol", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "a")
		got := g.Flatten("#x#")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol twice", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "a")
		got := g.Flatten("#x##x#")
		want := "aa"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol before and after local assignment", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "a")
		got := g.Flatten("#x#[x:b]#x#")
		want := "ab"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns literals assigned and read from two symbols", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "a")
		g.PushRules("y", "b")
		got := g.Flatten("#y# #x#")
		want := "b a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it reads the inline value before the context value", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "2")
		g.PushRules("y", "#x#")
		got := g.Flatten("[x:1]#y#")
		want := "1"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it reads and evaluates the values in order", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("y", "[y:#x#]#x#")
		g.PushRules("x", "[x:1]#y#[x:2]")
		got := g.Flatten("#y##y##x#")
		want := "212"
		// Expands to setting y = x = 1, x = 2, put x, put y, put x
		// [y: [x:1] [y:#x#]#x# [x:2] ]#x# #y# #x#
		/* So:
		1. #y# -> [y:#x#]#x# -- resolve y*
			a. #y# -> [y:[x:1]#y#[x:2]]#x# -- resolve x*
			b. #y# -> [y:(x=1)#y#[x:2]]#x# -- assign 1 to x
			c. #y# -> [y:(x=1)[y:#x#]#x#[x:2]]#x# -- resolve y†
			d. #y# -> [y:(x=1)(y=x{1})#x#[x:2]]#x# -- y = x (which is 1)
			e. #y# -> [y:(x=1)(y=x{1})1[x:2]]#x# -- resolve x† (which is 1)
			f. #y# -> [y:(x=1)(y=x{1})1(x=2)]#x# -- assign 2 to x
			g. #y# -> [y:1]#x# -- finished resolving x*
			h. #y# -> (y=1)#x# -- assign 1 to y
			i. #y# -> (y=1)2 -- resolve x
			j. #y# -> 2 -- finished resolving y*
		2. #y# -> 1 -- resolve y (which is 1 from 1.h)
		3. #x# -> 2 -- resolve x (which is 2 from 1.f)
		*/
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

func TestFlattenPop(t *testing.T) {
	t.Run("it returns the original value of a symbol after it's popped", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "a", "b")
		got := g.Flatten("[x:c]#x#[x:POP]#x#[x:POP]#x#")
		want := "cba"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
	t.Run("it ignores pop action when stack is empty", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "a")
		got := g.Flatten("#x#[x:POP]#x#")
		want := "aa"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

/**
Rule select
*/
func TestFlattenRuleSelect(t *testing.T) {
	t.Run("it selects a random literal rule to push", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "[y:a,b,c]")
		g.Rand = func(n int) int { return 0 }
		got := g.Flatten("#x##y#")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
		g.Rand = func(n int) int { return 2 }
		got = g.Flatten("#x##y#")
		want = "c"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
	t.Run("it selects a random literal or symbol rule to push", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "[y:a,#b#]")
		g.PushRules("b", "🥝")
		g.Rand = func(n int) int { return 0 }
		got := g.Flatten("#x##y#")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
		g.Rand = func(n int) int { return 1 }
		got = g.Flatten("#x##y#")
		want = "🥝"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
	t.Run("it selects a random rule to push which also pushes", func(t *testing.T) {
		g := NewGrammar()
		g.PushRules("x", "[y:a,#b#]")
		g.PushRules("b", "[c:🥝]")
		g.PushRules("c", "🏔")
		g.Rand = func(n int) int { return 0 }
		got := g.Flatten("#c##x##y##c#")
		want := "🏔a🏔"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
		g.Rand = func(n int) int { return 1 }
		got = g.Flatten("#c##x##y##c#")
		want = "🏔🥝"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

/**
#num# = ''// missing key
[num:1]#num# = 1 // assign literal
[num:1]#num##num# = 11
[num:1][animal:dog]#num# #animal# = 1 dog
[num:2][count:#num#]#count# = 2  // assign key
[num:1][num:#num#]#num# = 1 // self assignment
// unclosed tag
// unclosed action
// escaping

Recursion
[num:1][count:#num#][num:#count#]#num# = 1 // overriding
num:#count#; [count:#num#]#count# = ((count)) // missing key

CBDQ compat
// @incomplete: add option to enable?
[num:1][count:#num#,2]#count# = 1 | 2 // multiple assignment to key
[num:2][count:#num#]#[num:3]count# = 3  // assign literal inside tag
[two:4][num:1][count:#num#]#[num:#two#]count# = 4 // assign key inside tag
[two:4][num:1][count:#num#]#[num:3,#two#]count# = 3 | 4 // assign key inside tag
[num:1]#[num:#num#]num# = 1 // self assignment inside tag

Pop
[num:1]#num#[num:2]#num#[num:POP]#num# = 121 // Alternat sytanx for pop [num:] or [num] or [:num]?

Modifiers
e.g.
[words:two words]#words.sentencecase# = Two words
[words:words]#words.singliase# = word
[word:word]#word.plural# = words
[word:word]#word.upcase# = WORD
[word:WORD]#word.upcase# = word
Unknown modifier
// order of resolving when applying modifiers
count::1
num::#count#[num:3]
#num.join(#num#)# = 13

*/
