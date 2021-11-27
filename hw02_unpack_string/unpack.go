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
	var state previousSymbolState = INITIAL
	var needLast bool = true
	for _, char := range source {
		needLast = true
		switch state {
		case INITIAL:
			if unicode.IsDigit(char) {
				return "", ErrInvalidString
			} else if char == '\\' {
				state = ESCAPE
			} else {
				state = LETTER
			}
		case DIGIT:
			if unicode.IsDigit(char) {
				return "", ErrInvalidString
			} else if char == '\\' {
				sb.WriteRune(previous)
				state = ESCAPE
			} else {
				state = LETTER
			}
		case LETTER:
			if unicode.IsDigit(char) {
				number, err := strconv.Atoi(string(char))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(previous), number))
				state = DIGIT
				/* we need to add the last symbol to the result in all situations
				except when the last operation on sequence was unpacking */
				needLast = false
			} else if char == '\\' {
				sb.WriteRune(previous)
				state = ESCAPE
			} else {
				sb.WriteRune(previous)
				state = LETTER
			}
		case ESCAPE:
			if unicode.IsDigit(char) {
				state = LETTER
			} else if char == '\\' {
				state = LETTER
			} else {
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
