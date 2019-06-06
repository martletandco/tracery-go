package tracery

import (
	"regexp"
)

var actionRe = regexp.MustCompile(`^\[(.*?):(.*?)\]`)
var plainRe = regexp.MustCompile(`^([^\[#]+)`)
var tagRe = regexp.MustCompile(`^#([^#]+)#`)

func parse(input string) Rule {

	rules := []Rule{}
	var index []int
	for {
		if len(input) == 0 {
			break
		}
		index = actionRe.FindStringIndex(input)
		if index != nil {
			rule, newIndex := parseAction(input)
			rules = append(rules, rule)
			input = input[newIndex:]
			continue
		}
		index = tagRe.FindStringIndex(input)
		if index != nil {
			rule, newIndex := parseTag(input)
			rules = append(rules, rule)
			input = input[newIndex:]
			continue
		}
		index = plainRe.FindStringIndex(input)
		if index != nil {
			rule, newIndex := parseLiteral(input)
			rules = append(rules, rule)
			input = input[newIndex:]
			continue
		}
	}
	if len(rules) == 1 {
		return rules[0]
	}
	return ListRule{rules: rules}
}

func parseAction(input string) (Rule, int) {
	index := actionRe.FindStringIndex(input)
	match := actionRe.FindStringSubmatch(input[index[0]:index[1]])
	value := parse(match[2])
	return PushOp{key: match[1], value: value}, index[1]
}

func parseTag(input string) (Rule, int) {
	index := tagRe.FindStringIndex(input)
	match := tagRe.FindStringSubmatch(input[index[0]:index[1]])
	return SymbolValue{key: match[1]}, index[1]
}

func parseLiteral(input string) (Rule, int) {
	index := plainRe.FindStringIndex(input)
	match := plainRe.FindStringSubmatch(input[index[0]:index[1]])
	return LiteralValue{value: match[1]}, index[1]
}
