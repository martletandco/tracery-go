package tracery

import (
	"encoding/json"
	"fmt"
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
		}

		for _, tt := range tests {
			var set RuleSet
			if err := json.Unmarshal([]byte(tt.input), &set); err != nil {
				t.Errorf("String(%v): encountered error: %v", tt.input, err)
			}
			if fmt.Sprintf("%v", set) != fmt.Sprintf("%v", tt.expected) {
				t.Errorf("String(%v): expected %v, set %v", tt.input, tt.expected, set)
			}
		}
	})
}
