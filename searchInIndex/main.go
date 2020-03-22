package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/polisgo2020/search-AnBrusn/index"
)

func readUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func readIndexFromFile(indexPath string) (index.Index, error) {
	data, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}
	var index = map[string][]index.FileWithFreq{}
	if er := json.Unmarshal(data, &index); er != nil {
		return nil, err
	}
	return index, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("There must be 1 arguments: path to the index file")
	}

	invertedIndex, er := readIndexFromFile(os.Args[1])
	if er != nil {
		log.Fatal(er)
	}
	searchResults, er := invertedIndex.FindInIndex(readUserInput())
	if er != nil {
		log.Fatal(er)
	}
	if len(searchResults) == 0 {
		fmt.Println("No results")
	} else {
		for _, el := range searchResults {
			fmt.Printf("%s (%d words were found)\n", el.Filename, el.Freq)
		}
	}
}
