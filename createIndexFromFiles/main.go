package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/polisgo2020/search-AnBrusn/index"
)

var (
	Re = regexp.MustCompile(`[^\w]+`)
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

func readFile(ctx context.Context, cancel context.CancelFunc, path string, filename string,
	dataChan chan<- [2]string, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		errChan <- err
		cancel()
		return
	}
	words := Re.Split(string(data), -1)
	for _, word := range words {
		select {
		case <-ctx.Done():
			return
		default:
			dataChan <- [2]string{word, filename}
		}
	}
}

func createFromDirectory(dirname string) (index.Index, error) {
	invertedIndex := make(index.Index)
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	dataChan := make(chan [2]string)
	defer close(errChan)
	defer close(dataChan)

	ctx, cancel := context.WithCancel(context.Background())

	go listener(ctx, invertedIndex, dataChan, errChan)

	wgForFilesReading := &sync.WaitGroup{}
	wgForFilesReading.Add(len(files))
	go func() {
		wgForFilesReading.Wait()
		cancel()
	}()

	for _, currentFile := range files {
		path := filepath.Join(dirname, currentFile.Name())
		go readFile(ctx, cancel, path, currentFile.Name(), dataChan, errChan, wgForFilesReading)
	}

	for {
		select {
		case err := <-errChan:
			return nil, err
		case <-ctx.Done():
			return invertedIndex, nil
		}
	}
}

func listener(ctx context.Context, invertedIndex index.Index, dataChan chan [2]string, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case wordInfo := <-dataChan:
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
