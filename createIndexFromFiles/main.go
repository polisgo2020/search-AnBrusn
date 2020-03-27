package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

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

func readFile(path string, filename string, dataChan chan<- [2]string, errChan chan<- error,
	wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		errChan <- err
	}
	re := regexp.MustCompile(`[^\w]+`)
	words := re.Split(string(data), -1)
	for _, word := range words {
		mutex.Lock()
		dataChan <- [2]string{word, filename}
		mutex.Unlock()
	}
}

func createFromDirectory(dirname string) (index.Index, error) {
	invertedIndex := make(index.Index)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return invertedIndex, err
	}

	errChan := make(chan error)
	dataChan := make(chan [2]string)
	mutex := &sync.Mutex{}

	wgForListener := &sync.WaitGroup{}
	wgForListener.Add(1)
	go listener(invertedIndex, dataChan, errChan, wgForListener)

	wgForFilesReading := &sync.WaitGroup{}
	wgForFilesReading.Add(len(files))
	go func() {
		wgForFilesReading.Wait()
		close(dataChan)
	}()

	for _, currentFile := range files {
		path := filepath.Join(dirname, currentFile.Name())
		go readFile(path, currentFile.Name(), dataChan, errChan, wgForFilesReading, mutex)
	}

	wgForListener.Wait()
	if err, ok := <-errChan; ok {
		close(errChan)
		return invertedIndex, err
	}
	return invertedIndex, nil
}

func listener(invertedIndex index.Index, dataChan chan [2]string, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-errChan:
			close(dataChan)
			return
		case wordInfo, ok := <-dataChan:
			if !ok {
				close(errChan)
				return
			}
			if err := invertedIndex.AddToken(wordInfo[0], wordInfo[1]); err != nil {
				errChan <- err
			}
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("There must be 2 arguments: path to the folder with input files and output file path")
	}

	invertedIndex, err := createFromDirectory(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := writeInvertedIndex(os.Args[2], invertedIndex); err != nil {
		log.Fatal(err)
	}
}
