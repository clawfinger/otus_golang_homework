package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomainsOptimized(r, domain)
}

func countDomainsOptimized(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	searchString := "." + domain

	bufReader := bufio.NewReaderSize(r, 10000)
	for {
		line, readErr := bufReader.ReadString(10)
		var user User
		if err := jsoniter.Unmarshal([]byte(line), &user); err != nil {
			return result, err
		}

		matched := strings.Contains(user.Email, searchString)

		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}

		if readErr == io.EOF {
			break
		}
	}
	return result, nil
}
