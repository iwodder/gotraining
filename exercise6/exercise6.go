package exercise7

import (
	"regexp"
	"unicode"
)

func CamelCase(s string) int {
	if s == "" {
		return 0
	}
	re := regexp.MustCompile("[A-Z]")
	words := re.Split(s, -1)
	return len(words)
}

func CaesarCipher(s string, k int) string {
	if s == "" {
		return s
	}
	ret := make([]byte, len(s))
	for i, char := range s {
		if unicode.IsLetter(char) {
			val := uint8(char) + (uint8(k) % 26)
			if val > 'z' || (val > 'Z' && val < 'a') {
				val -= 26
			}
			ret[i] = val
		} else {
			ret[i] = uint8(char)
		}

	}
	return string(ret)
}
