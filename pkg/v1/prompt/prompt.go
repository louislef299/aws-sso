package prompt

import (
	"log"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
)

const LINEPREFIX = "#"

func Select(label string, toSelect []string, searcher func(input string, index int) bool) (int, string) {
	prompt := promptui.Select{
		Label:             label,
		Items:             toSelect,
		Size:              30,
		Searcher:          searcher,
		StartInSearchMode: true,
	}
	index, value, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	return index, value
}

func Prompt(label string, dfault string) string {
	prompt := promptui.Prompt{
		Label:     label,
		Default:   dfault,
		AllowEdit: false,
	}
	val, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func FuzzySearchWithPrefixAnchor(itemsToSelect []string) func(input string, index int) bool {
	return func(input string, index int) bool {
		role := itemsToSelect[index]
		if strings.HasPrefix(input, LINEPREFIX) {
			return strings.HasPrefix(role, input)
		} else {
			if fuzzy.MatchFold(input, role) {
				return true
			}
		}
		return false
	}
}
