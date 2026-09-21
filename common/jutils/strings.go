package jutils

// CutBytes splits a string into two parts: the first part contains at most max bytes,
// and the second part contains the remaining bytes.
// It ensures that multi-byte characters are not split in the middle.
func CutBytes(s string, max int) (before, after string) {
	if len(s) <= max {
		return s, ""
	}

	runes := []rune(s)

	var currentLength, i int
	for ; i < len(runes); i++ {
		runeLength := len(string(runes[i]))
		if currentLength+runeLength > max {
			break
		}
		currentLength += runeLength
	}

	return string(runes[:i]), string(runes[i:])
}
