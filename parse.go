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
		case RightBracket:
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
	rawValue := scanner.Next()
	if rawValue.Value == "POP" {
		return PopOp{key: key}
	}
	value := parse(rawValue.Value)
	return PushOp{key: key, value: value}
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
		case Text:
			texts = append(texts, token.Value)
		default:
			break Loop
		}
	}
	value := strings.Join(texts, ",")
	return LiteralValue{value: value}
}
