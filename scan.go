package tracery

import ()
import "unicode/utf8"

// Token found while scanning rules
type Token struct {
	Type  Type
	Value string
}

// Type of emitted token
type Type int

// any text then #symbol# or even [else:something] something #else.wrap(')#?
// LexText LexOpeningOcto LexInside LexIdent LexInside LexClosingOcto LexText
// Text(any then then ) Octo Identifier(symbol) Octo Text( or even ) RightBracket Identifier(else) Colon

// ident is any chars excluding control and whitespace?
// maybe we don't have idents at the scanner/lexer level

// any -> openOcto, leftBracket, colon, comma, period, closeOcto, rightBracket, rightParen, text
// openOcto -> openOcto, leftBracket, ident
// leftBracket -> openOcto, ident
// colon -> text, openOcto, ident(POP)
// comma -> text, openOcto

// text -> [...control chars, space]
//

const (
	EOF Type = iota // End of file/line
	Error
	// NewLine
	Text         // Any non-control grammar value, could be an ident for a symbol
	Identifier   // Symbol key
	WhiteSpace   // \n\t\s @enhance: get a more complete definition of whitespace
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
	return &Scanner{input: input, state: lexText, tokens: []Token{}}
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
	s.tokens = append(s.tokens, Token{Type: t, Value: value})
	s.start = s.end
}

func lexText(s *Scanner) stateFunc {
	var nextState stateFunc = nil
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
			nextState = lexOpeningOcto
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
		case r == ' ':
			fallthrough
		case r == '\n':
			fallthrough
		case r == '\t':
			nextState = lexChar(WhiteSpace)
			break Loop

		// Escaped chars
		case r == '\\':
			s.consume()
			if s.next() == eof {
				break Loop
			}
			s.consume()
		default:
			s.consume()
		}
	}

	// Because lexText doubles as our entry we might not have any text here
	if s.end > s.start {
		s.emit(Text)
	}

	return nextState
}

func lexChar(t Type) stateFunc {
	return func(s *Scanner) stateFunc {
		s.consume()
		s.emit(t)
		return lexText
	}
}

func lexOpeningOcto(s *Scanner) stateFunc {
	// consume opening #
	s.consume()
	s.emit(Octo)

	return lexIdent
}

func lexIdent(s *Scanner) stateFunc {
	var nextState stateFunc = nil
	// Keep grabbing runes until we hit a control character or EOF
	// @cleanup: Some type of consume until significant char
	// const index = strings.IndexAny("[]().,:#\n\t ")
Loop:
	for {
		r := s.next()
		switch {
		// eof
		case r == eof:
			break Loop
		// tag
		case r == '#':
			nextState = lexClosingOcto
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
		case r == ' ':
			fallthrough
		case r == '\n':
			fallthrough
		case r == '\t':
			nextState = lexChar(WhiteSpace)
			break Loop

		// Escaped chars
		case r == '\\':
			s.consume()
			if s.next() == eof {
				break Loop
			}
			s.consume()
		default:
			s.consume()
		}
	}

	if s.end > s.start {
		s.emit(Identifier)
	}

	return nextState
}

func lexClosingOcto(s *Scanner) stateFunc {
	s.consume()
	s.emit(Octo)
	return lexText
}
