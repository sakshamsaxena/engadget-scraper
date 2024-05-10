package processor

import (
	"regexp"

	"github.com/sakshamsaxena/engadget-scraper/cache"
)

func sanitize(words []string) []string {
	validWords := make([]string, 0)
	for _, word := range words {
		if isValid(word) && cache.CheckWordBank(word) {
			validWords = append(validWords, word)
		}
	}
	return validWords
}

func isValid(word string) bool {
	// Word should have 3 or more characters.
	if len(word) < 3 {
		return false
	}

	// Word should have alphabetical characters.
	// NOTE: I've included special characters as well so that,
	// for example, if "Schrödinger" is one of the popular word
	// it should not get missed. To support this search, I've used
	// Unicode letter class (\p{L}) and accents class (\p{Mn}).
	// Ref: https://www.compart.com/en/unicode/category
	matches := regexp.MustCompile(`[\p{L}\p{Mn}]+`).FindAllString(word, -1)
	return len(matches) == 1 // entire word should match
}
