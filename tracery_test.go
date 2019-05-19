package tracery

import "testing"

/**
Literals
*/
func TestFlattenLiterals(t *testing.T) {
	t.Run("it returns empty when given an empty rule", func(t *testing.T) {
		var g Grammar
		got := g.Flatten("")
		want := ""
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a plain value when given a plain value", func(t *testing.T) {
		var g Grammar
		got := g.Flatten("a")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a plain value when given a plain unicode value", func(t *testing.T) {
		var g Grammar
		got := g.Flatten("🌻")
		want := "🌻"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns an escaped tag", func(t *testing.T) {
		var g Grammar
		got := g.Flatten(`\#notakey\#`)
		want := "#notakey#"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns an escaped action", func(t *testing.T) {
		var g Grammar
		got := g.Flatten(`\[not:an,action\]`)
		want := "[not:an,action]"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

/**
Assignment and read (inline)
*/
func TestFlattenAssignmentAndReadInline(t *testing.T) {
	// @enhance: should return error or warning when configured to
	t.Run("it returns wrapped symbol when given a non-assigned symbol", func(t *testing.T) {
		var g Grammar
		got := g.Flatten("#x#")
		want := "((x))"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol", func(t *testing.T) {
		var g Grammar
		got := g.Flatten("[x:a]#x#")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol twice", func(t *testing.T) {
		var g Grammar
		got := g.Flatten("[x:a]#x##x#")
		want := "aa"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol before and after assignment", func(t *testing.T) {
		var g Grammar
		got := g.Flatten("#x#[x:a]#x#")
		want := "((x))a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns literals assigned and read from two symbols", func(t *testing.T) {
		var g Grammar
		got := g.Flatten("[x:a][y:b]#y# #x#")
		want := "b a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	// @incomplete: test for pop (unassignment)
}

func TestFlattenAssignmentAndReadContext(t *testing.T) {
	t.Run("it returns a literal assigned and read from a symbol", func(t *testing.T) {
		var g Grammar
		g.PushRules("x", []string{"a"})
		got := g.Flatten("#x#")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol twice", func(t *testing.T) {
		var g Grammar

		g.PushRules("x", []string{"a"})
		got := g.Flatten("#x##x#")
		want := "aa"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a symbol before and after local assignment", func(t *testing.T) {
		var g Grammar
		g.PushRules("x", []string{"a"})
		got := g.Flatten("#x#[x:b]#x#")
		want := "ab"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns literals assigned and read from two symbols", func(t *testing.T) {
		var g Grammar
		g.PushRules("x", []string{"a"})
		g.PushRules("y", []string{"b"})
		got := g.Flatten("#y# #x#")
		want := "b a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	// @incomplete: test for pop (unassignment)
}

/**
Rule select
*/
// func TestFlattenRuleSelect(t *testing.T) {

// }

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
num:[#count#]; [count:#num#]#count# = ?? // ?? need to add a depth limit?

CBDQ compat
[num:1][count:#num#,2]#count# = 1 | 2 // multiple assignment to key
[num:2][count:#num#]#[num:3]count# = 3  // assign literal inside tag
[two:4][num:1][count:#num#]#[num:#two#]count# = 4 // assign key inside tag
[two:4][num:1][count:#num#]#[num:3,#two#]count# = 3 | 4 // assign key inside tag
[num:1]#[num:#num#]num# = 1 // self assignment inside tag

Unassignment
[num:1]#num#[num:2]#num#[num:POP]#num# = 121 // Alternat sytanx for pop [num:] or [num] or [:num]?

Modifiers
e.g.
[words:two words]#words.sentencecase# = Two words
[words:words]#words.singliase# = word
[word:word]#word.plural# = words
[word:word]#word.upcase# = WORD
[word:WORD]#word.upcase# = word

*/
