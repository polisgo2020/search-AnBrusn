package index

import (
	"strings"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

type FileWithFreq struct {
	Filename string
	Freq     int
}

type Index map[string][]FileWithFreq

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
