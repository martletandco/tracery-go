package basic

import (
	"strings"
	"unicode/utf8"
)

func isVowel(char string) bool {
	return char == "a" || char == "e" || char == "i" || char == "o" || char == "u"
}

func AppendIndefArticle(input string, params []string) string {
	if len(input) == 0 {
		return ""
	}

	_, size := utf8.DecodeRuneInString(input)
	first := input[:size]
	if isVowel(first) {
		return "an " + input
	}
	return "a " + input
}

func CapitaliseFirst(input string, params []string) string {
	_, size := utf8.DecodeRuneInString(input)
	return strings.ToUpper(input[:size]) + input[size:]
}

func Replace(input string, params []string) string {
	search, replacement := params[0], params[1]
	return strings.Replace(input, search, replacement, -1)
}
