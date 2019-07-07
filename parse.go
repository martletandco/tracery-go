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
	for {
		rawValue := scanner.Next()
		if rawValue.Type == RightBracket {
			break
		}
		if rawValue.Value == "POP" {
			// Consume closing ]
			scanner.Next()
			return PopOp{key: key}
		}
		rule := parse(rawValue.Value)
		rules = append(rules, rule)
	}
	return PushOp{key: key, value: RandomRule{rules: rules}}
}

func parseTag(scanner *Scanner) Rule {
	// @cleanup: whole lot of assume the input is valid here
	// Consume opening #
	scanner.Next()
	key := scanner.Next().Value
	// Consume closing #
	scanner.Next()
	return SymbolValue{key: key}
}

func parseLiteral(scanner *Scanner) Rule {
	var texts []string

	var token Token
Loop:
	for {
		token = scanner.Next()
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
			texts = append(texts, token.Value)
		}
	}
	value := strings.Join(texts, "")
	return LiteralValue{value: value}
}
