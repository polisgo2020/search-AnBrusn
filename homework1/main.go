package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"polisAnnaBrusnitsyna/utils"
	"regexp"
)

type FileWithFreq struct {
	Filename string
	Freq int
}

func writeInvertedIndex(outputFile string, invertedIndexes map[string][]FileWithFreq) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	indexes, err := json.Marshal(invertedIndexes)
	if err != nil {
		return err
	}
	if _, err = file.Write(indexes); err != nil {
		return err
	}
	if err = file.Close();  err != nil {
		return err
	}
	return nil
}

func createInvertedIndex(dirname string, outputFile string) error {
	var invertedIndexes = make(map[string][]FileWithFreq)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, currentFile := range files {
		data, err := ioutil.ReadFile(dirname + "\\" + currentFile.Name())
		if err != nil {
			return err
		}
		re := regexp.MustCompile(`[^\w]+`)
		wordsInCurrentFile := re.Split(string(data), -1)
		for _, word := range wordsInCurrentFile {
			stemWord := utils.Stem(word)
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
	return writeInvertedIndex(outputFile, invertedIndexes)
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("There must be 2 arguments: path to the folder with input files and output file path")
	}

	err := createInvertedIndex(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}
