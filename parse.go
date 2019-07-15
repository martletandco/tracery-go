package tracery

import (
	"strings"
)

func parse(input string) Rule {
	rules := []Rule{}
	var rule Rule
	scanner := newScanner(input)
	for {
		token := scanner.Peek()
		if token.Type == EOF {
			break
		}
		switch token.Type {
		case LeftBracket:
			rule = parseAction(scanner)
		case Octo:
			rule = parseTag(scanner)
		default:
			rule = parseLiteral(scanner)
		}

		rules = append(rules, rule)
	}
	if len(rules) == 0 {
		return LiteralValue{""}
	}
	if len(rules) == 1 {
		return rules[0]
	}
	return ListRule{rules: rules}
}

func parseAction(scanner *Scanner) Rule {
	// @cleanup: whole lot of assume the input is valid here
	// Consume opening [
	scanner.Next()
	keyToken := scanner.Next()
	key := keyToken.Value
	// Consume :
	scanner.Next()
	var rules []Rule

	var ruleValue []string
	for {
		rawValue := scanner.Next()
		if rawValue.Type == RightBracket {
			rule := parse(strings.Join(ruleValue, ""))
			rules = append(rules, rule)
			break
		}
		if rawValue.Value == "POP" {
			// Consume closing ]
			scanner.Next()
			return PopOp{key: key}
		}
		if rawValue.Type == Comma {
			rule := parse(strings.Join(ruleValue, ""))
			rules = append(rules, rule)
			ruleValue = ruleValue[:0]
			continue
		}

		ruleValue = append(ruleValue, rawValue.Value)
	}
	if len(rules) == 0 {
		return PushOp{key: key, value: LiteralValue{value: ""}}
	}
	if len(rules) == 1 {
		return PushOp{key: key, value: rules[0]}
	}

	return PushOp{key: key, value: RandomRule{rules: rules}}
}

func parseTag(scanner *Scanner) Rule {
	// @cleanup: whole lot of assume the input is valid here
	// Consume opening #
	scanner.Next()
	key := scanner.Next().Value
	var mods []SymbolModifier

	var ruleValue []string
	for {
		rawValue := scanner.Next()
		if rawValue.Type == Octo {
			break
		}
		if rawValue.Type == Period {
			mod := parseModifier(scanner)
			mods = append(mods, mod)
			continue
		}

		ruleValue = append(ruleValue, rawValue.Value)
	}

	return SymbolValue{key: key, modifiers: mods}
}

func parseModifier(scanner *Scanner) SymbolModifier {
	key := scanner.Next().Value

	if scanner.Peek().Type != LeftParen {
		return SymbolModifier{key: key}
	}
	// Consume (
	scanner.Next()

	var rules []Rule

	var ruleValue []string
	for {
		rawValue := scanner.Next()
		if rawValue.Type == RightParen {
			rule := parse(strings.Join(ruleValue, ""))
			rules = append(rules, rule)
			break
		}
		if rawValue.Type == Comma {
			rule := parse(strings.Join(ruleValue, ""))
			rules = append(rules, rule)
			ruleValue = ruleValue[:0]
			continue
		}

		ruleValue = append(ruleValue, rawValue.Value)
	}

	return SymbolModifier{key: key, params: rules}
}

func parseLiteral(scanner *Scanner) Rule {
	var texts []string

	var token Token
Loop:
	for {
		token = scanner.Peek()
		switch token.Type {
		case LeftBracket:
			fallthrough
		case Octo:
			fallthrough
		case EOF:
			fallthrough
		case Error:
			break Loop
		default:
			texts = append(texts, scanner.Next().Value)
		}
	}
	value := strings.Join(texts, "")
	return LiteralValue{value: value}
}
