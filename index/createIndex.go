package index

import (
	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type FileWithFreq struct {
	Filename string
	Freq     int
}

func CreateInvertedIndex(dirname string) (map[string][]FileWithFreq, error) {
	var invertedIndexes = make(map[string][]FileWithFreq)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`[^\w]+`)
	for _, currentFile := range files {
		data, err := ioutil.ReadFile(filepath.Join(dirname, currentFile.Name()))
		if err != nil {
			return nil, err
		}
		wordsInCurrentFile := re.Split(string(data), -1)
		for _, word := range wordsInCurrentFile {
			stemWord := stopwords.CleanString(word, "en", true)
			if len(stemWord) > 0 {
				stemWord, err = snowball.Stem(word, "english", true)
				if err != nil {
					return nil, err
				}
				isFound := false
				for i, el := range invertedIndexes[stemWord] {
					if el.Filename == currentFile.Name() {
						invertedIndexes[stemWord][i].Freq++
						isFound = true
						break
					}
				}
				if !isFound {
					invertedIndexes[stemWord] = append(invertedIndexes[stemWord], FileWithFreq{currentFile.Name(), 1})
				}
			}
		}
	}
	return invertedIndexes, nil
}
