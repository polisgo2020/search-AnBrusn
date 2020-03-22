package index

import (
	"regexp"
	"sort"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
	"github.com/polisgo2020/search-AnBrusn/utils"
)

func getTokensFromInput(inpStr string) ([]string, error) {
	var userInput []string
	re := regexp.MustCompile(`[^\w]+`)
	wordsInInput := re.Split(inpStr, -1)
	for _, word := range wordsInInput {
		stemWord := stopwords.CleanString(word, "en", true)
		if len(stemWord) > 0 {
			var err error
			stemWord, err = snowball.Stem(word, "english", true)
			if err != nil {
				return nil, err
			}
			userInput = utils.AppendIfMissing(userInput, stemWord)
		}
	}
	return userInput, nil
}

func (index Index) FindInIndex(userInput string) ([]FileWithFreq, error) {
	inputTokens, err := getTokensFromInput(userInput)
	if err != nil {
		return nil, err
	}
	var searchResults []FileWithFreq
	for i, token := range inputTokens {
		filesWithToken, ok := index[token]
		if !ok {
			break
		}
		if i == 0 {
			searchResults = filesWithToken
		} else {
			searchResults = getIntersection(searchResults, filesWithToken)
		}
	}
	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].Freq > searchResults[j].Freq
	})
	return searchResults, nil
}

func getIntersection(slice1, slice2 []FileWithFreq) (intersection []FileWithFreq) {
	m := make(map[string]int)
	for _, el := range slice1 {
		m[el.Filename] = el.Freq
	}
	for _, el := range slice2 {
		if num, ok := m[el.Filename]; ok {
			intersection = append(intersection, FileWithFreq{el.Filename, el.Freq + num})
		}
	}
	return
}
