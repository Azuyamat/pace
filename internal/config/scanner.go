package config

type Scanner struct {
	input        string
	position     int
	readPosition int
	char         byte
	line         int
	column       int
}

func NewScanner(input string, position, readPosition int, char byte, line, column int) *Scanner {
	return &Scanner{
		input:        input,
		position:     position,
		readPosition: readPosition,
		char:         char,
		line:         line,
		column:       column,
	}
}

func (s *Scanner) ReadChar() (byte, int, int) {
	if s.readPosition >= len(s.input) {
		s.char = 0
	} else {
		s.char = s.input[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition++
	s.column++
	return s.char, s.position, s.readPosition
}

func (s *Scanner) PeekChar() byte {
	if s.readPosition >= len(s.input) {
		return 0
	}
	return s.input[s.readPosition]
}

func (s *Scanner) ScanIdentifier() string {
	position := s.position
	for isLetter(s.char) {
		s.ReadChar()
	}
	return s.input[position:s.position]
}

func (s *Scanner) ScanNumber() string {
	position := s.position
	for isDigit(s.char) {
		s.ReadChar()
	}
	return s.input[position:s.position]
}

func (s *Scanner) ScanString() string {
	position := s.position + 1
	for {
		s.ReadChar()
		if s.char == '"' || s.char == 0 {
			break
		}
	}
	return s.input[position:s.position]
}

func (s *Scanner) ScanMultilineString() string {
	// Skip the opening triple quotes
	s.ReadChar() // Skip second quote
	s.ReadChar() // Skip third quote
	s.ReadChar() // Move to first char of content

	position := s.position

	for {
		if s.char == 0 {
			break
		}

		// Check for closing triple quotes
		if s.char == '"' && s.PeekChar() == '"' && s.position+2 < len(s.input) && s.input[s.position+2] == '"' {
			result := s.input[position:s.position]
			s.ReadChar() // Skip first closing quote
			s.ReadChar() // Skip second closing quote
			// Third will be consumed by caller
			return result
		}

		if s.char == '\n' {
			s.line++
			s.column = 0
		}
		s.ReadChar()
	}

	return s.input[position:s.position]
}

func (s *Scanner) ScanComment() string {
	position := s.position
	for s.char != '\n' && s.char != 0 {
		s.ReadChar()
	}
	return s.input[position:s.position]
}

func (s *Scanner) SkipWhitespace() {
	for isWhitespace(s.char) {
		s.column++
		s.ReadChar()
	}
}

func (s *Scanner) GetState() (position, readPosition int, char byte, line, column int) {
	return s.position, s.readPosition, s.char, s.line, s.column
}

func (s *Scanner) SetState(position, readPosition int, char byte, line, column int) {
	s.position = position
	s.readPosition = readPosition
	s.char = char
	s.line = line
	s.column = column
}
