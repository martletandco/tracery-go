/**

The scanner's job is to turn a input string into a list of tokens for the parser to deal with

Does:
- Escape chars

Does not:
- Enforce gramatic rules e.g. is fine with #a.#('?]

*/

package tracery

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Token found while scanning rules
type Token struct {
	Type  Type
	Value string
}

// Type of emitted token
type Type int

const (
	EOF Type = iota // End of file/line
	Error
	Word         // Any non-control grammar value, could be an ident for a symbol
	WhiteSpace   // \n \t \s, etc
	LeftBracket  // [
	RightBracket // ]
	LeftParen    // (
	RightParen   // )
	BackStroke   // \
	Colon        // :
	Comma        // ,
	Octo         // #
	Period       // .
)

// Pretend to be a rune but signal instead
const eof = -1

type stateFunc func(*Scanner) stateFunc

type Scanner struct {
	input  string
	start  int
	end    int
	size   int
	state  stateFunc
	tokens []Token
}

func newScanner(input string) *Scanner {
	return &Scanner{input: input, state: lexAny, tokens: []Token{}}
}

func (s *Scanner) Peek() Token {
	for {
		if len(s.tokens) > 0 {
			return s.tokens[0]
		}
		if s.state == nil {
			return Token{Type: EOF}
		}
		s.state = s.state(s)
	}
}

func (s *Scanner) Next() Token {
	token := s.Peek()

	if len(s.tokens) > 0 {
		s.tokens = s.tokens[1:]
	}

	return token
}

func (s *Scanner) next() rune {
	if s.end == len(s.input) {
		return eof
	}
	r, size := utf8.DecodeRuneInString(s.input[s.end:])
	s.size = size
	return r
}
func (s *Scanner) consume() {
	s.end += s.size
	s.size = 0
}

func (s *Scanner) emit(t Type) {
	value := s.input[s.start:s.end]
	if t == Word {
		// BackStroke is only used in words as an escape, so we are cleaning up here
		value = strings.Replace(value, `\`, "", int(-1))
		// ignore empty values
		if len(value) == 0 {
			s.start = s.end
			return
		}
	}
	s.tokens = append(s.tokens, Token{Type: t, Value: value})
	s.start = s.end
}

func lexAny(s *Scanner) stateFunc {
	var nextState stateFunc
	// Keep grabbing runes until we hit a control character or EOF
	// @cleanup: use strings.Index to find control chars without a loop
Loop:
	for {
		r := s.next()
		switch {
		// eof
		case r == eof:
			break Loop
		// tag
		case r == '#':
			nextState = lexChar(Octo)
			break Loop
		// action
		case r == '[':
			// nextState = lexAction
			nextState = lexChar(LeftBracket)
			break Loop
		case r == ']':
			// nextState = lexAction
			nextState = lexChar(RightBracket)
			break Loop

		// Sometimes these are control chars
		// left paren
		case r == '(':
			nextState = lexChar(LeftParen)
			break Loop
		// right paren
		case r == ')':
			nextState = lexChar(RightParen)
			break Loop
		// colon
		case r == ':':
			nextState = lexChar(Colon)
			break Loop
		// comma
		case r == ',':
			nextState = lexChar(Comma)
			break Loop
		// period
		case r == '.':
			nextState = lexChar(Period)
			break Loop

		// Whitespace
		case unicode.IsSpace(r):
			nextState = lexChar(WhiteSpace)
			break Loop

		// Escaped chars
		case r == '\\':
			s.consume()
			switch s.next() {
			case eof:
				break Loop
			// Escaped backstroke
			case '\\':
				nextState = lexChar(BackStroke)
				break Loop
			}
			s.consume()
		default:
			s.consume()
		}
	}

	// Because lexAny doubles as our entry we might not have any text here
	if s.end > s.start {
		s.emit(Word)
	}

	return nextState
}

func lexChar(t Type) stateFunc {
	return func(s *Scanner) stateFunc {
		s.consume()
		s.emit(t)
		return lexAny
	}
}
