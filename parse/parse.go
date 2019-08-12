package parse

import (
	"strings"

	"github.com/martletandco/tracery-go/exec"
	"github.com/martletandco/tracery-go/scan"
)

func String(input string) exec.Operation {
	rules := []exec.Operation{}
	var rule exec.Operation
	scanner := scan.New(input)
	for {
		token := scanner.Peek()
		if token.Type == scan.EOF {
			break
		}
		switch token.Type {
		case scan.LeftBracket:
			rule = parseAction(scanner)
		case scan.Octo:
			rule = parseTag(scanner)
		default:
			rule = parseLiteral(scanner)
		}

		rules = append(rules, rule)
	}
	if len(rules) == 0 {
		return exec.NewLiteral("")
	}
	if len(rules) == 1 {
		return rules[0]
	}
	return exec.NewConcat(rules)
}

func parseAction(s *scan.Scanner) exec.Operation {
	// @cleanup: whole lot of assume the input is valid here
	// Consume opening [
	s.Next()
	keyToken := s.Next()
	key := keyToken.Value
	// Consume :
	s.Next()
	var rules []exec.Operation

	var ruleValue []string
	for {
		rawValue := s.Next()
		if rawValue.Type == scan.RightBracket {
			rule := String(strings.Join(ruleValue, ""))
			rules = append(rules, rule)
			break
		}
		if rawValue.Value == "POP" {
			// Consume closing ]
			s.Next()
			return exec.NewPop(key)
		}
		if rawValue.Type == scan.Comma {
			rule := String(strings.Join(ruleValue, ""))
			rules = append(rules, rule)
			ruleValue = ruleValue[:0]
			continue
		}

		ruleValue = append(ruleValue, rawValue.Value)
	}
	if len(rules) == 0 {
		return exec.NewPush(key, exec.NewLiteral(""))
	}
	if len(rules) == 1 {
		return exec.NewPush(key, rules[0])
	}

	return exec.NewPush(key, exec.NewSelect(rules))
}

func parseTag(s *scan.Scanner) exec.Operation {
	// @cleanup: whole lot of assume the input is valid here
	// Consume opening #
	s.Next()
	key := s.Next().Value
	var mods []exec.ModCall

	var ruleValue []string
	for {
		rawValue := s.Next()
		if rawValue.Type == scan.Octo {
			break
		}
		if rawValue.Type == scan.Period {
			mod := parseModifier(s)
			mods = append(mods, mod)
			continue
		}

		ruleValue = append(ruleValue, rawValue.Value)
	}

	return exec.NewSymbolWithMods(key, mods)
}

func parseModifier(s *scan.Scanner) exec.ModCall {
	key := s.Next().Value

	if s.Peek().Type != scan.LeftParen {
		return exec.NewModCallZero(key)
	}
	// Consume (
	s.Next()

	var rules []exec.Operation

	var ruleValue []string
	for {
		rawValue := s.Next()
		if rawValue.Type == scan.RightParen {
			rule := String(strings.Join(ruleValue, ""))
			rules = append(rules, rule)
			break
		}
		if rawValue.Type == scan.Comma {
			rule := String(strings.Join(ruleValue, ""))
			rules = append(rules, rule)
			ruleValue = ruleValue[:0]
			continue
		}

		ruleValue = append(ruleValue, rawValue.Value)
	}

	return exec.NewModCall(key, rules)
}

func parseLiteral(s *scan.Scanner) exec.Operation {
	var texts []string

	var token scan.Token
Loop:
	for {
		token = s.Peek()
		switch token.Type {
		case scan.LeftBracket:
			fallthrough
		case scan.Octo:
			fallthrough
		case scan.EOF:
			fallthrough
		case scan.Error:
			break Loop
		default:
			texts = append(texts, s.Next().Value)
		}
	}
	value := strings.Join(texts, "")
	return exec.NewLiteral(value)
}
