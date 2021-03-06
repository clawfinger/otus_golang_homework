package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordData struct {
	word  string
	count int
}

var matcher = regexp.MustCompile(`\p{L}+-*\p{L}*`)

func Top10(source string) []string {
	const NUMBER = 10
	dict := make(map[string]int)
	splitted := strings.Fields(source)
	for _, word := range splitted {
		matched := matcher.FindString(strings.ToLower(word))
		if len(matched) > 0 {
			dict[matched]++
		}
	}
	prepslice := make([]wordData, 0, len(dict))

	for word, count := range dict {
		prepslice = append(prepslice, wordData{word, count})
	}

	sort.Slice(prepslice, func(i, j int) bool {
		if prepslice[i].count == prepslice[j].count {
			return prepslice[i].word < prepslice[j].word
		}
		return prepslice[i].count > prepslice[j].count
	})

	end := NUMBER
	if len(prepslice) < NUMBER {
		end = len(prepslice)
	}

	res := make([]string, 0, end)

	for _, entry := range prepslice[:end] {
		res = append(res, entry.word)
	}
	return res
}
