package config

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '_'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t'
}

func isNewline(char byte) bool {
	return char == '\n'
}

func isCarriageReturn(char byte) bool {
	return char == '\r'
}
