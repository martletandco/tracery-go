package tracery

import "testing"

func testRuleEq(a, b Rule) bool {
	switch v := a.(type) {
	case ListRule:
		return testListRuleEq(v, b)
	case RandomRule:
		return testRandomRuleEq(v, b)
	case PushOp:
		return testPushOpEq(v, b)
	default:
		// If this is panicing we might need a type specific comparison type
		return a == b
	}
}

func testListRuleEq(a ListRule, b Rule) bool {
	if bl, ok := b.(ListRule); ok {
		return testRulesEq(a.rules, bl.rules)
	}
	return false
}

func testRandomRuleEq(a RandomRule, b Rule) bool {
	if br, ok := b.(RandomRule); ok {
		return testRulesEq(a.rules, br.rules)
	}
	return false
}

func testPushOpEq(a PushOp, b Rule) bool {
	if br, ok := b.(PushOp); ok {
		return testRuleEq(a.value, br.value)
	}
	return false
}

func testRulesEq(a, b []Rule) bool {
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
		expected Rule
	}{
		{"", LiteralValue{value: ""}},
		{"a", LiteralValue{value: "a"}},
		{"A complete sentence, oh my.", LiteralValue{value: "A complete sentence, oh my."}},
		{"12,3456.7890", LiteralValue{value: "12,3456.7890"}},
		{"🥝", LiteralValue{value: "🥝"}},
		{`\[\]`, LiteralValue{value: "[]"}},
		{`\#`, LiteralValue{value: "#"}},
		{`\\`, LiteralValue{value: `\`}},
		{`\#sym\#`, LiteralValue{value: "#sym#"}},
		{`\[key:literal\]`, LiteralValue{value: "[key:literal]"}},
	}

	for _, tt := range tests {
		actual := parse(tt.input)
		if !testRuleEq(actual, tt.expected) {
			t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}

func TestParseSymbols(t *testing.T) {
	t.Run("valid inputs", func(t *testing.T) {
		var tests = []struct {
			input    string
			expected Rule
		}{
			{"#sym#", SymbolValue{key: "sym"}},
			{"#symBol#", SymbolValue{key: "symBol"}},
			{"#sym_bol#", SymbolValue{key: "sym_bol"}},
			{"#🥝#", SymbolValue{key: "🥝"}},
			// @incomplete: Are modifiers a rule or are they a list on SymbolValue?
			// {"#sym.mod#", SymbolValue{key: "sym"}},
			// {"#sym.mod.mod#", SymbolValue{key: "sym"}},
			// {"#sym.mod.mod.mod.mod.mod.mod#", SymbolValue{key: "sym"}},
			// {"#sym.mod(param)", SymbolValue{key: "sym"}},
		}

		for _, tt := range tests {
			actual := parse(tt.input)
			if !testRuleEq(actual, tt.expected) {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		}
	})
	t.Run("inputs which should error", func(t *testing.T) {
		var tests = []struct {
			input    string
			expected Rule
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
		}

		for _, tt := range tests {
			actual := parse(tt.input)
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
			expected Rule
		}{
			// @incomplete: repeat all with action or symbol inplace of lit
			{"[act:lit]", PushOp{key: "act", value: LiteralValue{value: "lit"}}},
			{"[🥝:lit]", PushOp{key: "🥝", value: LiteralValue{value: "lit"}}},
			{"[act:lit,lit]", PushOp{key: "act", value: RandomRule{rules: []Rule{LiteralValue{value: "lit"}, LiteralValue{value: "lit"}}}}},
			{`[act:lit\,eral]`, PushOp{key: "act", value: LiteralValue{value: "lit,eral"}}},
			{"[act:POP]", PopOp{key: "act"}},
		}

		for _, tt := range tests {
			actual := parse(tt.input)
			if !testRuleEq(actual, tt.expected) {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		}
	})
	t.Run("inputs which should error", func(t *testing.T) {
		var tests = []struct {
			input    string
			expected Rule
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
			actual := parse(tt.input)
			if !testRuleEq(actual, tt.expected) {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		}
	})
}

/**
combinations
#[act:lit]sym# -- compat
#[act:lit]sym.mod(value)# -- compat
errs
#[act:lit\]sym# -- found action open inside symbol with no close, did you mean '#[act:lit]sym#'?

extensions
[#sym#:rule] -- push to dynamic symbol
(rule|rule|rule) or #rule|rule|rule#
#sym.mod(#bol#)# -- symbol as modifier param
[act:lit|lit] -- different separator

questions
[act:lit,POP] -- legal?
[act:\POP] -- can you escape POP?
[p:\POP][act:#p#] -- does this pop?

examples
push
push pop
nesting

*/
