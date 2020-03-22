package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/polisgo2020/search-AnBrusn/index"
)

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func writeInvertedIndex(outputFile string, invertedIndexes index.Index) error {
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

func readFromDirectory(dirname string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	dataFromFiles := make(map[string]string)
	for _, currentFile := range files {
		data, err := ioutil.ReadFile(filepath.Join(dirname, currentFile.Name()))
		if err != nil {
			return nil, err
		}
		dataFromFiles[currentFile.Name()] = string(data)
	}
	return dataFromFiles, nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("There must be 2 arguments: path to the folder with input files and output file path")
	}

	data, err := readFromDirectory(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	invertedIndex, err := index.CreateInvertedIndex(data)
	if err != nil {
		log.Fatal(err)
	}
	if err := writeInvertedIndex(os.Args[2], invertedIndex); err != nil {
		log.Fatal(err)
	}
}
