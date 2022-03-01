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

	bufReader := bufio.NewScanner(r)
	for bufReader.Scan() {
		line := bufReader.Text()
		found := strings.Contains(line, searchString)

		if found {
			var user User
			if err := jsoniter.Unmarshal([]byte(line), &user); err != nil {
				return result, err
			}

			result[strings.ToLower(strings.Split(user.Email, "@")[1])]++
		}
	}
	return result, nil
}
