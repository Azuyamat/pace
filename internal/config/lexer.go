package config

type Lexer struct {
	scanner *Scanner
	input   string
}

func NewLexer(input string) *Lexer {
	scanner := NewScanner(input, 0, 0, 0, 1, 1)
	scanner.ReadChar()
	return &Lexer{
		scanner: scanner,
		input:   input,
	}
}

func (l *Lexer) GetInput() string {
	return l.input
}

func (l *Lexer) Position() (line, column int) {
	_, _, _, line, column = l.scanner.GetState()
	return line, column
}

func (l *Lexer) isAtEnd() bool {
	_, _, char, _, _ := l.scanner.GetState()
	return char == 0
}

func (l *Lexer) NextToken() Token {
	l.scanner.SkipWhitespace()

	_, _, char, line, column := l.scanner.GetState()

	var token Token

	switch char {
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
	case '$':
		token = l.makeSingleCharToken(TOKEN_DOLLAR, line, column)
	case '(':
		token = l.makeSingleCharToken(TOKEN_LPAREN, line, column)
	case ')':
		token = l.makeSingleCharToken(TOKEN_RPAREN, line, column)
	case '"':
		// Check for triple-quoted string
		if l.scanner.PeekChar() == '"' {
			_, _, nextChar, _, _ := l.scanner.GetState()
			l.scanner.ReadChar() // Consume second quote
			if l.scanner.PeekChar() == '"' {
				token = l.scanMultilineString(line, column)
			} else {
				// Two quotes in a row, treat as empty string followed by quote
				l.scanner.SetState(l.scanner.position-1, l.scanner.readPosition-1, nextChar, line, column)
				token = l.scanString(line, column)
			}
		} else {
			token = l.scanString(line, column)
		}
	case '#':
		token = l.scanComment(line, column)
	case '\n':
		token = l.scanNewline(line, column)
	case '\r':
		l.scanner.ReadChar()
		return l.NextToken()
	case 0:
		token = NewEOFToken(line, column)
		l.scanner.ReadChar()
		return token
	default:
		if isLetter(char) {
			return l.scanIdentifier(line, column)
		} else if isDigit(char) {
			return l.scanNumber(line, column)
		} else {
			token = l.makeSingleCharToken(TOKEN_ILLEGAL, line, column)
		}
	}

	return token
}

func (l *Lexer) makeSingleCharToken(tokenType TokenType, line, column int) Token {
	_, _, char, _, _ := l.scanner.GetState()
	token := newToken(tokenType, char, line, column)
	l.scanner.ReadChar()
	return token
}

func (l *Lexer) scanIdentifier(line, column int) Token {
	literal := l.scanner.ScanIdentifier()
	tokenType := lookupIdent(literal)
	return NewTokenWithLiteral(tokenType, literal, line, column)
}

func (l *Lexer) scanNumber(line, column int) Token {
	literal := l.scanner.ScanNumber()
	return NewTokenWithLiteral(TOKEN_NUMBER, literal, line, column)
}

func (l *Lexer) scanString(line, column int) Token {
	literal := l.scanner.ScanString()
	l.scanner.ReadChar()
	return NewTokenWithLiteral(TOKEN_STRING, literal, line, column)
}

func (l *Lexer) scanMultilineString(line, column int) Token {
	literal := l.scanner.ScanMultilineString()
	l.scanner.ReadChar()
	return NewTokenWithLiteral(TOKEN_MULTILINE_STRING, literal, line, column)
}

func (l *Lexer) scanComment(line, column int) Token {
	literal := l.scanner.ScanComment()
	return NewTokenWithLiteral(TOKEN_COMMENT, literal, line, column)
}

func (l *Lexer) scanNewline(line, column int) Token {
	_, _, char, _, _ := l.scanner.GetState()
	token := newToken(TOKEN_NEWLINE, char, line, column)

	_, _, _, currentLine, _ := l.scanner.GetState()
	l.scanner.SetState(l.scanner.position, l.scanner.readPosition, l.scanner.char, currentLine+1, 0)
	l.scanner.ReadChar()
	return token
}

func lookupIdent(ident string) TokenType {
	switch ident {
	case "true", "false":
		return TOKEN_BOOLEAN
	}
	return TOKEN_IDENTIFIER
}
