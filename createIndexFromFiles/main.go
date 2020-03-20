package main

import (
	"encoding/json"
	"github.com/polisgo2020/search-AnBrusn/index"
	"log"
	"os"
)

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func writeInvertedIndex(outputFile string, invertedIndexes map[string][]index.FileWithFreq) error {
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

func main() {
	if len(os.Args) < 3 {
		log.Fatal("There must be 2 arguments: path to the folder with input files and output file path")
	}

	index, err := index.CreateInvertedIndex(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := writeInvertedIndex(os.Args[2], index); err != nil {
		log.Fatal(err)
	}
}
