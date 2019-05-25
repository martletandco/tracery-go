package tracery

import "testing"

/**
Literals
*/
func TestParseLiterals(t *testing.T) {
	var tests = []struct {
		input    string
		expected Rule
	}{
		{"a", LiteralValue{value: "a"}},
		{"A complete sentence, oh my.", LiteralValue{value: "A complete sentence, oh my."}},
		{"12,3456.7890", LiteralValue{value: "12,3456.7890"}},
		{"🥝", LiteralValue{value: "🥝"}},
		// @incomplete: these tests cause the current 'parse' func to loop forever
		// {`\[\]`, LiteralValue{value: "[]"}},
		// {`\#`, LiteralValue{value: "#"}},
		{`\\`, LiteralValue{value: `\`}},
		{`\#sym\#`, LiteralValue{value: "#sym#"}},
		{`\[key:literal\]`, LiteralValue{value: "[key:literal]"}},
	}

	for _, tt := range tests {
		actual := parse(tt.input)
		if actual != tt.expected {
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
			if actual != tt.expected {
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
			// #sym\# -- found symbol open with no close, did you mean '\#sym\#'?
			// \#sym# -- found symbol open with no close, did you mean '\#sym\#'?
			// #sym bol# -- symbols cannot contain spaces, did you mean '#sym_bol#'? -- stretch goal: suggest var name base on other key usage, e.g. snake, kebab, or cammel case
			// #sym\.mod# -- symbols cannot contain periods, did you mean '#sym.mod#' or '#sym_mod'?
			// #sym.# -- symbols cannot contain periods, did you mean to add a modifier or '#sym#'?
		}

		for _, tt := range tests {
			actual := parse(tt.input)
			if actual != tt.expected {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, tt.expected, actual)
			}
		}
	})
}

/**
actions
(repeat all with action or symbol inplace of lit
[act:lit]
[🥝:lit]
[act:lit,lit] // [act:lit|lit] compat
[act:lit,POP] ??
[act:\POP]
errs
\[act:lit] -- found action close without matching action open
[act\:lit] -- found action with no rule, did you mean to escape ':'?
[act:lit\] -- found action open with no matching close
[act] -- found action with no rule
[:lit] -- found push action with no symbol
[act:pop] -- ?? warning?

combinations
#[act:lit]sym# -- compat
errs
#[act:lit\]sym# -- found action open inside symbol with no close, did you mean '#[act:lit]sym#'?




extensions
[#sym#:rule] -- push to dynamic symbol
(rule|rule|rule) or #rule|rule|rule#
#sym.mod(#bol#)# -- symbol as modifier param


examples
push
push pop
nesting


*/
