package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"polisAnnaBrusnitsyna/utils"
	"regexp"
	"sort"
)

type FileWithFreq struct {
	Filename string
	Freq int
}

func readUserInput() []string {
	var userInput []string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inpStr := scanner.Text()
	re := regexp.MustCompile(`[^\w]+`)
	wordsInInput := re.Split(inpStr, -1)
	for _, word := range wordsInInput {
		stemWord := utils.Stem(word)
		userInput = utils.AppendIfMissing(userInput, stemWord)
	}
	return userInput
}

func readIndexFromFile(indexPath string) (map[string][]FileWithFreq, error) {
	data, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}
	var index = map[string][]FileWithFreq{}
	if er := json.Unmarshal(data, &index); er != nil {
		return nil, err
	}
	return index, nil
}

func findInIndex(index map[string][]FileWithFreq, userInput []string) error {
	var searchResults []FileWithFreq
	for i, token := range userInput {
		filesWithToken, ok := index[token]
		if !ok {
			fmt.Println("No results")
			return nil
		}
		if i == 0 {
			searchResults = filesWithToken
		} else {
			searchResults = getIntersection(searchResults, filesWithToken)
		}
	}
	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].Freq > searchResults[j].Freq
	})
	for _, el := range searchResults {
		fmt.Printf("%s (%d words were found)\n", el.Filename, el.Freq)
	}
	return nil
}

func getIntersection(slice1, slice2 []FileWithFreq) (intersection []FileWithFreq) {
	m := make(map[string]int)
	for _, el := range slice1 {
		m[el.Filename] = el.Freq
	}
	for _, el := range slice2 {
		if num, ok := m[el.Filename]; ok {
			intersection = append(intersection, FileWithFreq{el.Filename, el.Freq + num})
		}
	}
	return intersection
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("There must be 1 arguments: path to the index file")
	}

	userInput := readUserInput()
	index, er := readIndexFromFile(os.Args[1])
	if er != nil {
		log.Fatal(er)
	}
	if er := findInIndex(index, userInput); er != nil {
		log.Fatal(er)
	}
}
