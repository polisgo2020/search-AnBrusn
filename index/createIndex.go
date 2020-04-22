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
	"context"
	"fmt"
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

// FileWithFreq contains the name of the file and frequency og the token in it.
type FileWithFreq struct {
	Filename string
	Freq     int
}

type Index struct {
	Data     map[string][]FileWithFreq
	ErrChan  chan error
	DataChan chan [2]string
}

func NewIndex() Index {
	return Index{
		Data:     make(map[string][]FileWithFreq),
		ErrChan:  make(chan error),
		DataChan: make(chan [2]string),
	}
}

// GetTokenFromWord extracts token
func GetTokenFromWord(word string) (string, error) {
	stemWord := stopwords.CleanString(word, "en", true)
	stemWord = strings.TrimSpace(stemWord)
	return snowball.Stem(stemWord, "english", true)
}

// AddToken adds tokens in inverted index.
func (index Index) AddToken(word string, filename string) error {
	stemWord, err := GetTokenFromWord(word)
	if err != nil {
		return fmt.Errorf("can not extract token from word %w", err)
	}
	if len(stemWord) > 0 {
		isFound := false
		for i, el := range index.Data[stemWord] {
			if el.Filename == filename {
				index.Data[stemWord][i].Freq++
				isFound = true
				break
			}
		}
		if !isFound {
			index.Data[stemWord] = append(index.Data[stemWord], FileWithFreq{filename, 1})
		}
	}
	return nil
}

// Listener listens channel of words and adds them in index
func (index Index) Listener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case wordInfo := <-index.DataChan:
			if err := index.AddToken(wordInfo[0], wordInfo[1]); err != nil {
				index.ErrChan <- err
			}
		}
	}
}
