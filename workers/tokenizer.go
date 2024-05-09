package workers

import (
	"github.com/sakshamsaxena/engadget-scraper/cache"
	"regexp"
)

func tokenize(input []string) []string {
	ans := make([]string, 0)
	for _, s := range input {
		if check(s) && cache.CheckBank(s) {
			ans = append(ans, s)
		}
	}
	return ans
}

func check(str string) bool {
	if len(str) < 3 {
		return false
	}
	// unicode letter class and accents class
	var re = regexp.MustCompile(`[\p{L}\p{Mn}]+`)
	all := re.FindAllString(str, -1)
	if len(all) > 1 {
		return false // no exact match
	}
	return true
}
