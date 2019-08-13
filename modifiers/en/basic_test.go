package basic

import "testing"

func TestAppendIndefArticle(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"apple", "an apple"},
		{"emu", "an emu"},
		{"inter island ferry", "an inter island ferry"},
		{"old shoe", "an old shoe"},
		{"umbrella", "an umbrella"},
		{"one-year-old ham", "a one-year-old ham"},
		{"united group", "a united group"},

		{"helicopter", "a helicopter"},
		{"hour nap", "an hour nap"},

		{"7", "a 7"},
		{"🥝", "a 🥝"},
		{"kiwi 🥝", "a kiwi 🥝"},
	}

	for _, tt := range tests {
		actual := AppendIndefArticle(tt.input)
		if actual != tt.expected {
			t.Errorf("AppendIndefArticle(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}

func TestCapitalise(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"🥝fruit", "🥝fruit"},
		{"Kiwi 🥝", "Kiwi 🥝"},
	}

	for _, tt := range tests {
		actual := CapitaliseFirst(tt.input)
		if actual != tt.expected {
			t.Errorf("CapitaliseFirst(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}

func TestReplace(t *testing.T) {
	var tests = []struct {
		input       string
		search      string
		replacement string
		expected    string
	}{
		{"", "", "", ""},
		{"a", "a", "b", "b"},
		{"hello", "he", "o", "ollo"},
	}

	for _, tt := range tests {
		actual := Replace(tt.input, tt.search, tt.replacement)
		if actual != tt.expected {
			t.Errorf("Replace(%v): expected %v, actual %v", tt.input, tt.expected, actual)
		}
	}
}
