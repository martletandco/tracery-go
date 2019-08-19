package parse

import (
	"strings"

	"github.com/martletandco/tracery-go/exec"
	"github.com/martletandco/tracery-go/scan"
)

func String(input string) exec.Operation {
	ops := []exec.Operation{}
	var op exec.Operation
	scanner := scan.New(input)
	for {
		token := scanner.Peek()
		if token.Type == scan.EOF {
			break
		}
		switch token.Type {
		case scan.LeftBracket:
			op = parseAction(scanner)
		case scan.Octo:
			op = parseTag(scanner)
		default:
			op = parseLiteral(scanner)
		}

		ops = append(ops, op)
	}
	if len(ops) == 0 {
		return exec.NewLiteral("")
	}
	if len(ops) == 1 {
		return ops[0]
	}
	return exec.NewConcat(ops)
}

// Strings takes a list of inputs and always returns a single operation
// More than one rule will return a Select (similar to a multi-push rule)
// Less then one will return an empty Literal
func Strings(inputs []string) exec.Operation {
	ops := []exec.Operation{}

	for _, input := range inputs {
		op := String(input)
		ops = append(ops, op)
	}

	if len(ops) == 0 {
		return exec.NewLiteral("")
	}
	if len(ops) == 1 {
		return ops[0]
	}

	return exec.NewSelect(ops)
}

func parseAction(s *scan.Scanner) exec.Operation {
	// @cleanup: whole lot of assume the input is valid here
	// Consume opening [
	s.Next()
	keyToken := s.Next()
	key := keyToken.Value
	// Consume :
	s.Next()
	var ops []exec.Operation

	var ruleParts []string
	for {
		rawValue := s.Next()
		if rawValue.Type == scan.RightBracket {
			op := String(strings.Join(ruleParts, ""))
			ops = append(ops, op)
			break
		}
		if rawValue.Value == "POP" {
			// Consume closing ]
			s.Next()
			return exec.NewPop(key)
		}
		if rawValue.Type == scan.Comma {
			op := String(strings.Join(ruleParts, ""))
			ops = append(ops, op)
			ruleParts = ruleParts[:0]
			continue
		}

		ruleParts = append(ruleParts, rawValue.Value)
	}
	if len(ops) == 0 {
		return exec.NewPush(key, exec.NewLiteral(""))
	}
	if len(ops) == 1 {
		return exec.NewPush(key, ops[0])
	}

	return exec.NewPush(key, exec.NewSelect(ops))
}

func parseTag(s *scan.Scanner) exec.Operation {
	// @cleanup: whole lot of assume the input is valid here
	// Consume opening #
	s.Next()
	key := s.Next().Value
	var mods []exec.ModCall

	var ruleParts []string
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

		ruleParts = append(ruleParts, rawValue.Value)
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

	var ops []exec.Operation

	var ruleParts []string
	for {
		rawValue := s.Next()
		if rawValue.Type == scan.RightParen {
			op := String(strings.Join(ruleParts, ""))
			ops = append(ops, op)
			break
		}
		if rawValue.Type == scan.Comma {
			op := String(strings.Join(ruleParts, ""))
			ops = append(ops, op)
			ruleParts = ruleParts[:0]
			continue
		}

		ruleParts = append(ruleParts, rawValue.Value)
	}

	return exec.NewModCall(key, ops)
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
