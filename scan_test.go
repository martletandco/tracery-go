package tracery

import "testing"

func TestScanSingle(t *testing.T) {
	var tests = []struct {
		input    string
		expected []Token
	}{
		{"", []Token{Token{Type: EOF, Value: ""}}},
		{"a", []Token{Token{Type: Text, Value: "a"}, Token{Type: EOF, Value: ""}}},
		{"🥝", []Token{Token{Type: Text, Value: "🥝"}, Token{Type: EOF, Value: ""}}},
		{" ", []Token{Token{Type: WhiteSpace, Value: " "}, Token{Type: EOF, Value: ""}}},
		{"\n", []Token{Token{Type: WhiteSpace, Value: "\n"}, Token{Type: EOF, Value: ""}}},
		{"\t", []Token{Token{Type: WhiteSpace, Value: "\t"}, Token{Type: EOF, Value: ""}}},
		{"[", []Token{Token{Type: LeftBracket, Value: "["}, Token{Type: EOF, Value: ""}}},
		{"]", []Token{Token{Type: RightBracket, Value: "]"}, Token{Type: EOF, Value: ""}}},
		{"(", []Token{Token{Type: LeftParen, Value: "("}, Token{Type: EOF, Value: ""}}},
		{")", []Token{Token{Type: RightParen, Value: ")"}, Token{Type: EOF, Value: ""}}},
		{":", []Token{Token{Type: Colon, Value: ":"}, Token{Type: EOF, Value: ""}}},
		{",", []Token{Token{Type: Comma, Value: ","}, Token{Type: EOF, Value: ""}}},
		{"#", []Token{Token{Type: Octo, Value: "#"}, Token{Type: EOF, Value: ""}}},
		{".", []Token{Token{Type: Period, Value: "."}, Token{Type: EOF, Value: ""}}},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.input)
		for _, expected := range tt.expected {
			actual := scanner.Next()
			if actual != expected {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, expected, actual)
				break
			}
		}
	}
}

func TestScanSingleEscaped(t *testing.T) {
	var tests = []struct {
		input    string
		expected []Token
	}{
		{`\a`, []Token{Token{Type: Text, Value: `a`}, Token{Type: EOF, Value: ""}}},
		{`\🥝`, []Token{Token{Type: Text, Value: `🥝`}, Token{Type: EOF, Value: ""}}},
		{`\[`, []Token{Token{Type: Text, Value: `[`}, Token{Type: EOF, Value: ""}}},
		{`\]`, []Token{Token{Type: Text, Value: `]`}, Token{Type: EOF, Value: ""}}},
		{`\(`, []Token{Token{Type: Text, Value: `(`}, Token{Type: EOF, Value: ""}}},
		{`\)`, []Token{Token{Type: Text, Value: `)`}, Token{Type: EOF, Value: ""}}},
		{`\\`, []Token{Token{Type: BackStroke, Value: `\`}, Token{Type: EOF, Value: ""}}},
		{`\:`, []Token{Token{Type: Text, Value: `:`}, Token{Type: EOF, Value: ""}}},
		{`\,`, []Token{Token{Type: Text, Value: `,`}, Token{Type: EOF, Value: ""}}},
		{`\#`, []Token{Token{Type: Text, Value: `#`}, Token{Type: EOF, Value: ""}}},
		{`\.`, []Token{Token{Type: Text, Value: `.`}, Token{Type: EOF, Value: ""}}},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.input)
		for _, expected := range tt.expected {
			actual := scanner.Next()
			if actual != expected {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, expected, actual)
				break
			}
		}
	}
}

func TestScanText(t *testing.T) {
	var tests = []struct {
		input    string
		expected []Token
	}{
		{" a ", []Token{Token{Type: WhiteSpace, Value: " "}, Token{Type: Text, Value: "a"}, Token{Type: WhiteSpace, Value: " "}, Token{Type: EOF, Value: ""}}},
		{"a\nb", []Token{Token{Type: Text, Value: "a"}, Token{Type: WhiteSpace, Value: "\n"}, Token{Type: Text, Value: "b"}, Token{Type: EOF, Value: ""}}},
		{`\[wow\]`, []Token{Token{Type: LeftBracket, Value: "[wow]"}, Token{Type: EOF, Value: ""}}},
		{`\#wee\#`, []Token{Token{Type: Octo, Value: "#wee#"}, Token{Type: EOF, Value: ""}}},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.input)
		for _, expected := range tt.expected {
			actual := scanner.Next()
			if actual != expected {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, expected, actual)
				break
			}
		}
	}
}

func TestScanSymbol(t *testing.T) {
	var tests = []struct {
		input    string
		expected []Token
	}{
		{"#a#", []Token{Token{Type: Octo, Value: "#"}, Token{Type: Identifier, Value: "a"}, Token{Type: Octo, Value: "#"}, Token{Type: EOF, Value: ""}}},
		{"#a.b#", []Token{Token{Type: Octo, Value: "#"}, Token{Type: Identifier, Value: "a"}, Token{Type: Period, Value: "."}, Token{Type: Identifier, Value: "b"}, Token{Type: Octo, Value: "#"}, Token{Type: EOF, Value: ""}}},
		{"#a.b()#", []Token{Token{Type: Octo, Value: "#"}, Token{Type: Identifier, Value: "a"}, Token{Type: Period, Value: "."}, Token{Type: Identifier, Value: "b"}, Token{Type: LeftParen, Value: "("}, Token{Type: RightParen, Value: ")"}, Token{Type: Octo, Value: "#"}, Token{Type: EOF, Value: ""}}},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.input)
		for _, expected := range tt.expected {
			actual := scanner.Next()
			if actual != expected {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, expected, actual)
				break
			}
		}
	}
}
