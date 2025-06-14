package core

// Character gradient from dark to light
// Default characters: "@%#*+=:~-.  "

// Chars represents a mapping of 256 brightness levels to ASCII characters.
// The characters are ordered from darkest (low brightness) to lightest (high brightness).
type Chars [256]byte

// NewChars creates a new character set for brightness-to-ASCII conversion.
// Accepts only ASCII characters ordered from darkest to lightest.
//
// Returns default character set ("@%#*+=:~-.  ") if:
//   - Input string is empty
//   - Input contains non-ASCII characters
//
// Example:
//
//	chars := NewChars(" .:-=+*#%@") // Light to dark gradient
func NewChars(chars string) *Chars {
	if len(chars) == 0 {
		return DefaultChars()
	}

	// Only ascii chars
	if len(chars) != len([]rune(chars)) {
		return DefaultChars()
	}

	bytes := Chars{}

	charsLen := len(chars)
	for brightness := 0; brightness < 256; brightness++ {
		idx := brightness * (charsLen - 1) / 255
		bytes[brightness] = chars[idx]
	}

	return &bytes
}

// DefaultChars returns the default character set: "@%#*+=:~-.  "
func DefaultChars() *Chars {
	return defaultChars
}

// defaultChars contains the predefined default character set
var defaultChars = &Chars{
	64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
	37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37, 37,
	35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35,
	42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42,
	43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43, 43,
	61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61, 61,
	58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58, 58,
	126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126,
	45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45,
	46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46, 46,
	32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32,
}
