package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/polisgo2020/search-AnBrusn/index"
	"github.com/urfave/cli/v2"
)

var searchForm = []byte(`
<html>
	<body>
	<form action="/searchInIndex" method="get">
		Search: <input type="text" name="userInput">
		<input type="submit" value="Search">
	</form>
	</body>
</html>
`)

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

func readFile(ctx context.Context, path string, filename string,
	dataChan chan<- [2]string, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		errChan <- err
		return
	}
	words := strings.FieldsFunc(string(data), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
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
		go readFile(ctx, path, currentFile.Name(), dataChan, errChan, wgForFilesReading)
	}

	for {
		select {
		case err := <-errChan:
			cancel()
			return nil, err
		case <-ctx.Done():
			return invertedIndex, nil
		}
	}
}

func searchWithInputFromStdin(indexFile string) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userInput := scanner.Text()
	invertedIndex, err := readIndexFromFile(indexFile)
	if err != nil {
		return err
	}
	searchResults, err := invertedIndex.FindInIndex(userInput)
	if err != nil {
		return err
	}
	if len(searchResults) == 0 {
		fmt.Println("No results")
	} else {
		for _, el := range searchResults {
			fmt.Printf("%s (%d words were found)\n", el.Filename, el.Freq)
		}
	}
	return nil
}

func searchWithInputFromHttp(server string, indexFile string) error {
	srv := &http.Server{Addr: server}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(searchForm)
	})
	http.HandleFunc("/searchInIndex", func(w http.ResponseWriter, r *http.Request) {
		userInput := r.FormValue("userInput")
		invertedIndex, err := readIndexFromFile(indexFile)
		if err != nil {
			http.Error(w, "reading index error", http.StatusInternalServerError)
			return
		}
		searchResults, err := invertedIndex.FindInIndex(userInput)
		if err != nil {
			http.Error(w, "searching error", http.StatusInternalServerError)
			return
		}
		if len(searchResults) == 0 {
			fmt.Fprintln(w, "No results")
		} else {
			for _, el := range searchResults {
				fmt.Fprintf(w, "%s (%d words were found)\n", el.Filename, el.Freq)
			}
		}
	})
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func searchInIndex(c *cli.Context) error {
	if c.String("http") == "" {
		if err := searchWithInputFromStdin(c.String("index")); err != nil {
			return err
		}
	} else {
		if err := searchWithInputFromHttp(c.String("http"), c.String("index")); err != nil {
			return err
		}
	}
	return nil
}

func readIndexFromFile(indexPath string) (index.Index, error) {
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

func createIndexFromFiles(c *cli.Context) error {
	invertedIndex, err := createFromDirectory(c.String("dir"))
	if err != nil {
		return err
	}
	if err := writeInvertedIndex(c.String("index"), invertedIndex); err != nil {
		return err
	}
	return nil
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
	app := &cli.App{
		Name:  "Index",
		Usage: "Create index from directory and search in index",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "Create index from directory",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dir",
						Usage:    "Path to directory with input files",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "index",
						Usage:    "Path to output index file",
						Required: true,
					},
				},
				Action: createIndexFromFiles,
			},
			{
				Name:  "search",
				Usage: "Search in index",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "index",
						Usage:    "Path to index file",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "http",
						Usage: "Input from http",
					},
				},
				Action: searchInIndex,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
