package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	runes := []rune(str)
	var result strings.Builder
	var prev rune
	escaped := false

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if escaped {
			if r != '\\' && !unicode.IsDigit(r) {
				return "", ErrInvalidString
			}
			result.WriteRune(r)
			prev = r
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if unicode.IsDigit(r) {
			if prev == 0 {
				return "", ErrInvalidString
			}
			repeatCount := int(r - '0')
			if repeatCount > 0 {
				result.WriteString(strings.Repeat(string(prev), repeatCount-1))
			} else {
				currentStr := result.String()
				result.Reset()
				result.WriteString(currentStr[:len(currentStr)-1])
			}
			prev = 0
			continue
		}

		result.WriteRune(r)
		prev = r
	}

	if escaped {
		return "", ErrInvalidString
	}

	return result.String(), nil
}
