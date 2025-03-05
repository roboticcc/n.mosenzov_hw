package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	words := strings.Fields(input)
	frequency := make(map[string]int)

	for _, word := range words {
		if word == "-" {
			continue
		}

		normalized := normalizeWord(word)
		if normalized != "" {
			frequency[normalized]++
		}
	}

	uniqueWords := make([]string, 0, len(frequency))
	for word := range frequency {
		uniqueWords = append(uniqueWords, word)
	}

	sort.Slice(uniqueWords, func(i, j int) bool {
		if frequency[uniqueWords[i]] == frequency[uniqueWords[j]] {
			return uniqueWords[i] < uniqueWords[j]
		}
		return frequency[uniqueWords[i]] > frequency[uniqueWords[j]]
	})

	if len(uniqueWords) > 10 {
		return uniqueWords[:10]
	}
	return uniqueWords
}

func normalizeWord(word string) string {
	word = strings.ToLower(word)

	if strings.Count(word, "-") == len(word) {
		if len(word) > 1 {
			return word
		}
		return ""
	}

	word = strings.Trim(word, "!.,;:\"'()[] ")

	if len(word) == 0 {
		return ""
	}

	return word
}
