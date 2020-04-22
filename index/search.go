package index

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// GetTokensFromInput extracts tokens from user input.
func GetTokensFromInput(inpStr string) ([]string, error) {
	var userInput []string
	wordsInInput := strings.FieldsFunc(inpStr, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, word := range wordsInInput {
		stemWord, err := GetTokenFromWord(word)
		if err != nil {
			return nil, fmt.Errorf("can not extract token from word %w", err)
		}
		if len(stemWord) > 0 {
			userInput = AppendIfMissing(userInput, stemWord)
		}
	}
	return userInput, nil
}

// FindInIndex searches query over inverted index.
// Search results are ranged by amount of found tokens.
func (index Index) FindInIndex(userInput string) ([]FileWithFreq, error) {
	inputTokens, err := GetTokensFromInput(userInput)
	if err != nil {
		return nil, fmt.Errorf("can not get tokens from user input %w", err)
	}
	var searchResults []FileWithFreq
	for i, token := range inputTokens {
		filesWithToken, ok := index.Data[token]
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

// getIntersection of two slices of filenames and frequencies
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
