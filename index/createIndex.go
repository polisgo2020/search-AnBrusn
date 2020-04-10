/*
Package index implements inverted index, search over the built index.

Usage

New token can be added with AddToken function that extract token from word and add it to inverted index.
Example:

	err := i.AddToken(" word", "sourceFile.txt")

To search over the index use FindInIndex function.

	searchResults, err := i.FindInIndex("this is search query")

Search results are ranged by amount of found tokens.
*/
package index

import (
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

// FileWithFreq contains the name of the file and frequency og the token in it.
type FileWithFreq struct {
	Filename string
	Freq     int
}

type Index map[string][]FileWithFreq

// AddToken extracts token from word and adds it in inverted index.
func (index Index) AddToken(word string, filename string) error {
	stemWord := stopwords.CleanString(word, "en", true)
	stemWord = strings.TrimSpace(stemWord)
	if len(stemWord) > 0 {
		stemWord, err := snowball.Stem(stemWord, "english", true)
		if err != nil {
			return err
		}
		isFound := false
		for i, el := range index[stemWord] {
			if el.Filename == filename {
				index[stemWord][i].Freq++
				isFound = true
				break
			}
		}
		if !isFound {
			index[stemWord] = append(index[stemWord], FileWithFreq{filename, 1})
		}
	}
	return nil
}
