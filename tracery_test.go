package tracery

import "testing"

/**
Rule select
*/

func TestFlatten(t *testing.T) {
	t.Run("it returns empty when given an empty rule", func(t *testing.T) {
		got := Flatten("")
		want := ""
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a plain value when given a plain value", func(t *testing.T) {
		got := Flatten("a")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	/**
	Assignment and read
	*/
	t.Run("it returns empty when given a non-assigned key", func(t *testing.T) {
		got := Flatten("#x#")
		want := ""
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a key", func(t *testing.T) {
		got := Flatten("[x:a]#x#")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a key twice", func(t *testing.T) {
		got := Flatten("[x:a]#x##x#")
		want := "aa"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns a literal assigned and read from a key before and after assignment", func(t *testing.T) {
		got := Flatten("#x#[x:a]#x#")
		want := "a"
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("it returns literals assigned and read from two keys", func(t *testing.T) {
		got := Flatten("[x:a][y:b]#y# #x#")
		want := "b a"
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
