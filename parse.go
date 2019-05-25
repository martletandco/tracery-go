package tracery

import (
	"regexp"
)

func parse(input string) Rule {
	actionRe := regexp.MustCompile(`^\[(.*?):(.*?)\]`)
	plainRe := regexp.MustCompile(`^([^\[#]+)`)
	tagRe := regexp.MustCompile(`^#([^#]+)#`)

	rules := []Rule{}
	var index []int
	for {
		if len(input) == 0 {
			break
		}
		index = actionRe.FindStringIndex(input)
		if index != nil {
			match := actionRe.FindStringSubmatch(input[index[0]:index[1]])
			value := parse(match[2])
			rules = append(rules, PushOp{key: match[1], value: value})
			input = input[index[1]:]
			continue
		}
		index = tagRe.FindStringIndex(input)
		if index != nil {
			match := tagRe.FindStringSubmatch(input[index[0]:index[1]])
			rules = append(rules, SymbolValue{key: match[1]})
			input = input[index[1]:]
			continue
		}
		index = plainRe.FindStringIndex(input)
		if index != nil {
			match := plainRe.FindStringSubmatch(input[index[0]:index[1]])
			if match != nil {
				rules = append(rules, LiteralValue{value: match[1]})
			}
			input = input[index[1]:]
			continue
		}
	}
	if len(rules) == 1 {
		return rules[0]
	}
	return ListRule{rules: rules}
}
