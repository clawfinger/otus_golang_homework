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
	EXPECTING_LETTER unpackState = 0
	NORMAL           unpackState = 1
)

func Unpack(source string) (string, error) {
	if len(source) == 0 {
		return "", nil
	}
	var sb strings.Builder
	var previous rune
	var state unpackState = EXPECTING_LETTER
	for _, char := range source {
		switch state {
		case EXPECTING_LETTER:
			if unicode.IsDigit(char) {
				return "", ErrInvalidString
			} else {
				state = NORMAL
			}

		case NORMAL:
			if unicode.IsDigit(char) {
				currentAsString := string(char)
				currentAsDigit, err := strconv.Atoi(currentAsString)
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(previous), currentAsDigit))
				state = EXPECTING_LETTER
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
