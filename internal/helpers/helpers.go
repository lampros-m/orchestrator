package helpers

func IsOnlyLowercaseAndNumbersAndNotEmpty(s string) bool {
	if s == "" {
		return false
	}

	for _, letter := range s {
		// Check if letter is neither a lowercase letter nor a digit.
		if (letter < 'a' || letter > 'z') && (letter < '0' || letter > '9') {
			return false
		}
	}
	return true
}
