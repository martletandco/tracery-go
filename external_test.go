package tracery

import (
	"encoding/json"
	"testing"
)

// TestExternalGrammars are found examples from other implementations to help
// align the behaviour of this implemention with others
func TestExternalGrammars(t *testing.T) {
	assert := func(t *testing.T, got, want interface{}) {
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	}
	var tests = []struct {
		name     string
		origin   string
		input    string
		expected string
	}{
		{
			"comma workaround",
			"https://github.com/galaxykate/tracery/issues/20#issuecomment-220018871",
			`{
				"origin": "#defthing#The #WOTSIT#s #thing#",
				"defthing": [
						"[WOTSIT:colour][thing:#colour#]",
						"[WOTSIT:animal][thing:#animal#]"
				],
				"colour": "orange, blue and white",
				"animal": "unicorn, raven and sparrow"
			}`,
			"The colours orange, blue and white",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := grammarFromJSON(tt.input)
			if err != nil {
				t.Errorf("String(%v): encountered error: %v", tt.input, err)
			}
			got := g.Flatten("#origin#")
			assert(t, got, tt.expected)
		})
	}
}

func grammarFromJSON(input string) (Grammar, error) {
	g := NewGrammar()

	var set RuleSet
	if err := json.Unmarshal([]byte(input), &set); err != nil {
		return g, err
	}

	g.PushRuleSet(set)
	g.Rand = func(n int) int { return 0 }
	return g, nil
}
