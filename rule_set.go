package tracery

import (
	"encoding/json"
)

type RuleSet map[string]Rule

type Rule []string

// Handles `"rule"` and `["rule", "rule"]`
func (rs *Rule) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		return json.Unmarshal(b, (*[]string)(rs))
	}

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*rs = []string{s}
	return nil
}
