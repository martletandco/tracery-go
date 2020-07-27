package tracery

import (
	"encoding/json"
	"testing"
)

func TestRuleSetUnmarshal(t *testing.T) {
	t.Run("valid inputs", func(t *testing.T) {
		var tests = []struct {
			input    string
			expected RuleSet
		}{
			{`{}`, RuleSet{}},
			{`{"x": "a"}`, RuleSet{"x": Rule{"a"}}},
			{`{"x": ["a"]}`, RuleSet{"x": Rule{"a"}}},
			{`{"x": ["a", "b"]}`, RuleSet{"x": Rule{"a", "b"}}},
			{`{"x": "a", "y": "b"}`, RuleSet{"x": Rule{"a"}, "y": Rule{"b"}}},
			{`{"x": ["a"], "y": "b"}`, RuleSet{"x": Rule{"a"}, "y": Rule{"b"}}},
			{`{"x": ["a"], "y": ["b"]}`, RuleSet{"x": Rule{"a"}, "y": Rule{"b"}}},
		}

		for _, tt := range tests {
			var set RuleSet
			if err := json.Unmarshal([]byte(tt.input), &set); err != nil {
				t.Errorf("String(%v): encountered error: %v", tt.input, err)
			}
			if ruleSetEqual(set, tt.expected) == false {
				t.Errorf("String(%v): expected '%v', got '%v'", tt.input, tt.expected, set)
			}
		}
	})
}

func ruleSetEqual(a, b RuleSet) bool {
	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for k, va := range a {
		vb := b[k]
		if len(va) != len(vb) {
			return false
		}
		for i := range va {
			if va[i] != vb[i] {
				return false
			}
		}
	}

	return true
}
