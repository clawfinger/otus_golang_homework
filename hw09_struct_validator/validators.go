package hw09structvalidator

import (
	"regexp"
	"strconv"
	"strings"
)

func validateIntIn(v int, rule string) error {
	values := strings.Split(rule, ",")
	if len(values) == 0 {
		return errBadValidator
	}
	for _, entry := range values {
		intEntry, err := strconv.Atoi(entry)
		if err != nil {
			return errBadValidator
		}
		if intEntry == v {
			return nil
		}
	}

	return errIntNotInSet
}

func validateIntMin(v int, rule string) error {
	min, err := strconv.Atoi(rule)
	if err != nil {
		return errBadValidator
	}
	if v < min {
		return errIntLess
	}
	return nil
}

func validateIntMax(v int, rule string) error {
	max, err := strconv.Atoi(rule)
	if err != nil {
		return errBadValidator
	}
	if v > max {
		return errIntGreater
	}
	return nil
}

func validateStringIn(v string, rule string) error {
	values := strings.Split(rule, ",")
	if len(values) == 0 {
		return errBadValidator
	}
	for _, entry := range values {
		if v == entry {
			return nil
		}
	}
	return errStringNotInSet
}

func validateStringRegexp(v string, rule string) error {
	rExpr, err := regexp.Compile(rule)
	if err != nil {
		return errBadValidator
	}
	if rExpr.MatchString(v) {
		return nil
	}
	return errStringRegexp
}

func validateStringLen(v string, rule string) error {
	length, err := strconv.Atoi(rule)
	if err != nil {
		return errBadValidator
	}
	if len(v) != length {
		return errStringLen
	}
	return nil
}
