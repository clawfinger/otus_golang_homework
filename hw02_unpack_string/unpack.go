package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type previousSymbolState uint

const (
	INITIAL previousSymbolState = 0
	DIGIT   previousSymbolState = 1
	LETTER  previousSymbolState = 2
	ESCAPE  previousSymbolState = 3
)

func Unpack(source string) (string, error) {
	if len(source) == 0 {
		return "", nil
	}
	var sb strings.Builder
	var previous rune
	state := INITIAL
	needLast := true
	for _, char := range source {
		needLast = true
		switch state {
		case INITIAL:
			switch {
			case unicode.IsDigit(char):
				return "", ErrInvalidString
			case char == '\\':
				state = ESCAPE
			default:
				state = LETTER
			}
		case DIGIT:
			switch {
			case unicode.IsDigit(char):
				return "", ErrInvalidString
			case char == '\\':
				sb.WriteRune(previous)
				state = ESCAPE
			default:
				state = LETTER
			}
		case LETTER:
			switch {
			case unicode.IsDigit(char):
				number, err := strconv.Atoi(string(char))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(previous), number))
				state = DIGIT
				/* we need to add the last symbol to the result in all situations
				except when the last operation on sequence was unpacking */
				needLast = false
			case char == '\\':
				sb.WriteRune(previous)
				state = ESCAPE
			default:
				sb.WriteRune(previous)
				state = LETTER
			}
		case ESCAPE:
			switch {
			case unicode.IsDigit(char):
				state = LETTER
			case char == '\\':
				state = LETTER
			default:
				return "", ErrInvalidString
			}
		}
		previous = char
	}
	if needLast {
		sb.WriteRune(previous)
	}
	return sb.String(), nil
}
