package file

import (
	"encoding/json"
	"fmt"
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
		return index.Index{}, fmt.Errorf("can not read file %w", err)
	}
	var invertedIndex = index.NewIndex()
	if er := json.Unmarshal(data, &invertedIndex.Data); er != nil {
		return index.Index{}, fmt.Errorf("can not decode data %w", err)
	}
	return invertedIndex, nil
}

// WriteInvertedIndex encodes index with json and writes it in output file.
func WriteInvertedIndex(outputFile string, invertedIndexes index.Index) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("can not create output file %w", err)
	}
	defer closeFile(file)
	indexes, err := json.Marshal(invertedIndexes.Data)
	if err != nil {
		return fmt.Errorf("can not code data %w", err)
	}
	if _, err = file.Write(indexes); err != nil {
		return fmt.Errorf("can not write index %w", err)
	}
	return nil
}
