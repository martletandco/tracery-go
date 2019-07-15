package tracery

import "testing"

func TestScanSingle(t *testing.T) {
	var tests = []struct {
		input    string
		expected []Token
	}{
		{"", []Token{Token{Type: EOF, Value: ""}}},
		{"a", []Token{Token{Type: Word, Value: "a"}, Token{Type: EOF, Value: ""}}},
		{"🥝", []Token{Token{Type: Word, Value: "🥝"}, Token{Type: EOF, Value: ""}}},
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
		{`\a`, []Token{Token{Type: Word, Value: `a`}, Token{Type: EOF, Value: ""}}},
		{`\🥝`, []Token{Token{Type: Word, Value: `🥝`}, Token{Type: EOF, Value: ""}}},
		{`\[`, []Token{Token{Type: Word, Value: `[`}, Token{Type: EOF, Value: ""}}},
		{`\]`, []Token{Token{Type: Word, Value: `]`}, Token{Type: EOF, Value: ""}}},
		{`\(`, []Token{Token{Type: Word, Value: `(`}, Token{Type: EOF, Value: ""}}},
		{`\)`, []Token{Token{Type: Word, Value: `)`}, Token{Type: EOF, Value: ""}}},
		{`\\`, []Token{Token{Type: BackStroke, Value: `\`}, Token{Type: EOF, Value: ""}}},
		{`\:`, []Token{Token{Type: Word, Value: `:`}, Token{Type: EOF, Value: ""}}},
		{`\,`, []Token{Token{Type: Word, Value: `,`}, Token{Type: EOF, Value: ""}}},
		{`\#`, []Token{Token{Type: Word, Value: `#`}, Token{Type: EOF, Value: ""}}},
		{`\.`, []Token{Token{Type: Word, Value: `.`}, Token{Type: EOF, Value: ""}}},
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
		{" a ", []Token{Token{Type: WhiteSpace, Value: " "}, Token{Type: Word, Value: "a"}, Token{Type: WhiteSpace, Value: " "}, Token{Type: EOF, Value: ""}}},
		{"a b", []Token{Token{Type: Word, Value: "a"}, Token{Type: WhiteSpace, Value: " "}, Token{Type: Word, Value: "b"}, Token{Type: EOF, Value: ""}}},
		{"a\nb", []Token{Token{Type: Word, Value: "a"}, Token{Type: WhiteSpace, Value: "\n"}, Token{Type: Word, Value: "b"}, Token{Type: EOF, Value: ""}}},
		{`\[wow\]`, []Token{Token{Type: Word, Value: "[wow]"}, Token{Type: EOF, Value: ""}}},
		{`\#wee\#`, []Token{Token{Type: Word, Value: "#wee#"}, Token{Type: EOF, Value: ""}}},
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
	type TL []struct {
		Type
		string
	}
	var tests = []struct {
		input    string
		expected TL
	}{
		{"#a#", TL{{Octo, "#"}, {Word, "a"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b.c#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {Period, "."}, {Word, "c"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b()#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {LeftParen, "("}, {RightParen, ")"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b().c()#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {LeftParen, "("}, {RightParen, ")"}, {Period, "."}, {Word, "c"}, {LeftParen, "("}, {RightParen, ")"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b.c()#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {Period, "."}, {Word, "c"}, {LeftParen, "("}, {RightParen, ")"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b().c#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {LeftParen, "("}, {RightParen, ")"}, {Period, "."}, {Word, "c"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b(x)#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {LeftParen, "("}, {Word, "x"}, {RightParen, ")"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b(x).c(y)#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {LeftParen, "("}, {Word, "x"}, {RightParen, ")"}, {Period, "."}, {Word, "c"}, {LeftParen, "("}, {Word, "y"}, {RightParen, ")"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b.c(🥝)#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {Period, "."}, {Word, "c"}, {LeftParen, "("}, {Word, "🥝"}, {RightParen, ")"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b(z).c#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {LeftParen, "("}, {Word, "z"}, {RightParen, ")"}, {Period, "."}, {Word, "c"}, {Octo, "#"}, {EOF, ""}}},
		{"#a.b(#x#)#", TL{{Octo, "#"}, {Word, "a"}, {Period, "."}, {Word, "b"}, {LeftParen, "("}, {Octo, "#"}, {Word, "x"}, {Octo, "#"}, {RightParen, ")"}, {Octo, "#"}, {EOF, ""}}},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.input)
		for _, expected := range tt.expected {
			expectedToken := Token{Type: expected.Type, Value: expected.string}
			actual := scanner.Next()
			if actual != expectedToken {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, expectedToken, actual)
				break
			}
		}
	}
}

func TestScanAction(t *testing.T) {
	type TL []struct {
		Type
		string
	}
	var tests = []struct {
		input    string
		expected TL
	}{
		{"[a:b]", TL{{LeftBracket, "["}, {Word, "a"}, {Colon, ":"}, {Word, "b"}, {RightBracket, "]"}, {EOF, ""}}},
		{"[a:b,c]", TL{{LeftBracket, "["}, {Word, "a"}, {Colon, ":"}, {Word, "b"}, {Comma, ","}, {Word, "c"}, {RightBracket, "]"}, {EOF, ""}}},
		{"[a:#b#]", TL{{LeftBracket, "["}, {Word, "a"}, {Colon, ":"}, {Octo, "#"}, {Word, "b"}, {Octo, "#"}, {RightBracket, "]"}, {EOF, ""}}},
		{"[a:#b#,#c#]", TL{{LeftBracket, "["}, {Word, "a"}, {Colon, ":"}, {Octo, "#"}, {Word, "b"}, {Octo, "#"}, {Comma, ","}, {Octo, "#"}, {Word, "c"}, {Octo, "#"}, {RightBracket, "]"}, {EOF, ""}}},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.input)
		for _, expected := range tt.expected {
			expectedToken := Token{Type: expected.Type, Value: expected.string}
			actual := scanner.Next()
			if actual != expectedToken {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, expectedToken, actual)
				break
			}
		}
	}
}

func TestScanComplex(t *testing.T) {
	type TL []struct {
		Type
		string
	}
	var tests = []struct {
		input    string
		expected TL
	}{
		{"#[a:b]c.d(e,#f#)#", TL{{Octo, "#"}, {LeftBracket, "["}, {Word, "a"}, {Colon, ":"}, {Word, "b"}, {RightBracket, "]"}, {Word, "c"}, {Period, "."}, {Word, "d"}, {LeftParen, "("}, {Word, "e"}, {Comma, ","}, {Octo, "#"}, {Word, "f"}, {Octo, "#"}, {RightParen, ")"}, {Octo, "#"}, {EOF, ""}}},
	}

	for _, tt := range tests {
		scanner := newScanner(tt.input)
		for _, expected := range tt.expected {
			expectedToken := Token{Type: expected.Type, Value: expected.string}
			actual := scanner.Next()
			if actual != expectedToken {
				t.Errorf("parse(%v): expected %v, actual %v", tt.input, expectedToken, actual)
				break
			}
		}
	}
}
