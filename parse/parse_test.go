package parse

import (
	"fmt"
	"testing"

	"github.com/martletandco/tracery-go/exec"
)

func testRuleEq(a, b exec.Operation) bool {
	// Here we compare the operations by 'value'
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func testRulesEq(a, b []exec.Operation) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !testRuleEq(a[i], b[i]) {
			return false
		}
	}

	return true
}

/**
Literals
*/
func TestParseLiterals(t *testing.T) {
	var tests = []struct {
		input    string
		expected exec.Operation
	}{
		{"", exec.NewLiteral("")},
		{"a", exec.NewLiteral("a")},
		{"A complete sentence, oh my.", exec.NewLiteral("A complete sentence, oh my.")},
		{"12,3456.7890", exec.NewLiteral("12,3456.7890")},
		{"🥝", exec.NewLiteral("🥝")},
		{`\[\]`, exec.NewLiteral("[]")},
		{`\#`, exec.NewLiteral("#")},
		{`\\`, exec.NewLiteral(`\`)},
		{`\#sym\#`, exec.NewLiteral("#sym#")},
		{`\[key:literal\]`, exec.NewLiteral("[key:literal]")},
	}

	for _, tt := range tests {
		actual := String(tt.input)
		if !testRuleEq(actual, tt.expected) {
			t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}

func TestParseSymbols(t *testing.T) {
	t.Run("valid inputs", func(t *testing.T) {
		var tests = []struct {
			input    string
			expected exec.Operation
		}{
			{"#sym#", exec.NewSymbol("sym")},
			{"#symBol#", exec.NewSymbol("symBol")},
			{"#sym_bol#", exec.NewSymbol("sym_bol")},
			{"#🥝#", exec.NewSymbol("🥝")},
			{"#sym.mod#", exec.NewSymbolWithMods("sym", []exec.ModCall{exec.NewModCallZero("mod")})},
			{"#sym.mod.mod.mod#", exec.NewSymbolWithMods("sym", []exec.ModCall{exec.NewModCallZero("mod"), exec.NewModCallZero("mod"), exec.NewModCallZero("mod")})},
			{"#sym.mod(param)#", exec.NewSymbolWithMods("sym", []exec.ModCall{exec.NewModCall("mod", []exec.Operation{exec.NewLiteral("param")})})},
			{"#sym.mod(par,am)#", exec.NewSymbolWithMods("sym", []exec.ModCall{exec.NewModCall("mod", []exec.Operation{exec.NewLiteral("par"), exec.NewLiteral("am")})})},
		}

		for _, tt := range tests {
			actual := String(tt.input)
			if !testRuleEq(actual, tt.expected) {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		}
	})
	t.Run("inputs which should error", func(t *testing.T) {
		var tests = []struct {
			input    string
			expected exec.Operation
		}{
			// @incomplete: change func and test signature to include an error
			// #sym bol -- found symbol open with no close, did you mean '#sym#'?
			// #sym # -- found whitespace in symbol, did you mean '#sym#'?
			// # sym# -- found whitespace in symbol, did you mean '#sym#'?
			// #sym bol# -- found whitespace in symbol, did you mean '#sym_bol#'?
			// #sym\# -- found symbol open with no close, did you mean '\#sym\#'?
			// \#sym# -- found symbol open with no close, did you mean '\#sym\#'?
			// #sym bol# -- symbols cannot contain spaces, did you mean '#sym_bol#'? -- stretch goal: suggest var name base on other key usage, e.g. snake, kebab, or cammel case
			// #sym\.mod# -- symbols cannot contain periods, did you mean '#sym.mod#' or '#sym_mod'?
			// #sym.# -- symbols cannot contain periods, did you mean to add a modifier or '#sym#'?
			// #sym,sym# -- symbols cannot contain commas, did you mean '#sym-sym#'?
		}

		for _, tt := range tests {
			actual := String(tt.input)
			if !testRuleEq(actual, tt.expected) {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		}
	})
}

func TestParseActions(t *testing.T) {
	t.Run("valid inputs", func(t *testing.T) {
		var tests = []struct {
			input    string
			expected exec.Operation
		}{
			// @incomplete: repeat all with action or symbol inplace of lit
			{"[act:lit]", exec.NewPush("act", exec.NewLiteral("lit"))},
			{"[🥝:lit]", exec.NewPush("🥝", exec.NewLiteral("lit"))},
			{"[act:lit,lit]", exec.NewPush("act", exec.NewSelect([]exec.Operation{exec.NewLiteral("lit"), exec.NewLiteral("lit")}))},
			{`[act:lit\,eral]`, exec.NewPush("act", exec.NewLiteral("lit,eral"))},
			{"[act:POP]", exec.NewPop("act")},
			// @question: Can POP be escaped?
			// {"[act:\POP]", PushOp{key: "act", value: LiteralValue{value: "POP"}}},
		}

		for _, tt := range tests {
			actual := String(tt.input)
			if !testRuleEq(actual, tt.expected) {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		}
	})
	t.Run("inputs which should error", func(t *testing.T) {
		var tests = []struct {
			input    string
			expected exec.Operation
		}{
			// @incomplete: change func and test signature to include an error
			// @incomplete: whitespace errors around symbol
			// \[act:lit] -- found action close without matching action open
			// [act\:lit] -- found action with no rule, did you mean to escape ':'?
			// [act:lit\] -- found action open with no matching close
			// [act] -- found action with no rule
			// [:lit] -- found push action with no symbol
			// [act:pop] -- ?? warning?
		}

		for _, tt := range tests {
			actual := String(tt.input)
			if !testRuleEq(actual, tt.expected) {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		}
	})
}

// /**
// combinations
// #[act:lit]sym# -- compat
// #[act:lit]sym.mod(value)# -- compat
// errs
// #[act:lit\]sym# -- found action open inside symbol with no close, did you mean '#[act:lit]sym#'?

// extensions
// [#sym#:rule] -- push to dynamic symbol
// (rule|rule|rule) or #rule|rule|rule#
// #sym.mod(#bol#)# -- symbol as modifier param
// [act:lit|lit] -- different separator
// #a,b,c# -- random symbol

// questions
// [act:lit,POP] -- legal?
// [act:\POP] -- can you escape POP?
// [p:\POP][act:#p#] -- does this pop?

// examples
// push
// push pop
// nesting

// */
