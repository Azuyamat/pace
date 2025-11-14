package config

// Lexer tokenizes input text for the configuration parser
type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
	line         int
	column       int
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 1,
	}
	l.readChar()
	return l
}

// GetInput returns the original input string for error reporting
func (l *Lexer) GetInput() string {
	return l.input
}

func (l *Lexer) Position() (line, column int) {
	return l.line, l.column
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0 // EOF
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // EOF
	}
	return l.input[l.readPosition]
}

func (l *Lexer) isAtEnd() bool {
	return l.char == 0
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	line := l.line
	column := l.column

	var token Token

	switch l.char {
	case '{':
		token = l.makeSingleCharToken(TOKEN_LBRACE, line, column)
	case '}':
		token = l.makeSingleCharToken(TOKEN_RBRACE, line, column)
	case '[':
		token = l.makeSingleCharToken(TOKEN_LBRACKET, line, column)
	case ']':
		token = l.makeSingleCharToken(TOKEN_RBRACKET, line, column)
	case ',':
		token = l.makeSingleCharToken(TOKEN_COMMA, line, column)
	case '"':
		token = l.scanString(line, column)
	case '#':
		token = l.scanComment(line, column)
	case '\n':
		token = l.scanNewline(line, column)
	case '\r':
		l.readChar()
		return l.NextToken()
	case 0:
		token = NewEOFToken(line, column)
		l.readChar()
		return token
	default:
		if isLetter(l.char) {
			return l.scanIdentifier(line, column)
		} else if isDigit(l.char) {
			return l.scanNumber(line, column)
		} else {
			token = l.makeSingleCharToken(TOKEN_ILLEGAL, line, column)
		}
	}

	return token
}

func (l *Lexer) makeSingleCharToken(tokenType TokenType, line, column int) Token {
	token := newToken(tokenType, l.char, line, column)
	l.readChar()
	return token
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' {
		l.column++
		l.readChar()
	}
}

func (l *Lexer) scanIdentifier(line, column int) Token {
	literal := l.readIdentifier()
	tokenType := lookupIdent(literal)
	return NewTokenWithLiteral(tokenType, literal, line, column)
}

func (l *Lexer) scanNumber(line, column int) Token {
	literal := l.readNumber()
	return NewTokenWithLiteral(TOKEN_NUMBER, literal, line, column)
}

func (l *Lexer) scanString(line, column int) Token {
	literal := l.readString()
	l.readChar() // consume the closing quote
	return NewTokenWithLiteral(TOKEN_STRING, literal, line, column)
}

func (l *Lexer) scanComment(line, column int) Token {
	literal := l.readComment()
	return NewTokenWithLiteral(TOKEN_COMMENT, literal, line, column)
}

func (l *Lexer) scanNewline(line, column int) Token {
	token := newToken(TOKEN_NEWLINE, l.char, line, column)
	l.line++
	l.column = 0
	l.readChar()
	return token
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readComment() string {
	position := l.position
	for l.char != '\n' && l.char != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func lookupIdent(ident string) TokenType {
	switch ident {
	case "true", "false":
		return TOKEN_BOOLEAN
	}
	return TOKEN_IDENTIFIER
}
