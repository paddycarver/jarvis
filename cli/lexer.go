package cli

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

type tokenType int // used to distinguish between types of tokens

func (t tokenType) String() string {
	switch t {
	case tokenError:
		return "Error"
	case tokenEOF:
		return "EOF"
	case tokenCommand:
		return "Command"
	case tokenParam:
		return "Parameter"
	case tokenArg:
		return "Argument"
	}
	return "Unknown token"
}

const (
	tokenError tokenType = iota
	tokenEOF
	tokenCommand
	tokenParam
	tokenArg
)

const (
	eof = unicode.ReplacementChar // we'll use this to signal the end of the input. May cause bugs if a non-unicode character is entered.
)

var (
	argPrefix, _ = utf8.DecodeRune([]byte("-"))
)

// token is a single conceptual part of the command
type token struct {
	code  tokenType
	value string
}

func (t token) String() string {
	if t.code == tokenEOF {
		return "EOF"
	}
	return fmt.Sprintf("%s: %s", t.code, t.value)
}

// stateFunc is a recursive function definition, it's going to power our state machine
type stateFunc func(*lexer) stateFunc

// lexer keeps track of our state
type lexer struct {
	input  string
	start  int
	pos    int
	width  int
	state  stateFunc
	tokens chan token
}

// NewLexer exposes lexers to other packages
func NewLexer(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan token, 2), // tokens is a buffered channel; we'll be picking up tokens as we put them on, so only two need to be held at any given time
		state:  lexCommand,          // all input should begin with a command
	}
	return l
}

// nextToken steps the lexer to the next token in the input
func (l *lexer) nextToken() token {
	for {
		select {
		case token := <-l.tokens: // if we have a token, return it
			return token
		default:
			l.state = l.state(l) // otherwise, generate a token using the state machine
		}
	}
	panic("Should never be reached.")
}

// Lex tokenizes the entire input string
func (l *lexer) Lex() ([]token, error) {
	results := []token{}
	t := l.nextToken()
	for t.code > tokenEOF { // as long as the token is not an error or EOF, we want it
		results = append(results, t)
		t = l.nextToken()
	}
	if t.code == tokenError { // if there's an error token, throw an error
		return results, errors.New(t.value)
	}
	return results, nil
}

// emit is a helper function that will return tokens to the lexer
func (l *lexer) emit(code tokenType) {
	l.tokens <- token{code, l.input[l.start:l.pos]} // send the token to l.tokens for nextToken's select to pick up
	l.start = l.pos                                 // advance our position in the input
}

// next returns the next unicode rune in the input
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// ignore ignores everything since the last token was emitted, up to the current position
func (l *lexer) ignore() {
	l.start = l.pos
}

// backup steps back a single unicode rune
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns the next unicode rune without changing our position
func (l *lexer) peek() rune {
	r := l.next()
	l.backup() // we use this instead of having next call peek, as next needs the width return from utf8.DecodeRuneInString
	return r
}

// errorf is a helper function for returning errors while tokenizing
func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	// errors are just tokens with code of tokenError and value of the error message
	l.tokens <- token{
		tokenError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

// lexCommand is run when the next input is expected to be a command.
// It returns the next stateFunc that should be run to process the next input
func lexCommand(l *lexer) stateFunc {
	for {
		switch r := l.next(); {
		case r == eof: // we're at the end of the input
			if l.pos > l.start {
				l.emit(tokenCommand) // if we hit the end of input and there are runes entered, emit the command
			}
			l.emit(tokenEOF)
			return nil // terminate lexing
		case unicode.IsSpace(r): // spaces terminate commands
			l.backup() // we don't want to include the space rune in the token
			if l.pos > l.start {
				l.emit(tokenCommand) // only emit the command if there are runes in it
			} else {
				return l.errorf("No command entered.") // if there are no runes, emit an error instead
			}
			return lexSpace // we know the command is followed by a space, so return a stateFunc that knows how to deal with that
		}
	}
	panic("should never be reached.")
}

// lexSpace is run when the next input is expected to be or known to be a space.
// It figures out which stateFunc to hand off to next.
func lexSpace(l *lexer) stateFunc {
	l.next()   // step forward one rune (the space)
	l.ignore() // ignore the space
	switch r := l.next(); {
	case r == eof:
		l.emit(tokenEOF) // if our next rune is the end of input, emit the eoft token
		return nil       // stop the lexing process
	case unicode.IsSpace(r):
		return lexSpace // two spaces in a row, try again!
	case r == argPrefix:
		return lexArg // it has the argument prefix, it must be an argument!
	default:
		l.backup()      // we've already consumed the first character, but we need it, so back up
		return lexParam // it's not a command, a space, an argument... it must be a parameter!
	}
	panic("Should never be reached.")
}

// lexArg is run when the next input is expected to be or known to be an argument.
// It consumes the argument, then hands off to the next stateFunc.
// Arguments begin with argPrefix and are terminated by a space, unless a " is found.
// When a " is present in the argument, the argument is terminated by a ".
func lexArg(l *lexer) stateFunc {
	isTerminator := func(r rune) bool {
		return unicode.IsSpace(r)
	}
	l.ignore() // ignore the prefix
	if l.peek() == argPrefix {
		l.next()
		l.ignore() // ignore the prefix, if repeated (e.g. --)
	}
	for {
		switch r := l.next(); {
		case r == eof: // if we hit the end of input, we know we didn't find the terminator
			if !isTerminator(argPrefix) {
				// if argPrefix isn't the terminator, a space is the terminator, so just pretend we found it
				l.emit(tokenArg)
				l.emit(tokenEOF)
				return nil
			}
			return l.errorf("Argument ended prematurely, expected a close quote.")
		case isTerminator(r):
			if unicode.IsSpace(r) {
				l.backup() // we don't want the character in the token if it's a space
			}
			l.emit(tokenArg)
			if unicode.IsSpace(r) {
				// we know the next character is a space, so let's return a stateFunc to handle it
				return lexSpace
			}
			// if we're this far, the terminator character was "
			nextRune := l.peek()
			if nextRune == eof {
				l.emit(tokenEOF)
				return nil
			}
			if !unicode.IsSpace(nextRune) {
				l.next() // adjust our position, so debuggers can get the right position for the offending character
				return l.errorf("Bad argument syntax: expected a space or end of input after the closing quote.")
			}
			return lexSpace
    default:
      if r == '"' {
        isTerminator = func(r rune) bool {
          return r == '"'
        }
      }
		}
	}
	panic("Should never be reached.")
}

// lexParam is run when the next input is expected to be or known to be a parameter.
// It consumes the parameter, then hands off to the next stateFunc.
// Parameters that begin with a " are terminated by a ". All other parameters are terminated by a space.
func lexParam(l *lexer) stateFunc {
	isTerminator := func(r rune) bool {
		return unicode.IsSpace(r)
	}
	first := l.next()
	if first == '"' {
		isTerminator = func(r rune) bool {
			return r == '"'
		}
	} else {
		l.backup()
	}
	for {
		switch r := l.next(); {
		case r == eof:
			if !isTerminator(argPrefix) {
				// if the argPrefix isn't the terminator, a space is, so end of input is not an error
				l.emit(tokenParam)
				l.emit(tokenEOF)
				return nil
			}
			return l.errorf("Parameter ended prematurely, expected a close quote.")
		case isTerminator(r):
			l.backup()
			l.emit(tokenParam)
			if isTerminator('"') {
				l.next()
				l.ignore()
			}
			return lexSpace
		}
	}
	panic("Should never be reached.")
}
