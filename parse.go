package tracery

import (
	"regexp"
)

var actionRe = regexp.MustCompile(`^\[(.*?):(.*?)\]`)
var plainRe = regexp.MustCompile(`^([^\[#]+)`)
var tagRe = regexp.MustCompile(`^#([^#]+)#`)

func parse(input string) Rule {

	rules := []Rule{}
	var rule Rule
	var newIndex int
	for {
		if len(input) == 0 {
			break
		}
		nextChr := input[0:1] // @incomplete: doesn't handle unicode properly
		switch nextChr {
		case "[":
			rule, newIndex = parseAction(input)
		case "#":
			rule, newIndex = parseTag(input)
		default:
			rule, newIndex = parseLiteral(input)
		}

		rules = append(rules, rule)
		input = input[newIndex:]
	}
	if len(rules) == 1 {
		return rules[0]
	}
	return ListRule{rules: rules}
}

func parseAction(input string) (Rule, int) {
	index := actionRe.FindStringIndex(input)
	match := actionRe.FindStringSubmatch(input[index[0]:index[1]])
	key := match[1]
	rawValue := match[2]
	if rawValue == "POP" {
		return PopOp{key: key}, index[1]
	}
	value := parse(rawValue)
	return PushOp{key: key, value: value}, index[1]
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
