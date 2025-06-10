package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	runes := []rune(str)
	var result []rune
	var prev rune
	escaped := false

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if escaped {
			if r != '\\' && !unicode.IsDigit(r) {
				return "", ErrInvalidString
			}
			result = append(result, r)
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
			switch {
			case repeatCount > 0:
				for j := 0; j < repeatCount-1; j++ {
					result = append(result, prev)
				}
			case len(result) > 0:
				result = result[:len(result)-1]
			default:
				return "", ErrInvalidString
			}

			prev = 0
			continue
		}

		result = append(result, r)
		prev = r
	}

	if escaped {
		return "", ErrInvalidString
	}

	return string(result), nil
}
