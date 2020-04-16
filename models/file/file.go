package file

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/polisgo2020/search-AnBrusn/index"
	"github.com/rs/zerolog/log"
)

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Err(err)
	}
}

// ReadIndex reads and decodes inverted index.
func ReadIndex(indexPath string) (index.Index, error) {
	data, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}
	var invertedIndex = index.Index{}
	if er := json.Unmarshal(data, &invertedIndex); er != nil {
		return nil, err
	}
	return invertedIndex, nil
}

// WriteInvertedIndex encodes index with json and writes it in output file.
func WriteInvertedIndex(outputFile string, invertedIndexes index.Index) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer closeFile(file)
	indexes, err := json.Marshal(invertedIndexes)
	if err != nil {
		return err
	}
	if _, err = file.Write(indexes); err != nil {
		return err
	}
	return nil
}
