package compiler

import (
	"fmt"
	"log"
	"unicode/utf8"
)

// Lexer is the data structure that represents the transition state machine.
// it functions as container with input text string and other components, for
// example lex stores a SymbolTable to store the symbols/identifiers and a
// current state function pointer. It also keep tracks of current lex/scan
// head position (for roll back character should the state is mismatched).
//
// Most importantly the lexer outputs tokens into a channel, a buffered CSP
// message passing mechanism, to facilitate an asynchronous concurrent
// capability for future parser to consume.
type Lexer struct {
	S               *SymbolTable
	In              string        // Input string
	State           StateFunction // current state
	CurrentPosition int           // current position in the input
	StartPosition   int           // start position of this token
	Width           int           // width of last rune read from input
	Tokens          chan Token    // channel of scanned tokens
	LastPosition    int           // position of most recent token consumed by client
	//parenthesisDepth int // nesting depth of parenthesis
}

// LegalWords contains uppercase opening legal keywords, all other uppercase
// opening words are deemed illegal.
var LegalWords = []string{"Private", "Public", "Protected", "Static",
	"Primary", "Integer", "Exception", "Try"}

func (l *Lexer) buildLegalWords() {
	for _, v := range LegalWords {
		l.S.Process(v, Dynamic)
	}
}

// special characters and constant values
const (
	WhitespaceChars = " \t\r\n"
	EOF             = -1
	EscapeRune      = 126 // '~'
	QuoteRune       = 34  // '"'
	LparRune        = 40  // '('
	RparRune        = 41  // ')'
	SemicolonRune   = 59  // ';'
)

// Next returns the next rune in the input.
func (l *Lexer) Next() rune {
	if int(l.CurrentPosition) >= len(l.In) {
		l.Width = 0
		return EOF
	}
	r, w := utf8.DecodeRuneInString(l.In[l.CurrentPosition:])
	l.Width = w
	l.CurrentPosition += l.Width
	return r
}

// Peek returns but does not consume the next rune in the input.
func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

// Backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) Backup() {
	l.CurrentPosition -= l.Width
}

// Emit passes a token back to the client.
// the variadic e is the array of indices of EscapeRune within the token, they
// need to be escaped.
func (l *Lexer) Emit(t Token) {
	l.Tokens <- t
	l.Ignore()
}

// Ignore skips over the pending input before this point.
func (l *Lexer) Ignore() {
	l.StartPosition = l.CurrentPosition
}

// Errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.ConsumeToken.
func (l *Lexer) Errorf(format string, args ...interface{}) StateFunction {
	l.Tokens <- Token{TokenError, l.StartPosition,
		fmt.Sprintf(format, args...)}
	return nil
}

// Token presents a token or text string returned from the scanner.
type Token struct {
	T tokenType // the type of this Token
	P int       // the starting byte postion of this Token in the input string
	V string    // the value of this Token
}

type tokenType int

// Token Types
const (
	TokenError            tokenType = iota // error occurred;
	TokenInteger                           // integer constant
	TokenKeyword                           // uppercase keywords
	TokenIdentifier                        // lowercase identifier
	TokenLeftParenthesis                   // left parenthesis
	TokenRightParenthesis                  // right parenthesis
	TokenSemicolon                         // ;
	TokenText                              // string constant
)

// StateFunction is a function pointer definition, to replace switch-case
// implementation of the state machine.
type StateFunction func(*Lexer) StateFunction

// Run is a main processing routine which calls a state machine, implemented
// as a rolling function pointer, which reduces code complexity of a switch
// case implementation
func (l *Lexer) Run() {
	for l.State = startState; l.State != nil; {
		l.State = l.State(l)
	}
	close(l.Tokens)
}

func startState(l *Lexer) StateFunction {
	for r := l.Next(); r != EOF; r = l.Next() {
		if isWhitespaces(r) {
			l.Ignore()
			continue
		}
		if r == QuoteRune {
			// not rewinding
			return lexQuotedText
		}
		if r == LparRune {
			l.Emit(Token{TokenLeftParenthesis, l.StartPosition, "0"})
			continue
		}
		if r == RparRune {
			l.Emit(Token{TokenRightParenthesis, l.StartPosition, "0"})
			continue
		}
		if r == SemicolonRune {
			l.Emit(Token{TokenSemicolon, l.StartPosition, "0"})
			continue
		}
		if isNumber(r) {
			l.Backup()
			return lexInteger
		}
		if isCharacter(r) {
			l.Backup()
			return lexIdentifier
		}
		// if reaches here, error
		l.Emit(Token{TokenError, l.StartPosition,
			fmt.Sprintf("no matching state for rune %q", r)})
		l.Ignore()
	}
	// reached EOF
	return nil
}

func lexInteger(l *Lexer) StateFunction {
	var rst int
	const maxint = 65535
	for {
		r := l.Next()
		if r == EOF {
			break
		}
		if i := int(r - '0'); i >= 0 && i <= 9 {
			if rst <= 6553 && rst > -1 {
				rst *= 10
				if maxint-rst < i {
					rst = -1
				} else {
					rst += i
				}
			} else {
				rst = -1
				// not break, instead skip until transition
			}
		} else {
			// potentially a state transition here
			l.Backup()
			break
		}
	}
	l.Emit(Token{TokenInteger, l.StartPosition, fmt.Sprintf("%d", rst)})
	return startState
}

func lexIdentifier(l *Lexer) StateFunction {
	var flag FlagVar
	r := l.Next()
	if r >= 'A' && r <= 'Z' {
		flag = Static
	} else if r >= 'a' && r <= 'z' {
		flag = Dynamic
	} else {
		return l.Errorf("first character is not a a-zA-Z in identifier state")
	}
	for r := l.Next(); r != EOF && r >= 'a' && r <= 'z'; r = l.Next() {
	}
	l.Backup()
	v, e := l.S.Process(l.In[l.StartPosition:l.CurrentPosition], flag)
	if e != nil {
		log.Fatalf("Trie reports error: %v", e)
	}
	l.Emit(Token{TokenIdentifier, l.StartPosition, fmt.Sprintf("%d", v)})
	return startState
}

// BUG(n) test this properly against quoted text objects
func lexQuotedText(l *Lexer) StateFunction {
	var escapes []int // array of escape charactor indices
	for {
		r := l.Next()
		if r == EOF {
			// error
			break
		}
		if r == EscapeRune {
			//escape charactor, advancing 1 rune
			escapes = append(escapes, l.CurrentPosition-1)
			l.Next()
			continue
		}
		if r == QuoteRune {
			t := textToken(l, escapes)
			l.Emit(t)
			return startState
		}
	}
	l.Emit(Token{TokenError, l.StartPosition,
		fmt.Sprintf("not matched \" found for string token, position %d",
			l.CurrentPosition-2)})
	return startState
}

func textToken(l *Lexer, e []int) Token {
	value := ""
	curr := l.StartPosition + 1
	for _, idx := range e {
		if idx >= l.CurrentPosition {
			// panic because it should never happen
			panic("EscapeRune idx is larger than text end")
		}
		value += l.In[curr:idx]
		curr = idx + 1
	}
	value += l.In[curr : l.CurrentPosition-1] // minus the endQuote
	return Token{TokenText, l.StartPosition, value}
}

func isWhitespaces(r rune) bool {
	for _, v := range WhitespaceChars {
		if r == v {
			return true
		}
	}
	return false
}

func isNumber(r rune) bool {
	i := int(r - '0')
	return i <= 9 && i >= 0
}

func isCharacter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// NewLexer is the constructor for Lexer, which creates a symboltable and load
// LegalWords as a "static trie"
func NewLexer(input string) *Lexer {
	pointer := &Lexer{S: NewSymbolTable(), In: input,
		Tokens: make(chan Token)}
	pointer.buildLegalWords()
	go pointer.Run()
	return pointer
}
