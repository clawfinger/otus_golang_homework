package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type unpackState uint

const (
	EXPECTING_NON_DIGIT unpackState = 0
	NORMAL              unpackState = 1
)

func Unpack(source string) (string, error) {
	if len(source) == 0 {
		return "", nil
	}
	var sb strings.Builder
	var previous rune
	var state unpackState = EXPECTING_NON_DIGIT
	for _, char := range source {
		switch state {
		case EXPECTING_NON_DIGIT:
			if unicode.IsDigit(char) {
				return "", ErrInvalidString
			} else {
				state = NORMAL
			}
		case NORMAL:
			if unicode.IsDigit(char) {
				currentAsDigit, err := strconv.Atoi(string(char))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(previous), currentAsDigit))
				state = EXPECTING_NON_DIGIT
			} else {
				sb.WriteRune(previous)
			}
		}
		previous = char
	}
	if !unicode.IsDigit(previous) {
		sb.WriteRune(previous)
	}

	return sb.String(), nil
}
