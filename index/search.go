package index

import (
	"sort"
	"strings"
	"unicode"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

func getTokensFromInput(inpStr string) ([]string, error) {
	var userInput []string
	wordsInInput := strings.FieldsFunc(inpStr, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, word := range wordsInInput {
		stemWord := stopwords.CleanString(word, "en", true)
		stemWord = strings.TrimSpace(stemWord)
		if len(stemWord) > 0 {
			var err error
			stemWord, err = snowball.Stem(word, "english", true)
			if err != nil {
				return nil, err
			}
			userInput = AppendIfMissing(userInput, stemWord)
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
			return nil, nil
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
