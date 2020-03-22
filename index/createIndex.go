package index

import (
	"regexp"

	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
)

type FileWithFreq struct {
	Filename string
	Freq     int
}

type Index map[string][]FileWithFreq

func CreateInvertedIndex(data map[string]string) (Index, error) {
	var invertedIndexes = make(Index)
	re := regexp.MustCompile(`[^\w]+`)
	for fileName, fileData := range data {
		wordsInCurrentFile := re.Split(fileData, -1)
		for _, word := range wordsInCurrentFile {
			stemWord := stopwords.CleanString(word, "en", true)
			if len(stemWord) > 0 {
				stemWord, err := snowball.Stem(stemWord, "english", true)
				if err != nil {
					return nil, err
				}
				isFound := false
				for i, el := range invertedIndexes[stemWord] {
					if el.Filename == fileName {
						invertedIndexes[stemWord][i].Freq++
						isFound = true
						break
					}
				}
				if !isFound {
					invertedIndexes[stemWord] = append(invertedIndexes[stemWord], FileWithFreq{fileName, 1})
				}
			}
		}
	}
	return invertedIndexes, nil
}
