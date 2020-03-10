package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func createAndWriteInvertedIndex(dirname string, outputFile string) error {
	var invertedIndexes = make(map[string][]string)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, currentFile := range files {
		data, err := ioutil.ReadFile(dirname + "\\" + currentFile.Name())
		if err != nil {
			return err
		}
		wordsInCurrentFile := strings.Fields(string(data))
		for _, word := range wordsInCurrentFile {
			sort.Strings(invertedIndexes[word])
			if sort.SearchStrings(invertedIndexes[word], currentFile.Name()) >= len(invertedIndexes[word]) {
				invertedIndexes[word] = append(invertedIndexes[word], currentFile.Name())
			}
		}
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	indexes, err := json.Marshal(invertedIndexes)
	if err != nil {
		return err
	}
	_, err = file.Write(indexes)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("There must be 2 arguments: path to the folder with input files and output file path")
	}

	err := createAndWriteInvertedIndex(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}
